package storage

import (
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
	CreateUser(username, email, password string) (UserID, error)
	GetUser(id UserID) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	Login(email, password string) (UserID, error)

	// Session management
	CreateSession(id UserID, token string, expiresAt UnixTime) error
	GetSession(token string) (models.Session, error)
	DeleteSession(token string) error
	UpdateSessionExpiry(sessionID string, expiresAt UnixTime) error

	// Snippet management
	GetSnippets() ([]models.Snippet, error)
	GetSnippet(id SnippetID) (models.Snippet, error)
	CreateSnippet(snippet models.Snippet) (SnippetID, error)
	UpdateSnippet(snippet models.Snippet) error
	DeleteSnippet(id SnippetID) error
	ToggleLikeSnippet(id SnippetID, isLike bool) error

	// Database management
	Seed() error
	Close() error
}
