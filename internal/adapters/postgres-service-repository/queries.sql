-- name: SearchUsers :many
SELECT id, username, email FROM users
WHERE (sqlc.narg('email')::text IS NULL OR email = sqlc.narg('email'));

-- name: CreateUser :one
INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id, username, email;