package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/happytoolin/unolog"
	"github.com/happytoolin/unolog/integration/flow"
	"github.com/rofleksey/serega/internal/entity"
	"github.com/rofleksey/serega/internal/observability"
)

func RequestLogger(events *unolog.Runtime, metrics *observability.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			operation := flow.StartRequest(r.Context(), events, r.Method, "")
			requestContext := operation.Context()
			observability.Enrich(requestContext,
				entity.FieldEvent, entity.EventHTTPRequest,
				entity.FieldRequestID, RequestIDFrom(r.Context()),
			)
			r = r.WithContext(requestContext)

			response := &statusWriter{ResponseWriter: w}
			started := time.Now()

			metrics.HTTPStarted(requestContext, r.Method)

			next.ServeHTTP(response, r)

			status := response.status
			if status == 0 {
				status = http.StatusOK
			}

			duration := time.Since(started)

			route := chi.RouteContext(r.Context()).RoutePattern()
			if route == "" {
				route = "unmatched"
			}

			metrics.HTTPCompleted(requestContext, r.Method, route, status, duration)
			observability.Enrich(requestContext,
				entity.FieldOutcome, httpOutcome(status),
				entity.FieldHTTPSize, response.bytes,
			)
			flow.FinalizeRequest(operation, route, status, nil, nil)
		})
	}
}
