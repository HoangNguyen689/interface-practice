package sqlmigration

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/HoangNguyen689/interface-practice/pkg/postgresql"
)

type Migrator interface {
	Execute(ctx context.Context, version uint) error
}

type migrator struct {
	client             *postgresql.Client
	sourceDir          string
	databaseName       string
	migrationTableName string
	logger             *log.Logger
}

func NewMigrator(client *postgresql.Client, sourceDir, databaseName string, logger *log.Logger) Migrator {
	m := &migrator{
		client:             client,
		sourceDir:          fmt.Sprintf("file://%s", sourceDir),
		databaseName:       databaseName,
		migrationTableName: "schema_migrations",
		logger:             log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

	return m
}

func (m *migrator) Execute(ctx context.Context, version uint) (err error) {
	db := stdlib.OpenDBFromPool(m.client.Pool)

	driver, err := pgx.WithInstance(db, &pgx.Config{
		MigrationsTable: m.migrationTableName,
		DatabaseName:    m.databaseName,
	})
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithDatabaseInstance(m.sourceDir, m.databaseName, driver)
	if err != nil {
		return err
	}

	// Show current version.
	if err := showVersion(ctx, migrator, m.logger); err != nil {
		return err
	}

	// Start migrating.
	if version > 0 {
		m.logger.Println(fmt.Sprintf("start migrating to version %d", version))
		err = migrator.Migrate(version)
	} else {
		m.logger.Println("start migrating to latest version")
		err = migrator.Up()
	}

	if errors.Is(err, migrate.ErrNoChange) {
		m.logger.Println("there was no changes since the last apply")

		return nil
	}
	if err != nil {
		return err
	}

	m.logger.Println("successfully executed migration task")

	// Show new version.
	if err := showVersion(ctx, migrator, m.logger); err != nil {
		return err
	}

	return nil
}

func showVersion(ctx context.Context, m *migrate.Migrate, logger *log.Logger) error {
	curVer, dirty, err := m.Version()

	if errors.Is(err, migrate.ErrNilVersion) {
		logger.Println("no migration has been applied before, this will be the first one")

		return nil
	}

	if err != nil {
		logger.Println("failed to check current version")
		logger.Println(err)

		return err
	}

	logger.Println(fmt.Sprintf("current version: %d, dirty: %t", curVer, dirty))

	return nil
}
