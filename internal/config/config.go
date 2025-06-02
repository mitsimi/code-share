package config

import (
	"fmt"
	"log"
	"os"
)

// Config holds the application configuration
type Config struct {
	Environment string
	Port        string
	DBPath      string
	LogLevel    string
	Seed        bool
	JWTSecret   string
}

// New creates a new configuration
func New() (*Config, error) {
	// Get environment
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get database path
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/codeshare.db"
	}

	// Get log level
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	// Get seed flag
	seed := os.Getenv("SEED") == "true"

	// Get JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		if env == "production" {
			return nil, fmt.Errorf("JWT_SECRET environment variable is required in production")
		}
		// Only use default in development
		jwtSecret = "dev-secret-key"
		log.Printf("WARNING: Using default JWT secret in development environment. This is not secure for production use.")
	}

	return &Config{
		Environment: env,
		Port:        port,
		DBPath:      dbPath,
		LogLevel:    logLevel,
		Seed:        seed,
		JWTSecret:   jwtSecret,
	}, nil
}
