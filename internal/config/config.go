package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config holds the application configuration
type Config struct {
	Environment string `env:"GO_ENV" env-default:"development"`
	Port        string `env:"PORT" env-default:"8080"`
	DBPath      string `env:"DB_PATH" env-default:"data/codeshare.db"`
	LogLevel    string `env:"LOG_LEVEL" env-default:"info"`
	Seed        bool   `env:"SEED" env-default:"false"`
	JWTSecret   string `env:"JWT_SECRET"`
}

// New creates a new configuration
func New() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(".env", cfg)
	if err != nil {
		log.Printf("WARNING: Could not read .env file: %v", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading environment variables: %w", err)
	}

	if cfg.JWTSecret == "" {
		if cfg.Environment == "production" {
			return nil, fmt.Errorf("JWT_SECRET environment variable is required in production")
		}
		// Only use default in development
		cfg.JWTSecret = "dev-secret-key"
		log.Printf("WARNING: Using default JWT secret in development environment. This is not secure for production use.")
	}

	return cfg, nil
}
