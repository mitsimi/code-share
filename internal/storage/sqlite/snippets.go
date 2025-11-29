package sqlite

import (
	"context"
	"database/sql"

	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/domain"
	"mitsimi.dev/codeShare/internal/repository"
)

var _ repository.SnippetRepository = (*SnippetRepository)(nil)

type SnippetRepository struct {
	db *sql.DB
	q  *db.Queries
}

func NewSnippetRepository(dbConn *sql.DB) *SnippetRepository {
	return &SnippetRepository{
		db: dbConn,
		q:  db.New(dbConn),
	}
}

func (r *SnippetRepository) Create(ctx context.Context, snippet *domain.Snippet) error {
	_, err := r.q.CreateSnippet(ctx, db.CreateSnippetParams{
		ID:       snippet.ID,
		Title:    snippet.Title,
		Content:  snippet.Content,
		Language: snippet.Language,
		Author:   snippet.Author.ID,
	})
	return err
}

func (r *SnippetRepository) GetByID(ctx context.Context, snippetID, userID string) (*domain.Snippet, error) {
	snippet, err := r.q.GetSnippet(ctx, db.GetSnippetParams{
		SnippetID: snippetID,
		UserID:    userID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, repository.WrapError(err, "failed to get snippet")
	}

	var avatar *string
	if snippet.AuthorAvatar.Valid {
		avatar = &snippet.AuthorAvatar.String
	}

	return &domain.Snippet{
		ID:       snippet.ID,
		Title:    snippet.Title,
		Content:  snippet.Content,
		Language: snippet.Language,
		Author: &domain.User{
			ID:       snippet.AuthorID.String,
			Username: snippet.AuthorUsername.String,
			Email:    snippet.AuthorEmail.String,
			Avatar:   avatar,
		},
		CreatedAt: snippet.CreatedAt,
		UpdatedAt: snippet.UpdatedAt,
		Views:     int(snippet.Views),
		Likes:     int(snippet.Likes),
		IsLiked:   snippet.IsLiked == 1,
		IsSaved:   snippet.IsSaved == 1,
	}, nil
}

func (r *SnippetRepository) GetAllByAuthor(ctx context.Context, authorID, userID string) ([]*domain.Snippet, error) {
	snippets, err := r.q.GetSnippetsByAuthor(ctx, db.GetSnippetsByAuthorParams{
		AuthorID: authorID,
		UserID:   userID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, repository.WrapError(err, "failed to get snippets by author")
	}

	result := make([]*domain.Snippet, len(snippets))
	for i, snippet := range snippets {
		var avatar *string
		if snippet.AuthorAvatar.Valid {
			avatar = &snippet.AuthorAvatar.String
		}

		result[i] = &domain.Snippet{
			ID:       snippet.ID,
			Title:    snippet.Title,
			Content:  snippet.Content,
			Language: snippet.Language,
			Author: &domain.User{
				ID:       snippet.AuthorID.String,
				Username: snippet.AuthorUsername.String,
				Email:    snippet.AuthorEmail.String,
				Avatar:   avatar,
			},
			CreatedAt: snippet.CreatedAt,
			UpdatedAt: snippet.UpdatedAt,
			Views:     int(snippet.Views),
			Likes:     int(snippet.Likes),
			IsLiked:   snippet.IsLiked == 1,
			IsSaved:   snippet.IsSaved == 1,
		}
	}

	return result, nil
}

func (r *SnippetRepository) GetAll(ctx context.Context, userID string) ([]*domain.Snippet, error) {
	snippets, err := r.q.GetSnippets(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, repository.WrapError(err, "failed to get snippets by author")
	}

	result := make([]*domain.Snippet, len(snippets))
	for i, snippet := range snippets {
		var avatar *string
		if snippet.AuthorAvatar.Valid {
			avatar = &snippet.AuthorAvatar.String
		}

		result[i] = &domain.Snippet{
			ID:       snippet.ID,
			Title:    snippet.Title,
			Content:  snippet.Content,
			Language: snippet.Language,
			Author: &domain.User{
				ID:       snippet.AuthorID.String,
				Username: snippet.AuthorUsername.String,
				Email:    snippet.AuthorEmail.String,
				Avatar:   avatar,
			},
			CreatedAt: snippet.CreatedAt,
			UpdatedAt: snippet.UpdatedAt,
			Views:     int(snippet.Views),
			Likes:     int(snippet.Likes),
			IsLiked:   snippet.IsLiked == 1,
			IsSaved:   snippet.IsSaved == 1,
		}
	}

	return result, nil
}

func (r *SnippetRepository) Update(ctx context.Context, snippet *domain.Snippet) error {
	_, err := r.q.UpdateSnippet(ctx, db.UpdateSnippetParams{
		SnippetID: snippet.ID, // The ID of the snippet to update
		Title:     snippet.Title,
		Content:   snippet.Content,
		Language:  snippet.Language,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.ErrNotFound
		}
		return repository.WrapError(err, "failed to delete snippet")
	}
	return nil
}

func (r *SnippetRepository) Delete(ctx context.Context, id string) error {
	err := r.q.DeleteSnippet(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.ErrNotFound
		}
		return repository.WrapError(err, "failed to delete snippet")
	}
	return nil
}
