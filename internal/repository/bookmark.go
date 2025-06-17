package repository

import (
	"context"

	"mitsimi.dev/codeShare/internal/domain"
)

type BookmarkRepository interface {
	ToggleSave(ctx context.Context, userID, snippetID string, isSave bool) error
	GetSavedSnippets(ctx context.Context, userID string) ([]*domain.Snippet, error)
}
