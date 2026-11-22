-- +goose Up
CREATE TABLE inbox (
    event_id UUID PRIMARY KEY,
    stream TEXT NOT NULL,
    type TEXT NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE inbox;
