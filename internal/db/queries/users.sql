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

-- name: UpdateUserAvatar :one
UPDATE users
SET 
    avatar = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateUserInfo :one
UPDATE users
SET 
    username = ?,
    email = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET 
    password_hash = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
