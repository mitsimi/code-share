-- name: SaveSnippet :exec
INSERT OR IGNORE INTO user_saves (snippet_id, user_id)
VALUES (@snippet_id, @user_id);

-- name: DeleteSavedSnippet :exec
DELETE FROM user_saves
WHERE snippet_id = @snippet_id AND user_id = @user_id;

-- name: GetSavedSnippets :many
SELECT s.*, 
    CASE WHEN ul.user_id IS NOT NULL THEN 1 ELSE 0 END as is_liked,
    CASE WHEN us.user_id IS NOT NULL THEN 1 ELSE 0 END as is_saved,
    u.id AS author_id, 
    u.username AS author_username, 
    u.email AS author_email,
    u.avatar AS author_avatar
FROM snippets s
LEFT JOIN user_likes ul ON s.id = ul.snippet_id AND ul.user_id = @user_id
LEFT JOIN user_saves us ON s.id = us.snippet_id AND us.user_id = @user_id
LEFT JOIN users u ON s.author = u.id
WHERE us.user_id = @user_id
ORDER BY ul.created_at DESC;