package handler

import (
	"errors"
	"net/http"

	"github.com/rofleksey/serega/internal/entity"
	"github.com/rofleksey/serega/internal/observability"
	"github.com/rofleksey/serega/internal/usecase/account"
)

func (h *handler) currentUser(r *http.Request) (entity.User, error) {
	id := h.sessions.GetString(r.Context(), userIDKey)
	if id == "" {
		observability.Enrich(r.Context(), entity.FieldAuthMechanism, entity.AuthMechanismSession, entity.FieldAuthDecision, entity.DecisionDenied)
		return entity.User{}, account.ErrUserNotFound
	}

	currentAccount, err := h.accounts.User(r.Context(), id)

	decision := entity.DecisionAllowed
	if err != nil {
		decision = entity.DecisionDenied
	}

	observability.Enrich(r.Context(), entity.FieldAuthMechanism, entity.AuthMechanismSession, entity.FieldAuthDecision, decision, entity.FieldAccountID, id)

	return currentAccount, err
}

func (h *handler) requireUser(w http.ResponseWriter, r *http.Request) (entity.User, bool) {
	user, err := h.currentUser(r)
	if errors.Is(err, account.ErrUserNotFound) {
		writeError(w, r, http.StatusUnauthorized, "unauthorized", "authentication is required")
		return entity.User{}, false
	}

	if err != nil {
		h.writeInternalError(w, r, "account_error", "could not load account", err)
		return entity.User{}, false
	}

	return user, true
}
