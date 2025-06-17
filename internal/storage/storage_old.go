package storage

import (
	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/models"
)

// StorageOLD defines the interface for data storage
type StorageOLD interface {
	// User management
	CreateUser(username, email, password string) (db.User, error)
	GetUserByID(id UserID) (db.User, error)
	GetUserByUsername(username string) (db.User, error)
	GetUserByEmail(email string) (db.User, error)
	UpdateUser(userID UserID, username, email string) (db.User, error)
	UpdateUserAvatar(userID UserID, avatarURL string) (db.User, error)
	UpdateUserPassword(userID UserID, password string) error

	// Session management
	Login(email, password string) (UserID, error)
	CreateSession(id UserID, token string, refreshToken string, expiresAt UnixTime) error
	GetSession(token string) (models.Session, error)
	DeleteSession(token string) error
	DeleteExpiredSessions() error
	UpdateSessionExpiry(sessionID string, expiresAt UnixTime, refreshToken string) error

	// Snippet management
	GetSnippet(userID UserID, id SnippetID) (models.Snippet, error)
	GetSnippets(userID UserID) ([]models.Snippet, error)
	GetSnippetsByAuthor(userID, authorID UserID) ([]models.Snippet, error)

	CreateSnippet(snippet models.Snippet) (SnippetID, error)
	UpdateSnippet(snippet models.Snippet) error
	DeleteSnippet(id SnippetID) error

	// Like management
	ToggleLikeSnippet(userID UserID, id SnippetID, isLike bool) error
	GetLikedSnippets(userID UserID) ([]models.Snippet, error)

	// Bookmark management
	ToggleSaveSnippet(userID UserID, id SnippetID, isSave bool) error
	GetSavedSnippets(userID UserID) ([]models.Snippet, error)

	// Database management
	Seed() error
	Close() error
}
