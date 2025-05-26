package main

import (
	"codeShare/internal/config"
	"codeShare/internal/server"
	"codeShare/internal/storage"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize SQLite storage
	store, err := storage.NewSQLiteStorage(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize SQLite storage: %v", err)
	}
	defer store.Close()

	if cfg.Seed {
		log.Default().Println("Seeding")
		// Seed the database with sample data
		if err := store.Seed(); err != nil {
			log.Printf("Warning: Failed to seed database: %v", err)
		}
		log.Default().Println("Seeding done")
	}

	// Create and start server
	srv := server.New(store)
	if err := srv.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
