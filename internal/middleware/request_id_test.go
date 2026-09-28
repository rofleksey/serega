package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rofleksey/meg/httpx"

	httpapi "github.com/rofleksey/serega/internal/api"
	"github.com/rofleksey/serega/internal/api/generated"
	"github.com/rofleksey/serega/internal/observability"
)

func TestRequestIDCorrelatesValidationErrorsAndEvents(t *testing.T) {
	for _, test := range []struct {
		name    string
		inbound string
	}{
		{name: "missing"},
		{name: "invalid", inbound: "bad"},
		{name: "valid", inbound: "client-correlation-123"},
	} {
		t.Run(test.name, func(t *testing.T) {
			validate, err := httpapi.NewRequestBodyValidator(nil)
			if err != nil {
				t.Fatal(err)
			}

			var output bytes.Buffer

			var contextID string

			router := chi.NewRouter()
			router.Use(RequestID, RequestLogger(observability.NewEventRuntime(observability.NewLogger("json", &output)), nil))
			router.Post("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
				contextID = RequestIDFrom(r.Context())
				if !httpx.ValidRequestID(contextID) || r.Header.Get("X-Request-ID") != contextID {
					t.Fatalf("context ID = %q, downstream header = %q", contextID, r.Header.Get("X-Request-ID"))
				}

				if validate(w, r) {
					t.Fatal("unknown property unexpectedly passed schema validation")
				}
			})

			request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"username":"user","password":"password","unexpected":true}`))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("X-Request-ID", test.inbound)

			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			var body generated.ErrorResponse
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}

			if response.Code != http.StatusBadRequest || body.Error.Code != "invalid_request" {
				t.Fatalf("response = %d %s", response.Code, response.Body.String())
			}

			if response.Header().Get("X-Request-ID") != contextID || body.Error.RequestID != contextID {
				t.Fatalf("response ID = %q, error ID = %q, context ID = %q", response.Header().Get("X-Request-ID"), body.Error.RequestID, contextID)
			}

			if httpx.ValidRequestID(test.inbound) && contextID != test.inbound {
				t.Fatalf("valid request ID changed to %q", contextID)
			}

			var event struct {
				RequestID string `json:"request_id"`
			}
			if err := json.Unmarshal(output.Bytes(), &event); err != nil {
				t.Fatal(err)
			}

			if event.RequestID != contextID {
				t.Fatalf("event request ID = %q, context ID = %q", event.RequestID, contextID)
			}
		})
	}
}
