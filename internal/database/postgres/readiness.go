package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rofleksey/serega/internal/database/postgres/migrations"
)

// CheckReady is read-only. Goose's version helpers can initialize its metadata
// table, which must never happen in the serving process. This infrastructure
// query is deliberately outside application SQLC queries.
func CheckReady(ctx context.Context, pool *pgxpool.Pool) error {
	entries, err := migrations.Files.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read embedded migrations: %w", err)
	}

	if len(entries) == 0 {
		return errors.New("no embedded migrations")
	}

	var latest int64

	for _, entry := range entries {
		prefix, _, found := strings.Cut(entry.Name(), "_")
		if !found || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		version, parseErr := strconv.ParseInt(prefix, 10, 64)
		if parseErr != nil {
			return fmt.Errorf("invalid migration filename: %w", parseErr)
		}

		latest = max(latest, version)
	}

	var applied int64

	err = pool.QueryRow(ctx, `SELECT COALESCE(MAX(version_id), 0) FROM (
  SELECT DISTINCT ON (version_id) version_id, is_applied
  FROM goose_db_version ORDER BY version_id, id DESC
 ) versions WHERE is_applied`).Scan(&applied)
	if err != nil {
		return fmt.Errorf("schema is not ready; run serega migrate up: %w", err)
	}

	if applied != latest {
		return fmt.Errorf("schema version %d, expected %d; run serega migrate up", applied, latest)
	}

	return nil
}
