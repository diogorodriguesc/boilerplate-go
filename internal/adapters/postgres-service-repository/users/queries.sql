-- name: SearchUsers :many
SELECT id, username, email FROM users
WHERE 
    (sqlc.narg('email')::text IS NULL OR email = sqlc.narg('email')) AND
    (sqlc.narg('username')::text IS NULL OR username = sqlc.narg('username'));

-- name: CreateUser :one
INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id, username, email;

-- name: GetUser :one
SELECT id, username, email FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT id, username, email FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: DeleteUser :execrows
DELETE FROM users WHERE id = $1;