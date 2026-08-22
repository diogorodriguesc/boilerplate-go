-- name: GetUserByEmail :one
SELECT id, username, email FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id, username, email;