# Workflow Engine — Low-Level Design (layout v1)

## Purpose

Define the **core interfaces** first: a **Workflow** with lifecycle hooks and a **Task** catalog with CRUD. Execution (running each task) comes later; tasks are persisted in PostgreSQL (phase 1).

## Interfaces

### Task (`task/` package)

```
Task
  ID() string
  Name() string
  Status() Status

  Run() error
  Pause() error
  End() error
  Fail() error
```

#### Task status

| Status | Meaning |
|--------|---------|
| `pending` | Created; not started yet |
| `running` | In progress |
| `paused` | Started but paused; `Run()` resumes |
| `ended` | Finished successfully (terminal) |
| `error` | Failed (terminal) |

### Task repository

```
Create(spec Spec) (Task, error)
Edit(id, spec Spec) (Task, error)
Delete(id) error
Get(id) (Task, error)
List() ([]Task, error)
```

### Workflow (`engine/` package)

```
Workflow
  ID() string
  Status() WorkflowStatus
  TaskIDs() []string

  Start() error
  Add(taskID string) error
  Cancel() error
  Pause() error
```

## Workflow status

| Status | Meaning |
|--------|---------|
| `draft` | Created; tasks can be added |
| `running` | Started (executor will hook here later) |
| `paused` | Started but paused; can `Start()` again to resume |
| `cancelled` | Terminal; no further changes |

## Layout rules (v1)

- **Add** — only when status is `draft` or `paused`; task must exist in repository.
- **Start** — from `draft` or `paused` → `running`.
- **Pause** — only from `running` → `paused`.
- **Cancel** — from `draft`, `running`, or `paused` → `cancelled`.

## Package layout

```
task/
  status.go, task.go, repository.go, lifecycle.go
  memory/          in-memory Repository (tests, no DATABASE_URL)
  postgres/        sqlc-backed Repository
engine/
  status.go        WorkflowStatus
  workflow.go      Workflow, BasicWorkflow
internal/db/
  migrations/      goose SQL (001_tasks — single table)
  queries/         sqlc input
  sqlc/            generated queries
  postgres.go      pgx pool
  migrate.go       goose Up
```

## Persistence (phase 1)

Single `tasks` table: `id` (UUID), `name`, `status`, timestamps. Queries via **sqlc**; connections via **pgx** pool.

- `DATABASE_URL` set → migrate + `task/postgres` store
- unset → `task/memory` store

Local Postgres: `docker compose up -d` then  
`DATABASE_URL=postgres://jobee:jobee@localhost:5432/jobee?sslmode=disable`

Regenerate sqlc after query changes: `sqlc generate`

## Phase 2 (planned)

Split into `task_definitions` + `task_executions` and workflow tables. Migrate existing `tasks` rows; orphan/GC policy applies to definitions when workflows are persisted.

## Not in this layout

- Workflow persistence, workers, retries, API/HTTP layer
