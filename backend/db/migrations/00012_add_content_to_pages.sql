-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS pages
ADD content TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS pages
DROP COLUMN content;
-- +goose StatementEnd
