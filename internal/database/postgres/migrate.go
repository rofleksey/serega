// Package postgres owns explicit PostgreSQL lifecycle operations.
package postgres

import (
	"context"
	"errors"
)

// MigrateUp applies the versioned schema migrations. It is intentionally called
// only by the migrate command; serving Serega never changes a production schema.
func MigrateUp(ctx context.Context, databaseURL string) error {
	if databaseURL == "" {
		return errors.New("DATABASE_URL is required for migrations")
	}

	return migrateURL(ctx, databaseURL)
}
