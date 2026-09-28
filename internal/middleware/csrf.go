package middleware

import (
	"net/http"

	"github.com/justinas/nosurf"
	httpapi "github.com/rofleksey/serega/internal/api"
)

func CSRF(next http.Handler, secureCookies bool) http.Handler {
	csrf := nosurf.New(next)
	// nosurf assumes HTTPS by default, including behind a TLS-terminating proxy.
	// Match the explicit cookie configuration for local plain HTTP.
	csrf.SetIsTLSFunc(func(r *http.Request) bool { return secureCookies || r.TLS != nil })
	csrf.SetBaseCookie(http.Cookie{ //nolint:gosec // Secure defaults true; explicit local HTTP configuration may disable it.
		Name:     "serega_csrf",
		Path:     "/",
		Secure:   secureCookies,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	csrf.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpapi.WriteError(w, http.StatusForbidden, "csrf_failed", "CSRF token is missing or invalid", RequestIDFrom(r.Context()))
	}))

	return csrf
}
