-- +goose Up
-- +goose StatementBegin

-- Projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Basic info
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT,
    project_date DATE,
    
    -- Status (reusing existing page_status enum if it exists, or create it)
    status TEXT DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    
    -- Project-specific fields
    client_name TEXT,
    project_year INTEGER,
    project_url TEXT,
    project_status TEXT DEFAULT 'completed' CHECK (project_status IN ('completed', 'in-progress', 'planned')),

    -- Media
    featured_image_id UUID REFERENCES media(id) ON DELETE SET NULL,
    -- need to create media table

    -- Metadata
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    
    -- Constraints
    CONSTRAINT project_slug_not_empty CHECK (trim(slug) <> ''),
    CONSTRAINT project_title_not_empty CHECK (trim(title) <> '')
);

-- Unique constraint only for non-deleted projects
CREATE UNIQUE INDEX idx_projects_unique_slug ON projects(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_slug ON projects(slug);
CREATE INDEX idx_projects_status ON projects(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_author ON projects(author_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_published_at ON projects(published_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_deleted_at ON projects(deleted_at);

-- Junction table: projects <-> tags (many-to-many)
CREATE TABLE project_tags (
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (project_id, tag_id)
);

CREATE INDEX idx_project_tags_project ON project_tags(project_id);
CREATE INDEX idx_project_tags_tag ON project_tags(tag_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS project_tags;
DROP TABLE IF EXISTS projects;

-- +goose StatementEnd
