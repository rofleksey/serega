package frontend

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.NotFound(w, r)

		return
	}

	requested := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
	if requested != "." && requested != "" {
		if info, err := fs.Stat(h.files, requested); err == nil && !info.IsDir() {
			h.fileServer.ServeHTTP(w, r)

			return
		}
	}

	indexRequest := r.Clone(r.Context())
	indexRequest.URL.Path = "/"
	h.fileServer.ServeHTTP(w, indexRequest)
}
