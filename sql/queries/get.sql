-- name: ListUsers :many
SELECT * FROM USERS;

-- name: AddUser :one
INSERT INTO USERS (username, hash)
VALUES ($1, $2)
RETURNING *;

-- name: GetUser :one
SELECT * FROM USERS
WHERE username = $1;