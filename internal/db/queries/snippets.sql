-- name: GetSnippets :many
SELECT 
    s.*,
    CASE WHEN us.user_id IS NOT NULL THEN 1 ELSE 0 END as is_saved,
    CASE WHEN ul.user_id IS NOT NULL THEN 1 ELSE 0 END as is_liked,
    u.id AS author_id, 
    u.username AS author_username, 
    u.email AS author_email,
    u.avatar AS author_avatar
FROM snippets s
LEFT JOIN user_likes ul ON s.id = ul.snippet_id AND ul.user_id = @user_id
LEFT JOIN user_saves us ON s.id = us.snippet_id AND us.user_id = @user_id
LEFT JOIN users u ON s.author = u.id
ORDER BY s.created_at DESC;

-- name: GetSnippetsByAuthor :many
SELECT 
    s.*,
    CASE WHEN us.user_id IS NOT NULL THEN 1 ELSE 0 END as is_saved,
    CASE WHEN ul.user_id IS NOT NULL THEN 1 ELSE 0 END as is_liked,
    u.id AS author_id, 
    u.username AS author_username, 
    u.email AS author_email,
    u.avatar AS author_avatar
FROM snippets s
LEFT JOIN user_likes ul ON s.id = ul.snippet_id AND ul.user_id = @user_id
LEFT JOIN user_saves us ON s.id = us.snippet_id AND us.user_id = @user_id
LEFT JOIN users u ON s.author = u.id
WHERE s.author = @author_id
ORDER BY s.created_at DESC;

-- name: GetSnippet :one
SELECT 
    s.*,
    CASE WHEN us.user_id IS NOT NULL THEN 1 ELSE 0 END as is_saved,
    CASE WHEN ul.user_id IS NOT NULL THEN 1 ELSE 0 END as is_liked,
    u.id AS author_id, 
    u.username AS author_username, 
    u.email AS author_email,
    u.avatar AS author_avatar
FROM snippets s
LEFT JOIN user_saves us ON s.id = us.snippet_id AND us.user_id = @user_id
LEFT JOIN user_likes ul ON s.id = ul.snippet_id AND ul.user_id = @user_id
LEFT JOIN users u ON s.author = u.id
WHERE s.id = @snippet_id;

-- name: CreateSnippet :one
INSERT INTO snippets (
    id,
    title,
    content,
    language,
    author
) VALUES (
    ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateSnippet :one
UPDATE snippets
SET 
    title = @title,
    content = @content,
    language = @language,
    updated_at = CURRENT_TIMESTAMP
WHERE id = @snippet_id
RETURNING *;

-- name: DeleteSnippet :exec
DELETE FROM snippets
WHERE id = ?;

-- name: IncrementViews :exec
UPDATE snippets
SET views = views + 1
WHERE id = @snippet_id;

-- name: CheckRecentView :one
SELECT 
    snippet_id,
    viewer_identifier,
    last_viewed_at,
    (strftime('%s', 'now') - strftime('%s', last_viewed_at)) as seconds_since_last_view
FROM snippet_views 
WHERE snippet_id = @snippet_id 
AND viewer_identifier = @viewer_identifier
LIMIT 1;

-- name: RecordView :exec
INSERT OR REPLACE INTO snippet_views (
    snippet_id, 
    viewer_identifier, 
    ip_address, 
    last_viewed_at, 
    view_count
) VALUES (
    @snippet_id,
    @viewer_identifier,
    @ip_address,
    CURRENT_TIMESTAMP,
    CASE 
        WHEN EXISTS (
            SELECT 1 FROM snippet_views 
            WHERE snippet_id = @snippet_id 
            AND viewer_identifier = @viewer_identifier
        ) THEN (
            SELECT view_count + 1 FROM snippet_views 
            WHERE snippet_id = @snippet_id 
            AND viewer_identifier = @viewer_identifier
        )
        ELSE 1
    END
);

-- name: CleanupOldViews :exec
DELETE FROM snippet_views 
WHERE last_viewed_at < datetime('now', '-30 days');