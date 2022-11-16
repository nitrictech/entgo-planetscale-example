package migrations

import (
	"embed"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var migrationFiles embed.FS

func MigrationFS() (source.Driver, error) {
	return iofs.New(migrationFiles, ".")
}
