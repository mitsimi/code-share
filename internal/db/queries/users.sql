-- name: CreateUser :one
INSERT INTO users (
    id,
    username,
    email,
    password_hash
) VALUES (
    ?, ?, ?, ?
) 
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = ?;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ?;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;