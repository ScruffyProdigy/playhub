# Database Migrations

This document describes the database migration system for PlayHub.

## Overview

PlayHub uses a custom migration system built on top of [golang-migrate](https://github.com/golang-migrate/migrate) to manage database schema changes. The system supports both programmatic and CLI-based migration management.

## Migration Files

Migrations are stored in the `backend/migrations/` directory and follow the naming convention:
- `{version}_{description}.up.sql` - Migration to apply changes
- `{version}_{description}.down.sql` - Migration to rollback changes

### Current Migrations

- `000001_initial_schema.up.sql` - Creates the initial database schema including:
  - Users table with magic link authentication support
  - Magic links table for secure email-based login
  - Games table for available games
  - Game queues table for matchmaking
  - Game sessions table for active games
  - Digital goods table for trading system
  - User inventory table for owned items

## CLI Usage

The migration CLI tool is located at `backend/cmd/migrate/main.go` and can be used with the following commands:

### Using Make Commands

```bash
# Run all pending migrations
make migrate-up

# Rollback the last migration
make migrate-down

# Check current migration version
make migrate-version

# Force a specific migration version (use with caution)
make migrate-force VERSION=1
```

### Direct CLI Usage

```bash
# Run migrations up
go run ./cmd/migrate -action=up

# Rollback last migration
go run ./cmd/migrate -action=down

# Run specific number of steps
go run ./cmd/migrate -action=steps -steps=2

# Check current version
go run ./cmd/migrate -action=version

# Force version (use with caution)
go run ./cmd/migrate -action=force -version=1
```

## Programmatic Usage

The migration system can also be used programmatically:

```go
import "github.com/scruffyprodigy/playhub/internal/migrate"

// Create migrator
migrator, err := migrate.NewMigrator(db)
if err != nil {
    return err
}
defer migrator.Close()

// Run all pending migrations
if err := migrator.Up(); err != nil {
    return err
}
```

## Database Schema

### Users Table
- `id` - UUID primary key
- `email` - Unique email address
- `username` - Unique username
- `display_name` - User's display name
- `avatar_url` - Optional avatar image URL
- `is_active` - Whether the user account is active
- `is_verified` - Whether the user has verified their email
- `last_login_at` - Timestamp of last login
- `created_at` - Account creation timestamp
- `updated_at` - Last update timestamp

### Magic Links Table
- `id` - UUID primary key
- `user_id` - Foreign key to users table (nullable for new users)
- `email` - Email address for the magic link
- `token` - Unique token for the magic link
- `expires_at` - Expiration timestamp
- `used_at` - When the link was used (nullable)
- `created_at` - Creation timestamp

### Games Table
- `id` - UUID primary key
- `name` - Game name
- `description` - Game description
- `version` - Game version
- `min_players` - Minimum players required
- `max_players` - Maximum players allowed
- `estimated_duration_minutes` - Estimated game duration
- `category` - Game category
- `status` - Game status (active, inactive, maintenance)
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Game Queues Table
- `id` - UUID primary key
- `game_id` - Foreign key to games table
- `user_id` - Foreign key to users table
- `status` - Queue status (waiting, matched, cancelled, expired)
- `priority` - Queue priority
- `preferences` - JSON preferences
- `joined_at` - When user joined queue
- `matched_at` - When user was matched (nullable)
- `expires_at` - Queue expiration (nullable)

### Game Sessions Table
- `id` - UUID primary key
- `game_id` - Foreign key to games table
- `status` - Session status (active, completed, cancelled)
- `started_at` - Session start time
- `ended_at` - Session end time (nullable)
- `session_data` - JSON session data

### Game Session Participants Table
- `id` - UUID primary key
- `session_id` - Foreign key to game_sessions table
- `user_id` - Foreign key to users table
- `joined_at` - When user joined session
- `left_at` - When user left session (nullable)
- `role` - User's role in the session

### Digital Goods Table
- `id` - UUID primary key
- `name` - Item name
- `description` - Item description
- `category` - Item category
- `rarity` - Item rarity (common, uncommon, rare, epic, legendary)
- `game_id` - Foreign key to games table
- `is_tradeable` - Whether item can be traded
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### User Inventory Table
- `id` - UUID primary key
- `user_id` - Foreign key to users table
- `good_id` - Foreign key to digital_goods table
- `quantity` - Number of items owned
- `acquired_at` - When item was acquired

## Environment Variables

- `DATABASE_URL` - PostgreSQL connection string (required)

## Kubernetes Integration

The migration system is integrated with Kubernetes through the `k8s/jobs/migration.yaml` job, which runs migrations automatically during deployment.

## Testing

Migration tests are located in `backend/internal/migrate/migrate_test.go` and can be run with:

```bash
go test ./internal/migrate -v
```

Note: Tests require a `DATABASE_URL` environment variable to be set.

## Best Practices

1. **Always create both up and down migrations** - Every migration should be reversible
2. **Test migrations thoroughly** - Test both up and down migrations in development
3. **Use transactions** - Wrap complex migrations in transactions when possible
4. **Backup before major changes** - Always backup production data before running migrations
5. **Version control** - Keep migration files in version control
6. **Sequential numbering** - Use sequential version numbers for migrations
7. **Descriptive names** - Use clear, descriptive names for migration files

## Troubleshooting

### Migration Stuck in "Dirty" State
If a migration fails and leaves the database in a dirty state:

```bash
# Check current version and dirty status
make migrate-version

# Force the version to the last successful migration
make migrate-force VERSION=1
```

### Connection Issues
Ensure the `DATABASE_URL` environment variable is set correctly:

```bash
export DATABASE_URL="postgres://user:password@host:port/database?sslmode=disable"
```

### Permission Issues
Ensure the database user has sufficient permissions to create tables, indexes, and manage the schema_migrations table.
