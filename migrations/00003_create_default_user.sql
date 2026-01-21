-- +goose Up
-- +goose StatementBegin
INSERT INTO users (id, first_name, last_name, email, password_hash, created_at, updated_at)
VALUES (
  gen_random_uuid(),
  'admin',
  'user',
  'admin@example.com',
  '$2a$10$TYELnOfNjpWcGJAsjUYTM.w8l4Hr8MLufk0lvL/r2oW1LMDzVQdT6',
  NOW(),
  NOW()
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE email = 'admin@example.com';
-- +goose StatementEnd
