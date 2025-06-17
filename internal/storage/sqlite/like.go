package sqlite

import (
	"context"
	"database/sql"
	"errors"

	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/domain"
	"mitsimi.dev/codeShare/internal/repository"
	"mitsimi.dev/codeShare/internal/storage"
)

var _ repository.LikeRepository = (*LikeRepository)(nil)

type LikeRepository struct {
	db *sql.DB
	q  *db.Queries
}

func NewLikeRepository(dbConn *sql.DB) *LikeRepository {
	return &LikeRepository{
		db: dbConn,
		q:  db.New(dbConn),
	}
}

func (r *LikeRepository) ToggleLike(ctx context.Context, userID, snippetID string, isLike bool) error {
	_, err := r.q.GetSnippet(ctx, db.GetSnippetParams{
		SnippetID: snippetID})
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.ErrNotFound
		}
		return repository.WrapError(err, "failed to get snippet")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return repository.WrapError(err, "failed to begin transaction")
	}
	defer tx.Rollback()
	qtx := r.q.WithTx(tx)

	exists, err := r.checkLikeExists(ctx, qtx, userID, snippetID)
	if err != nil {
		return err
	}

	if isLike && !exists { // If not already liked, add the like
		if err := qtx.LikeSnippet(ctx, db.LikeSnippetParams{
			SnippetID: snippetID,
			UserID:    userID,
		}); err != nil {
			return err
		}
		if err := qtx.IncrementLikesCount(ctx, snippetID); err != nil {
			return err
		}

	} else if !isLike && exists { // If already liked, remove the like
		if err := qtx.DeleteLike(ctx, db.DeleteLikeParams{
			SnippetID: snippetID,
			UserID:    userID,
		}); err != nil {
			return err
		}
		if err := qtx.DecrementLikesCount(ctx, snippetID); err != nil {
			return err
		}

	}

	return tx.Commit()
}

func (r *LikeRepository) GetLikedSnippets(ctx context.Context, userID string) ([]*domain.Snippet, error) {
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
			UpdatedAt: snippet.UpdatedAt,
			Likes:     int(snippet.Likes),
			IsLiked:   snippet.IsLiked == 1,
			IsSaved:   snippet.IsSaved == 1,
		}
	}

	return result, nil
}

func (r *LikeRepository) checkLikeExists(ctx context.Context, qtx *db.Queries, userID storage.UserID, snippetID storage.SnippetID) (bool, error) {
	exists, err := qtx.CheckLikeExists(ctx, db.CheckLikeExistsParams{
		SnippetID: snippetID,
		UserID:    userID,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	return exists == 1, nil
}
