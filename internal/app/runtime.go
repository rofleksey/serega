package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rofleksey/serega/internal/config"
	databasepostgres "github.com/rofleksey/serega/internal/database/postgres"
	"github.com/rofleksey/serega/internal/entity"
	"github.com/rofleksey/serega/internal/observability"
)

// NewRuntime assembles the transport lifecycle around an arbitrary database
// boundary. Production Postgres composition is explicit in Open.
func NewRuntime(cfg config.Config, db Database, logger *slog.Logger) *Runtime {
	return assembleRuntime(cfg, db, logger, observability.Disabled(Name), nil)
}

// Handler exposes the assembled public surface for focused tests.
func (r *Runtime) Handler() http.Handler { return r.server.Handler }

// Serve blocks until the server exits or its context is cancelled. Cancellation
// first drains HTTP work, then closes the database pool.
func (r *Runtime) Serve(ctx context.Context) error {
	defer r.closeResources()

	errs := make(chan error, 1)
	go func() { errs <- r.server.ListenAndServe() }()

	select {
	case err := <-errs:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), r.shutdownTimeout)
		defer cancel()

		if err := r.server.Shutdown(shutdownCtx); err != nil {
			_ = r.server.Close()
			return err
		}

		return nil
	}
}

// NewLogger creates the process logger. JSON is the deployment default; text
// is useful for local runs.
func NewLogger(format string, output io.Writer) *slog.Logger {
	return observability.NewLogger(format, output)
}

// Open validates database readiness before the process accepts HTTP traffic.
func Open(ctx context.Context, cfg config.Config, logger *slog.Logger) (*Runtime, error) {
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if err := databasepostgres.CheckReady(ctx, pool); err != nil {
		pool.Close()
		return nil, err
	}

	provider, observabilityErr := observability.Open(ctx, Name)
	if observabilityErr != nil {
		logger.Error("metrics initialization failed",
			entity.FieldEvent, entity.EventMetricsInit,
			entity.FieldOutcome, entity.OutcomeFailed,
			entity.FieldErrorKind, fmt.Sprintf("%T", observabilityErr),
		)

		provider = observability.Disabled(Name)
	}

	return newPostgresRuntime(cfg, pool, logger, provider), nil
}
