-- +goose Up
-- +goose StatementBegin
-- Header block data
CREATE TABLE block_header (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    block_id UUID UNIQUE NOT NULL REFERENCES blocks(id) ON DELETE CASCADE,
    heading TEXT NOT NULL,
    subheading TEXT,
    level VARCHAR(10) DEFAULT 'h2',
    
    CONSTRAINT valid_level CHECK (level IN ('h1', 'h2', 'h3', 'h4', 'h5', 'h6')),
    CONSTRAINT heading_not_empty CHECK (trim(heading) <> '')
);

-- Indexes for performance
CREATE INDEX idx_block_header_block_id ON block_header(block_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS block_header CASCADE;
-- +goose StatementEnd
