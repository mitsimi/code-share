package storage

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"mitsimi.dev/codeShare/internal/auth"
	"mitsimi.dev/codeShare/internal/storage"

	db "mitsimi.dev/codeShare/internal/db/sqlc"
	_ "modernc.org/sqlite"
)

// CreateUser creates a new user
func (s *SQLiteStorage) CreateUser(username, email, password string) (db.User, error) {
	// Hash the password
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return db.User{}, err
	}

	user, err := s.q.CreateUser(s.ctx, db.CreateUserParams{
		ID:           uuid.NewString(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

// GetUserByID gets a user by ID
func (s *SQLiteStorage) GetUserByID(id storage.UserID) (db.User, error) {
	user, err := s.q.GetUser(s.ctx, string(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.User{}, errors.New("user not found")
		}
		return db.User{}, err
	}
	return user, nil
}

// GetUserByUsername gets a user by username
func (s *SQLiteStorage) GetUserByUsername(username string) (db.User, error) {
	user, err := s.q.GetUserByUsername(s.ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.User{}, errors.New("user not found")
		}
		return db.User{}, err
	}
	return user, nil
}

// GetUserByEmail gets a user by email
func (s *SQLiteStorage) GetUserByEmail(email string) (db.User, error) {
	user, err := s.q.GetUserByEmail(s.ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.User{}, errors.New("user not found")
		}
		return db.User{}, err
	}
	return user, nil
}

// UpdateUser updates a user info
func (s *SQLiteStorage) UpdateUser(userID storage.UserID, username, email string) (db.User, error) {
	user, err := s.q.UpdateUserInfo(s.ctx, db.UpdateUserInfoParams{
		ID:       userID,
		Username: username,
		Email:    email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.User{}, errors.New("user not found")
		}
		return db.User{}, err
	}
	return user, nil
}

// UpdateUserAvatar updates a user's avatar URL
func (s *SQLiteStorage) UpdateUserAvatar(userID storage.UserID, avatarURL string) (db.User, error) {
	user, err := s.q.UpdateUserAvatar(s.ctx, db.UpdateUserAvatarParams{
		ID:     userID,
		Avatar: sql.NullString{String: avatarURL, Valid: avatarURL != ""},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.User{}, errors.New("user not found")
		}
		return db.User{}, err
	}
	return user, nil
}

// UpdateUserPassword updates a user's password
func (s *SQLiteStorage) UpdateUserPassword(userID storage.UserID, password string) error {
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	return s.q.UpdateUserPassword(s.ctx, db.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: passwordHash,
	})
}
