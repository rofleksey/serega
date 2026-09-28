package handler

import (
	"net/http"

	"github.com/rofleksey/serega/internal/api/generated"
)

func (h *handler) ListCards(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireUser(w, r); !ok {
		return
	}

	cards, err := h.board.List(r.Context())
	if err != nil {
		h.writeBoardError(w, r, err)
		return
	}

	response := make([]generated.Card, 0, len(cards))
	for _, card := range cards {
		response = append(response, cardResponse(card))
	}

	h.writeJSON(w, r, http.StatusOK, generated.CardList{Cards: response})
}
