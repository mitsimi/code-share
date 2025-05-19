package main

import (
	"codeShare/frontend"
	"log"
	"mime"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func init() {
	// Add correct MIME types for modern web files
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".mjs", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
}



func main() {
	// Create a new Chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.CleanPath)
	r.Use(middleware.GetHead)

	// Create a file server handler for the embedded dist directory
	fs := http.FS(frontend.DistDirFS)

	// Handle all routes
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		serveStaticAsset(w, r, fs)
	})

	})

	// Start the server
	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}

