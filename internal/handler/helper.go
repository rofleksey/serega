package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	apiadapter "github.com/rofleksey/serega/internal/api"
	"github.com/rofleksey/serega/internal/api/generated"
	"github.com/rofleksey/serega/internal/entity"
	"github.com/rofleksey/serega/internal/middleware"
	"github.com/rofleksey/serega/internal/observability"
)

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		observability.Enrich(r.Context(), entity.FieldValidation, entity.DecisionRejected)
		writeError(w, r, http.StatusBadRequest, "invalid_request", "request body is invalid")

		return false
	}

	observability.Enrich(r.Context(), entity.FieldValidation, entity.DecisionAccepted)

	return true
}

func (h *handler) writeJSON(w http.ResponseWriter, r *http.Request, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		h.logError(r, "response_encode_failed", err)

		return
	}
}

func (h *handler) writeInternalError(w http.ResponseWriter, r *http.Request, code, message string, err error) {
	h.logError(r, code, err)
	writeError(w, r, http.StatusInternalServerError, code, message)
}

func (h *handler) logError(r *http.Request, code string, err error) {
	// Database errors can contain user input. Record a stable code and error type,
	// leaving the request-wide event as the single owner of log emission.
	observability.Enrich(r.Context(), entity.FieldErrorKind, fmt.Sprintf("%T", err))
	observability.RecordError(r.Context(), errors.New(code))
}

// writeError is the single HTTP error boundary for every handler. It keeps the
// generated transport model out of individual endpoints and always attaches
// the request ID in the same way.
func writeError(w http.ResponseWriter, r *http.Request, status int, code, message string, fields ...generated.FieldError) {
	if status == http.StatusUnauthorized && code == "unauthorized" {
		if _, err := r.Cookie(middleware.SessionCookieName); err == nil {
			status = sessionExpiredHTTPStatus
			code = "session_expired"
			message = "browser session expired; sign in again"
		}
	}

	observability.Enrich(r.Context(),
		entity.FieldAPIErrorCode, code,
		entity.FieldHTTPStatusCode, status,
	)
	apiadapter.WriteError(w, status, code, message, r.Header.Get("X-Request-ID"), fields...)
}

func (h *handler) validateJSONRequest(w http.ResponseWriter, r *http.Request) bool {
	return h.validateBody == nil || h.validateBody(w, r)
}
