-- +goose Up
-- +goose StatementBegin

-- Page types enum
CREATE TYPE page_type AS ENUM ('generic', 'project', 'blog');

-- Page status enum
CREATE TYPE page_status AS ENUM ('draft', 'published', 'archived');

-- Main pages table
CREATE TABLE pages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Basic info
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    page_type page_type NOT NULL DEFAULT 'generic',
    
    -- Status
    status page_status DEFAULT 'draft',
    
    -- Media
    featured_image_id UUID REFERENCES media(id) ON DELETE SET NULL,
    
    -- Metadata
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    
    -- Constraints
    CONSTRAINT slug_not_empty CHECK (trim(slug) <> ''),
    CONSTRAINT title_not_empty CHECK (trim(title) <> '')
);

-- Unique constraint only for non-deleted pages
CREATE UNIQUE INDEX idx_pages_unique_slug 
ON pages(slug) 
WHERE deleted_at IS NULL;

-- Other indexes for performance
CREATE INDEX idx_pages_slug ON pages(slug);
CREATE INDEX idx_pages_status ON pages(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_pages_type ON pages(page_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_pages_author ON pages(author_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_pages_published ON pages(published_at) WHERE published_at IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_pages_deleted ON pages(deleted_at) WHERE deleted_at IS NOT NULL;

-- Trigger for updated_at
CREATE TRIGGER update_pages_updated_at
    BEFORE UPDATE ON pages
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_pages_updated_at ON pages;
DROP TABLE IF EXISTS pages CASCADE;
DROP TYPE IF EXISTS page_status;
DROP TYPE IF EXISTS page_type;

-- +goose StatementEnd
