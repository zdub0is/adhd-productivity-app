-- +goose Up
CREATE TABLE checkins (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mood       SMALLINT NOT NULL,
    energy     SMALLINT NOT NULL,
    note       TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT checkins_mood_range CHECK (mood BETWEEN 1 AND 5),
    CONSTRAINT checkins_energy_range CHECK (energy BETWEEN 1 AND 5)
);

CREATE INDEX checkins_created_at_idx ON checkins (created_at);

-- +goose Down
DROP TABLE checkins;
