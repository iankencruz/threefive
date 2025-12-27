-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS block_feature (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  block_id UUID UNIQUE NOT NULL REFERENCES blocks(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  heading TEXT NOT NULL,
  subheading TEXT NOT NULL,
  image_id UUID REFERENCES media(id) ON DELETE SET NULL
);

CREATE INDEX idx_block_feature_block_id ON block_feature(block_id);
CREATE INDEX idx_block_feature_image ON block_feature(title) WHERE image_id IS NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS block_feature CASCADE;
-- +goose StatementEnd
