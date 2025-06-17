package sqlite

import (
	"context"
	"database/sql"
	"errors"

	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/domain"
	"mitsimi.dev/codeShare/internal/repository"
)

var _ repository.BookmarkRepository = (*BookmarkRepository)(nil)

type BookmarkRepository struct {
	db *sql.DB
	q  *db.Queries
}

func NewBookmarkRepository(dbConn *sql.DB) *BookmarkRepository {
	return &BookmarkRepository{
		db: dbConn,
		q:  db.New(dbConn),
	}
}

func (r *BookmarkRepository) ToggleSave(ctx context.Context, userID, snippetID string, isSave bool) error {
	_, err := r.q.GetSnippet(ctx, db.GetSnippetParams{
		SnippetID: snippetID})
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.ErrNotFound
		}
		return repository.WrapError(err, "failed to get snippet")
	}

	if isSave {
		err := r.q.SaveSnippet(ctx, db.SaveSnippetParams{
			SnippetID: snippetID,
			UserID:    userID,
		})
		if err != nil {
			return err
		}
	} else {
		err := r.q.DeleteSavedSnippet(ctx, db.DeleteSavedSnippetParams{
			SnippetID: snippetID,
			UserID:    userID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *BookmarkRepository) GetSavedSnippets(ctx context.Context, userID string) ([]*domain.Snippet, error) {
	snippets, err := r.q.GetLikedSnippets(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, repository.WrapError(err, "failed to get liked snippets")
	}

	result := make([]*domain.Snippet, len(snippets))
	for i, snippet := range snippets {
		result[i] = &domain.Snippet{
			ID:        snippet.ID,
			Title:     snippet.Title,
			Content:   snippet.Content,
			Language:  snippet.Language,
			AuthorID:  snippet.Author,
			CreatedAt: snippet.CreatedAt,
			UpdatedAt: snippet.UpdatedAt}
	}

	return result, nil
}
