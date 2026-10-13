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

	w.Header().Set("Cache-Control", "no-cache")

	requested := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
	if requested != "." && requested != "" {
		if info, err := fs.Stat(h.files, requested); err == nil && !info.IsDir() {
			if strings.HasPrefix(requested, "assets/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}

			h.fileServer.ServeHTTP(w, r)

			return
		}
	}

	if strings.HasPrefix(requested, "assets/") || path.Ext(requested) != "" {
		http.NotFound(w, r)

		return
	}

	indexRequest := r.Clone(r.Context())
	indexRequest.URL.Path = "/"
	h.fileServer.ServeHTTP(w, indexRequest)
}
