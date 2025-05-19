package server

import (
	"codeShare/frontend"
	"codeShare/internal/api"
	"codeShare/internal/storage"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server represents the application server
type Server struct {
	router  *chi.Mux
	storage storage.Storage
}

// New creates a new server instance
func New(storage storage.Storage) *Server {
	s := &Server{
		router:  chi.NewRouter(),
		storage: storage,
	}
	s.setupMiddleware()
	s.setupRoutes()
	return s
}

// setupMiddleware configures the server middleware
func (s *Server) setupMiddleware() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.CleanPath)
	s.router.Use(middleware.GetHead)
}

// setupRoutes configures the server routes
func (s *Server) setupRoutes() {
	// Create a file server handler for the embedded dist directory
	fs := http.FileServer(http.FS(frontend.DistDirFS))

	// Handle static assets
	s.router.HandleFunc("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		ext := filepath.Ext(path)
		mimeType := mime.TypeByExtension(ext)

		if mimeType != "" {
			w.Header().Set("Content-Type", mimeType)
		}
		fs.ServeHTTP(w, r)
	})

	// Handle favicon
	s.router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})

	// API routes
	s.router.Route("/api/snippets", func(r chi.Router) {
		handler := api.NewSnippetHandler(s.storage)
		r.Get("/", handler.GetSnippets)
		r.Post("/", handler.CreateSnippet)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetSnippet)
			r.Put("/", handler.UpdateSnippet)
			r.Delete("/", handler.DeleteSnippet)
			r.Patch("/like", handler.ToggleLikeSnippet)
		})
	})

	// Handle all other routes by serving index.html
	s.router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		indexFile, err := frontend.DistDirFS.Open("index.html")
		if err != nil {
			http.Error(w, "Error loading index.html", http.StatusInternalServerError)
			return
		}
		defer indexFile.Close()

		http.ServeContent(w, r, "index.html", time.Now(), indexFile.(io.ReadSeeker))
	})
}

// Start starts the server
func (s *Server) Start(port string) error {
	log.Printf("Server starting on http://localhost%s", port)
	return http.ListenAndServe(port, s.router)
} 