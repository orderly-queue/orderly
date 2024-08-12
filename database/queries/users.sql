-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users(id, name, email, password)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = EXTRACT(epoch FROM NOW())
WHERE id = $1;
