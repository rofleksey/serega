package middleware

import (
	"bufio"
	"net"
)

type hijackingWriter struct {
	connection net.Conn
	peer       net.Conn
	reader     *bufio.ReadWriter
}
