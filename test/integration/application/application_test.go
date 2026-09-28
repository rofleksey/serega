//go:build integration

package application_test

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

const requestIDHeader = "X-Request-ID"

func TestApplicationBlackBox(t *testing.T) {
	application := startApplication(t)

	t.Run("public health and stable API errors", func(t *testing.T) {
		health := request(t, application.client, http.MethodGet, application.baseURL+"/healthz", nil, nil)

		if health.StatusCode != http.StatusNoContent {
			t.Fatalf("health status = %d, want %d", health.StatusCode, http.StatusNoContent)
		}

		missing := request(t, application.client, http.MethodGet, application.baseURL+"/api/v1/does-not-exist", nil, map[string]string{requestIDHeader: "integration-request-0001"})

		body := requireJSONResponse(t, missing, http.StatusNotFound)
		if missing.Header.Get(requestIDHeader) != "integration-request-0001" || errorCode(body) != "not_found" {
			t.Fatalf("unexpected missing response: headers=%v body=%v", missing.Header, body)
		}
	})

	t.Run("authentication and invalid input fail through public API", func(t *testing.T) {
		response1 := request(t, application.client, http.MethodGet, application.baseURL+"/api/v1/auth/me", nil, nil)

		body1 := requireJSONResponse(t, response1, http.StatusUnauthorized)
		if code := errorCode(body1); code != "unauthorized" {
			t.Fatalf("unauthenticated code = %q", code)
		}

		response2 := request(t, application.client, http.MethodGet, application.baseURL+"/api/v1/auth/me", nil, map[string]string{"Cookie": "serega_session=expired"})

		body2 := requireJSONResponse(t, response2, 419)
		if code := errorCode(body2); code != "session_expired" {
			t.Fatalf("expired-session code = %q", code)
		}

		body3 := requireJSONResponse(t, request(t, application.client, http.MethodPost, application.baseURL+"/api/v1/auth/login", bytes.NewBufferString("not-json"), map[string]string{requestIDHeader: "integration-request-0003"}), http.StatusUnsupportedMediaType)
		if code := errorCode(body3); code != "unsupported_media_type" {
			t.Fatalf("invalid login code = %q", code)
		}

		var record map[string]any

		waitFor(t, 3*time.Second, func() bool {
			record = findLogRecord(application.logRecords(), "http.request.completed", "integration-request-0003")
			return record != nil
		}, "route-enriched HTTP wide event")

		if record["http.route"] != "/api/v1/auth/login" || record["http.status"] != float64(http.StatusUnsupportedMediaType) {
			t.Fatalf("invalid-input event = %#v; logs=%s", record, describeLogs(application.logLines()))
		}
	})

	t.Run("HTTP wide event is correlated and schema stable", func(t *testing.T) {
		response := request(t, application.client, http.MethodGet, application.baseURL+"/api/v1/not-real/12345678-1234-1234-1234-123456789abc", nil, map[string]string{
			"Authorization": "Bearer must-never-appear",
			requestIDHeader: "integration-request-0002",
		})

		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusNotFound)
		}

		var record map[string]any

		waitFor(t, 3*time.Second, func() bool {
			record = findLogRecord(application.logRecords(), "http.request.completed", "integration-request-0002")
			return record != nil
		}, "correlated HTTP wide event")

		for field, expected := range map[string]any{
			"outcome": "rejected", "http.status": float64(http.StatusNotFound), "http.method": "GET",
			"http.path": "", "http.route": "unmatched",
			"op.domain": "http", "op.name": "unmatched", "op.outcome": "success", "api.error.code": "not_found",
		} {
			if record[field] != expected {
				t.Errorf("%s = %#v, want %#v; logs=%s", field, record[field], expected, describeLogs(application.logLines()))
			}
		}

		if _, ok := record["duration_ms"]; !ok {
			t.Error("wide event is missing duration_ms")
		}

		if strings.Contains(strings.Join(application.logLines(), "\n"), "must-never-appear") {
			t.Error("authorization secret appeared in logs")
		}
	})

	t.Run("OTLP metrics export without blocking requests", func(t *testing.T) {
		for range 3 {
			response := request(t, application.client, http.MethodGet, application.baseURL+"/healthz", nil, nil)
			_, _ = io.Copy(io.Discard, response.Body)

			if response.StatusCode != http.StatusNoContent {
				t.Fatalf("health status = %d", response.StatusCode)
			}
		}

		var latest metricRequest

		waitFor(t, 8*time.Second, func() bool {
			requests := application.metrics.snapshot()
			if len(requests) == 0 {
				return false
			}

			latest = requests[len(requests)-1]

			return true
		}, "OTLP metric export")

		if latest.Path != "/v1/metrics" || latest.ContentType != "application/x-protobuf" || len(latest.Body) == 0 {
			t.Fatalf("unexpected metric request: path=%q content-type=%q body=%d", latest.Path, latest.ContentType, len(latest.Body))
		}
	})

	t.Run("metric export failure does not block application work", func(t *testing.T) {
		before := len(application.metrics.snapshot())

		application.metrics.status.Store(http.StatusServiceUnavailable)
		defer application.metrics.status.Store(http.StatusOK)

		waitFor(t, 5*time.Second, func() bool {
			response := request(t, application.client, http.MethodGet, application.baseURL+"/healthz", nil, nil)

			if response.StatusCode != http.StatusNoContent {
				t.Fatalf("health status during exporter failure = %d", response.StatusCode)
			}

			for _, exported := range application.metrics.snapshot()[before:] {
				if exported.ResponseStatus == http.StatusServiceUnavailable {
					return true
				}
			}

			return false
		}, "simulated OTLP exporter failure")
	})
}
