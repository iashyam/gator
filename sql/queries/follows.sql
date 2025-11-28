-- name: CreateFeedFollow :one
with something as (INSERT INTO feed_follows (id, created_At, updated_At, feed_id, user_id)
VALUES($1, $2, $3, $4, $5)
RETURNING *) 

SELECT
    something.*,
    feeds.name as feed_name, 
    users.name as user_name
FROM something
INNER JOIN feeds ON feeds.id = something.feed_id
INNER JOIN users ON users.id = something.user_id;

-- name: GetFeedFollowsForUser :many
SELECT * FROM feed_follows WHERE user_id=$1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE user_id=$1 AND feed_id=$2;

