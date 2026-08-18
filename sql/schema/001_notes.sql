-- +goose Up

CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL
);

-- +goose Down
DROP TABLE notes;