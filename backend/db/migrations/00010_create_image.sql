
-- +goose Up
CREATE TABLE IF NOT EXISTS image_block (
    block_id UUID PRIMARY KEY REFERENCES blocks(id) ON DELETE CASCADE  DEFAULT gen_random_uuid() ,
    media_id UUID  REFERENCES media(id) ON DELETE CASCADE DEFAULT gen_random_uuid(),
    alt_text TEXT DEFAULT '',
    align TEXT NOT NULL DEFAULT 'center',
    object_fit TEXT NOT NULL DEFAULT 'cover'
);

-- +goose Down
DROP TABLE IF EXISTS image_block;

