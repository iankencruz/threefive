-- +goose Up
-- +goose StatementBegin

-- Media table with S3 support
CREATE TABLE media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- File metadata
    filename TEXT NOT NULL,              -- Generated filename: YYYYMMDD-<uuid>.ext (e.g., "20240115-abc123.jpg")
    original_filename TEXT NOT NULL,     -- User's original filename
    mime_type TEXT NOT NULL,             -- e.g., "image/jpeg", "video/mp4", "image/png"
    file_size BIGINT NOT NULL,           -- File size in bytes
    
    -- Image/Video dimensions
    width INTEGER,                       -- Original width in pixels
    height INTEGER,                      -- Original height in pixels
    duration INTEGER,                    -- For videos (in seconds)
    
    -- Storage configuration
    storage_type TEXT NOT NULL DEFAULT 's3' CHECK (storage_type IN ('local', 's3')),
    
    -- S3 keys (pattern: media/<size>/<year>/<filename>)
    -- Example: media/original/2024/20240115-abc123.jpg
    s3_bucket TEXT,                      -- Bucket name (e.g., "threefive-media")
    s3_region TEXT,                      -- Region/endpoint (e.g., "us-east-1" or Vultr endpoint)
    original_key TEXT,                   -- media/original/2024/20240115-abc123.jpg
    large_key TEXT,                      -- media/large/2024/20240115-abc123.jpg (1920px)
    medium_key TEXT,                     -- media/medium/2024/20240115-abc123.jpg (1024px)
    thumbnail_key TEXT,                  -- media/thumbnail/2024/20240115-abc123.jpg (300px)
    
    -- Accessibility & SEO
    alt_text TEXT,                       -- Alt text for images/videos
    
    -- Metadata
    uploaded_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,              -- Soft delete (media only)
    
    CONSTRAINT filename_not_empty CHECK (trim(filename) <> '')
);

-- Indexes for media table
CREATE INDEX idx_media_uploaded_by ON media(uploaded_by) WHERE deleted_at IS NULL;
CREATE INDEX idx_media_created_at ON media(created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_media_mime_type ON media(mime_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_media_storage_type ON media(storage_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_media_deleted_at ON media(deleted_at);
CREATE INDEX idx_media_filename ON media(filename) WHERE deleted_at IS NULL;

-- Media Relations table (polymorphic many-to-many)
-- Links media to projects, blogs, pages
CREATE TABLE media_relations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    media_id UUID NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    entity_type TEXT NOT NULL,           -- 'project', 'blog', 'page'
    entity_id UUID NOT NULL,             -- ID of project/blog/page
    relation_type TEXT NOT NULL DEFAULT 'gallery' CHECK (relation_type IN ('gallery', 'featured', 'content')),
    sort_order INTEGER DEFAULT 0,        -- For ordering gallery images (0-based)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Ensure unique relationship per media-entity-type combination
    UNIQUE(media_id, entity_type, entity_id, relation_type)
);

-- Indexes for media_relations
CREATE INDEX idx_media_relations_media ON media_relations(media_id);
CREATE INDEX idx_media_relations_entity ON media_relations(entity_type, entity_id);
CREATE INDEX idx_media_relations_entity_type ON media_relations(entity_type, entity_id, relation_type);
CREATE INDEX idx_media_relations_sort ON media_relations(entity_type, entity_id, sort_order);

-- Comments for documentation
COMMENT ON TABLE media IS 'Stores all media files (images, videos) with S3 keys for different sizes';
COMMENT ON COLUMN media.filename IS 'Generated filename format: YYYYMMDD-<uuid>.ext';
COMMENT ON COLUMN media.original_key IS 'S3 key pattern: media/original/<year>/<filename>';
COMMENT ON COLUMN media.large_key IS 'S3 key for large version (1920px width)';
COMMENT ON COLUMN media.medium_key IS 'S3 key for medium version (1024px width)';
COMMENT ON COLUMN media.thumbnail_key IS 'S3 key for thumbnail (300px width) or video thumbnail from FFmpeg';
COMMENT ON COLUMN media.duration IS 'Video duration in seconds (NULL for images)';

COMMENT ON TABLE media_relations IS 'Polymorphic relationship between media and entities (projects, blogs, pages)';
COMMENT ON COLUMN media_relations.relation_type IS 'Type: gallery (gallery images), featured (hero/featured image), content (inline content)';
COMMENT ON COLUMN media_relations.sort_order IS 'Display order for gallery images (0-based index)';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS media_relations CASCADE;
DROP TABLE IF EXISTS media CASCADE;

-- +goose StatementEnd
