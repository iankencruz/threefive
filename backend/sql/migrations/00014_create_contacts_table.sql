-- +goose Up
-- +goose StatementBegin

CREATE TABLE contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL CHECK (trim(name) <> ''),
    email TEXT NOT NULL CHECK (trim(email) <> ''),
    subject TEXT,
    message TEXT NOT NULL CHECK (trim(message) <> ''),
    status TEXT NOT NULL DEFAULT 'new' CHECK (status IN ('new', 'read', 'archived')),
    ip_address INET,
    user_agent TEXT,
    email_sent BOOLEAN NOT NULL DEFAULT false,
    email_error TEXT,
    email_sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Indexes for better performance
CREATE INDEX idx_contacts_status ON contacts(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_contacts_created_at ON contacts(created_at);
CREATE INDEX idx_contacts_deleted_at ON contacts(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX idx_contacts_email ON contacts(email);
CREATE INDEX idx_contacts_email_sent ON contacts(email_sent) WHERE email_sent = false AND deleted_at IS NULL;

-- Trigger for automatic updated_at timestamp
CREATE TRIGGER update_contacts_updated_at
    BEFORE UPDATE ON contacts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_contacts_updated_at ON contacts;
DROP TABLE IF EXISTS contacts;

-- +goose StatementEnd
