// Package middleware contains reusable HTTP middleware for Serega.
package middleware

import (
	"net/http"
)

type contextKey uint8

const clientIPContextKey contextKey = iota

type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}
