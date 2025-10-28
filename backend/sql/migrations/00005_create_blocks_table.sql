-- +goose Up
-- +goose StatementBegin

-- Base blocks table (polymorphic)
CREATE TABLE blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Constraint for valid block types
    CONSTRAINT valid_block_type CHECK (type IN ('hero', 'richtext', 'header', 'gallery'))
);







-- Indexes for performance
CREATE INDEX idx_blocks_page_id ON blocks(page_id);
CREATE INDEX idx_blocks_sort ON blocks(page_id, sort_order);
CREATE INDEX idx_blocks_type ON blocks(type);



-- Trigger for updated_at
CREATE TRIGGER update_blocks_updated_at
    BEFORE UPDATE ON blocks
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;
DROP TABLE IF EXISTS blocks CASCADE;

-- +goose StatementEnd
