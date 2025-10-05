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
    CONSTRAINT valid_block_type CHECK (type IN ('hero', 'richtext', 'header'))
);

-- Hero block data
CREATE TABLE block_hero (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    block_id UUID UNIQUE NOT NULL REFERENCES blocks(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    subtitle TEXT,
    image_id UUID REFERENCES media(id) ON DELETE SET NULL,
    cta_text TEXT,
    cta_url TEXT,
    
    CONSTRAINT title_not_empty CHECK (trim(title) <> '')
);

-- Richtext block data
CREATE TABLE block_richtext (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    block_id UUID UNIQUE NOT NULL REFERENCES blocks(id) ON DELETE CASCADE,
    content TEXT NOT NULL
);

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
CREATE INDEX idx_blocks_page_id ON blocks(page_id);
CREATE INDEX idx_blocks_sort ON blocks(page_id, sort_order);
CREATE INDEX idx_blocks_type ON blocks(type);

CREATE INDEX idx_block_hero_block_id ON block_hero(block_id);
CREATE INDEX idx_block_hero_image ON block_hero(image_id) WHERE image_id IS NOT NULL;

CREATE INDEX idx_block_richtext_block_id ON block_richtext(block_id);

CREATE INDEX idx_block_header_block_id ON block_header(block_id);

-- Trigger for updated_at
CREATE TRIGGER update_blocks_updated_at
    BEFORE UPDATE ON blocks
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;
DROP TABLE IF EXISTS block_header CASCADE;
DROP TABLE IF EXISTS block_richtext CASCADE;
DROP TABLE IF EXISTS block_hero CASCADE;
DROP TABLE IF EXISTS blocks CASCADE;

-- +goose StatementEnd
