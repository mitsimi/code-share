package storage

import (
	"errors"
	"sync"
	"time"

	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/models"

	"github.com/google/uuid"
)

var _ Storage = (*MemoryStorage)(nil)

// MemoryStorage implements Storage interface with in-memory storage
type MemoryStorage struct {
	snippets map[string]models.Snippet
	users    map[string]models.User
	sessions map[string]models.Session
	mu       sync.RWMutex
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
	snippets := make(map[string]models.Snippet)
	for _, snippet := range SampleSnippets {
		snippets[snippet.ID] = snippet
	}
	return &MemoryStorage{
		snippets: snippets,
		users:    make(map[string]models.User),
		sessions: make(map[string]models.Session),
	}
}

// CreateUser creates a new user
func (s *MemoryStorage) CreateUser(username, email, password string) (db.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if username or email already exists
	for _, user := range s.users {
		if user.Username == username || user.Email == email {
			return db.User{}, errors.New("username or email already exists")
		}
	}

	user := db.User{
		ID:           uuid.NewString(),
		Username:     username,
		Email:        email,
		PasswordHash: password,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	s.users[user.ID] = models.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
	return user, nil
}

// GetUserByID gets a user by ID
func (s *MemoryStorage) GetUserByID(id UserID) (db.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[string(id)]
	if !exists {
		return db.User{}, errors.New("user not found")
	}
	return db.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

// GetUserByUsername gets a user by username
func (s *MemoryStorage) GetUserByUsername(username string) (db.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Username == username {
			return db.User{
				ID:           user.ID,
				Username:     user.Username,
				Email:        user.Email,
				PasswordHash: user.PasswordHash,
				CreatedAt:    user.CreatedAt,
				UpdatedAt:    user.UpdatedAt,
			}, nil
		}
	}
	return db.User{}, errors.New("user not found")
}

// GetUserByEmail gets a user by email
func (s *MemoryStorage) GetUserByEmail(email string) (db.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Email == email {
			return db.User{
				ID:           user.ID,
				Username:     user.Username,
				Email:        user.Email,
				PasswordHash: user.PasswordHash,
				CreatedAt:    user.CreatedAt,
				UpdatedAt:    user.UpdatedAt,
			}, nil
		}
	}
	return db.User{}, errors.New("user not found")
}

// Login authenticates a user
func (s *MemoryStorage) Login(email, password string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Email == email && user.PasswordHash == password { // In a real app, use proper password comparison
			return user.ID, nil
		}
	}
	return "", errors.New("invalid credentials")
}

// CreateSession creates a new session
func (s *MemoryStorage) CreateSession(userID string, token string, refreshToken string, expiresAt UnixTime) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session := models.Session{
		ID:           uuid.NewString(),
		UserID:       userID,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		CreatedAt:    time.Now(),
	}
	s.sessions[token] = session
	return nil
}

// GetSession gets a session by token
func (s *MemoryStorage) GetSession(token string) (models.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[token]
	if !exists {
		return models.Session{}, errors.New("session not found")
	}
	return session, nil
}

// DeleteSession deletes a session
func (s *MemoryStorage) DeleteSession(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, token)
	return nil
}

// DeleteExpiredSessions deletes expired sessions
func (s *MemoryStorage) DeleteExpiredSessions() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for token, session := range s.sessions {
		if session.ExpiresAt < time.Now().Unix() {
			delete(s.sessions, token)
		}
	}
	return nil
}

// UpdateSessionExpiry updates the expiration time and refresh token of a session
func (s *MemoryStorage) UpdateSessionExpiry(token string, expiresAt UnixTime, refreshToken string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[token]
	if !exists {
		return errors.New("session not found")
	}

	session.ExpiresAt = expiresAt
	session.RefreshToken = refreshToken
	s.sessions[token] = session
	return nil
}

func (s *MemoryStorage) GetSnippets(id UserID) ([]models.Snippet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snippets := make([]models.Snippet, 0, len(s.snippets))
	for _, snippet := range s.snippets {
		snippets = append(snippets, snippet)
	}
	return snippets, nil
}

func (s *MemoryStorage) GetSnippet(userUd, id string) (models.Snippet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snippet, exists := s.snippets[id]
	if !exists {
		return models.Snippet{}, errors.New("snippet not found")
	}
	return snippet, nil
}

func (s *MemoryStorage) CreateSnippet(snippet models.Snippet) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	snippet.ID = uuid.NewString()
	snippet.CreatedAt = time.Now()
	snippet.UpdatedAt = time.Now()
	snippet.Likes = 0
	snippet.IsLiked = false

	s.snippets[snippet.ID] = snippet
	return snippet.ID, nil
}

func (s *MemoryStorage) UpdateSnippet(snippet models.Snippet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.snippets[snippet.ID]; !exists {
		return errors.New("snippet not found")
	}

	snippet.UpdatedAt = time.Now()
	s.snippets[snippet.ID] = snippet
	return nil
}

func (s *MemoryStorage) DeleteSnippet(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.snippets[id]; !exists {
		return errors.New("snippet not found")
	}

	delete(s.snippets, id)
	return nil
}

func (s *MemoryStorage) ToggleLikeSnippet(userID, id string, isLike bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	snippet, exists := s.snippets[id]
	if !exists {
		return errors.New("snippet not found")
	}

	if isLike {
		if !snippet.IsLiked {
			snippet.Likes++
			snippet.IsLiked = true
		}
	} else {
		if snippet.IsLiked {
			snippet.Likes--
			snippet.IsLiked = false
		}
	}
	s.snippets[id] = snippet
	return nil
}

// Close implements the Storage interface
func (s *MemoryStorage) Close() error {
	return nil
}

// Seed populates the storage with sample data
func (s *MemoryStorage) Seed() error {
	return nil
}

// UpdateUser updates a user info
func (s *MemoryStorage) UpdateUser(userID UserID, username, email string) (db.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[string(userID)]
	if !exists {
		return db.User{}, errors.New("user not found")
	}

	user.Username = username
	user.Email = email
	user.UpdatedAt = time.Now()
	s.users[string(userID)] = user

	return db.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

// UpdateUserAvatar updates a user's avatar URL
func (s *MemoryStorage) UpdateUserAvatar(userID UserID, avatarURL string) (db.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[string(userID)]
	if !exists {
		return db.User{}, errors.New("user not found")
	}

	user.Avatar = avatarURL
	user.UpdatedAt = time.Now()
	s.users[string(userID)] = user
	return db.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

// UpdateUserPassword updates a user's password
func (s *MemoryStorage) UpdateUserPassword(userID UserID, password string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[string(userID)]
	if !exists {
		return errors.New("user not found")
	}

	user.PasswordHash = password // In a real app, this would be hashed
	user.UpdatedAt = time.Now()
	s.users[string(userID)] = user
	return nil
}
