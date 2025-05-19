package storage

import (
	"codeShare/internal/models"
)

// Storage defines the interface for data storage
type Storage interface {
	// Snippet operations
	GetSnippets() ([]models.Snippet, error)
	GetSnippet(id string) (models.Snippet, error)
	CreateSnippet(snippet models.Snippet) error
	UpdateSnippet(snippet models.Snippet) error
	DeleteSnippet(id string) error
	ToggleLikeSnippet(id string) error
} 