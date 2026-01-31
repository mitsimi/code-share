-- name: LikeSnippet :exec
INSERT INTO user_likes (snippet_id, user_id)
VALUES (?, ?);

-- name: IncrementLikesCount :exec
UPDATE snippets 
SET likes = likes + 1 
WHERE id = ?;

-- name: DeleteLike :exec
DELETE FROM user_likes
WHERE snippet_id = ? AND user_id = ?;

-- name: DecrementLikesCount :exec
UPDATE snippets 
SET likes = likes - 1 
WHERE id = ?;

-- name: UpdateLikesCount :exec
UPDATE snippets
SET likes = (
    SELECT COUNT(*)
    FROM user_likes
    WHERE snippet_id = ?
)
WHERE id = ?;

-- name: CheckLikeExists :one
SELECT 1 as exists_flag
FROM user_likes 
WHERE snippet_id = ? AND user_id = ?
LIMIT 1;

-- name: GetLikedSnippets :many
SELECT s.*, 
    CASE WHEN us.user_id IS NOT NULL THEN 1 ELSE 0 END as is_saved,
    CASE WHEN ul.user_id IS NOT NULL THEN 1 ELSE 0 END as is_liked,
    u.id AS author_id, 
    u.username AS author_username, 
    u.email AS author_email,
    u.avatar AS author_avatar
FROM snippets s
LEFT JOIN user_likes ul ON s.id = ul.snippet_id
LEFT JOIN user_saves us ON s.id = us.snippet_id AND us.user_id = @user_id
LEFT JOIN users u ON s.author = u.id
WHERE ul.user_id = @user_id
ORDER BY ul.created_at DESC;