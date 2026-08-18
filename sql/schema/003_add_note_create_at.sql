-- +goose Up

ALTER TABLE notes ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT now();

-- +goose Down
ALTER TABLE notes DROP COLUMN created_at;