package main

import (
	"context"
	"log"
	"path/filepath"

	"github.com/gofreego/goutils/databases"
	"github.com/gofreego/goutils/databases/connections/sql/pgsql"
	"github.com/gofreego/goutils/databases/migrations/sql"
)

func main() {
	// Connect to your database
	db, err := pgsql.GetConnection(&pgsql.Config{
		Host:     "localhost",
		Port:     5432,
		Username: "admin",
		Password: "******",
		DBName:   "*****",
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Get the path to migration files
	migrationPath := filepath.Join(".", "sql/")

	// Create a new migrator instance
	migrator, err := sql.NewPostgresMigrator(db, &sql.Config{
		Path:         migrationPath,
		DBType:       databases.Postgres,
		Action:       sql.ActionUp,
		ForceVersion: 0,
	})
	if err != nil {
		log.Fatal("Failed to create migrator:", err)
	}
	defer migrator.Close()

	// Run migrations
	ctx := context.Background()
	if err := migrator.Run(ctx); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Check current version
	version, dirty, err := migrator.Version()
	if err != nil {
		log.Fatal("Failed to get version:", err)
	}
	log.Printf("Current migration version after migrate to 2: %d, dirty: %t", version, dirty)
}
