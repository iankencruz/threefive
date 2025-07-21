-- +goose Up
CREATE TABLE IF NOT EXISTS heading_block (
    block_id UUID PRIMARY KEY REFERENCES blocks(id) ON DELETE CASCADE  DEFAULT gen_random_uuid() ,
    title TEXT NOT NULL,
    description TEXT
);

-- +goose Down
DROP TABLE IF EXISTS heading_block;
