

-- +goose Up
CREATE TABLE IF NOT EXISTS media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url TEXT NOT NULL,
    thumbnail_url TEXT,
    type TEXT NOT NULL, -- e.g. "image", "video", "embed"
    is_public BOOLEAN DEFAULT FALSE,
    title TEXT,
    alt_text TEXT,
    mime_type TEXT,
    file_size INT,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS media;
