-- +goose Up
-- +goose StatementBegin
-- Create the function if it doesn't exist
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';


CREATE TABLE seo (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Polymorphic reference: 'page', 'project', 'blog'
    entity_type TEXT NOT NULL CHECK (entity_type IN ('page', 'project', 'blog')),
    entity_id   UUID NOT NULL,

    -- Basic SEO
    seo_title       TEXT,
    seo_description TEXT,

    -- Open Graph
    og_title       TEXT,
    og_description TEXT,
    og_image_id    UUID REFERENCES media(id) ON DELETE SET NULL,

    -- Advanced
    canonical_url TEXT,
    robots_index  BOOLEAN NOT NULL DEFAULT true,
    robots_follow BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- One SEO record per entity
    CONSTRAINT seo_entity_unique UNIQUE (entity_type, entity_id)
);

CREATE INDEX idx_seo_entity ON seo (entity_type, entity_id);

-- Auto-update updated_at
CREATE TRIGGER update_seo_updated_at
    BEFORE UPDATE ON seo
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_seo_updated_at ON seo;
DROP TABLE IF EXISTS seo;

-- +goose StatementEnd
