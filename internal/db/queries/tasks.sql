-- name: CreateTask :one
INSERT INTO tasks (name, status)
VALUES ($1, 'pending')
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1;

-- name: UpdateTaskName :one
UPDATE tasks
SET name = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateTaskStatus :one
UPDATE tasks
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;

-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY created_at, id;
