package board

import (
	"context"

	"github.com/rofleksey/serega/internal/entity"
	"github.com/rofleksey/serega/internal/observability"
)

func (s *Service) Delete(ctx context.Context, cardID string, version int64) error {
	observability.Enrich(ctx, entity.FieldUseCase, entity.OperationDeleteCard)

	if err := validateIdentity(cardID, version); err != nil {
		observeResult(ctx, err)
		return err
	}

	observability.Enrich(ctx, entity.FieldLogCardID, cardID)

	err := s.store.DeleteCard(ctx, cardID, version)
	observeResult(ctx, err)

	return err
}
