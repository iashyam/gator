-- name: CreateUser :one
INSERT INTO users (id, created_At, updated_At, name)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
select * from users where name=$1 limit 1;

-- name: GetUserNameByID :one
select * from users where id=$1 limit 1;

-- name: Reset :exec
DELETE FROM users;

-- name: ListUsers :many
SELECT name FROM users;