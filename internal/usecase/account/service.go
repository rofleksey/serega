// Package account implements authentication and account provisioning.
package account

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/rofleksey/serega/internal/entity"
	"github.com/rofleksey/serega/internal/observability"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidUsername    = errors.New("username must be between 1 and 64 characters")
	ErrUsernameTaken      = errors.New("username already exists")
	ErrUserNotFound       = errors.New("user not found")
)

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) Authenticate(ctx context.Context, username, password string) (User, error) {
	observability.Enrich(ctx, entity.FieldUseCase, entity.OperationAuthenticate)

	user, err := s.store.FindUserByUsername(ctx, strings.TrimSpace(username))
	if err != nil {
		if !errors.Is(err, ErrUserNotFound) {
			return User{}, fmt.Errorf("find account: %w", err)
		}

		return User{}, ErrInvalidCredentials
	}

	if !CheckPassword(user.PasswordHash, password) {
		observability.Enrich(ctx, entity.FieldDecision, entity.DecisionDenied)
		return User{}, ErrInvalidCredentials
	}

	observability.Enrich(ctx, entity.FieldDecision, entity.DecisionAllowed, entity.FieldAccountID, user.ID)

	return user, nil
}

func (s *Service) User(ctx context.Context, id string) (User, error) {
	return s.store.FindUserByID(ctx, id)
}

// CreateUser adds one account without changing any other user's credentials or sessions.
func (s *Service) CreateUser(ctx context.Context, username, password string) (User, error) {
	username = strings.TrimSpace(username)
	if !utf8.ValidString(username) || utf8.RuneCountInString(username) < 1 || utf8.RuneCountInString(username) > 64 {
		return User{}, ErrInvalidUsername
	}

	hash, err := HashPassword(password)
	if err != nil {
		return User{}, fmt.Errorf("%w: %w", ErrInvalidPassword, err)
	}

	return s.store.CreateUser(ctx, username, hash)
}
