package board

import (
	"context"

	"github.com/rofleksey/serega/internal/entity"
	"github.com/rofleksey/serega/internal/observability"
)

func (s *Service) Update(ctx context.Context, userID, cardID string, input UpdateInput) (Card, error) {
	observability.Enrich(ctx, entity.FieldUseCase, entity.OperationUpdateCard)

	if err := validateIdentity(cardID, input.Version); err != nil {
		observeResult(ctx, err)
		return Card{}, err
	}

	observability.Enrich(ctx, entity.FieldLogCardID, cardID)

	input.Title = normalizeTitle(input.Title)
	if err := validateText(input.Title, input.Description); err != nil {
		observeResult(ctx, err)
		return Card{}, err
	}

	switch input.Status {
	case entity.CardStatusTodo, entity.CardStatusDoing, entity.CardStatusDone:
	default:
		err := &ValidationError{Field: entity.FieldStatus, Message: "must be todo, doing, or done"}
		observeResult(ctx, err)

		return Card{}, err
	}

	card, err := s.store.UpdateCard(ctx, userID, cardID, input)
	observeResult(ctx, err)

	return card, err
}
