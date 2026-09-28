// Package handler adapts HTTP requests to application use cases.
package handler

import (
	"log/slog"
	"net/http"

	"github.com/rofleksey/serega/internal/middleware"
)

const userIDKey = "user_id"

// NewHandler composes the HTTP transport over explicit application services.
func NewHandler(dependencies HandlerDependencies) http.Handler {
	logger := dependencies.Logger
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	h := &handler{
		logger: logger, accounts: dependencies.Accounts, board: dependencies.Board,
		sessions: middleware.NewSessionManager(dependencies.SessionStore, dependencies.SecureCookies),
	}
	mux := newGeneratedRouter(h)

	return h.sessions.LoadAndSave(middleware.CSRF(mux, dependencies.SecureCookies))
}
