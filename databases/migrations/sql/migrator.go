package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gofreego/goutils/databases"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator interface {
	Migrate(ctx context.Context) error
	Rollback(ctx context.Context) error
	MigrateTo(ctx context.Context, version uint) error
	Version() (uint, bool, error)
	Force(int) error
	Close() error
}

// NewMigrator creates a new Migrator instance.
// db is the primary database connection to use for migrations.
// path is the file system path to the migration files.
// dbType specifies the database type (postgres, mysql, etc.).
func NewMigrator(db *sql.DB, path string, dbType databases.DatabaseName) (Migrator, error) {
	var driver database.Driver
	var err error
	var databaseName string

	switch dbType {
	case databases.Postgres:
		driver, err = postgres.WithInstance(db, &postgres.Config{})
		databaseName = "postgres"
	case databases.MySQL:
		driver, err = mysql.WithInstance(db, &mysql.Config{})
		databaseName = "mysql"
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create %s driver: %w", dbType, err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", path),
		databaseName, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return &migratorImpl{
		db:      db,
		path:    path,
		migrate: m,
	}, nil
}

// NewPostgresMigrator creates a new Migrator instance for PostgreSQL.
// This is a convenience function for backward compatibility.
func NewPostgresMigrator(db *sql.DB, path string) (Migrator, error) {
	return NewMigrator(db, path, databases.Postgres)
}

// NewMySQLMigrator creates a new Migrator instance for MySQL.
// This is a convenience function.
func NewMySQLMigrator(db *sql.DB, path string) (Migrator, error) {
	return NewMigrator(db, path, databases.MySQL)
}

type migratorImpl struct {
	db      *sql.DB
	path    string
	migrate *migrate.Migrate
}

// Force will force the migration to a specific version.
// it will false the dirty state
func (m *migratorImpl) Force(version int) error {
	if err := m.migrate.Force(version); err != nil {
		return fmt.Errorf("failed to force migration to version %d: %w", version, err)
	}
	return nil
}

// Rollback reverts the migration to the previous version.
func (m *migratorImpl) Rollback(ctx context.Context) error {
	if err := m.migrate.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}
	return nil
}

// Migrate applies all available migrations.
func (m *migratorImpl) Migrate(ctx context.Context) error {
	if err := m.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

// MigrateTo migrates the database to a specific version.
func (m *migratorImpl) MigrateTo(ctx context.Context, version uint) error {
	if err := m.migrate.Migrate(version); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}
	return nil
}

// Version returns the current migration version and dirty state.
func (m *migratorImpl) Version() (uint, bool, error) {
	version, dirty, err := m.migrate.Version()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}
	return version, dirty, nil
}

// Close will close the db connections
func (m *migratorImpl) Close() error {
	sourceErr, dbErr := m.migrate.Close()
	if sourceErr != nil {
		return fmt.Errorf("failed to close source: %w", sourceErr)
	}
	if dbErr != nil {
		return fmt.Errorf("failed to close database: %w", dbErr)
	}
	return nil
}
