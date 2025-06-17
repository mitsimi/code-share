package storage

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	db "mitsimi.dev/codeShare/internal/db/sqlc"
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

func (s *SQLiteStorage) GetSnippet(userID storage.UserID, id string) (models.Snippet, error) {
	snippet, err := s.q.GetSnippet(s.ctx, db.GetSnippetParams{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Snippet{}, errors.New("snippet not found")
		}
		return models.Snippet{}, err
	}

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

func (s *SQLiteStorage) DeleteSnippet(id storage.SnippetID) error {
	return s.q.DeleteSnippet(s.ctx, id)
}

func (s *SQLiteStorage) ToggleLikeSnippet(userID storage.UserID, id storage.SnippetID, isLike bool) error {
	// Check if snippet exists first
	_, err := s.q.GetSnippet(s.ctx, db.GetSnippetParams{
		UserID: userID,
		ID:     id,
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
		exists, err := s.checkLikeExists(qtx, userID, id)
		if err != nil {
			return err
		}
		// If not already liked, add the like
		if !exists {
			if err := qtx.LikeSnippet(s.ctx, db.LikeSnippetParams{
				SnippetID: id,
				UserID:    userID,
			}); err != nil {
				return err
			}
			if err := qtx.IncrementLikesCount(s.ctx, id); err != nil {
				return err
			}
		}
	} else {
		// Check if already unliked
		exists, err := s.checkLikeExists(qtx, userID, id)
		if err != nil {
			return err
		}
		// If liked, remove the like
		if exists {
			if err := qtx.DeleteLike(s.ctx, db.DeleteLikeParams{
				SnippetID: id,
				UserID:    userID,
			}); err != nil {
				return err
			}
			if err := qtx.DecrementLikesCount(s.ctx, id); err != nil {
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
