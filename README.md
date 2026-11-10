# jobee

Workflow orchestrator skeleton (Go).

## Tasks

- Domain: `task/` — CRUD + lifecycle (`pending` → `running` → …)
- Stores: `task/memory` (default), `task/postgres` when `DATABASE_URL` is set
- DB: PostgreSQL, migrations (goose), queries (sqlc)

```bash
docker compose up -d
export DATABASE_URL=postgres://jobee:jobee@localhost:5432/jobee?sslmode=disable
go run .
sqlc generate   # after editing internal/db/queries
```