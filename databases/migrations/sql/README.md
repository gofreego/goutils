# SQL Migrator

A Go library for database migrations using [golang-migrate/migrate](https://github.com/golang-migrate/migrate).

## Features

- Support for PostgreSQL and MySQL databases
- Migration up and down operations
- Rollback functionality
- Migration to specific versions
- Version tracking and dirty state detection
- Clean resource management

## Installation

```bash
go get github.com/gofreego/goutils/databases/migrations/sql
```

## Usage

### Basic Usage with PostgreSQL

```go
package main

import (
    "context"
    "database/sql"
    "log"
    
    sqlmigrations "github.com/gofreego/goutils/databases/migrations/sql"
    _ "github.com/lib/pq"
)

func main() {
    // Connect to database
    db, err := sql.Open("postgres", "postgres://user:password@localhost/dbname?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create migrator
    migrator, err := sqlmigrations.NewPostgresMigrator(db, "./migrations")
    if err != nil {
        log.Fatal(err)
    }
    defer migrator.Close()

    // Run migrations
    ctx := context.Background()
    if err := migrator.Migrate(ctx); err != nil {
        log.Fatal(err)
    }
}
```

### Using with MySQL

```go
migrator, err := sqlmigrations.NewMySQLMigrator(db, "./migrations")
```

### Generic Usage

```go
migrator, err := sqlmigrations.NewMigrator(db, "./migrations", sqlmigrations.MySQL)
```

## Migration Files

Migration files should follow the naming convention:
- `{version}_{description}.up.sql` - for applying migrations
- `{version}_{description}.down.sql` - for rolling back migrations

Example:
- `000001_create_users_table.up.sql`
- `000001_create_users_table.down.sql`

## API Reference

### Interface

```go
type Migrator interface {
    Migrate(ctx context.Context) error           // Apply all pending migrations
    Rollback(ctx context.Context) error          // Rollback one migration
    MigrateTo(ctx context.Context, version uint) error // Migrate to specific version
    Version() (uint, bool, error)                // Get current version and dirty state
    Close() error                                // Clean up resources
}
```

### Database Types

```go
const (
    Postgres DatabaseType = "postgres"
    MySQL    DatabaseType = "mysql"
)
```

### Constructor Functions

- `NewMigrator(db *sql.DB, path string, dbType DatabaseType) (Migrator, error)`
- `NewPostgresMigrator(db *sql.DB, path string) (Migrator, error)`
- `NewMySQLMigrator(db *sql.DB, path string) (Migrator, error)`

## Examples

See the `example/` directory for a complete working example with sample migration files.

## Error Handling

The migrator properly handles:
- `migrate.ErrNoChange` - when no migrations need to be applied
- Database connection issues
- Invalid migration files
- Version conflicts

## Dependencies

- [golang-migrate/migrate/v4](https://github.com/golang-migrate/migrate)
- Database drivers (postgres, mysql)
