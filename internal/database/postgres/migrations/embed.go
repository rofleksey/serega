// Package migrations embeds the PostgreSQL Goose migrations shipped with Serega.
package migrations

import "embed"

// Files contains migration SQL files rooted at this package directory.
//
//go:embed *.sql
var Files embed.FS
