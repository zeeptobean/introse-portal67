-- +goose Up
ALTER TABLE enrollments ADD COLUMN version BIGINT NOT NULL DEFAULT 1 CHECK (version > 0);
UPDATE enrollments SET version = 2 WHERE status = 'cancelled';

CREATE TABLE outbox (
    seq BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    event_id UUID NOT NULL UNIQUE,
    stream TEXT NOT NULL,
    type TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at TIMESTAMPTZ
);
CREATE INDEX outbox_unpublished_idx ON outbox (seq) WHERE published_at IS NULL;

-- +goose Down
-- Downgrade the application first and keep this migration applied if events
-- still need to be published or replayed.
DROP TABLE outbox;
ALTER TABLE enrollments DROP COLUMN version;
