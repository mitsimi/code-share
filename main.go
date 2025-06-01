package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mitsimi.dev/codeShare/internal/config"
	"mitsimi.dev/codeShare/internal/logger"
	"mitsimi.dev/codeShare/internal/server"
	"mitsimi.dev/codeShare/internal/storage"

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

	if cfg.Seed {
		logger.Debug("Seeding database")
		// Seed the database with sample data
		if err := store.Seed(); err != nil {
			logger.Error("Warning: Failed to seed database", zap.Error(err))
		}
		logger.Debug("Seeding finished successfully")
	}

	// Create server
	srv := server.New(store, cfg.JWTSecret)

	// Channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(":"+cfg.Port, cfg.Environment); err != nil {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	logger.Info("Server started successfully. Press Ctrl+C to shutdown gracefully")

	// Wait for interrupt signal
	sig := <-sigChan
	logger.Info("Received shutdown signal", zap.String("signal", sig.String()))

	// Create a context with timeout for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Perform graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Graceful shutdown failed", zap.Error(err))
	}

	// Close storage connection
	if err := store.Close(); err != nil {
		logger.Error("Failed to close storage", zap.Error(err))
	} else {
		logger.Info("Storage closed successfully")
	}

	logger.Info("Application exited gracefully")
}
