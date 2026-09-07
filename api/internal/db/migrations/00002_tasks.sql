-- +goose Up
CREATE TABLE tasks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT 'todo',
    priority    TEXT NOT NULL DEFAULT 'normal',
    due_date    DATE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT tasks_status_check CHECK (status IN ('todo', 'in_progress', 'done', 'archived')),
    CONSTRAINT tasks_priority_check CHECK (priority IN ('low', 'normal', 'high'))
);

CREATE INDEX tasks_status_idx ON tasks (status);

-- +goose Down
DROP TABLE tasks;
