---- +goose Up
-- +goose StatementBegin

-- Tags table (shared between projects and blogs)
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT tag_name_not_empty CHECK (trim(name) <> ''),
    CONSTRAINT tag_slug_not_empty CHECK (trim(slug) <> '')
);

-- Unique constraint for tag slugs
CREATE UNIQUE INDEX idx_tags_unique_slug ON tags(slug);
CREATE INDEX idx_tags_name ON tags(name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tags;

-- +goose StatementEnd
