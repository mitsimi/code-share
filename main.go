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
	sqlite "mitsimi.dev/codeShare/internal/storage/sqlite"

	_ "github.com/mattn/go-sqlite3"
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

	// Initialize SQLite database
	sqliteStorage, err := sqlite.New(cfg.DBPath)
	if err != nil {
		logger.Fatal("Failed to initialize SQLite storage", zap.Error(err))
	}
	defer sqliteStorage.Close()

	// Initialize repositories
	snippets := sqlite.NewSnippetRepository(sqliteStorage.DB())
	likes := sqlite.NewLikeRepository(sqliteStorage.DB())
	bookmarks := sqlite.NewBookmarkRepository(sqliteStorage.DB())
	users := sqlite.NewUserRepository(sqliteStorage.DB())
	sessions := sqlite.NewSessionRepository(sqliteStorage.DB())
	views := sqlite.NewViewRepository(sqliteStorage.DB())

	// Create storage instance
	storage := storage.NewStorage(snippets, likes, bookmarks, users, sessions)

	if cfg.Seed {
		logger.Debug("Seeding database")
		if err := storage.SeedSampleData(context.Background()); err != nil {
			logger.Error("Failed to seed database", zap.Error(err))
		} else {
			logger.Debug("Seeding finished successfully")
		}
	}

	// Create server with repositories
	srv := server.New(
		snippets,
		likes,
		bookmarks,
		users,
		sessions,
		views,
		cfg.JWTSecret,
	)

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

	logger.Info("Application exited gracefully")
}
