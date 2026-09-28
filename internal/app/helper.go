package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rofleksey/serega/internal/config"
	"github.com/rofleksey/serega/internal/handler"
	"github.com/rofleksey/serega/internal/observability"
	"github.com/rofleksey/serega/internal/store/postgres"
	"github.com/rofleksey/serega/internal/usecase/account"
	"github.com/rofleksey/serega/internal/usecase/board"
)

func newPostgresRuntime(cfg config.Config, pool *pgxpool.Pool, logger *slog.Logger, provider *observability.Provider) *Runtime {
	if err := provider.Metrics.RegisterDatabasePool(func() observability.DatabasePoolSnapshot {
		stats := pool.Stat()
		return observability.DatabasePoolSnapshot{Acquired: int64(stats.AcquiredConns()), Idle: int64(stats.IdleConns()), Maximum: int64(stats.MaxConns())}
	}); err != nil {
		logger.Error("database pool metrics registration failed", "error", err)
	}

	store := postgres.New(pool)
	sessions := pgxstore.NewWithCleanupInterval(pool, 5*time.Minute)
	api := handler.NewHandler(handler.HandlerDependencies{
		Logger: logger, SessionStore: sessions, Accounts: account.NewService(store), Board: board.NewService(store), SecureCookies: cfg.SecureCookies,
	})
	runtime := assembleRuntime(cfg, pool, logger, provider, api)
	runtime.stopSessions = sessions.StopCleanup

	return runtime
}

func assembleRuntime(cfg config.Config, db Database, logger *slog.Logger, provider *observability.Provider, api http.Handler) *Runtime {
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	if provider == nil {
		provider = observability.Disabled(Name)
	}

	return &Runtime{
		server: &http.Server{
			Addr: cfg.HTTPAddr,
			Handler: NewHTTPHandler(HTTPDependencies{
				TrustedProxies: cfg.TrustedProxies, Database: db, Logger: logger, API: api,
				Metrics: provider.Metrics, Events: observability.NewEventRuntime(logger),
			}),
			ReadHeaderTimeout: 10 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second,
		},
		database: db, shutdownTimeout: cfg.ShutdownTimeout, logger: logger, observability: provider,
	}
}

func (r *Runtime) closeResources() {
	if r.stopSessions != nil {
		r.stopSessions()
	}

	if r.database != nil {
		r.database.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.observability.Shutdown(ctx); err != nil {
		r.logger.Error("metrics shutdown failed", "error", err)
	}
}
