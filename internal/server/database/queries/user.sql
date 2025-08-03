-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES($1, $2)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUsernameById :one
SELECT username FROM users
WHERE id = $1;

-- name: GetUsers :many
SELECT username, id FROM users;
