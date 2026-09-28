package postgres

import "github.com/jackc/pgx/v5/pgxpool"

// Store is a PostgreSQL-backed store. It owns no schema migration side effect.
type Store struct {
	pool *pgxpool.Pool
}
