package handler

import (
	"net/http"

	"github.com/rofleksey/serega/internal/api/generated"
)

func (h *handler) DeleteCard(w http.ResponseWriter, r *http.Request, cardID string, params generated.DeleteCardParams) {
	if _, ok := h.requireUser(w, r); !ok {
		return
	}

	if err := h.board.Delete(r.Context(), cardID, params.Version); err != nil {
		h.writeBoardError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
