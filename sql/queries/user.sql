-- name: GetUsers :exec
SELECT * FROM "users";

-- name: CreateUser :one
INSERT INTO users(
  id, 
  created_at, 
  name, 
  email, 
  auth_provider, 
  password_hash
) VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *; 

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET name = COALESCE($1, name)
WHERE id = $2
RETURNING *;

-- name: CountUsersByEmail :one
SELECT COUNT(*) FROM users WHERE email = $1;