
-- +goose Up
CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    user_agent TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Add the foreign key constraint
ALTER TABLE sessions
ADD CONSTRAINT fk_sessions_user_id
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE TABLE session_data (
    token TEXT NOT NULL REFERENCES sessions(token) ON DELETE CASCADE,
    key TEXT NOT NULL,
    value TEXT,
    PRIMARY KEY (token, key)
);

-- +goose Down
DROP TABLE IF EXISTS session_data;
DROP TABLE IF EXISTS sessions;

