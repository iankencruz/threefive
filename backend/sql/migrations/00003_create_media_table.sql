-- +goose Up
-- +goose StatementBegin

-- Media storage types
CREATE TYPE storage_type AS ENUM ('local', 's3');

-- Media table
CREATE TABLE media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    filename VARCHAR(255) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    size_bytes BIGINT NOT NULL,
    width INTEGER,
    height INTEGER,
    storage_type storage_type NOT NULL DEFAULT 'local',
    storage_path TEXT NOT NULL,
    s3_bucket VARCHAR(255),
    s3_key TEXT,
    s3_region VARCHAR(50),
    url TEXT,
    thumbnail_url TEXT,
    uploaded_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_media_uploaded_by ON media(uploaded_by);
CREATE INDEX idx_media_created_at ON media(created_at);
CREATE INDEX idx_media_mime_type ON media(mime_type);
CREATE INDEX idx_media_storage_type ON media(storage_type);
CREATE INDEX idx_media_deleted_at ON media(deleted_at) WHERE deleted_at IS NULL;

-- Media relationships table (polymorphic)
CREATE TABLE media_relations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    media_id UUID NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL, -- 'project', 'page', 'gallery', etc
    entity_id UUID NOT NULL,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for relations
CREATE INDEX idx_media_relations_media_id ON media_relations(media_id);
CREATE INDEX idx_media_relations_entity ON media_relations(entity_type, entity_id);
CREATE UNIQUE INDEX idx_media_relations_unique ON media_relations(media_id, entity_type, entity_id);

-- Trigger for updated_at
CREATE TRIGGER update_media_updated_at
    BEFORE UPDATE ON media
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_media_updated_at ON media;
DROP TABLE IF EXISTS media_relations;
DROP TABLE IF EXISTS media;
DROP TYPE IF EXISTS storage_type;

-- +goose StatementEnd
