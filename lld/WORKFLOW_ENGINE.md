# Workflow Engine — Low-Level Design (layout v1)

## Purpose

Define the **core interfaces** first: a **Workflow** with lifecycle hooks and a **Task** catalog with CRUD. Execution (running each task) comes later; this document is the skeleton.

## Interfaces

### Task

```
Task
  ID() string
  Name() string
```

### TaskRepository

```
Create(spec TaskSpec) (Task, error)
Edit(id, spec TaskSpec) (Task, error)
Delete(id) error
Get(id) (Task, error)
```

### Workflow

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
engine/
  status.go       WorkflowStatus
  task.go         Task, TaskSpec, TaskRepository
  task_memory.go  MemoryTaskStore
  workflow.go     Workflow, BasicWorkflow, NewWorkflow
```

## Reference implementation

- `MemoryTaskStore` — in-memory tasks
- `BasicWorkflow` — holds ordered task IDs + status; uses `TaskRepository` to validate `Add`

## Not in this layout

- Task execution, context, retries, persistence, API/HTTP layer
