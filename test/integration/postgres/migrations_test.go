//go:build integration

package postgres_test

import (
	"testing"

	databasepostgres "github.com/rofleksey/serega/internal/database/postgres"
)

func TestReadinessNeverMigratesAndExplicitMigrationsRoundTrip(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	pool := newDatabase(t)
	if err := databasepostgres.CheckReady(ctx, pool); err == nil {
		t.Fatal("empty database reported ready")
	}

	var tables int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM pg_tables WHERE schemaname = 'public'").Scan(&tables); err != nil {
		t.Fatal(err)
	}

	if tables != 0 {
		t.Fatalf("readiness created %d tables", tables)
	}

	if err := databasepostgres.MigrateUpPool(ctx, pool); err != nil {
		t.Fatal(err)
	}

	if err := databasepostgres.MigrateUpPool(ctx, pool); err != nil {
		t.Fatalf("idempotent migration = %v", err)
	}

	if err := databasepostgres.CheckReady(ctx, pool); err != nil {
		t.Fatal(err)
	}

	if err := databasepostgres.MigrateDownPool(ctx, pool); err != nil {
		t.Fatal(err)
	}

	if err := databasepostgres.CheckReady(ctx, pool); err == nil {
		t.Fatal("rolled-back database reported ready")
	}

	if err := pool.QueryRow(ctx, "SELECT count(*) FROM pg_tables WHERE schemaname = 'public' AND tablename IN ('users','sessions','cards')").Scan(&tables); err != nil {
		t.Fatal(err)
	}

	if tables != 0 {
		t.Fatalf("rollback left %d application tables", tables)
	}

	if err := databasepostgres.MigrateUpPool(ctx, pool); err != nil {
		t.Fatal(err)
	}

	if err := databasepostgres.CheckReady(ctx, pool); err != nil {
		t.Fatal(err)
	}
}
