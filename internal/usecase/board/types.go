package board

import (
	"context"
	"errors"

	"github.com/rofleksey/serega/internal/entity"
)

var (
	ErrNotFound = errors.New("card not found")
	ErrConflict = errors.New("card changed; refresh and try again")
)

type Card = entity.Card

type CreateInput struct{ Title, Description string }
type UpdateInput struct {
	Title, Description, Status string
	Version                    int64
}

// ValidationError identifies a domain field without depending on HTTP types.
type ValidationError struct{ Field, Message string }

func (e *ValidationError) Error() string { return e.Field + ": " + e.Message }

// Store implements atomic compare-and-swap for updates and deletes.
type Store interface {
	ListCards(context.Context) ([]Card, error)
	CreateCard(context.Context, string, CreateInput) (Card, error)
	UpdateCard(context.Context, string, string, UpdateInput) (Card, error)
	DeleteCard(context.Context, string, int64) error
}

type Service struct{ store Store }
