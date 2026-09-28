// Package board owns the rules for Serega's single shared kanban board.
package board

import (
	"context"

	"github.com/rofleksey/serega/internal/entity"
	"github.com/rofleksey/serega/internal/observability"
)

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) List(ctx context.Context) ([]Card, error) {
	observability.Enrich(ctx, entity.FieldUseCase, entity.OperationListCards)
	cards, err := s.store.ListCards(ctx)
	observeResult(ctx, err)
	observability.Enrich(ctx, entity.FieldResultCount, len(cards))

	return cards, err
}
