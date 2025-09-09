# SQL Migrator

A command-line tool for managing SQL database migrations using the [golang-migrate](https://github.com/golang-migrate/migrate) library. This tool supports PostgreSQL and MySQL databases and provides a simple interface for running database schema migrations.

## Features

- **Multi-database support**: PostgreSQL and MySQL
- **Multiple migration actions**: Up, Down, Migrate to specific version, Force version, and Version check
- **Configuration-based**: YAML configuration file for easy management
- **Safe migrations**: Built on top of the battle-tested golang-migrate library
- **Graceful shutdown**: Proper cleanup and connection handling

## Installation

### From Source

```bash
go install github.com/gofreego/goutils/cmd/sql-migrator@latest
```

### Build Locally

```bash
git clone https://github.com/gofreego/goutils.git
cd goutils/cmd/sql-migrator
go build -o sql-migrator .
```

## Usage

### Basic Usage

```bash
# Use default config file (./migrator.yaml)
./sql-migrator

# Specify custom config file
./sql-migrator /path/to/your/migrator.yaml
```

### Configuration

Create a `migrator.yaml` configuration file:

```yaml
Repository:
  name: postgres  # or mysql
  postgres:
    primary:
      host: localhost
      port: 5432
      username: your_username
      password: your_password
      dbname: your_database
      sslmode: disable
      # Optional: Additional connection parameters
      # max_open_conns: 10
      # max_idle_conns: 5
      # conn_max_lifetime: 3600s

Migrator:
  path: "./migrations"       # Path to migration files directory
  action: up                 # Migration action to perform
  force_version: 0          # Version number for force/migrate_to actions
```

#### MySQL Configuration Example

```yaml
Repository:
  name: mysql
  mysql:
    primary:
      host: localhost
      port: 3306
      username: your_username
      password: your_password
      dbname: your_database
      # Optional: Additional MySQL parameters
      # charset: utf8mb4
      # parseTime: true
      # loc: Local

Migrator:
  path: "./migrations"
  action: up
  force_version: 0
```

## Migration Actions

### Available Actions

| Action | Description | Configuration |
|--------|-------------|---------------|
| `up` | Apply all pending migrations | `action: up` |
| `down` | Rollback the last migration | `action: down` |
| `migrate_to` | Migrate to a specific version | `action: migrate_to`<br>`force_version: <target_version>` |
| `force` | Force database to a specific version (fixes dirty state) | `action: force`<br>`force_version: <version>` |
| `version` | Show current database version and dirty state | `action: version` |

### Examples

#### Apply All Migrations
```yaml
Migrator:
  path: "./migrations"
  action: up
```

#### Rollback Last Migration
```yaml
Migrator:
  path: "./migrations"
  action: down
```

#### Migrate to Specific Version
```yaml
Migrator:
  path: "./migrations"
  action: migrate_to
  force_version: 5
```

#### Force Version (Fix Dirty State)
```yaml
Migrator:
  path: "./migrations"
  action: force
  force_version: 3
```

#### Check Current Version
```yaml
Migrator:
  path: "./migrations"
  action: version
```

## Migration Files

Migration files must follow the golang-migrate naming convention:

```
{version}_{description}.up.sql    # Forward migration
{version}_{description}.down.sql  # Reverse migration
```

### Example Migration Files

#### `000001_create_users_table.up.sql`
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

#### `000001_create_users_table.down.sql`
```sql
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

#### `000002_create_posts_table.up.sql`
```sql
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
```

#### `000002_create_posts_table.down.sql`
```sql
DROP INDEX IF EXISTS idx_posts_user_id;
DROP TABLE IF EXISTS posts;
```

## Directory Structure

Organize your project with the following structure:

```
your-project/
├── migrator.yaml              # Configuration file
├── migrations/                # Migration files directory
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_posts_table.up.sql
│   └── 000002_create_posts_table.down.sql
└── sql-migrator              # Binary (after build)
```

## Error Handling

### Common Issues and Solutions

1. **Dirty State**: If a migration fails partway through, the database may be in a "dirty" state.
   ```yaml
   # Fix by forcing to a known good version
   Migrator:
     action: force
     force_version: 1  # Last known good version
   ```

2. **Connection Issues**: Verify database connectivity and credentials in your configuration.

3. **Missing Migration Files**: Ensure all migration files are present and follow the correct naming convention.

4. **Permission Issues**: Verify the database user has the necessary permissions to create/modify schema objects.

## Best Practices

1. **Version Control**: Always commit migration files to version control
2. **Backup**: Create database backups before running migrations in production
3. **Test Migrations**: Test both up and down migrations in development environments
4. **Incremental Changes**: Keep migrations small and focused on single changes
5. **Naming Convention**: Use descriptive names for migrations
6. **Review Process**: Implement code review for migration files

## Environment-Specific Configurations

### Development
```yaml
Repository:
  name: postgres
  postgres:
    primary:
      host: localhost
      port: 5432
      username: dev_user
      password: dev_password
      dbname: myapp_dev
      sslmode: disable
```

### Production
```yaml
Repository:
  name: postgres
  postgres:
    primary:
      host: prod-db.example.com
      port: 5432
      username: app_user
      password: ${DB_PASSWORD}  # Use environment variables
      dbname: myapp_prod
      sslmode: require
```

## Integration with CI/CD

### Docker Example
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o sql-migrator ./cmd/sql-migrator

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/sql-migrator .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/migrator.yaml .
CMD ["./sql-migrator"]
```

### GitHub Actions Example
```yaml
name: Run Migrations
on:
  push:
    branches: [main]

jobs:
  migrate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      
      - name: Build migrator
        run: go build -o sql-migrator ./cmd/sql-migrator
      
      - name: Run migrations
        run: ./sql-migrator ./config/production.yaml
        env:
          DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
```

## Troubleshooting

### Enable Debug Logging

Set log level in your application configuration or environment:

```bash
LOG_LEVEL=debug ./sql-migrator
```

### Checking Migration Status

Use the version action to check current state:

```yaml
Migrator:
  action: version
```

This will output:
```
Database version: 5, dirty: false
```

### Manual Intervention

In case of critical issues, you can directly access the migration tracking table:

```sql
-- PostgreSQL
SELECT * FROM schema_migrations;

-- MySQL  
SELECT * FROM schema_migrations;
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.

## Related Documentation

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [GoUtils Database Connections](../../databases/connections/)
- [GoUtils Configuration](../../configutils/)
- [GoUtils Logger](../../logger/)
