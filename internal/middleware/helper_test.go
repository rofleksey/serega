package middleware

import (
	"bufio"
	"net"
	"net/http"
	"testing"
)

func newHijackingWriter(t *testing.T) *hijackingWriter {
	t.Helper()

	connection, peer := net.Pipe()

	return &hijackingWriter{
		connection: connection,
		peer:       peer,
		reader:     bufio.NewReadWriter(bufio.NewReader(connection), bufio.NewWriter(connection)),
	}
}

func (writer *hijackingWriter) Header() http.Header {
	return make(http.Header)
}

func (writer *hijackingWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func (writer *hijackingWriter) WriteHeader(_ int) {}

func (writer *hijackingWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return writer.connection, writer.reader, nil
}

func (writer *hijackingWriter) closePeer() {
	_ = writer.peer.Close()
}
