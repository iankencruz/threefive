-- +goose Up
-- +goose StatementBegin

-- Page status enum (keep this)
CREATE TYPE page_status AS ENUM ('draft', 'published', 'archived');

-- Main pages table (simplified - only for generic pages)
CREATE TABLE pages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Basic info
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    
    -- Status
    status page_status DEFAULT 'draft',
    
    -- Media
    featured_image_id UUID REFERENCES media(id) ON DELETE SET NULL,
    
    -- Metadata
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
CREATE INDEX idx_pages_published ON pages(published_at) WHERE published_at IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_pages_deleted ON pages(deleted_at) WHERE deleted_at IS NOT NULL;

-- Trigger for updated_at
CREATE TRIGGER update_pages_updated_at
    BEFORE UPDATE ON pages
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- Soft Delete Cleanup Triggers
-- ============================================

-- Function to cleanup blocks when page is soft-deleted
CREATE OR REPLACE FUNCTION cleanup_page_blocks()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL THEN
        DELETE FROM blocks 
        WHERE entity_type = 'page' AND entity_id = NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_cleanup_page_blocks
    AFTER UPDATE OF deleted_at ON pages
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_page_blocks();

-- Function to cleanup media_relations when page is soft-deleted
CREATE OR REPLACE FUNCTION cleanup_page_media_relations()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL THEN
        DELETE FROM media_relations 
        WHERE entity_type = 'page' AND entity_id = NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_cleanup_page_media_relations
    AFTER UPDATE OF deleted_at ON pages
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_page_media_relations();

-- Function to cleanup SEO when page is soft-deleted
CREATE OR REPLACE FUNCTION cleanup_page_seo()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL THEN
        DELETE FROM seo 
        WHERE entity_type = 'page' AND entity_id = NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_cleanup_page_seo
    AFTER UPDATE OF deleted_at ON pages
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_page_seo();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS trigger_cleanup_page_seo ON pages;
DROP TRIGGER IF EXISTS trigger_cleanup_page_media_relations ON pages;
DROP TRIGGER IF EXISTS trigger_cleanup_page_blocks ON pages;
DROP FUNCTION IF EXISTS cleanup_page_seo();
DROP FUNCTION IF EXISTS cleanup_page_media_relations();
DROP FUNCTION IF EXISTS cleanup_page_blocks();
DROP TRIGGER IF EXISTS update_pages_updated_at ON pages;
DROP TABLE IF EXISTS pages CASCADE;
DROP TYPE IF EXISTS page_status;

-- +goose StatementEnd
