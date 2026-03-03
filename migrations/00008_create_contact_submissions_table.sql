-- +goose Up
-- +goose StatementBegin
CREATE TABLE contact_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(100) NOT NULL,
    last_name  VARCHAR(100) NOT NULL,
    email      VARCHAR(256) NOT NULL,
    subject    VARCHAR(100) NOT NULL,
    message    TEXT NOT NULL,
    -- email delivery tracking
    email_sent      BOOLEAN NOT NULL DEFAULT FALSE,
    email_attempts  INT NOT NULL DEFAULT 0,
    email_last_attempted_at TIMESTAMPTZ,
    email_error     TEXT,
    -- general
    read_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ);


CREATE INDEX idx_contact_submissions_created_at ON contact_submissions(created_at DESC);
CREATE INDEX idx_contact_submissions_unsent ON contact_submissions(email_sent, email_attempts)
    WHERE deleted_at IS NULL AND email_sent = FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS contact_submissions;
-- +goose StatementEnd
