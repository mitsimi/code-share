package sqlite

import (
	"context"
	"database/sql"

	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/domain"
	"mitsimi.dev/codeShare/internal/repository"
)

var _ repository.SessionRepository = (*SessionRepository)(nil)

type SessionRepository struct {
	db *sql.DB
	q  *db.Queries
}

func NewSessionRepository(dbConn *sql.DB) *SessionRepository {
	return &SessionRepository{
		db: dbConn,
		q:  db.New(dbConn),
	}
}

func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	_, err := r.q.CreateSession(ctx, db.CreateSessionParams{
		ID:           session.ID,
		UserID:       session.UserID,
		Token:        session.Token,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt,
	})
	if err != nil {
		return repository.WrapError(err, "failed to create session")
	}
	return nil
}

func (r *SessionRepository) GetByToken(ctx context.Context, token string) (*domain.Session, error) {
	session, err := r.q.GetSession(ctx, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, repository.WrapError(err, "failed to get session")
	}

	return &domain.Session{
		ID:           session.ID,
		UserID:       session.UserID,
		Token:        session.Token,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}, nil
}

func (r *SessionRepository) Delete(ctx context.Context, token string) error {
	err := r.q.DeleteSession(ctx, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.ErrNotFound
		}
		return repository.WrapError(err, "failed to delete session")
	}
	return nil
}

func (r *SessionRepository) DeleteExpired(ctx context.Context) error {
	err := r.q.DeleteExpiredSessions(ctx)
	if err != nil {
		return repository.WrapError(err, "failed to delete expired sessions")
	}
	return nil
}

func (r *SessionRepository) UpdateExpiry(ctx context.Context, token string, expiresAt domain.UnixTime, refreshToken string) error {
	// First get the session to verify it exists
	session, err := r.GetByToken(ctx, token)
	if err != nil {
		return err
	}

	err = r.q.UpdateSessionExpiry(ctx, db.UpdateSessionExpiryParams{
		Token:        session.Token,
		ExpiresAt:    expiresAt,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return repository.WrapError(err, "failed to update session expiry")
	}
	return nil
}
