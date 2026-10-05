-- +goose Up
CREATE TABLE users (id UUID Primary key, created_at TIMESTAMP NOT NULL, updated_at TIMESTAMP NOT NULL, name TEXT UNIQUE NOT NULL);

-- +goose Down
DROP TABLE
users;