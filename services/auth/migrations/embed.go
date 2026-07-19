// Package migrations embeds the raw SQL migration files so the compiled
// auth service binary can self-migrate on startup without external tooling.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
