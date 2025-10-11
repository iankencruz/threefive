-- +goose Up
-- +goose StatementBegin
-- Richtext block data
CREATE TABLE block_richtext (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    block_id UUID UNIQUE NOT NULL REFERENCES blocks(id) ON DELETE CASCADE,
    content TEXT NOT NULL
);



-- Indexes for performance
CREATE INDEX idx_block_richtext_block_id ON block_richtext(block_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS block_richtext CASCADE;
-- +goose StatementEnd
