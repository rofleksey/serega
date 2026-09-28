package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	databasepostgres "github.com/rofleksey/serega/internal/database/postgres"
	"github.com/rofleksey/serega/internal/store/postgres"
	"github.com/rofleksey/serega/internal/usecase/account"
)

// CreateUser composes account provisioning without exposing persistence to CLI.
func CreateUser(ctx context.Context, databaseURL, username, password string) error {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("open database pool: %w", err)
	}
	defer pool.Close()

	if err := databasepostgres.CheckReady(ctx, pool); err != nil {
		return err
	}

	if _, err := account.NewService(postgres.New(pool)).CreateUser(ctx, username, password); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}
