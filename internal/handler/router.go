package handler

import (
	"net/http"

	apiadapter "github.com/rofleksey/serega/internal/api"
	"github.com/rofleksey/serega/internal/api/generated"
	"github.com/rofleksey/serega/internal/entity"
	"github.com/rofleksey/serega/internal/observability"
)

var _ generated.ServerInterface = (*handler)(nil)

// observeHandler attaches stable route patterns to the request's wide event.
func observeHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enrichHandler(r)
		next.ServeHTTP(w, r)
	})
}

func enrichHandler(r *http.Request) {
	if r != nil && r.Pattern != "" {
		observability.Enrich(r.Context(), entity.FieldHandler, r.Pattern)
	}
}

func newGeneratedRouter(h *handler) http.Handler {
	validateBody, err := apiadapter.NewRequestBodyValidator(h.logger)
	if err != nil {
		panic(err)
	}

	h.validateBody = validateBody

	return generated.HandlerWithOptions(h, generated.StdHTTPServerOptions{
		BaseURL:     "/api",
		Middlewares: []generated.MiddlewareFunc{observeHandler},
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, _ error) {
			enrichHandler(r)
			writeError(w, r, http.StatusBadRequest, "invalid_request", "request parameters are invalid")
		},
	})
}
