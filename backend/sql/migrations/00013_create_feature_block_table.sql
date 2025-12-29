-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS block_feature (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  block_id UUID UNIQUE NOT NULL REFERENCES blocks(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  heading TEXT NOT NULL,
  subheading TEXT NOT NULL
);

-- Feature images are linked via media_relations with:
-- entity_type = 'block_feature' and entity_id = block_feature.id

CREATE INDEX idx_block_feature_block_id ON block_feature(block_id);

-- Trigger function to clean up media_relations when feature block is deleted
CREATE OR REPLACE FUNCTION cleanup_feature_block_media_relations()
RETURNS TRIGGER AS $$
BEGIN
    -- Delete all media_relations entries for this feature block
    DELETE FROM media_relations
    WHERE entity_type = 'block_feature' AND entity_id = OLD.id;
    
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Trigger to cleanup media relations on feature block deletion
CREATE TRIGGER trigger_cleanup_feature_block_media_relations
    BEFORE DELETE ON block_feature
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_feature_block_media_relations();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_cleanup_feature_block_media_relations ON block_feature;
DROP FUNCTION IF EXISTS cleanup_feature_block_media_relations();
DROP TABLE IF EXISTS block_feature CASCADE;
-- +goose StatementEnd
