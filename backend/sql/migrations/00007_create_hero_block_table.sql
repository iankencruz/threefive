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
    
    CONSTRAINT title_not_empty CHECK (trim(title) <> '')
);

-- Indexes for performance
CREATE INDEX idx_block_hero_block_id ON block_hero(block_id);
CREATE INDEX idx_block_hero_image ON block_hero(image_id) WHERE image_id IS NOT NULL;




-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS block_hero CASCADE;
-- +goose StatementEnd
