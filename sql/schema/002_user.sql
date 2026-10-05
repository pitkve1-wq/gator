-- +goose Up
CREATE TABLE feeds (id UUID Primary key, created_at TIMESTAMP NOT NULL, updated_at TIMESTAMP NOT NULL, name TEXT NOT NULL, url TEXT UNIQUE NOT NULL, user_id UUID NOT NULL references users (id) ON DELETE CASCADE);

-- +goose Down
DROP TABLE
feeds;