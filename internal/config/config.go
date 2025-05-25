package config

import (
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	DBPath string
	Port   string
}

// New creates a new Config instance
func New() (*Config, error) {
	// Get database path from environment variable or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = filepath.Join(".", "data", "sqlite.db")
	}

	// Ensure the database directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, err
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DBPath: dbPath,
		Port:   port,
	}, nil
}
