-- name: GetUsers :exec
SELECT * FROM "users";

-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, name)
VALUES ($1, $2, $3, $4)
RETURNING *; 

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;