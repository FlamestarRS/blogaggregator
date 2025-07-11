-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users 
WHERE name = $1;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: ListUsers :many
SELECT * FROM users;

-- name: GetUserByID :one
SELECT name FROM users
WHERE id = $1;