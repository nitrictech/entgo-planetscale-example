package main

import (
	"context"
	"database/sql"
	"time"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/go-sql-driver/mysql"
	sqlmysql "github.com/go-sql-driver/mysql"
	gomigrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/nitrictech/entgo-planetscale-example/ent"
	"github.com/nitrictech/entgo-planetscale-example/ent/migrate"
	"github.com/nitrictech/entgo-planetscale-example/ent/migrate/migrations"
)

func createMigration(name string) error {
	if name == "" {
		return errors.New("migration name is required. Use: 'go run ./cmd/migration <name>'")
	}

	dir, err := atlas.NewLocalDir("ent/migrate/migrations")
	if err != nil {
		return errors.WithMessage(err, "failed creating atlas migration directory")
	}

	opts := []schema.MigrateOption{
		schema.WithDir(dir),                         // provide migration directory
		schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
		schema.WithDialect(dialect.MySQL),           // Ent dialect to use
		schema.WithForeignKeys(false),               // planetscale uses https://vitess.io/ that requires foreign keys off
		schema.WithDropColumn(true),
	}

	// Generate migrations using Atlas support for MySQL (note the Ent dialect option passed above).
	return migrate.NamedDiff(context.TODO(), "mysql://root:pass@localhost:3306/deploy-test", name, opts...)
}

func mysqlConnectAndMigrate(dsn string, migrate bool) (*ent.Client, error) {
	pDSN, err := sqlmysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	pDSN.ParseTime = true
	pDSN.Loc = time.Local

	if pDSN.Params == nil {
		pDSN.Params = map[string]string{}
	}

	pDSN.Params["tls"] = "true"
	pDSN.Params["charset"] = "utf8mb4"
	if migrate {
		pDSN.Params["multiStatements"] = "true"
	}

	dsn = pDSN.FormatDSN()

	db, err := sql.Open(dialect.MySQL, dsn)
	if err != nil {
		return nil, err
	}

	if migrate {
		d, err := migrations.MigrationFS()
		if err != nil {
			return nil, errors.WithMessage(err, "iofs.New")
		}

		m, err := gomigrate.NewWithSourceInstance("iofs", d, "mysql://"+dsn)
		if err != nil {
			return nil, errors.WithMessage(err, "NewWithSourceInstance")
		}

		if err := m.Up(); err != nil {
			if !errors.Is(err, gomigrate.ErrNoChange) {
				return nil, errors.WithMessage(err, "db migrations update")
			}
		}
	}

	return ent.NewClient(
		ent.Driver(entsql.NewDriver(
			dialect.MySQL,
			entsql.Conn{ExecQuerier: db})),
		ent.Log(logrus.Info)), nil
}
