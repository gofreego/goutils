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
	Version() (uint, bool, error)
	Run(ctx context.Context) error
	Close() error
}

type Action string

const (
	ActionUp        Action = "up"
	ActionDown      Action = "down"
	ActionMigrateTo Action = "migrate_to"
	ActionForce     Action = "force"
	ActionVersion   Action = "version"
)

type Config struct {
	Path         string                 `yaml:"Path" json:"path"`
	DBType       databases.DatabaseName `yaml:"DBType" json:"dbType"`
	Action       Action                 `yaml:"Action" json:"action"`
	ForceVersion int                    `yaml:"ForceVersion" json:"forceVersion"`
}

// NewMigrator creates a new Migrator instance.
// db is the primary database connection to use for migrations.
// path is the file system path to the migration files.
// dbType specifies the database type (postgres, mysql, etc.).
func NewMigrator(db *sql.DB, cfg *Config) (Migrator, error) {
	var driver database.Driver
	var err error
	var databaseName string

	switch cfg.DBType {
	case databases.Postgres:
		driver, err = postgres.WithInstance(db, &postgres.Config{})
		databaseName = "postgres"
	case databases.MySQL:
		driver, err = mysql.WithInstance(db, &mysql.Config{})
		databaseName = "mysql"
	default:
		return nil, fmt.Errorf("unsupported database type: %s, expected: %v", cfg.DBType, []databases.DatabaseName{databases.Postgres, databases.MySQL})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create %s driver: %w", cfg.DBType, err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.Path),
		databaseName, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return &migratorImpl{
		db:      db,
		cfg:     cfg,
		migrate: m,
	}, nil
}

// NewPostgresMigrator creates a new Migrator instance for PostgreSQL.
// This is a convenience function for backward compatibility.
func NewPostgresMigrator(db *sql.DB, cfg *Config) (Migrator, error) {
	return NewMigrator(db, cfg)
}

// NewMySQLMigrator creates a new Migrator instance for MySQL.
// This is a convenience function.
func NewMySQLMigrator(db *sql.DB, cfg *Config) (Migrator, error) {
	return NewMigrator(db, cfg)
}

type migratorImpl struct {
	db      *sql.DB
	cfg     *Config
	migrate *migrate.Migrate
}

// Run implements Migrator.
func (m *migratorImpl) Run(ctx context.Context) error {
	switch m.cfg.Action {
	case ActionUp:
		return m.up()
	case ActionDown:
		return m.down()
	case ActionMigrateTo:
		return m.migrateTo(uint(m.cfg.ForceVersion))
	case ActionForce:
		return m.force(m.cfg.ForceVersion)
	default:
		return fmt.Errorf("unknown action: %s, expected one of: %v", m.cfg.Action, []Action{ActionUp, ActionDown, ActionMigrateTo, ActionForce})
	}
}

// Force will force the migration to a specific version.
// it will false the dirty state
func (m *migratorImpl) force(version int) error {
	if err := m.migrate.Force(version); err != nil {
		return fmt.Errorf("failed to force migration to version %d: %w", version, err)
	}
	return nil
}

// down reverts the migration to the previous version.
func (m *migratorImpl) down() error {
	if err := m.migrate.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to down migration: %w", err)
	}
	return nil
}

// Migrate applies all available migrations.
func (m *migratorImpl) up() error {
	if err := m.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

// MigrateTo migrates the database to a specific version.
func (m *migratorImpl) migrateTo(version uint) error {
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
