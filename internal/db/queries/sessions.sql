-- name: CreateSession :one
INSERT INTO sessions (
    id,
    user_id,
    token,
    refresh_token,
    expires_at
) VALUES (
    ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE token = ? LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE token = ?;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at < unixepoch();

-- name: UpdateSessionExpiry :exec
UPDATE sessions
SET expires_at = ?,
    refresh_token = ?
WHERE token = ?;