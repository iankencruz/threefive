
-- +goose Up
CREATE TABLE IF NOT EXISTS richtext_block (
    block_id UUID PRIMARY KEY REFERENCES blocks(id) ON DELETE CASCADE  DEFAULT gen_random_uuid() ,
    html TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS richtext_block;

