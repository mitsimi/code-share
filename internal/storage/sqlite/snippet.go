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
			IsSaved:   snippet.IsSaved == 1,
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
		IsSaved:   snippet.IsSaved == 1,
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

func (s *SQLiteStorage) GetSnippetsByAuthor(userID, authorID storage.UserID) ([]models.Snippet, error) {

	snippets, err := s.q.GetSnippetsByAuthor(s.ctx, db.GetSnippetsByAuthorParams{
		UserID:   userID,
		AuthorID: authorID,
	})
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
			IsSaved:   snippet.IsSaved == 1,
		}
	}

	return result, nil
}
