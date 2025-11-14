-- +goose Up
-- +goose StatementBegin

-- Project status enum
CREATE TYPE project_status AS ENUM ('completed', 'ongoing', 'archived');

-- Projects table (portfolio/gallery showcase)
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Basic info
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT, -- Short description/summary
    project_date DATE, -- Date of project (year, or specific date)
    
    -- Status
    status page_status DEFAULT 'draft', -- Reuse page_status enum
    
    -- Project-specific fields
    client_name TEXT,
    project_year INTEGER,
    project_url TEXT,
    technologies JSONB DEFAULT '[]'::jsonb,
    project_status project_status DEFAULT 'completed',
    
    -- Media
    featured_image_id UUID REFERENCES media(id) ON DELETE SET NULL,
    
    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    
    -- Constraints
    CONSTRAINT project_slug_not_empty CHECK (trim(slug) <> ''),
    CONSTRAINT project_title_not_empty CHECK (trim(title) <> ''),
    CONSTRAINT project_year_valid CHECK (
        project_year IS NULL OR 
        (project_year >= 1900 AND project_year <= 2100)
    )
);

-- Unique constraint only for non-deleted projects
CREATE UNIQUE INDEX idx_projects_unique_slug 
ON projects(slug) 
WHERE deleted_at IS NULL;

-- Other indexes for performance
CREATE INDEX idx_projects_slug ON projects(slug);
CREATE INDEX idx_projects_status ON projects(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_published ON projects(published_at) WHERE published_at IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_projects_deleted ON projects(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX idx_projects_date ON projects(project_date) WHERE project_date IS NOT NULL;
CREATE INDEX idx_projects_year ON projects(project_year) WHERE project_year IS NOT NULL;

-- Trigger for updated_at
CREATE TRIGGER update_projects_updated_at
    BEFORE UPDATE ON projects
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- Soft Delete Cleanup Triggers
-- ============================================

-- Function to cleanup media_relations when project is soft-deleted
CREATE OR REPLACE FUNCTION cleanup_project_media_relations()
RETURNS TRIGGER AS $$
BEGIN
    -- When project is soft-deleted, unlink all media relations
    IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL THEN
        DELETE FROM media_relations 
        WHERE entity_type = 'project' AND entity_id = NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to cleanup media_relations on project soft-delete
CREATE TRIGGER trigger_cleanup_project_media_relations
    AFTER UPDATE OF deleted_at ON projects
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_project_media_relations();

-- Function to cleanup SEO when project is soft-deleted
CREATE OR REPLACE FUNCTION cleanup_project_seo()
RETURNS TRIGGER AS $$
BEGIN
    -- When project is soft-deleted, delete SEO record
    IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL THEN
        DELETE FROM seo 
        WHERE entity_type = 'project' AND entity_id = NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to cleanup SEO on project soft-delete
CREATE TRIGGER trigger_cleanup_project_seo
    AFTER UPDATE OF deleted_at ON projects
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_project_seo();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS trigger_cleanup_project_seo ON projects;
DROP TRIGGER IF EXISTS trigger_cleanup_project_media_relations ON projects;
DROP FUNCTION IF EXISTS cleanup_project_seo();
DROP FUNCTION IF EXISTS cleanup_project_media_relations();
DROP TRIGGER IF EXISTS update_projects_updated_at ON projects;
DROP TABLE IF EXISTS projects CASCADE;
DROP TYPE IF EXISTS project_status;

-- +goose StatementEnd
