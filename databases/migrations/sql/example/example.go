package main

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"

	sqlmigrations "github.com/gofreego/goutils/databases/migrations/sql"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Connect to your database
	db, err := sql.Open("postgres", "postgres://user:password@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Get the path to migration files
	migrationPath := filepath.Join(".", "sql")

	// Create a new migrator instance
	migrator, err := sqlmigrations.NewPostgresMigrator(db, migrationPath)
	if err != nil {
		log.Fatal("Failed to create migrator:", err)
	}
	defer migrator.Close()

	// Run migrations
	ctx := context.Background()
	if err := migrator.Migrate(ctx); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Check current version
	version, dirty, err := migrator.Version()
	if err != nil {
		log.Fatal("Failed to get version:", err)
	}

	log.Printf("Current migration version: %d, dirty: %t", version, dirty)

	// Example of rolling back one step
	// if err := migrator.Rollback(ctx); err != nil {
	//     log.Fatal("Failed to rollback:", err)
	// }

	// Example of migrating to a specific version
	// if err := migrator.MigrateTo(ctx, 2); err != nil {
	//     log.Fatal("Failed to migrate to version 2:", err)
	// }
}
