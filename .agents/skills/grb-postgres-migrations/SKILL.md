---
name: grb-postgres-migrations
description: Use when changing the database schema, the seed, a gorm store or a migration in the Go React boilerplate. Covers embedded golang-migrate files, gorm boundaries, UUID keys, seed idempotency and the relevant commands.
---

# PostgreSQL, gorm and migrations

PostgreSQL is the only database. The adapter is
`internal/adapter/postgres`. gorm rows never leave that package.

## Schema changes

1. Add a new pair of files in
   `internal/adapter/postgres/migrations/`:
   - `000002_short_name.up.sql`
   - `000002_short_name.down.sql`
2. Never edit an applied migration.
3. Use snake_case table and column names.
4. Use `uuid` primary keys generated in Go with `entity.NewUUID()`.
5. Add indexes for foreign keys and common filters.
6. Use `timestamptz` and defaults.
7. Update the gorm row struct in `models.go`.
8. Update the store methods.
9. Run `make db-migrate` or `make db-reset`.

## Rules

- gorm stays inside `internal/adapter/postgres`, `internal/xmodule` and `cmd`.
- Store rows convert to entities at the boundary.
- Use `port.ErrUniqueViolation` for duplicate-key failures.
- Use transactions for multi-row writes.
- Use `context.Context` on every query.
- Never scan directly into an HTTP DTO.

## Seed

- `internal/adapter/postgres/seed.go` is idempotent.
- Permissions and roles are upserted.
- Demo users are inserted only when absent.
- Demo password is `password`.
- The seeded state is 14 users and 3 roles.
- Run `make db-seed` after a role or permission change.

## Commands

```bash
make db-up        # Postgres and Mailpit
make db-migrate   # apply migrations
make db-seed      # permissions, roles, demo users
make db-reset     # destructive rebuild
```

The server also applies migrations at boot. `cmd/migrate` exists for release
pipelines.

## Testing a schema change

1. Run `make db-reset`.
2. Sign in as `superadmin@example.com`.
3. Exercise the affected screen.
4. Run `make test` and the relevant Playwright spec.
5. Confirm the seed still produces the documented counts.

## Useful checks

```sql
SELECT count(*) FROM users;
SELECT count(*) FROM roles;
SELECT key, value FROM settings;
```
