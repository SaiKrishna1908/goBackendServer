package migrations

import (
	"embed"
	_ "embed"
)

//go:embed *.sql
var FS embed.FS
