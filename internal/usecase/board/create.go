package board

import (
	"context"

	"github.com/rofleksey/serega/internal/entity"
	"github.com/rofleksey/serega/internal/observability"
)

func (s *Service) Create(ctx context.Context, userID string, input CreateInput) (Card, error) {
	observability.Enrich(ctx, entity.FieldUseCase, entity.OperationCreateCard)

	input.Title = normalizeTitle(input.Title)
	if err := validateText(input.Title, input.Description); err != nil {
		observeResult(ctx, err)
		return Card{}, err
	}

	card, err := s.store.CreateCard(ctx, userID, input)
	observeResult(ctx, err)

	return card, err
}
