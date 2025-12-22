-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS block_about (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  block_id UUID UNIQUE NOT NULL REFERENCES blocks(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  heading TEXT NOT NULL,
  subheading TEXT NOT NULL
);

CREATE INDEX idx_block_about_block_id ON block_about(block_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS block_about CASCADE;
-- +goose StatementEnd
