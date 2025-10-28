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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS block_gallery CASCADE;
-- +goose StatementEnd
