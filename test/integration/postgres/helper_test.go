//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	databasepostgres "github.com/rofleksey/serega/internal/database/postgres"
	storepg "github.com/rofleksey/serega/internal/store/postgres"
)

func newDatabase(tb testing.TB) *pgxpool.Pool {
	tb.Helper()

	if integrationPostgres == nil {
		tb.Fatal("PostgreSQL integration harness was not initialized")
	}

	ctx := context.Background()
	name := fmt.Sprintf("serega_test_%d_%d", time.Now().UnixNano(), integrationPostgres.sequence.Add(1))

	identifier := pgx.Identifier{name}.Sanitize()
	if _, err := integrationPostgres.admin.Exec(ctx, "CREATE DATABASE "+identifier); err != nil {
		tb.Fatal(err)
	}

	config := integrationPostgres.base.Copy()
	config.ConnConfig.Database = name

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		tb.Fatal(err)
	}

	tb.Cleanup(func() {
		pool.Close()

		if _, err := integrationPostgres.admin.Exec(context.Background(), "DROP DATABASE "+identifier+" WITH (FORCE)"); err != nil {
			tb.Error(err)
		}
	})

	return pool
}

func newMigratedStore(tb testing.TB) (*pgxpool.Pool, *storepg.Store) {
	tb.Helper()

	pool := newDatabase(tb)
	if err := databasepostgres.MigrateUpPool(context.Background(), pool); err != nil {
		tb.Fatal(err)
	}

	return pool, storepg.New(pool)
}
