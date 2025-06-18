package repository

import (
	"context"

	"mitsimi.dev/codeShare/internal/domain"
)

type LikeRepository interface {
	ToggleLike(ctx context.Context, userID, snippetID string, isLike bool) error
	GetLikedSnippets(ctx context.Context, userID string) ([]*domain.Snippet, error)
}
