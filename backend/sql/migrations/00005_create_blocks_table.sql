-- +goose Up
-- +goose StatementBegin

-- Base blocks table (polymorphic - works with pages, projects, blogs)
CREATE TABLE blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Polymorphic reference
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    
    -- Block data
    type VARCHAR(50) NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT valid_block_type CHECK (type IN ('hero', 'richtext', 'header', 'gallery')),
    CONSTRAINT valid_entity_type CHECK (entity_type IN ('page', 'project',  'blog'))
);

-- Indexes for performance
CREATE INDEX idx_blocks_entity ON blocks(entity_type, entity_id);
CREATE INDEX idx_blocks_entity_sort ON blocks(entity_type, entity_id, sort_order);
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
