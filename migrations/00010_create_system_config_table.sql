-- +goose Up
CREATE TABLE IF NOT EXISTS system_config (
  config_code TEXT PRIMARY KEY NOT NULL,
  value TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);



-- +goose Down
DROP TABLE IF EXISTS system_config;
