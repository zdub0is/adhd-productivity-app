-- +goose Up
CREATE TABLE journal_entries (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title      TEXT NOT NULL DEFAULT '',
    body       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX journal_entries_created_at_idx ON journal_entries (created_at);

-- +goose Down
DROP TABLE journal_entries;
