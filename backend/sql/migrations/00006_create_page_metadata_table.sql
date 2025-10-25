-- +goose Up
-- +goose StatementBegin

-- SEO data (one-to-one with pages)
CREATE TABLE page_seo (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id UUID UNIQUE NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    
    -- Basic SEO
    meta_title TEXT,
    meta_description TEXT,
    
    -- Open Graph
    og_title TEXT,
    og_description TEXT,
    og_image_id UUID REFERENCES media(id) ON DELETE SET NULL,
    
    -- Advanced SEO
    canonical_url TEXT,
    robots_index BOOLEAN DEFAULT true,
    robots_follow BOOLEAN DEFAULT true,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Project-specific data (for project pages)
CREATE TABLE page_project_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id UUID UNIQUE NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    
    client_name TEXT,
    project_year INTEGER,
    project_url TEXT,
    technologies JSONB DEFAULT '[]'::jsonb,
    project_status VARCHAR(50) DEFAULT 'completed',
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT project_year_valid CHECK (
        project_year IS NULL OR 
        (project_year >= 1900 AND project_year <= 2100)
    ),
    CONSTRAINT project_status_valid CHECK (
        project_status IN ('completed', 'ongoing', 'archived')
    )
);

-- Blog-specific data (for blog pages)
CREATE TABLE page_blog_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id UUID UNIQUE NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    
    excerpt TEXT,
    reading_time INTEGER,
    is_featured BOOLEAN DEFAULT false,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT reading_time_valid CHECK (
        reading_time IS NULL OR reading_time > 0
    )
);

-- Indexes
CREATE INDEX idx_page_seo_page_id ON page_seo(page_id);
CREATE INDEX idx_page_project_data_page_id ON page_project_data(page_id);
CREATE INDEX idx_page_project_data_year ON page_project_data(project_year) WHERE project_year IS NOT NULL;
CREATE INDEX idx_page_blog_data_page_id ON page_blog_data(page_id);
CREATE INDEX idx_page_blog_data_featured ON page_blog_data(is_featured) WHERE is_featured = true;

-- Triggers for updated_at
CREATE TRIGGER update_page_seo_updated_at
    BEFORE UPDATE ON page_seo
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_page_project_data_updated_at
    BEFORE UPDATE ON page_project_data
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_page_blog_data_updated_at
    BEFORE UPDATE ON page_blog_data
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_page_seo_updated_at ON page_seo;
DROP TRIGGER IF EXISTS update_page_project_data_updated_at ON page_project_data;
DROP TRIGGER IF EXISTS update_page_blog_data_updated_at ON page_blog_data;

DROP TABLE IF EXISTS page_blog_data CASCADE;
DROP TABLE IF EXISTS page_project_data CASCADE;
DROP TABLE IF EXISTS page_seo CASCADE;

-- +goose StatementEnd
