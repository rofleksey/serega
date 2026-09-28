//go:build integration

package application_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const (
	contentTypeHeader = "Content-Type"
	originHeader      = "Origin"
	csrfHeader        = "X-Csrf-Token"
	jsonContentType   = "application/json"
)

type card struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Version     int    `json:"version"`
	CreatedBy   struct {
		Username string `json:"username"`
	} `json:"createdBy"`
	UpdatedBy struct {
		Username string `json:"username"`
	} `json:"updatedBy"`
}

type browser struct {
	client *http.Client
	csrf   string
}

func TestSharedBoardAcrossAccounts(t *testing.T) {
	app := startApplication(t)
	createAccount(t, app, "alice")
	alice := signIn(t, app, "alice")
	// Provisioning Bob must preserve Alice's account AND live session.
	createAccount(t, app, "bob")
	bob := signIn(t, app, "bob")

	response := boardRequest(t, app, alice, http.MethodPost, "/api/v1/cards", `{"title":"Write the reference","description":"Keep it small"}`)

	first := readCard(t, response, http.StatusCreated)
	if first.CreatedBy.Username != "alice" || first.Status != "todo" || first.Version != 1 {
		t.Fatalf("created card = %+v", first)
	}

	list := requireJSONResponse(t, boardRequest(t, app, bob, http.MethodGet, "/api/v1/cards", ""), http.StatusOK)

	cards, ok := list["cards"].([]any)
	if !ok || len(cards) != 1 {
		t.Fatalf("shared board = %+v", list)
	}

	path := "/api/v1/cards/" + first.ID
	update := `{"title":"Write the reference","description":"Reviewed together","status":"doing","version":1}`

	current := readCard(t, boardRequest(t, app, bob, http.MethodPatch, path, update), http.StatusOK)
	if current.UpdatedBy.Username != "bob" || current.Version != 2 {
		t.Fatalf("updated card = %+v", current)
	}

	requireJSONResponse(t, boardRequest(t, app, alice, http.MethodPatch, path, update), http.StatusConflict)
	requireJSONResponse(t, boardRequest(t, app, alice, http.MethodDelete, path+"?version=1", ""), http.StatusConflict)
	requireJSONResponse(t, boardRequest(t, app, alice, http.MethodPost, "/api/v1/cards", `{"title":"  ","description":""}`), http.StatusBadRequest)

	noCSRF := request(t, alice.client, http.MethodPost, app.baseURL+"/api/v1/cards", strings.NewReader(`{"title":"blocked","description":""}`), map[string]string{contentTypeHeader: jsonContentType, originHeader: app.baseURL})
	requireJSONResponse(t, noCSRF, http.StatusForbidden)
	crossOrigin := request(t, alice.client, http.MethodPost, app.baseURL+"/api/v1/cards", strings.NewReader(`{"title":"blocked","description":""}`), map[string]string{contentTypeHeader: jsonContentType, originHeader: "https://untrusted.example", csrfHeader: alice.csrf})
	requireJSONResponse(t, crossOrigin, http.StatusForbidden)

	outsider := &http.Client{Timeout: 2 * time.Second}
	requireJSONResponse(t, request(t, outsider, http.MethodGet, app.baseURL+"/api/v1/cards", nil, nil), http.StatusUnauthorized)

	deleted := boardRequest(t, app, alice, http.MethodDelete, fmt.Sprintf("%s?version=%d", path, current.Version), "")
	requireNoContent(t, deleted)
	requireJSONResponse(t, boardRequest(t, app, bob, http.MethodPatch, path, update), http.StatusNotFound)
	requireNoContent(t, boardRequest(t, app, alice, http.MethodPost, "/api/v1/auth/logout", ""))
	requireJSONResponse(t, boardRequest(t, app, alice, http.MethodGet, "/api/v1/cards", ""), http.StatusUnauthorized)
	// Logging out one user never logs another out.
	requireJSONResponse(t, boardRequest(t, app, bob, http.MethodGet, "/api/v1/cards", ""), http.StatusOK)

	for _, line := range app.logLines() {
		if strings.Contains(line, "test-password-123") || strings.Contains(line, "Reviewed together") {
			t.Fatal("private content leaked into logs")
		}
	}
}

func createAccount(t *testing.T, app *applicationHarness, username string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()

	command := exec.CommandContext(ctx, app.binary, "user", "create", "--username", username, "--password-stdin")
	command.Env = app.environment

	command.Stdin = strings.NewReader("test-password-123\n")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("provision account: %v %s", err, output)
	}
}

func signIn(t *testing.T, app *applicationHarness, username string) browser {
	t.Helper()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	user := browser{client: &http.Client{Jar: jar, Timeout: 3 * time.Second}}
	token := requireJSONResponse(t, request(t, user.client, http.MethodGet, app.baseURL+"/api/v1/auth/csrf", nil, nil), http.StatusOK)
	user.csrf, _ = token["token"].(string)
	data := fmt.Sprintf(`{"username":%q,"password":"test-password-123"}`, username)
	requireNoContent(t, boardRequest(t, app, user, http.MethodPost, "/api/v1/auth/login", data))

	return user
}

func boardRequest(t *testing.T, app *applicationHarness, user browser, method, path, body string) *responseSnapshot {
	t.Helper()
	return request(t, user.client, method, app.baseURL+path, strings.NewReader(body), map[string]string{contentTypeHeader: jsonContentType, csrfHeader: user.csrf, originHeader: app.baseURL})
}

func readCard(t *testing.T, response *responseSnapshot, status int) card {
	t.Helper()

	if response.StatusCode != status {
		body, _ := io.ReadAll(response.Body)
		t.Fatalf("status = %d, want %d: %s", response.StatusCode, status, body)
	}

	var value card
	if err := json.NewDecoder(response.Body).Decode(&value); err != nil {
		t.Fatal(err)
	}

	return value
}

func requireNoContent(t *testing.T, response *responseSnapshot) {
	t.Helper()

	if response.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(response.Body)
		t.Fatalf("status = %d: %s", response.StatusCode, body)
	}
}
