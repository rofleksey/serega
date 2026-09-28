package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	storesqlc "github.com/rofleksey/serega/internal/store/postgres/sqlc/generated"
	"github.com/rofleksey/serega/internal/usecase/account"
)

var _ account.Store = (*Store)(nil)

func (s *Store) FindUserByID(ctx context.Context, id string) (account.User, error) {
	row, err := s.generated(ctx).FindAuthUserByID(ctx, id)
	return accountRow(row), accountError(err)
}

func (s *Store) FindUserByUsername(ctx context.Context, username string) (account.User, error) {
	row, err := s.generated(ctx).FindUserByUsername(ctx, username)
	return accountRow(row), accountError(err)
}

func (s *Store) CreateUser(ctx context.Context, username, passwordHash string) (account.User, error) {
	row, err := s.generated(ctx).CreateUser(ctx, storesqlc.CreateUserParams{Username: username, PasswordHash: passwordHash})
	return accountRow(row), accountError(err)
}

func accountRow(row storesqlc.User) account.User {
	return account.User{ID: row.ID, Username: row.Username, PasswordHash: row.PasswordHash, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt}
}

func accountError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return account.ErrUserNotFound
	}

	var constraint *pgconn.PgError
	if errors.As(err, &constraint) && constraint.Code == "23505" {
		return account.ErrUsernameTaken
	}

	return err
}
