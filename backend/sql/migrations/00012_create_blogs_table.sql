-- +goose Up
-- +goose StatementBegin

-- Blogs table (first-class domain entity)
CREATE TABLE blogs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Basic info
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    
    -- Status
    status page_status DEFAULT 'draft', -- Reuse page_status enum
    
    -- Blog-specific fields
    excerpt TEXT,
    reading_time INTEGER,
    is_featured BOOLEAN DEFAULT false,
    
    -- Media
    featured_image_id UUID REFERENCES media(id) ON DELETE SET NULL,
    
    -- Metadata
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    
    -- Constraints
    CONSTRAINT blog_slug_not_empty CHECK (trim(slug) <> ''),
    CONSTRAINT blog_title_not_empty CHECK (trim(title) <> ''),
    CONSTRAINT reading_time_valid CHECK (
        reading_time IS NULL OR reading_time > 0
    )
);

-- Unique constraint only for non-deleted blogs
CREATE UNIQUE INDEX idx_blogs_unique_slug 
ON blogs(slug) 
WHERE deleted_at IS NULL;

-- Other indexes for performance
CREATE INDEX idx_blogs_slug ON blogs(slug);
CREATE INDEX idx_blogs_status ON blogs(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_blogs_author ON blogs(author_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_blogs_published ON blogs(published_at) WHERE published_at IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_blogs_deleted ON blogs(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX idx_blogs_featured ON blogs(is_featured) WHERE is_featured = true;

-- Trigger for updated_at
CREATE TRIGGER update_blogs_updated_at
    BEFORE UPDATE ON blogs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- Soft Delete Cleanup Triggers
-- ============================================

-- Function to cleanup blocks when blog is soft-deleted
CREATE OR REPLACE FUNCTION cleanup_blog_blocks()
RETURNS TRIGGER AS $$
BEGIN
    -- When blog is soft-deleted, delete all associated blocks
    IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL THEN
        DELETE FROM blocks 
        WHERE entity_type = 'blog' AND entity_id = NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to cleanup blocks on blog soft-delete
CREATE TRIGGER trigger_cleanup_blog_blocks
    AFTER UPDATE OF deleted_at ON blogs
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_blog_blocks();

-- Function to cleanup media_relations when blog is soft-deleted
CREATE OR REPLACE FUNCTION cleanup_blog_media_relations()
RETURNS TRIGGER AS $$
BEGIN
    -- When blog is soft-deleted, unlink all media relations
    IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL THEN
        DELETE FROM media_relations 
        WHERE entity_type = 'blog' AND entity_id = NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to cleanup media_relations on blog soft-delete
CREATE TRIGGER trigger_cleanup_blog_media_relations
    AFTER UPDATE OF deleted_at ON blogs
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_blog_media_relations();

-- Function to cleanup SEO when blog is soft-deleted
CREATE OR REPLACE FUNCTION cleanup_blog_seo()
RETURNS TRIGGER AS $$
BEGIN
    -- When blog is soft-deleted, delete SEO record
    IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL THEN
        DELETE FROM seo 
        WHERE entity_type = 'blog' AND entity_id = NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to cleanup SEO on blog soft-delete
CREATE TRIGGER trigger_cleanup_blog_seo
    AFTER UPDATE OF deleted_at ON blogs
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_blog_seo();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS trigger_cleanup_blog_seo ON blogs;
DROP TRIGGER IF EXISTS trigger_cleanup_blog_media_relations ON blogs;
DROP TRIGGER IF EXISTS trigger_cleanup_blog_blocks ON blogs;
DROP FUNCTION IF EXISTS cleanup_blog_seo();
DROP FUNCTION IF EXISTS cleanup_blog_media_relations();
DROP FUNCTION IF EXISTS cleanup_blog_blocks();
DROP TRIGGER IF EXISTS update_blogs_updated_at ON blogs;
DROP TABLE IF EXISTS blogs CASCADE;

-- +goose StatementEnd
