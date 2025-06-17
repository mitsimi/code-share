package storage

import (
	"database/sql"
	"errors"

	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/models"
	"mitsimi.dev/codeShare/internal/storage"
)

func (s *SQLiteStorage) ToggleLikeSnippet(userID storage.UserID, snippetID storage.SnippetID, isLike bool) error {
	// Check if snippet exists first
	_, err := s.q.GetSnippet(s.ctx, db.GetSnippetParams{
		UserID:    userID,
		SnippetID: snippetID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("snippet not found")
		}
		return err
	}

	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.q.WithTx(tx)

	// Like or unlike the snippet
	if isLike {
		exists, err := s.checkLikeExists(qtx, userID, snippetID)
		if err != nil {
			return err
		}
		// If not already liked, add the like
		if !exists {
			if err := qtx.LikeSnippet(s.ctx, db.LikeSnippetParams{
				SnippetID: snippetID,
				UserID:    userID,
			}); err != nil {
				return err
			}
			if err := qtx.IncrementLikesCount(s.ctx, snippetID); err != nil {
				return err
			}
		}
	} else {
		// Check if already unliked
		exists, err := s.checkLikeExists(qtx, userID, snippetID)
		if err != nil {
			return err
		}
		// If liked, remove the like
		if exists {
			if err := qtx.DeleteLike(s.ctx, db.DeleteLikeParams{
				SnippetID: snippetID,
				UserID:    userID,
			}); err != nil {
				return err
			}
			if err := qtx.DecrementLikesCount(s.ctx, snippetID); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (s *SQLiteStorage) GetLikedSnippets(userID storage.UserID) ([]models.Snippet, error) {
	snippets, err := s.q.GetLikedSnippets(s.ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	result := make([]models.Snippet, len(snippets))
	for i, snippet := range snippets {
		result[i] = models.Snippet{
			ID:        snippet.ID,
			Title:     snippet.Title,
			Content:   snippet.Content,
			Language:  snippet.Language,
			Author:    snippet.AuthorUsername,
			CreatedAt: snippet.CreatedAt,
			UpdatedAt: snippet.UpdatedAt,
			Likes:     int(snippet.Likes),
			IsLiked:   snippet.IsLiked == 1,
			IsSaved:   snippet.IsSaved == 1,
		}
	}

	return result, nil
}

// checkLikeExists checks if a user has liked a snippet
func (s *SQLiteStorage) checkLikeExists(qtx *db.Queries, userID storage.UserID, snippetID storage.SnippetID) (bool, error) {
	exists, err := qtx.CheckLikeExists(s.ctx, db.CheckLikeExistsParams{
		SnippetID: snippetID,
		UserID:    userID,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	return exists == 1, nil
}
