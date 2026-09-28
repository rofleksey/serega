package middleware

import (
	"testing"

	"github.com/alexedwards/scs/v2/memstore"
)

func TestNewSessionManagerUsesSecurePersistentCookies(t *testing.T) {
	sessions := NewSessionManager(memstore.New(), true)

	if sessions.Cookie.Name != SessionCookieName || !sessions.Cookie.Secure || !sessions.Cookie.HttpOnly || !sessions.Cookie.Persist {
		t.Fatalf("cookie settings = %+v", sessions.Cookie)
	}
}
