// Package frontend serves the immutable Vite build embedded in Serega.
package frontend

import "embed"

// assets contains the immutable production frontend built by Vite.
//
//go:embed all:dist
var assets embed.FS
