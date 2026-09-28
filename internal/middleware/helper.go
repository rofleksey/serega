package middleware

import (
	"bufio"
	"mime"
	"net"
	"net/http"

	httpapi "github.com/rofleksey/serega/internal/api"
	"github.com/rofleksey/serega/internal/entity"
)

func unsafeMethod(method string) bool {
	return method != http.MethodGet && method != http.MethodHead && method != http.MethodOptions
}

func trusted(peer net.IP, networks []*net.IPNet) bool {
	for _, network := range networks {
		if network.Contains(peer) {
			return true
		}
	}

	return false
}

func jsonContentType(value string) bool {
	mediaType, _, err := mime.ParseMediaType(value)

	return err == nil && mediaType == "application/json"
}

func writeError(w http.ResponseWriter, status int, code, message, requestID string) {
	httpapi.WriteError(w, status, code, message, requestID)
}

func (w *statusWriter) WriteHeader(status int) {
	if w.status == 0 {
		w.status = status
	}

	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}

	n, err := w.ResponseWriter.Write(data)
	w.bytes += n

	return n, err
}

func (w *statusWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *statusWriter) Flush() {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}

	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *statusWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}

	return hijacker.Hijack()
}

func httpOutcome(status int) string {
	switch {
	case status >= http.StatusInternalServerError:
		return entity.OutcomeFailed
	case status >= http.StatusBadRequest:
		return entity.OutcomeRejected
	default:
		return entity.OutcomeSucceeded
	}
}
