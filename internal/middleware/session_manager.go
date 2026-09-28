package middleware

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const SessionCookieName = "serega_session"

func NewSessionManager(store scs.Store, secureCookies bool) *scs.SessionManager {
	sessions := scs.New()
	sessions.Store = store
	sessions.HashTokenInStore = true
	sessions.IdleTimeout = 30 * time.Minute
	sessions.Lifetime = 12 * time.Hour
	sessions.Cookie.Name = SessionCookieName
	sessions.Cookie.Path = "/"
	sessions.Cookie.HttpOnly = true
	sessions.Cookie.Secure = secureCookies
	sessions.Cookie.SameSite = http.SameSiteStrictMode
	sessions.Cookie.Persist = true

	return sessions
}
