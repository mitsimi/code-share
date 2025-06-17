package storage

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/logger"
	"mitsimi.dev/codeShare/internal/models"
	"mitsimi.dev/codeShare/internal/storage"
)

func (s *SQLiteStorage) GetSnippets(userID storage.UserID) ([]models.Snippet, error) {
	snippets, err := s.q.GetSnippets(s.ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]models.Snippet, len(snippets))
	for i, snippet := range snippets {
		result[i] = models.Snippet{
			ID:        snippet.ID,
			Title:     snippet.Title,
			Content:   snippet.Content,
			Language:  snippet.Language,
			Author:    snippet.AuthorUsername.String, // Handle sql.NullString
			CreatedAt: snippet.CreatedAt,
			UpdatedAt: snippet.UpdatedAt,
			Likes:     int(snippet.Likes),
			IsLiked:   snippet.IsLiked == 1,
		}
	}

	return result, nil
}

func (s *SQLiteStorage) GetSnippet(userID storage.UserID, snippetID string) (models.Snippet, error) {
	snippet, err := s.q.GetSnippet(s.ctx, db.GetSnippetParams{
		UserID:    userID,
		SnippetID: snippetID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Snippet{}, errors.New("snippet not found")
		}
		return models.Snippet{}, err
	}

	isLiked, err := s.checkLikeExists(s.q, userID, snippetID)
	if err != nil {
		return models.Snippet{}, err
	}

	logger.Debug("THE FUCK IS THIS????!??!?!?!?!?",
		zap.String("id", snippet.ID),
		zap.String("title", snippet.Title),
		zap.Bool("is_liked", isLiked),
		zap.Bool("DB isLiked", snippet.IsLiked == 1),
	)

	return models.Snippet{
		ID:        snippet.ID,
		Title:     snippet.Title,
		Content:   snippet.Content,
		Language:  snippet.Language,
		Author:    snippet.AuthorUsername.String, // Handle sql.NullString
		CreatedAt: snippet.CreatedAt,
		UpdatedAt: snippet.UpdatedAt,
		Likes:     int(snippet.Likes),
		IsLiked:   snippet.IsLiked == 1,
	}, nil
}

func (s *SQLiteStorage) CreateSnippet(snippet models.Snippet) (storage.SnippetID, error) {
	// Get the user by ID to verify it exists
	_, err := s.q.GetUser(s.ctx, snippet.Author)
	if err != nil {
		return "", errors.New("author not found")
	}

	result, err := s.q.CreateSnippet(s.ctx, db.CreateSnippetParams{
		ID:       uuid.NewString(),
		Title:    snippet.Title,
		Content:  snippet.Content,
		Language: snippet.Language,
		Author:   snippet.Author, // Use the user ID directly
	})
	if err != nil {
		return "", err
	}

	return result.ID, nil
}

func (s *SQLiteStorage) UpdateSnippet(snippet models.Snippet) error {
	_, err := s.q.UpdateSnippet(s.ctx, db.UpdateSnippetParams{
		Title:    snippet.Title,
		Content:  snippet.Content,
		Language: snippet.Language,
		ID:       snippet.ID,
	})
	return err
}

func (s *SQLiteStorage) DeleteSnippet(snippetID storage.SnippetID) error {
	return s.q.DeleteSnippet(s.ctx, snippetID)
}

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
