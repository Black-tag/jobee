-- +goose Up
-- +goose StatementBegin
CREATE TYPE task_status AS ENUM (
    'pending', 'running', 'paused', 'ended', 'error'
);
-- +goose StatementEnd

CREATE TABLE tasks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL CHECK (name <> ''),
    status      task_status NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_tasks_status ON tasks (status);
CREATE INDEX idx_tasks_name ON tasks (name);
