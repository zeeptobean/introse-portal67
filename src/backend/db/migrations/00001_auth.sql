-- +goose Up
CREATE TABLE sessions (
    token TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

CREATE TABLE users(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email TEXT NOT NULL UNIQUE CHECK (email = lower(btrim(email)) AND email <> ''),
    password_hash TEXT NOT NULL CHECK (password_hash <> ''),
    institutional_id TEXT UNIQUE,
    full_name TEXT NOT NULL CHECK (btrim(full_name) <> ''),
    date_of_birth DATE NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('student', 'lecturer', 'staff')),
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now()
);
-- +goose Down
DROP TABLE sessions;
DROP TABLE users;
