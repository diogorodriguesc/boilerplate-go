-- name: GetUserByEmail :one
SELECT id, username, email FROM users WHERE email = $1;