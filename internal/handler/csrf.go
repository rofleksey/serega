package handler

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func (h *handler) GetCsrfToken(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, r, http.StatusOK, map[string]string{"token": nosurf.Token(r)})
}
