package handler

import (
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/rofleksey/serega/internal/usecase/account"
	"github.com/rofleksey/serega/internal/usecase/board"
)

const sessionExpiredHTTPStatus = 419

type HandlerDependencies struct {
	Logger        *slog.Logger
	SessionStore  scs.Store
	Accounts      *account.Service
	Board         *board.Service
	SecureCookies bool
}

type handler struct {
	logger       *slog.Logger
	accounts     *account.Service
	sessions     *scs.SessionManager
	board        *board.Service
	validateBody func(http.ResponseWriter, *http.Request) bool
}
