-- +goose Up
-- +goose StatementBegin
-- Gallery block data
CREATE TABLE IF NOT EXISTS block_gallery (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    block_id UUID UNIQUE NOT NULL REFERENCES blocks(id) ON DELETE CASCADE,
    title TEXT
);

-- Gallery images are linked via media_relations with:
-- entity_type = 'block_gallery' and entity_id = block_gallery.id

CREATE INDEX idx_block_gallery_block_id ON block_gallery(block_id);

-- Trigger function to clean up media_relations when gallery block is deleted
CREATE OR REPLACE FUNCTION cleanup_gallery_block_media_relations()
RETURNS TRIGGER AS $$
BEGIN
    -- Delete all media_relations entries for this gallery block
    DELETE FROM media_relations
    WHERE entity_type = 'block_gallery' AND entity_id = OLD.id;
    
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Trigger to cleanup media relations on gallery block deletion
CREATE TRIGGER trigger_cleanup_gallery_block_media_relations
    BEFORE DELETE ON block_gallery
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_gallery_block_media_relations();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_cleanup_gallery_block_media_relations ON block_gallery;
DROP FUNCTION IF EXISTS cleanup_gallery_block_media_relations();
DROP TABLE IF EXISTS block_gallery CASCADE;
-- +goose StatementEnd
