package server

import (
	"context"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mitsimi.dev/codeShare/frontend"
	"mitsimi.dev/codeShare/internal/api"
	"mitsimi.dev/codeShare/internal/logger"
	"mitsimi.dev/codeShare/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

// Server represents the application server
type Server struct {
	router     *chi.Mux
	httpServer *http.Server
	storage    storage.Storage
	logger     *zap.Logger
	secretKey  string
}

// New creates a new server instance
func New(storage storage.Storage, secretKey string) *Server {
	s := &Server{
		router:    chi.NewRouter(),
		storage:   storage,
		logger:    logger.Log,
		secretKey: secretKey,
	}
	s.setupMiddleware()
	s.setupRoutes()
	s.startSessionCleanup()
	return s
}

// setupMiddleware configures the server middleware
func (s *Server) setupMiddleware() {
	// Add request ID to context
	s.router.Use(middleware.RequestID)

	// Use our structured logger
	s.router.Use(logger.RequestLogger)

	// Other middleware
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

	// Create auth middleware
	authMiddleware := api.NewAuthMiddleware(s.storage, s.secretKey)

	// Auth routes
	s.router.Route("/api/auth", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{
				"http://localhost:3000",         // Development
				"https://codeshare.mitsimi.dev", // Production
			},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))

		handler := api.NewAuthHandler(s.storage, s.secretKey)
		r.Post("/signup", handler.Signup)
		r.Post("/login", handler.Login)
		r.Post("/logout", handler.Logout)
		r.Post("/refresh", handler.RefreshToken)
		r.Get("/me", handler.GetCurrentUser)
	})

	// API routes
	s.router.Route("/api/snippets", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{
				"http://localhost:3000",         // Development
				"https://codeshare.mitsimi.dev", // Production
			},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))

		handler := api.NewSnippetHandler(s.storage)

		// Public routes
		r.Group(func(r chi.Router) {
			r.Get("/", handler.GetSnippets)
			r.Get("/{id}", handler.GetSnippet)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Post("/", handler.CreateSnippet)
			r.Put("/{id}", handler.UpdateSnippet)
			r.Delete("/{id}", handler.DeleteSnippet)
			r.Patch("/{id}/like", handler.ToggleLikeSnippet)
		})
	})

	// Only serve static files if SERVE_STATIC is set to "true"
	if !(strings.ToLower(os.Getenv("SERVE_STATIC")) == "false") {
		// Handle all other routes by serving index.html
		s.router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
			indexFile, err := frontend.DistDirFS.Open("index.html")
			if err != nil {
				s.logger.Error("failed to load index.html",
					zap.Error(err),
					zap.String("request_id", middleware.GetReqID(r.Context())),
				)
				http.Error(w, "Error loading index.html", http.StatusInternalServerError)
				return
			}
			defer indexFile.Close()

			http.ServeContent(w, r, "index.html", time.Now(), indexFile.(io.ReadSeeker))
		})
	}
}

// startSessionCleanup starts a background goroutine to periodically clean up expired sessions
func (s *Server) startSessionCleanup() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour) // Run cleanup every hour
		defer ticker.Stop()

		for range ticker.C {
			if err := s.storage.DeleteExpiredSessions(); err != nil {
				s.logger.Error("Failed to delete expired sessions", zap.Error(err))
			} else {
				s.logger.Debug("Successfully cleaned up expired sessions")
			}
		}
	}()
}

// Start starts the server
func (s *Server) Start(port, env string) error {
	s.httpServer = &http.Server{
		Addr:    port,
		Handler: s.router,
		// Add reasonable timeouts
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	s.logger.Info("server starting",
		zap.String("port", port),
		zap.String("environment", env),
		zap.String("serve_static", os.Getenv("SERVE_STATIC")),
	)

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server...")

	if s.httpServer == nil {
		return nil
	}

	// Shutdown the HTTP server gracefully
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("server shutdown failed", zap.Error(err))
		return err
	}

	s.logger.Info("server shutdown completed successfully")
	return nil
}
