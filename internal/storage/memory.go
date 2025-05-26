package storage

import (
	"codeShare/internal/models"
	"errors"
	"sync"
	"time"

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
	for _, snippet := range sampleSnippets {
		snippets[snippet.ID] = snippet
	}
	return &MemoryStorage{
		snippets: snippets,
		users:    make(map[string]models.User),
		sessions: make(map[string]models.Session),
	}
}

// CreateUser creates a new user
func (s *MemoryStorage) CreateUser(username, email, password string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if username or email already exists
	for _, user := range s.users {
		if user.Username == username || user.Email == email {
			return "", errors.New("username or email already exists")
		}
	}

	user := models.User{
		ID:           uuid.NewString(),
		Username:     username,
		Email:        email,
		PasswordHash: password, // In a real app, this should be hashed
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	s.users[user.ID] = user
	return user.ID, nil
}

// GetUser gets a user by ID
func (s *MemoryStorage) GetUser(id string) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

// GetUserByUsername gets a user by username
func (s *MemoryStorage) GetUserByUsername(username string) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Username == username {
			return user, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

// GetUserByEmail gets a user by email
func (s *MemoryStorage) GetUserByEmail(email string) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}
	return models.User{}, errors.New("user not found")
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
func (s *MemoryStorage) CreateSession(userID string, token string, expiresAt UnixTime) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session := models.Session{
		ID:        uuid.NewString(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
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

func (s *MemoryStorage) GetSnippets() ([]models.Snippet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snippets := make([]models.Snippet, 0, len(s.snippets))
	for _, snippet := range s.snippets {
		snippets = append(snippets, snippet)
	}
	return snippets, nil
}

func (s *MemoryStorage) GetSnippet(id string) (models.Snippet, error) {
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

func (s *MemoryStorage) ToggleLikeSnippet(id string, isLike bool) error {
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
