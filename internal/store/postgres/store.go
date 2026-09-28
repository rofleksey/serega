// Package postgres adapts SQLC-generated PostgreSQL queries to domain ports.
// Queries stay generated; transactions and domain error mapping remain handwritten.
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	storesqlc "github.com/rofleksey/serega/internal/store/postgres/sqlc/generated"
	"github.com/rofleksey/serega/internal/usecase/port"
)

var _ port.Transaction = (*Store)(nil)

// New binds generated queries to an already-open connection pool.
func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

// Close releases database connections owned by the caller's Store.
func (s *Store) Close() {
	s.pool.Close()
}

// WithinTransaction implements the application transaction port.
func (s *Store) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	if _, exists := ctx.Value(transactionKey{}).(pgx.Tx); exists {
		return fn(ctx)
	}

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(context.WithValue(ctx, transactionKey{}, tx)); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *Store) queryer(ctx context.Context) queryer {
	if tx, ok := ctx.Value(transactionKey{}).(pgx.Tx); ok {
		return tx
	}

	return s.pool
}

// generated binds sqlc's typed queries to the same pool-or-transaction
// boundary used by every adapter operation.
func (s *Store) generated(ctx context.Context) *storesqlc.Queries {
	return storesqlc.New(s.queryer(ctx))
}

type transactionKey struct{}

type queryer interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}
