package handler

import (
	"net/http"

	"github.com/rofleksey/serega/internal/api/generated"
)

func (h *handler) Logout(w http.ResponseWriter, r *http.Request, _ generated.LogoutParams) {
	if _, ok := h.requireUser(w, r); !ok {
		return
	}

	if err := h.sessions.Destroy(r.Context()); err != nil {
		h.writeInternalError(w, r, "session_error", "could not destroy session", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
