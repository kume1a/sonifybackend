-- name: GetUsers :exec
SELECT * FROM "users";

-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, name, email)
VALUES ($1, $2, $3, $4, $5)
RETURNING *; 

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET name = COALESCE(@name, name)
WHERE id = @id
RETURNING *;
