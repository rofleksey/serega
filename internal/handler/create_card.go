package handler

import (
	"net/http"

	"github.com/rofleksey/serega/internal/api/generated"
	"github.com/rofleksey/serega/internal/usecase/board"
)

func (h *handler) CreateCard(w http.ResponseWriter, r *http.Request, _ generated.CreateCardParams) {
	user, ok := h.requireUser(w, r)
	if !ok || !h.validateJSONRequest(w, r) {
		return
	}

	var request generated.CreateCardRequest
	if !decodeJSON(w, r, &request) {
		return
	}

	card, err := h.board.Create(r.Context(), user.ID, board.CreateInput{Title: request.Title, Description: request.Description})
	if err != nil {
		h.writeBoardError(w, r, err)
		return
	}

	h.writeJSON(w, r, http.StatusCreated, cardResponse(card))
}
