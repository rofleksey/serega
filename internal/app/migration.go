package app

import (
	"context"

	databasepostgres "github.com/rofleksey/serega/internal/database/postgres"
)

// MigrateUp runs the explicit administrative schema migration command.
func MigrateUp(ctx context.Context, databaseURL string) error {
	return databasepostgres.MigrateUp(ctx, databaseURL)
}
