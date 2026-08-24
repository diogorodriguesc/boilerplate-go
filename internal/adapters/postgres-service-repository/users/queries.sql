-- name: SearchUsers :many
SELECT id, username, email FROM users
WHERE 
    (sqlc.narg('email')::text IS NULL OR email = sqlc.narg('email')) AND
    (sqlc.narg('username')::text IS NULL OR username = sqlc.narg('username'));

-- name: CreateUser :one
INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id, username, email;

-- name: GetUser :one
SELECT id, username, email FROM users WHERE id = $1;