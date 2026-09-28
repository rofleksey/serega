// Package port contains interfaces required by application use cases.
package port

import (
	"context"
)

// Transaction runs fn atomically. Repositories receive the supplied context so
// calls inside fn use the same transaction without exposing SQL to use cases.
type Transaction interface {
	WithinTransaction(ctx context.Context, fn func(context.Context) error) error
}
