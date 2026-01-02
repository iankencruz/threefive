-- +goose Up
-- +goose StatementBegin
-- Hero block data
CREATE TABLE IF NOT EXISTS block_hero (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    block_id UUID UNIQUE NOT NULL REFERENCES blocks(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    subtitle TEXT,
    image_id UUID REFERENCES media(id) ON DELETE SET NULL,
    cta_text TEXT,
    cta_url TEXT,
    navigation_bar BOOLEAN DEFAULT FALSE,
    
    CONSTRAINT title_not_empty CHECK (trim(title) <> '')
);

-- Indexes for performance
CREATE INDEX idx_block_hero_block_id ON block_hero(block_id);
CREATE INDEX idx_block_hero_image ON block_hero(image_id) WHERE image_id IS NOT NULL;

-- Trigger function to clean up media_relations when hero block is deleted
CREATE OR REPLACE FUNCTION cleanup_hero_block_media_relations()
RETURNS TRIGGER AS $$
BEGIN
    -- Delete all media_relations entries for this hero block
    DELETE FROM media_relations
    WHERE entity_type = 'block_hero' AND entity_id = OLD.id;
    
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Trigger to cleanup media relations on hero block deletion
CREATE TRIGGER trigger_cleanup_hero_block_media_relations
    BEFORE DELETE ON block_hero
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_hero_block_media_relations();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_cleanup_hero_block_media_relations ON block_hero;
DROP FUNCTION IF EXISTS cleanup_hero_block_media_relations();
DROP TABLE IF EXISTS block_hero CASCADE;
-- +goose StatementEnd
