package storage

import (
	"database/sql"
	"errors"

	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/models"
	"mitsimi.dev/codeShare/internal/storage"
)

func (s *SQLiteStorage) ToggleSaveSnippet(userID storage.UserID, snippetID storage.SnippetID, isSave bool) error {
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

	if isSave {
		err := s.q.SaveSnippet(s.ctx, db.SaveSnippetParams{
			SnippetID: snippetID,
			UserID:    userID,
		})
		if err != nil {
			return err
		}
	} else {
		err := s.q.DeleteSavedSnippet(s.ctx, db.DeleteSavedSnippetParams{
			SnippetID: snippetID,
			UserID:    userID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLiteStorage) GetSavedSnippets(userID storage.UserID) ([]models.Snippet, error) {
	snippets, err := s.q.GetSavedSnippets(s.ctx, userID)
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
			Author:    snippet.AuthorUsername.String,
			CreatedAt: snippet.CreatedAt,
			UpdatedAt: snippet.UpdatedAt,
			Likes:     int(snippet.Likes),
			IsLiked:   snippet.IsLiked == 1,
			IsSaved:   snippet.IsSaved == 1,
		}
	}

	return result, nil
}
