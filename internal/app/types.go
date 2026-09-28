package app

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/happytoolin/unolog"
	"github.com/rofleksey/serega/internal/observability"
)

type Database interface {
	Ping(context.Context) error
	Close()
}

type HTTPDependencies struct {
	TrustedProxies []*net.IPNet
	Database       Database
	Logger         *slog.Logger
	API            http.Handler
	Frontend       http.Handler
	Metrics        *observability.Metrics
	Events         *unolog.Runtime
}

// Runtime owns the HTTP server and resources of one serving process.
type Runtime struct {
	server          *http.Server
	database        Database
	shutdownTimeout time.Duration
	logger          *slog.Logger
	observability   *observability.Provider
	stopSessions    func()
}
