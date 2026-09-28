package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"

	"github.com/rofleksey/serega/internal/api/generated"
)

// NewRequestBodyValidator returns a request-body validation boundary for
// handlers that have already applied Serega's authentication, CSRF, request-size,
// and JSON syntax policy. It reports whether the concrete handler may continue
// and always restores the body for normal decoding.
func NewRequestBodyValidator(logger *slog.Logger) (func(http.ResponseWriter, *http.Request) bool, error) {
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	document, err := generated.GetSpec()
	if err != nil {
		return nil, fmt.Errorf("load embedded OpenAPI contract: %w", err)
	}

	stripRequestParameters(document)

	validator := nethttpmiddleware.OapiRequestValidatorWithOptions(document, &nethttpmiddleware.Options{
		Prefix:               "/api",
		DoNotValidateServers: true,
		Options: openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		},
		ErrorHandlerWithOpts: func(ctx context.Context, validationErr error, w http.ResponseWriter, r *http.Request, options nethttpmiddleware.ErrorHandlerOpts) {
			if options.StatusCode >= http.StatusInternalServerError {
				logger.ErrorContext(ctx, "OpenAPI request validation failed",
					slog.String("error_code", "request_validation_failed"),
					slog.String("error_kind", fmt.Sprintf("%T", validationErr)),
					slog.String("request_id", r.Header.Get("X-Request-ID")),
				)
				WriteError(w, http.StatusInternalServerError, "internal_error", "request validation failed", r.Header.Get("X-Request-ID"))

				return
			}

			WriteError(w, http.StatusBadRequest, "invalid_request", "request does not match API schema", r.Header.Get("X-Request-ID"))
		},
	})
	checked := validator(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		if state, ok := r.Context().Value(validationStateKey{}).(*validationState); ok {
			state.accepted = true
		}
	}))

	return func(w http.ResponseWriter, r *http.Request) bool {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.WarnContext(r.Context(), "request body could not be read for OpenAPI validation",
				slog.String("error_code", "request_body_read_failed"),
				slog.String("error_kind", fmt.Sprintf("%T", err)),
				slog.String("request_id", r.Header.Get("X-Request-ID")),
			)

			return true
		}

		defer func() { r.Body = io.NopCloser(bytes.NewReader(body)) }()

		state := &validationState{}
		r.Body = io.NopCloser(bytes.NewReader(body))
		checked.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), validationStateKey{}, state)))

		return state.accepted
	}, nil
}
