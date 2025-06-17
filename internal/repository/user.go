package repository

import (
	"context"

	"mitsimi.dev/codeShare/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.UserCreation) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdateAvatar(ctx context.Context, userID, avatarURL string) error
	UpdatePassword(ctx context.Context, userID, password string) error
}
