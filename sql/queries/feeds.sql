-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT f.id, f.created_at, f.updated_at, f.name, f.url, f.user_id,
       u.name AS user_name
FROM feeds f
JOIN users u ON f.user_id = u.id;