package account

import (
	"context"

	"github.com/rofleksey/serega/internal/entity"
)

type User = entity.User

// Store is the persistence port required by account use cases.
type Store interface {
	FindUserByUsername(context.Context, string) (User, error)
	FindUserByID(context.Context, string) (User, error)
	CreateUser(context.Context, string, string) (User, error)
}

// Service owns credential validation and multi-user account creation.
type Service struct{ store Store }
