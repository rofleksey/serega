package frontend

import (
	"io/fs"
	"net/http"
)

type spaHandler struct {
	files      fs.FS
	fileServer http.Handler
}
