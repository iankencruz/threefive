-- +goose Up
-- +goose StatementBegin


CREATE TABLE IF NOT EXISTS blogs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    cover_image_id UUID REFERENCES media(id) ON DELETE SET NULL,
    seo_description TEXT DEFAULT '',
    seo_title TEXT DEFAULT '',
    canonical_url TEXT DEFAULT '',
    is_draft BOOLEAN DEFAULT TRUE,
    is_published BOOLEAN DEFAULT FALSE,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS blogs;
-- +goose StatementEnd
