-- +goose Up
create TABLE posts (id UUID primary key, created_at TIMESTAMP NOT NULL, updated_at TIMESTAMP NOT NULL, title TEXT not NULL, url TEXT unique not NULL, description TEXT not NULL, published_at TIMESTAMP not NULL, feed_id UUID NOT NULL references feeds (id) ON DELETE CASCADE);

-- +goose Down
DROP TABLE
posts;