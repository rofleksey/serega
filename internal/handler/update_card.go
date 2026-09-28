package handler

import (
	"net/http"

	"github.com/rofleksey/serega/internal/api/generated"
	"github.com/rofleksey/serega/internal/usecase/board"
)

func (h *handler) UpdateCard(w http.ResponseWriter, r *http.Request, cardID string, _ generated.UpdateCardParams) {
	user, ok := h.requireUser(w, r)
	if !ok || !h.validateJSONRequest(w, r) {
		return
	}

	var request generated.UpdateCardRequest
	if !decodeJSON(w, r, &request) {
		return
	}

	card, err := h.board.Update(r.Context(), user.ID, cardID, board.UpdateInput{Title: request.Title, Description: request.Description, Status: string(request.Status), Version: request.Version})
	if err != nil {
		h.writeBoardError(w, r, err)
		return
	}

	h.writeJSON(w, r, http.StatusOK, cardResponse(card))
}
