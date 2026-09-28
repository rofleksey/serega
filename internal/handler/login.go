package handler

import (
	"errors"
	"net/http"

	"github.com/rofleksey/serega/internal/usecase/account"

	"github.com/rofleksey/serega/internal/api/generated"
)

func (h *handler) Login(w http.ResponseWriter, r *http.Request, _ generated.LoginParams) {
	if !h.validateJSONRequest(w, r) {
		return
	}

	var request generated.LoginRequest
	if !decodeJSON(w, r, &request) {
		return
	}

	currentAccount, err := h.accounts.Authenticate(r.Context(), request.Username, request.Password)
	if err != nil {
		if !errors.Is(err, account.ErrInvalidCredentials) {
			h.writeInternalError(w, r, "login_error", "could not sign in", err)
			return
		}

		writeError(w, r, http.StatusUnauthorized, "invalid_credentials", "invalid username or password")

		return
	}

	if err := h.sessions.RenewToken(r.Context()); err != nil {
		h.writeInternalError(w, r, "session_error", "could not create session", err)
		return
	}

	h.sessions.Put(r.Context(), userIDKey, currentAccount.ID)
	w.WriteHeader(http.StatusNoContent)
}
