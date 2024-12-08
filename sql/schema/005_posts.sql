-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title VARCHAR(255) NULL,
    url VARCHAR(255) UNIQUE NOT NULL,
    description TEXT NULL,
    published_at TIMESTAMP NULL,
    feed_id UUID REFERENCES feeds (id) ON DELETE CASCADE NOT NULL
);

-- +goose Down
DROP TABLE posts;