-- +goose Up

ALTER TABLE users ADD COLUMN auth_version BIGINT UNSIGNED NOT NULL DEFAULT 1;

-- +goose Down

ALTER TABLE users DROP COLUMN auth_version;
