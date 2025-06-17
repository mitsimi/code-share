package repository

import (
	"context"

	"mitsimi.dev/codeShare/internal/domain"
)

type SnippetRepository interface {
	Create(ctx context.Context, snippet *domain.Snippet) error
	GetByID(ctx context.Context, id string, userID string) (*domain.Snippet, error)
	GetAll(ctx context.Context, userID string) ([]*domain.Snippet, error)
	GetAllByAuthor(ctx context.Context, authorID string, userID string) ([]*domain.Snippet, error)
	Update(ctx context.Context, snippet *domain.Snippet) error
	Delete(ctx context.Context, id string) error
}
