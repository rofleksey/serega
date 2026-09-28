package handler

import "net/http"

func (h *handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	currentAccount, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	h.writeJSON(w, r, http.StatusOK, userResponse(currentAccount))
}
