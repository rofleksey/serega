package frontend

import (
	"io/fs"
	"net/http"
)

// NewHandler serves the frontend embedded in the current executable.
func NewHandler() http.Handler {
	files, err := fs.Sub(assets, "dist")
	if err != nil {
		panic("open embedded frontend: " + err.Error())
	}

	return &spaHandler{
		files:      files,
		fileServer: http.FileServer(http.FS(files)),
	}
}
