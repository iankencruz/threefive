-- +goose Up
-- +goose StatementBegin

-- SEO data (polymorphic - works with pages, projects, blogs)
CREATE TABLE seo (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Polymorphic reference
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    
    -- Basic SEO
    meta_title TEXT,
    meta_description TEXT,
    
    -- Open Graph
    og_title TEXT,
    og_description TEXT,
    og_image_id UUID REFERENCES media(id) ON DELETE SET NULL,
    
    -- Advanced SEO
    canonical_url TEXT,
    robots_index BOOLEAN DEFAULT true,
    robots_follow BOOLEAN DEFAULT true,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT valid_seo_entity_type CHECK (entity_type IN ('page','project', 'blog'))
);

-- Unique constraint: one SEO record per entity
CREATE UNIQUE INDEX idx_seo_entity_unique ON seo(entity_type, entity_id);

-- Indexes
CREATE INDEX idx_seo_entity ON seo(entity_type, entity_id);

-- Trigger for updated_at
CREATE TRIGGER update_seo_updated_at
    BEFORE UPDATE ON seo
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_seo_updated_at ON seo;
DROP TABLE IF EXISTS seo CASCADE;

-- +goose StatementEnd
