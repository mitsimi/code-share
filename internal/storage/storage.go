package storage

import (
	"context"

	"mitsimi.dev/codeShare/internal/domain"
	"mitsimi.dev/codeShare/internal/repository"
)

type (
	UserID    = string
	SnippetID = string

	UnixTime = int64
)

type Storage struct {
	snippets  repository.SnippetRepository
	likes     repository.LikeRepository
	bookmarks repository.BookmarkRepository
	users     repository.UserRepository
	sessions  repository.SessionRepository
}

func NewStorage(
	snippets repository.SnippetRepository,
	likes repository.LikeRepository,
	bookmarks repository.BookmarkRepository,
	users repository.UserRepository,
	sessions repository.SessionRepository,
) *Storage {
	return &Storage{
		snippets:  snippets,
		likes:     likes,
		bookmarks: bookmarks,
		users:     users,
		sessions:  sessions,
	}
}

// Snippet operations
func (s *Storage) CreateSnippet(ctx context.Context, snippet *domain.Snippet) error {
	return s.snippets.Create(ctx, snippet)
}

func (s *Storage) GetSnippet(ctx context.Context, id, userID string) (*domain.Snippet, error) {
	return s.snippets.GetByID(ctx, id, userID)
}

func (s *Storage) GetSnippets(ctx context.Context, userID string) ([]*domain.Snippet, error) {
	return s.snippets.GetAll(ctx, userID)
}

func (s *Storage) UpdateSnippet(ctx context.Context, snippet *domain.Snippet) error {
	return s.snippets.Update(ctx, snippet)
}

func (s *Storage) DeleteSnippet(ctx context.Context, id string) error {
	return s.snippets.Delete(ctx, id)
}

// Like operations
func (s *Storage) GetLikedSnippets(ctx context.Context, userID string) ([]*domain.Snippet, error) {
	return s.likes.GetLikedSnippets(ctx, userID)
}

func (s *Storage) ToggleLikeSnippet(ctx context.Context, userID, snippetID string, isLike bool) error {
	return s.likes.ToggleLike(ctx, userID, snippetID, isLike)
}

// Bookmark operations
func (s *Storage) GetSavedSnippets(ctx context.Context, userID string) ([]*domain.Snippet, error) {
	return s.bookmarks.GetSavedSnippets(ctx, userID)
}

func (s *Storage) ToggleSaveSnippet(ctx context.Context, userID, snippetID string, isSave bool) error {
	return s.bookmarks.ToggleSave(ctx, userID, snippetID, isSave)
}

// User operations
func (s *Storage) CreateUser(ctx context.Context, user *domain.UserCreation) error {
	return s.users.Create(ctx, user)
}

func (s *Storage) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.users.GetByUsername(ctx, username)
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.users.GetByEmail(ctx, email)
}

func (s *Storage) UpdateUser(ctx context.Context, user *domain.User) error {
	return s.users.Update(ctx, user)
}

// Session operations
func (s *Storage) CreateSession(ctx context.Context, session *domain.Session) error {
	return s.sessions.Create(ctx, session)
}

func (s *Storage) GetSession(ctx context.Context, token string) (*domain.Session, error) {
	return s.sessions.GetByToken(ctx, token)
}

func (s *Storage) DeleteSession(ctx context.Context, token string) error {
	return s.sessions.Delete(ctx, token)
}

func (s *Storage) DeleteExpiredSessions(ctx context.Context) error {
	return s.sessions.DeleteExpired(ctx)
}

func (s *Storage) UpdateSessionExpiry(ctx context.Context, sessionID string, expiresAt int64, refreshToken string) error {
	return s.sessions.UpdateExpiry(ctx, sessionID, expiresAt, refreshToken)
}
