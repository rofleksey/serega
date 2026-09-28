package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rofleksey/serega/internal/entity"
	"github.com/rofleksey/serega/internal/observability"
	storesqlc "github.com/rofleksey/serega/internal/store/postgres/sqlc/generated"
	"github.com/rofleksey/serega/internal/usecase/board"
)

var _ board.Store = (*Store)(nil)

func (s *Store) ListCards(ctx context.Context) ([]board.Card, error) {
	observability.Enrich(ctx, entity.FieldAdapter, entity.AdapterPostgres, entity.FieldAdapterOperation, entity.OperationListCards)

	rows, err := s.generated(ctx).ListCards(ctx)
	if err != nil {
		return nil, err
	}

	cards := make([]board.Card, 0, len(rows))
	for _, row := range rows {
		cards = append(cards, cardRow(storesqlc.FindCardRow(row)))
	}

	return cards, nil
}

func (s *Store) CreateCard(ctx context.Context, userID string, input board.CreateInput) (board.Card, error) {
	observability.Enrich(ctx, entity.FieldAdapter, entity.AdapterPostgres, entity.FieldAdapterOperation, entity.OperationCreateCard)

	var card board.Card

	err := s.WithinTransaction(ctx, func(ctx context.Context) error {
		queries := s.generated(ctx)

		id, err := queries.CreateCard(ctx, storesqlc.CreateCardParams{Title: input.Title, Description: input.Description, UserID: userID})
		if err != nil {
			return err
		}

		row, err := queries.FindCard(ctx, id)
		card = cardRow(row)

		return err
	})

	return card, err
}

func (s *Store) UpdateCard(ctx context.Context, userID, cardID string, input board.UpdateInput) (board.Card, error) {
	observability.Enrich(ctx, entity.FieldAdapter, entity.AdapterPostgres, entity.FieldAdapterOperation, entity.OperationUpdateCard)

	var card board.Card

	err := s.WithinTransaction(ctx, func(ctx context.Context) error {
		queries := s.generated(ctx)
		if err := lockCardVersion(ctx, queries, cardID, input.Version); err != nil {
			return err
		}

		count, err := queries.UpdateCard(ctx, storesqlc.UpdateCardParams{ID: cardID, Title: input.Title, Description: input.Description,
			Status: input.Status, Version: input.Version, UpdatedBy: userID})
		if err != nil {
			return err
		}

		if count != 1 {
			return board.ErrConflict
		}

		row, err := queries.FindCard(ctx, cardID)
		card = cardRow(row)

		return err
	})

	return card, err
}

func (s *Store) DeleteCard(ctx context.Context, cardID string, version int64) error {
	observability.Enrich(ctx, entity.FieldAdapter, entity.AdapterPostgres, entity.FieldAdapterOperation, entity.OperationDeleteCard)

	return s.WithinTransaction(ctx, func(ctx context.Context) error {
		queries := s.generated(ctx)
		if err := lockCardVersion(ctx, queries, cardID, version); err != nil {
			return err
		}

		count, err := queries.DeleteCard(ctx, storesqlc.DeleteCardParams{ID: cardID, Version: version})
		if err != nil {
			return err
		}

		if count != 1 {
			return board.ErrConflict
		}

		return nil
	})
}

// Holding the row lock makes the missing/stale decision and write one operation.
func lockCardVersion(ctx context.Context, queries *storesqlc.Queries, cardID string, expected int64) error {
	version, err := queries.LockCard(ctx, cardID)
	if errors.Is(err, pgx.ErrNoRows) {
		return board.ErrNotFound
	}

	if err != nil {
		return err
	}

	if version != expected {
		return board.ErrConflict
	}

	return nil
}

func cardRow(row storesqlc.FindCardRow) board.Card {
	return board.Card{ID: row.ID, Title: row.Title, Description: row.Description, Status: row.Status, Version: row.Version,
		CreatedBy: entity.User{ID: row.CreatedBy, Username: row.CreatorUsername}, UpdatedBy: entity.User{ID: row.UpdatedBy, Username: row.EditorUsername},
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt}
}
