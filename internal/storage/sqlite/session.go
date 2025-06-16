package storage

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"mitsimi.dev/codeShare/internal/auth"
	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/models"
	"mitsimi.dev/codeShare/internal/storage"
)

// Login authenticates a user and returns their ID
func (s *SQLiteStorage) Login(email, password string) (storage.UserID, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return "", auth.ErrInvalidCredentials
	}

	if !auth.CheckPasswordHash(password, user.PasswordHash) {
		return "", auth.ErrInvalidCredentials
	}

	return user.ID, nil
}

// CreateSession creates a new session for a user
func (s *SQLiteStorage) CreateSession(userID string, token string, refreshToken string, expiresAt storage.UnixTime) error {
	_, err := s.q.CreateSession(s.ctx, db.CreateSessionParams{
		ID:           uuid.NewString(),
		UserID:       userID,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	})
	return err
}

// GetSession gets a session by token
func (s *SQLiteStorage) GetSession(token string) (models.Session, error) {
	session, err := s.q.GetSession(s.ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Session{}, errors.New("session not found")
		}
		return models.Session{}, err
	}

	return models.Session{
		ID:           session.ID,
		UserID:       session.UserID,
		Token:        session.Token,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}, nil
}

// DeleteSession deletes a session by token
func (s *SQLiteStorage) DeleteSession(token string) error {
	return s.q.DeleteSession(s.ctx, token)
}

// DeleteExpiredSessions deletes all expired sessions
func (s *SQLiteStorage) DeleteExpiredSessions() error {
	return s.q.DeleteExpiredSessions(s.ctx)
}

// UpdateSessionExpiry updates the expiry time and refresh token of a session
func (s *SQLiteStorage) UpdateSessionExpiry(token string, expiresAt storage.UnixTime, refreshToken string) error {
	// First get the session to verify it exists
	session, err := s.GetSession(token)
	if err != nil {
		return err
	}

	err = s.q.UpdateSessionExpiry(s.ctx, db.UpdateSessionExpiryParams{
		Token:        session.Token,
		ExpiresAt:    expiresAt,
		RefreshToken: refreshToken,
	})
	return err
}
