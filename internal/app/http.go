package app

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	httpapi "github.com/rofleksey/serega/internal/api"
	"github.com/rofleksey/serega/internal/entity"
	"github.com/rofleksey/serega/internal/frontend"
	"github.com/rofleksey/serega/internal/middleware"
	"github.com/rofleksey/serega/internal/observability"
)

// NewHTTPHandler builds the same-origin surface. It deliberately contains no
// CORS middleware: the browser UI and API share this one origin.
func NewHTTPHandler(dependencies HTTPDependencies) http.Handler {
	logger := dependencies.Logger
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	events := dependencies.Events
	if events == nil {
		events = observability.NewEventRuntime(logger)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.TrustedProxy(dependencies.TrustedProxies))
	router.Use(middleware.RequestLogger(events, dependencies.Metrics))
	router.Use(middleware.RecoverAPI(logger))

	router.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	router.Get("/readyz", func(w http.ResponseWriter, request *http.Request) {
		if dependencies.Database == nil {
			logger.ErrorContext(request.Context(), "readiness check failed",
				slog.String("error_code", "database_unavailable"),
				slog.String("request_id", middleware.RequestIDFrom(request.Context())),
			)
			httpapi.WriteError(w, http.StatusServiceUnavailable, "not_ready", "database is not ready", middleware.RequestIDFrom(request.Context()))

			return
		}

		if err := dependencies.Database.Ping(request.Context()); err != nil {
			logger.ErrorContext(request.Context(), "readiness check failed",
				slog.String("error_code", "database_ping_failed"),
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFrom(request.Context())),
			)
			httpapi.WriteError(w, http.StatusServiceUnavailable, "not_ready", "database is not ready", middleware.RequestIDFrom(request.Context()))

			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	if dependencies.API != nil {
		router.With(middleware.NoStore, middleware.RateLimit(12, time.Minute), middleware.RequestBodyLimit, middleware.RequestJSON).
			Method(http.MethodPost, "/api/v1/auth/login", dependencies.API)

		router.Group(func(apiRoutes chi.Router) {
			apiRoutes.Use(middleware.NoStore)
			apiRoutes.Use(middleware.RequestBodyLimit)
			apiRoutes.Use(middleware.RequestJSON)
			apiRoutes.Mount("/api/v1/auth", dependencies.API)
			apiRoutes.Mount("/api/v1/cards", dependencies.API)
		})
	}

	frontendHandler := dependencies.Frontend
	if frontendHandler == nil {
		frontendHandler = frontend.NewHandler()
	}

	router.NotFound(func(w http.ResponseWriter, request *http.Request) {
		if strings.HasPrefix(request.URL.Path, "/api/") {
			observability.Enrich(request.Context(), entity.FieldAPIErrorCode, entity.APIErrorNotFound)
			httpapi.WriteError(w, http.StatusNotFound, entity.APIErrorNotFound, "API endpoint not found", middleware.RequestIDFrom(request.Context()))

			return
		}

		frontendHandler.ServeHTTP(w, request)
	})

	return router
}
