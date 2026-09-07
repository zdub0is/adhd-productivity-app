-- +goose Up
CREATE TABLE daily_plans (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_date     DATE NOT NULL UNIQUE,
    intention     TEXT NOT NULL DEFAULT '',
    started_at    TIMESTAMPTZ,
    review_notes  TEXT NOT NULL DEFAULT '',
    reviewed_at   TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE daily_plan_tasks (
    plan_id  UUID NOT NULL REFERENCES daily_plans (id) ON DELETE CASCADE,
    task_id  UUID NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
    position INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (plan_id, task_id)
);

-- +goose Down
DROP TABLE daily_plan_tasks;
DROP TABLE daily_plans;
