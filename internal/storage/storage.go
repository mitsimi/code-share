package storage

import (
	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/models"
)

type (
	UserID    = string
	SnippetID = string

	UnixTime = int64
)

// Storage defines the interface for data storage
type Storage interface {
	// User management
	CreateUser(username, email, password string) (db.User, error)
	GetUser(id UserID) (db.User, error)
	GetUserByUsername(username string) (db.User, error)
	GetUserByEmail(email string) (db.User, error)
	Login(email, password string) (UserID, error)

	// Session management
	CreateSession(id UserID, token string, refreshToken string, expiresAt UnixTime) error
	GetSession(token string) (models.Session, error)
	DeleteSession(token string) error
	DeleteExpiredSessions() error
	UpdateSessionExpiry(sessionID string, expiresAt UnixTime, refreshToken string) error

	// Snippet management
	GetSnippets(userID UserID) ([]models.Snippet, error)
	GetSnippet(userID UserID, id SnippetID) (models.Snippet, error)
	CreateSnippet(snippet models.Snippet) (SnippetID, error)
	UpdateSnippet(snippet models.Snippet) error
	DeleteSnippet(id SnippetID) error
	ToggleLikeSnippet(userID UserID, id SnippetID, isLike bool) error

	// Database management
	Seed() error
	Close() error
}
