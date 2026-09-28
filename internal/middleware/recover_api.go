package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
)

func RecoverAPI(logger *slog.Logger) func(http.Handler) http.Handler {
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				recovered := recover()
				if recovered == nil {
					return
				}

				logger.ErrorContext(r.Context(), "HTTP handler panicked",
					slog.String("error_code", "handler_panic"),
					slog.String("panic", fmt.Sprint(recovered)),
					slog.String("request_id", RequestIDFrom(r.Context())),
					slog.String("stack", string(debug.Stack())),
				)

				if strings.HasPrefix(r.URL.Path, "/api/") {
					writeError(w, http.StatusInternalServerError, "internal_error", "internal server error", RequestIDFrom(r.Context()))

					return
				}

				http.Error(w, "internal server error", http.StatusInternalServerError)
			}()

			next.ServeHTTP(w, r)
		})
	}
}
