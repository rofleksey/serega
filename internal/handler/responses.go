package handler

import (
	"errors"
	"net/http"

	"github.com/rofleksey/serega/internal/api/generated"
	"github.com/rofleksey/serega/internal/entity"
	"github.com/rofleksey/serega/internal/usecase/board"
)

func userResponse(user entity.User) generated.User {
	return generated.User{ID: user.ID, Username: user.Username}
}

func cardResponse(card entity.Card) generated.Card {
	return generated.Card{ID: card.ID, Title: card.Title, Description: card.Description, Status: generated.CardStatus(card.Status), Version: card.Version,
		CreatedBy: userResponse(card.CreatedBy), UpdatedBy: userResponse(card.UpdatedBy), CreatedAt: card.CreatedAt, UpdatedAt: card.UpdatedAt}
}

func (h *handler) writeBoardError(w http.ResponseWriter, r *http.Request, err error) {
	var invalid *board.ValidationError
	switch {
	case errors.As(err, &invalid):
		writeError(w, r, http.StatusBadRequest, "invalid_request", "card is invalid", generated.FieldError{Field: invalid.Field, Message: invalid.Message})
	case errors.Is(err, board.ErrNotFound):
		writeError(w, r, http.StatusNotFound, "not_found", "card not found")
	case errors.Is(err, board.ErrConflict):
		writeError(w, r, http.StatusConflict, "card_conflict", "card changed; refresh and try again")
	default:
		h.writeInternalError(w, r, "board_error", "could not update the board", err)
	}
}
