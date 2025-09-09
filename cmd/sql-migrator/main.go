package main

import (
	"context"
	std_sql "database/sql"
	"fmt"
	"os"

	"github.com/gofreego/goutils/configutils"
	"github.com/gofreego/goutils/databases"
	"github.com/gofreego/goutils/databases/connections/sql"
	"github.com/gofreego/goutils/databases/connections/sql/pgsql"
	migrator "github.com/gofreego/goutils/databases/migrations/sql"
	"github.com/gofreego/goutils/logger"
)

type Config struct {
	Repository sql.Config      `yaml:"Repository" json:"repository"`
	Migrator   migrator.Config `yaml:"Migrator" json:"migrator"`
}

type SQLMigrator struct {
	migrator migrator.Migrator
	cfg      migrator.Config
}

func NewSQLMigrator(cfg *Config) *SQLMigrator {
	var db *std_sql.DB
	var err error
	switch cfg.Repository.Name {
	case databases.Postgres:
		db, err = pgsql.GetConnection(&cfg.Repository.Postgres.Primary)
		if err != nil {
			panic("failed to get Postgres connection, err:" + err.Error())
		}
	default:
		panic(fmt.Sprintf("unsupported database for migration: %s, expected: %s", cfg.Repository.Name, databases.Postgres))
	}

	cfg.Migrator.DBType = cfg.Repository.Name

	migrator, err := migrator.NewMigrator(db, &cfg.Migrator)
	if err != nil {
		panic(fmt.Sprintf("failed to create SQL migrator, err: %s", err.Error()))
	}
	return &SQLMigrator{migrator: migrator, cfg: cfg.Migrator}
}

// Run implements apputils.Application.
func (s *SQLMigrator) Run(ctx context.Context) error {
	defer s.migrator.Close()

	var err error
	if s.cfg.Action == migrator.ActionVersion {
		version, dirty, err := s.migrator.Version()
		if err != nil {
			logger.Error(ctx, "Failed to get database version: %s", err.Error())
			return err
		}
		logger.Info(ctx, "Database version: %d, dirty: %t, press ctrl+c to exit", version, dirty)
		return nil
	}

	err = s.migrator.Run(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to migrate database: %s", err.Error())
		return err
	}
	version, dirty, err := s.migrator.Version()
	if err != nil {
		logger.Error(ctx, "Failed to get database version: %s", err.Error())
		return err
	}
	logger.Info(ctx, "Database version: %d, dirty: %t", version, dirty)
	logger.Info(ctx, "✅ Database migration completed successfully, press ctrl+c to exit")
	return nil
}

// Shutdown implements apputils.Application.
func (s *SQLMigrator) Shutdown(ctx context.Context) {
	// Close the migrator connection
	if err := s.migrator.Close(); err != nil {
		// Log the error but don't panic during shutdown
		logger.Warn(ctx, "Warning: failed to close migrator: %s", err.Error())
	}
}

func main() {
	configPath := "./migrator.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	var cfg Config
	ctx := context.Background()
	err := configutils.ReadConfig(ctx, configPath, &cfg)
	if err != nil {
		panic("failed to read config, from " + configPath + ", err: " + err.Error())
	}
	app := NewSQLMigrator(&cfg)
	if err := app.Run(ctx); err != nil {
		panic("failed to run SQL migrator, err: " + err.Error())
	}
}
