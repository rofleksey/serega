//go:build integration

package postgres_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rofleksey/serega/internal/entity"
	storepg "github.com/rofleksey/serega/internal/store/postgres"
	"github.com/rofleksey/serega/internal/usecase/account"
	"github.com/rofleksey/serega/internal/usecase/board"
)

func TestSharedBoardPersistsAcrossConnectionsAndGuardsConcurrentChanges(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	pool, store := newMigratedStore(t)
	accounts := account.NewService(store)

	const password = "example password for integration"

	alice, err := accounts.CreateUser(ctx, "alice", password)
	if err != nil {
		t.Fatal(err)
	}

	bob, err := accounts.CreateUser(ctx, "bob", password)
	if err != nil {
		t.Fatal(err)
	}

	if alice.ID == bob.ID {
		t.Fatal("accounts must be distinct")
	}

	if _, err := accounts.CreateUser(ctx, "alice", password); !errors.Is(err, account.ErrUsernameTaken) {
		t.Fatalf("duplicate username = %v", err)
	}

	if _, err := accounts.Authenticate(ctx, "alice", password); err != nil {
		t.Fatal(err)
	}

	if _, err := accounts.Authenticate(ctx, "bob", password); err != nil {
		t.Fatal(err)
	}

	service := board.NewService(store)

	card, err := service.Create(ctx, alice.ID, board.CreateInput{Title: "  First task  ", Description: "A shared task"})
	if err != nil {
		t.Fatal(err)
	}

	if card.Title != "First task" || card.Status != entity.CardStatusTodo || card.Version != 1 || card.CreatedBy.Username != alice.Username {
		t.Fatalf("created card = %#v", card)
	}

	secondPool, err := pgxpool.NewWithConfig(ctx, pool.Config().Copy())
	if err != nil {
		t.Fatal(err)
	}
	defer secondPool.Close()

	otherService := board.NewService(storepg.New(secondPool))

	cards, err := otherService.List(ctx)
	if err != nil || len(cards) != 1 || cards[0].ID != card.ID {
		t.Fatalf("second connection list = %#v, %v", cards, err)
	}

	// Two independent connections try to replace the same observed version.
	var workers sync.WaitGroup

	start := make(chan struct{})
	results := make(chan error, 2)

	for _, writer := range []*board.Service{service, otherService} {
		workers.Go(func() {
			<-start

			_, updateErr := writer.Update(ctx, bob.ID, card.ID, board.UpdateInput{Title: card.Title, Description: card.Description, Status: entity.CardStatusDoing, Version: card.Version})
			results <- updateErr
		})
	}

	close(start)
	workers.Wait()
	close(results)

	successes, conflicts := 0, 0

	for result := range results {
		switch {
		case result == nil:
			successes++
		case errors.Is(result, board.ErrConflict):
			conflicts++
		default:
			t.Fatalf("concurrent update = %v", result)
		}
	}

	if successes != 1 || conflicts != 1 {
		t.Fatalf("successes=%d conflicts=%d", successes, conflicts)
	}

	if err := service.Delete(ctx, card.ID, 1); !errors.Is(err, board.ErrConflict) {
		t.Fatalf("stale delete = %v", err)
	}

	cards, err = otherService.List(ctx)
	if err != nil || len(cards) != 1 {
		t.Fatalf("after stale delete = %#v, %v", cards, err)
	}

	updated := cards[0]
	if updated.Version != 2 || updated.Status != entity.CardStatusDoing || updated.CreatedBy.ID != alice.ID || updated.UpdatedBy.ID != bob.ID || updated.UpdatedAt.Before(updated.CreatedAt) {
		t.Fatalf("updated card = %#v", updated)
	}

	secondPool.Close()

	reconnect, err := pgxpool.NewWithConfig(ctx, pool.Config().Copy())
	if err != nil {
		t.Fatal(err)
	}
	defer reconnect.Close()

	persisted, err := board.NewService(storepg.New(reconnect)).List(ctx)
	if err != nil || len(persisted) != 1 || persisted[0].Version != 2 {
		t.Fatalf("reconnected cards = %#v, %v", persisted, err)
	}

	if err := service.Delete(ctx, card.ID, updated.Version); err != nil {
		t.Fatal(err)
	}

	if err := service.Delete(ctx, card.ID, updated.Version); !errors.Is(err, board.ErrNotFound) {
		t.Fatalf("missing delete = %v", err)
	}

	cards, err = service.List(ctx)
	if err != nil || cards == nil || len(cards) != 0 {
		t.Fatalf("empty board = %#v, %v", cards, err)
	}
}
