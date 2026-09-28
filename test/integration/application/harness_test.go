//go:build integration

package application_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const windowsOS = "windows"

// responseSnapshot owns a fully consumed response; request closes the network body.
type responseSnapshot struct {
	StatusCode int
	Header     http.Header
	Body       *bytes.Reader
}

type metricRequest struct {
	Path           string
	ContentType    string
	Body           []byte
	ResponseStatus int
}

type metricReceiver struct {
	status   atomic.Int32
	mu       sync.Mutex
	requests []metricRequest
}

func newMetricReceiver() *metricReceiver {
	receiver := &metricReceiver{}
	receiver.status.Store(http.StatusOK)

	return receiver
}

func (receiver *metricReceiver) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	body, _ := io.ReadAll(request.Body)
	status := int(receiver.status.Load())
	receiver.mu.Lock()
	receiver.requests = append(receiver.requests, metricRequest{
		Path: request.URL.Path, ContentType: request.Header.Get("Content-Type"), Body: body, ResponseStatus: status,
	})
	receiver.mu.Unlock()
	writer.WriteHeader(status)
}

func (receiver *metricReceiver) snapshot() []metricRequest {
	receiver.mu.Lock()
	defer receiver.mu.Unlock()

	return append([]metricRequest(nil), receiver.requests...)
}

type applicationHarness struct {
	baseURL     string
	client      *http.Client
	metrics     *metricReceiver
	mu          sync.Mutex
	logs        []string
	binary      string
	environment []string
}

func startApplication(t *testing.T) *applicationHarness {
	t.Helper()

	ctx := context.Background()

	container, err := testpostgres.Run(ctx, "postgres:17.6-bookworm",
		testpostgres.WithDatabase("serega"),
		testpostgres.WithUsername("serega"),
		testpostgres.WithPassword("serega"),
		testpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("start PostgreSQL application integration container: %v", err)
	}

	databaseURL, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)

		t.Fatalf("read PostgreSQL application integration connection string: %v", err)
	}

	receiver := newMetricReceiver()
	metricServer := httptest.NewServer(receiver)
	httpAddress := reserveAddress(t)
	root := repositoryRoot(t)

	binary := os.Getenv("SEREGA_BINARY")
	if binary == "" {
		binary = filepath.Join(root, "bin", executableName("serega"))
	}

	if _, err := os.Stat(binary); err != nil {
		metricServer.Close()

		_ = container.Terminate(ctx)

		t.Fatalf("built Serega binary is unavailable at %s: %v", binary, err)
	}

	environment := append(os.Environ(),
		"DATABASE_URL="+databaseURL,
		"SEREGA_HTTP_ADDR="+httpAddress,
		"SEREGA_LOG_FORMAT=json",
		"SEREGA_COOKIE_SECURE=false",
		"SEREGA_SHUTDOWN_TIMEOUT=5s",
		"OTEL_METRICS_EXPORTER=otlp",
		"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT="+metricServer.URL+"/v1/metrics",
		"OTEL_EXPORTER_OTLP_METRICS_PROTOCOL=http/protobuf",
		"OTEL_METRIC_EXPORT_INTERVAL=500",
		"OTEL_METRIC_EXPORT_TIMEOUT=1000",
	)

	migrationContext, cancelMigration := context.WithTimeout(ctx, 30*time.Second)
	defer cancelMigration()

	migrate := exec.CommandContext(migrationContext, binary, "migrate", "up")
	migrate.Dir = root

	migrate.Env = environment
	if output, err := migrate.CombinedOutput(); err != nil {
		metricServer.Close()

		_ = container.Terminate(ctx)

		t.Fatalf("migrate application integration database: %v\n%s", err, output)
	}

	command := exec.CommandContext(context.Background(), binary, "serve")
	command.Dir = root
	command.Env = environment

	stderr, err := command.StderrPipe()
	if err != nil {
		metricServer.Close()

		_ = container.Terminate(ctx)

		t.Fatalf("open Serega stderr: %v", err)
	}

	command.Stdout = io.Discard
	if err := command.Start(); err != nil {
		metricServer.Close()

		_ = container.Terminate(ctx)

		t.Fatalf("start Serega: %v", err)
	}

	harness := &applicationHarness{
		baseURL: "http://" + httpAddress,
		client:  &http.Client{Timeout: 2 * time.Second},
		metrics: receiver,
		binary:  binary, environment: environment,
	}

	logsDone := make(chan struct{})
	go func() {
		defer close(logsDone)

		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			harness.mu.Lock()
			harness.logs = append(harness.logs, scanner.Text())
			harness.mu.Unlock()
		}
	}()

	t.Cleanup(func() {
		requestsBeforeShutdown := len(receiver.snapshot())

		waitDone := make(chan error, 1)
		go func() { waitDone <- command.Wait() }()

		var waitErr error

		if runtime.GOOS == windowsOS {
			_ = command.Process.Kill()
			waitErr = <-waitDone
		} else {
			_ = command.Process.Signal(os.Interrupt)

			select {
			case waitErr = <-waitDone:
			case <-time.After(10 * time.Second):
				_ = command.Process.Kill()
				waitErr = <-waitDone

				t.Errorf("Serega did not stop within 10 seconds")
			}
		}

		if runtime.GOOS != windowsOS && waitErr != nil {
			t.Errorf("Serega did not shut down cleanly: %v", waitErr)
		}

		if runtime.GOOS != windowsOS {
			deadline := time.Now().Add(3 * time.Second)
			for time.Now().Before(deadline) && len(receiver.snapshot()) <= requestsBeforeShutdown {
				time.Sleep(50 * time.Millisecond)
			}

			if len(receiver.snapshot()) <= requestsBeforeShutdown {
				t.Errorf("Serega shutdown did not flush metrics")
			}
		}

		<-logsDone
		metricServer.Close()

		if err := container.Terminate(context.Background()); err != nil {
			t.Errorf("terminate PostgreSQL application integration container: %v", err)
		}
	})

	waitFor(t, 20*time.Second, func() bool {
		readiness, err := http.NewRequestWithContext(t.Context(), http.MethodGet, harness.baseURL+"/readyz", nil)
		if err != nil {
			t.Fatal(err)
		}

		response, err := harness.client.Do(readiness)
		if err != nil {
			return false
		}
		defer func() { _ = response.Body.Close() }()

		return response.StatusCode == http.StatusNoContent
	}, "Serega readiness")

	return harness
}

func (application *applicationHarness) logLines() []string {
	application.mu.Lock()
	defer application.mu.Unlock()

	return append([]string(nil), application.logs...)
}

func (application *applicationHarness) logRecords() []map[string]any {
	lines := application.logLines()

	records := make([]map[string]any, 0, len(lines))
	for _, line := range lines {
		var record map[string]any
		if json.Unmarshal([]byte(line), &record) == nil {
			records = append(records, record)
		}
	}

	return records
}

func repositoryRoot(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve application integration source location")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", ".."))
}

func executableName(name string) string {
	if runtime.GOOS == windowsOS {
		return name + ".exe"
	}

	return name
}

func reserveAddress(t *testing.T) string {
	t.Helper()

	var listenConfig net.ListenConfig

	listener, err := listenConfig.Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	address := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}

	return address
}

func waitFor(t *testing.T, timeout time.Duration, condition func() bool, description string) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}

		time.Sleep(50 * time.Millisecond)
	}

	t.Fatalf("timed out waiting for %s", description)
}

func requireJSONResponse(t *testing.T, response *responseSnapshot, status int) map[string]any {
	t.Helper()

	if response.StatusCode != status {
		body, _ := io.ReadAll(response.Body)
		t.Fatalf("status = %d, want %d: %s", response.StatusCode, status, strings.TrimSpace(string(body)))
	}

	var value map[string]any
	if err := json.NewDecoder(response.Body).Decode(&value); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	return value
}

func errorCode(value map[string]any) string {
	nested, _ := value["error"].(map[string]any)
	result, _ := nested["code"].(string)

	return result
}

func findLogRecord(records []map[string]any, eventName, requestID string) map[string]any {
	for _, record := range records {
		if record["event.name"] == eventName && record["request_id"] == requestID {
			return record
		}
	}

	return nil
}

func request(t *testing.T, client *http.Client, method, url string, body io.Reader, headers map[string]string) *responseSnapshot {
	t.Helper()

	httpRequest, err := http.NewRequestWithContext(t.Context(), method, url, body)
	if err != nil {
		t.Fatal(err)
	}

	for name, value := range headers {
		httpRequest.Header.Set(name, value)
	}

	response, err := client.Do(httpRequest)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}

	defer func() { _ = response.Body.Close() }()

	data, err := io.ReadAll(io.LimitReader(response.Body, 2<<20))
	if err != nil {
		t.Fatal(err)
	}

	return &responseSnapshot{StatusCode: response.StatusCode, Header: response.Header, Body: bytes.NewReader(data)}
}

func describeLogs(lines []string) string {
	if len(lines) > 10 {
		lines = lines[len(lines)-10:]
	}

	return fmt.Sprintf("%q", lines)
}
