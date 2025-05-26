package main

import (
	"codeShare/internal/config"
	"codeShare/internal/logger"
	"codeShare/internal/server"
	"codeShare/internal/storage"
	"log"

	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize logger with config
	if err := logger.Init(logger.Config{
		Environment: cfg.Environment,
		Level:       cfg.LogLevel,
	}); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	// Initialize SQLite storage
	store, err := storage.NewSQLiteStorage(cfg.DBPath)
	if err != nil {
		logger.Fatal("Failed to initialize SQLite storage", zap.Error(err))
	}
	defer store.Close()

	if cfg.Seed {
		logger.Debug("Seeding database")
		// Seed the database with sample data
		if err := store.Seed(); err != nil {
			logger.Error("Warning: Failed to seed database", zap.Error(err))
		}
		logger.Debug("Seeding finished successfully")
	}

	// Create and start server
	srv := server.New(store)
	if err := srv.Start(":"+cfg.Port, cfg.Environment); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
