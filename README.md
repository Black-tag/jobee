# jobee

Workflow orchestrator skeleton (Go).

## Tasks (PostgreSQL + sqlc)

- Domain: `task/` — CRUD + lifecycle on `Task` (`pending` → `running` → …)
- Store: `task/postgres` — sqlc queries in `internal/db/queries`, schema in `internal/db/migrations`
- Migrations: [goose](https://github.com/pressly/goose) via `internal/db/migrate.go`

### Connect to local Postgres (port 5432)

Create a database once (adjust user as needed):

```bash
createdb jobee
# or: psql -c 'CREATE DATABASE jobee;'
```

Default DSN (override with `DATABASE_URL`):

```text
postgres://postgres:postgres@127.0.0.1:5432/jobee?sslmode=disable
```

```bash
export DATABASE_URL='postgres://YOUR_USER@127.0.0.1:5432/jobee?sslmode=disable'
go run .
```

After changing SQL:

```bash
sqlc generate
```

### Task model (`tasks` table)

| Column       | Type        |
|-------------|-------------|
| id          | UUID (PK)   |
| name        | TEXT        |
| status      | `task_status` enum |
| created_at  | TIMESTAMPTZ |
| updated_at  | TIMESTAMPTZ |

Repository: `Create`, `Edit`, `Run`, `Delete`, `Get`, `List`.
