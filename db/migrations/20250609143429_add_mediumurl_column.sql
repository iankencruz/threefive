
-- +goose Up
ALTER TABLE media
ADD COLUMN medium_url TEXT;

-- +goose Down
ALTER TABLE media
DROP COLUMN medium_url;

