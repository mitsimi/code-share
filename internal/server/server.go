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
	"mitsimi.dev/codeShare/internal/logger"
	"mitsimi.dev/codeShare/internal/repository"
	"mitsimi.dev/codeShare/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Server represents the application server
type Server struct {
	router      *chi.Mux
	httpServer  *http.Server
	snippets    repository.SnippetRepository
	likes       repository.LikeRepository
	bookmarks   repository.BookmarkRepository
	users       repository.UserRepository
	sessions    repository.SessionRepository
	views       repository.ViewRepository
	viewTracker *services.ViewTracker
	logger      *zap.Logger
	secretKey   string
}

// New creates a new server instance
func New(
	snippets repository.SnippetRepository,
	likes repository.LikeRepository,
	bookmarks repository.BookmarkRepository,
	users repository.UserRepository,
	sessions repository.SessionRepository,
	views repository.ViewRepository,
	secretKey string,
) *Server {
	// Create view tracker
	viewTracker := services.NewViewTracker(views)

	s := &Server{
		router:      chi.NewRouter(),
		snippets:    snippets,
		likes:       likes,
		bookmarks:   bookmarks,
		users:       users,
		sessions:    sessions,
		views:       views,
		viewTracker: viewTracker,
		logger:      logger.Log,
		secretKey:   secretKey,
	}
	s.setupMiddleware()
	s.setupRoutes()
	s.startSessionCleanup()
	s.startViewCleanup()
	return s
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
	// Setup API routes
	s.router.Route("/api", s.setupAPIRoutes)

	// Only serve static files if SERVE_STATIC is set to "true"
	if !(strings.ToLower(os.Getenv("SERVE_STATIC")) == "false") {
		s.setupStaticRoutes()
	}
}

func (s *Server) setupStaticRoutes() {
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
