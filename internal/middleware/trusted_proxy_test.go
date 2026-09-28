package middleware

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrustedProxyAcceptsForwardingOnlyFromConfiguredNetwork(t *testing.T) {
	_, network, err := net.ParseCIDR("10.0.0.0/8")
	if err != nil {
		t.Fatal(err)
	}

	handler := TrustedProxy([]*net.IPNet{network})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ClientIP(r) != "203.0.113.8" {
			t.Fatalf("ClientIP() = %q", ClientIP(r))
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	request.RemoteAddr = "10.1.2.3:443"
	request.Header.Set("X-Forwarded-For", "203.0.113.8")
	handler.ServeHTTP(httptest.NewRecorder(), request)
}
