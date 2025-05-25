-- name: CreateUser :one
INSERT INTO users (
    id,
    username
) VALUES (
    ?, ?
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = ?;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ?;

-- name: GetSnippets :many
SELECT 
    s.*,
    CASE WHEN ul.user_id IS NOT NULL THEN 1 ELSE 0 END as is_liked,
    u.username as author_username
FROM snippets s
LEFT JOIN user_likes ul ON s.id = ul.snippet_id AND ul.user_id = ?
LEFT JOIN users u ON s.author = u.id
ORDER BY s.created_at DESC;

-- name: GetSnippet :one
SELECT 
    s.*,
    CASE WHEN ul.user_id IS NOT NULL THEN 1 ELSE 0 END as is_liked,
    u.username as author_username
FROM snippets s
LEFT JOIN user_likes ul ON s.id = ul.snippet_id AND ul.user_id = ?
LEFT JOIN users u ON s.author = u.id
WHERE s.id = ?;

-- name: CreateSnippet :one
INSERT INTO snippets (
    id,
    title,
    content,
    author
) VALUES (
    ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateSnippet :one
UPDATE snippets
SET 
    title = ?,
    content = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteSnippet :exec
DELETE FROM snippets
WHERE id = ?;

-- name: LikeSnippet :exec
INSERT OR IGNORE INTO user_likes (snippet_id, user_id)
VALUES (?, ?);

-- name: UnlikeSnippet :exec
DELETE FROM user_likes
WHERE snippet_id = ? AND user_id = ?;

-- name: UpdateLikesCount :exec
UPDATE snippets
SET likes = (
    SELECT COUNT(*)
    FROM user_likes
    WHERE snippet_id = ?
)
WHERE id = ?; 