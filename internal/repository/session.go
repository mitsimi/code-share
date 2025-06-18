package repository

import (
	"context"

	"mitsimi.dev/codeShare/internal/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByToken(ctx context.Context, token string) (*domain.Session, error)
	Delete(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
	UpdateExpiry(ctx context.Context, sessionID string, expiresAt domain.UnixTime, refreshToken string) error
}
