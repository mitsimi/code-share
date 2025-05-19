package main

import (
	"codeShare/internal/server"
	"codeShare/internal/storage"
	"log"
)

func main() {
	// Initialize storage
	store := storage.NewMemoryStorage()

	// Create and start server
	srv := server.New(store)
	if err := srv.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

