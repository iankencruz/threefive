
-- +goose Up
CREATE TABLE sessions (
    token TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS sessions;

