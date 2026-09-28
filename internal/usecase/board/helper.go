package board

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/rofleksey/serega/internal/entity"
	"github.com/rofleksey/serega/internal/observability"
)

func normalizeTitle(title string) string { return strings.TrimSpace(title) }

func validateText(title, description string) error {
	if !utf8.ValidString(title) || utf8.RuneCountInString(title) < 1 || utf8.RuneCountInString(title) > 200 {
		return &ValidationError{Field: entity.FieldTitle, Message: "must be between 1 and 200 characters"}
	}

	if !utf8.ValidString(description) || utf8.RuneCountInString(description) > 5000 {
		return &ValidationError{Field: entity.FieldDescription, Message: "must be at most 5000 characters"}
	}

	return nil
}

func validateIdentity(cardID string, version int64) error {
	parsed, err := uuid.Parse(cardID)
	if err != nil || parsed.String() != cardID {
		return &ValidationError{Field: entity.FieldCardID, Message: entity.ValidationMessageCanonicalUUID}
	}

	if version < 1 {
		return &ValidationError{Field: entity.FieldVersion, Message: "must be a positive integer"}
	}

	return nil
}

func observeResult(ctx context.Context, err error) {
	decision := entity.DecisionAccepted

	var invalid *ValidationError
	switch {
	case err == nil:
	case errors.As(err, &invalid):
		decision = entity.DecisionValidationFailed
	case errors.Is(err, ErrNotFound):
		decision = entity.DecisionTargetNotFound
	case errors.Is(err, ErrConflict):
		decision = entity.DecisionConflict
	default:
		decision = entity.DecisionPersistenceFailed
	}

	observability.Enrich(ctx, entity.FieldDecision, decision)
}
