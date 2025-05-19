package storage

import (
	"codeShare/internal/models"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MemoryStorage implements Storage interface with in-memory storage
type MemoryStorage struct {
	snippets map[string]models.Snippet
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
	}
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

func (s *MemoryStorage) CreateSnippet(snippet models.Snippet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	snippet.ID = uuid.New().String()
	snippet.CreatedAt = time.Now()
	snippet.UpdatedAt = time.Now()
	snippet.Likes = 0

	s.snippets[snippet.ID] = snippet
	return nil
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

func (s *MemoryStorage) ToggleLikeSnippet(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	snippet, exists := s.snippets[id]
	if !exists {
		return errors.New("snippet not found")
	}

	snippet.Likes++
	s.snippets[id] = snippet
	return nil
} 