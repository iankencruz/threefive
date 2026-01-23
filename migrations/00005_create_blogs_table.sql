-- +goose Up
-- +goose StatementBegin

-- Blogs table
CREATE TABLE blogs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Basic info
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    excerpt TEXT,
    
    -- Status
    status TEXT DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    
    -- Blog-specific fields
    reading_time INTEGER, -- estimated reading time in minutes
    is_featured BOOLEAN DEFAULT false,
    category_id UUID, -- optional category reference (if you want categories later)
    
    -- Media
    featured_image_id UUID REFERENCES media(id) ON DELETE SET NULL,
    
    -- Metadata
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    
    -- Constraints
    CONSTRAINT blog_slug_not_empty CHECK (trim(slug) <> ''),
    CONSTRAINT blog_title_not_empty CHECK (trim(title) <> '')
);

-- Unique constraint only for non-deleted blogs
CREATE UNIQUE INDEX idx_blogs_unique_slug ON blogs(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_blogs_slug ON blogs(slug);
CREATE INDEX idx_blogs_status ON blogs(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_blogs_author ON blogs(author_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_blogs_published_at ON blogs(published_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_blogs_deleted_at ON blogs(deleted_at);
CREATE INDEX idx_blogs_is_featured ON blogs(is_featured) WHERE deleted_at IS NULL AND is_featured = true;

-- Junction table: blogs <-> tags (many-to-many)
CREATE TABLE blog_tags (
    blog_id UUID NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (blog_id, tag_id)
);

CREATE INDEX idx_blog_tags_blog ON blog_tags(blog_id);
CREATE INDEX idx_blog_tags_tag ON blog_tags(tag_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS blog_tags;
DROP TABLE IF EXISTS blogs;

-- +goose StatementEnd
