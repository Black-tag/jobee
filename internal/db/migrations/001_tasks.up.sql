-- Phase 1: single tasks table (catalog + runtime status).
-- Phase 2: split into task_definitions + task_executions; migrate rows from here.

CREATE TYPE task_status AS ENUM (
    'pending', 'running', 'paused', 'ended', 'error'
);

CREATE TABLE tasks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL CHECK (name <> ''),
    status      task_status NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_tasks_status ON tasks (status);
CREATE INDEX idx_tasks_name ON tasks (name);
