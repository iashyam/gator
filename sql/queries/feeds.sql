-- name: AddFeed :one
INSERT INTO feeds (id, created_At, updated_At, name, url, user_id)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListFeeds :many
SELECT * from feeds;

-- name: GetFeedByName :one
SELECT * from feeds where name=$1 limit 1;

-- name: GetFeedByID :one
SELECT * from feeds where id=$1 limit 1;

-- name: GetFeedByURL :one
SELECT * from feeds where url=$1 limit 1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_At = $2, 
    last_fetched_at = $3
WHERE id=$1;

-- name: GetLastFetched :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST LIMIT 1;