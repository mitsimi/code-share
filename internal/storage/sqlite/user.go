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

// GetUser gets a user by ID
func (s *SQLiteStorage) GetUser(id storage.UserID) (db.User, error) {
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
