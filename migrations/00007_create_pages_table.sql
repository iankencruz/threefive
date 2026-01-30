-- +goose Up
-- +goose StatementBegin

-- Pages table (home, about, contact - static pages)
CREATE TABLE pages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Basic info
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    
    -- Page type (home, about, contact - these 3 are hardcoded)
    page_type TEXT NOT NULL CHECK (page_type IN ('home', 'about', 'contact')),
    
    -- Hero section (common to all pages)
    hero_media_id UUID REFERENCES media(id) ON DELETE SET NULL,
    header TEXT,
    sub_header TEXT,
    
    -- About page specific
    content TEXT, -- Rich text content
    content_image_id UUID REFERENCES media(id) ON DELETE SET NULL,
    cta_text TEXT,
    cta_link TEXT,
    
    -- Contact page specific
    email TEXT,
    social_links JSONB, -- {"twitter": "url", "linkedin": "url", "github": "url", "instagram": "url"}
    
    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    -- Constraints
    CONSTRAINT page_slug_not_empty CHECK (trim(slug) <> ''),
    CONSTRAINT page_title_not_empty CHECK (trim(title) <> ''),
    CONSTRAINT page_type_required CHECK (trim(page_type) <> '')
);

-- Unique constraint: one page per type (non-deleted)
CREATE UNIQUE INDEX idx_pages_unique_type ON pages(page_type) WHERE deleted_at IS NULL;

-- Indexes
CREATE INDEX idx_pages_slug ON pages(slug);
CREATE INDEX idx_pages_type ON pages(page_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_pages_deleted_at ON pages(deleted_at);

-- Featured projects for About page (max 3)
CREATE TABLE page_featured_projects (
    page_id UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    display_order INTEGER NOT NULL CHECK (display_order BETWEEN 0 AND 2), -- Max 3 projects (0, 1, 2)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (page_id, project_id)
);

CREATE INDEX idx_page_featured_projects_page ON page_featured_projects(page_id);
CREATE INDEX idx_page_featured_projects_order ON page_featured_projects(page_id, display_order);

-- Seed the 3 static pages
INSERT INTO pages (title, slug, page_type, header, sub_header) VALUES
    ('Home', 'home', 'home', 'Welcome to My Portfolio', 'Building amazing digital experiences'),
    ('About', 'about', 'about', 'About Me', 'Learn more about my journey'),
    ('Contact', 'contact', 'contact', 'Get In Touch', 'Let''s work together');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS page_featured_projects;
DROP TABLE IF EXISTS pages;

-- +goose StatementEnd
