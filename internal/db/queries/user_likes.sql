-- name: LikeSnippet :exec
INSERT OR IGNORE INTO user_likes (snippet_id, user_id)
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