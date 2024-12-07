-- name: CreateFeedFollow :one
WITH new_feed_follow AS (
INSERT INTO feed_follows (id, created_at, updated_at, feed_id, user_id)
VALUES ($1, $2, $3, $4, $5)
    RETURNING *
    )
SELECT nff.id, nff.created_at, nff.updated_at, nff.feed_id, nff.user_id,
       f.name AS feed_name, u.name AS user_name
FROM new_feed_follow nff
         JOIN feeds f ON nff.feed_id = f.id
         JOIN users u ON nff.user_id = u.id;

-- name: GetFeedFollowsForUser :many
SELECT ff.id, ff.created_at, ff.updated_at, ff.feed_id, ff.user_id,
       f.name AS feed_name, u.name AS user_name
FROM feed_follows ff
         JOIN feeds f ON ff.feed_id = f.id
         JOIN users u ON ff.user_id = u.id
WHERE ff.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE feed_id = $1
  AND user_id = $2;