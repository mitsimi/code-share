package server

import (
	"context"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"time"

	"mitsimi.dev/codeShare/frontend"
	"mitsimi.dev/codeShare/internal/api"
	"mitsimi.dev/codeShare/internal/logger"
	"mitsimi.dev/codeShare/internal/repository"
	"mitsimi.dev/codeShare/internal/services"
	ws "mitsimi.dev/codeShare/internal/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

// Server represents the application server
type Server struct {
	router             *chi.Mux
	httpServer         *http.Server
	repos              *repository.Container
	viewTracker        *services.ViewTracker
	wsHub              *ws.Hub
	logger             *zap.Logger
	secretKey          string
	serveStatic        bool
	corsAllowedOrigins []string
	devProxy           *DevProxy
}

// New creates a new server instance
func New(
	repos *repository.Container,
	secretKey string,
	serveStatic bool,
	corsAllowedOrigins []string,
) *Server {
	// Create view tracker
	viewTracker := services.NewViewTracker(repos.Views)
	wsHub := ws.NewHub()

	s := &Server{
		router:             chi.NewRouter(),
		repos:              repos,
		viewTracker:        viewTracker,
		wsHub:              wsHub,
		logger:             logger.Log,
		secretKey:          secretKey,
		serveStatic:        serveStatic,
		corsAllowedOrigins: corsAllowedOrigins,
	}

	// Setup Vite dev server proxy if not serving static files
	if !serveStatic {
		if err := s.setupDevProxy(); err != nil {
			s.logger.Fatal("Failed to setup development proxy", zap.Error(err))
		}
	}

	s.setupMiddleware()
	s.setupRoutes()
	s.startSessionCleanup()
	s.startViewCleanup()

	// Start the WebSocket hub
	go wsHub.Run()
	s.logger.Info("WebSocket hub started")

	return s
}

// setupDevProxy configures the development proxy
func (s *Server) setupDevProxy() error {
	devProxy, err := NewDevProxy("http://localhost:3000", s.logger)
	if err != nil {
		return err
	}
	s.devProxy = devProxy
	return nil
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
		zap.Bool("serve_static", s.serveStatic),
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
	// Create auth middleware
	authMiddleware := api.NewAuthMiddleware(s.repos.Users, s.repos.Sessions, s.secretKey)

	s.router.Use(authMiddleware.TryAttachUserID) // Attach user ID to context
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   s.corsAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Setup WebSocket routes (must be before API routes to avoid middleware conflicts)
	s.router.Route("/ws", func(r chi.Router) {
		s.setupWebSocketRoutes(r)
	})

	// Setup API routes
	s.router.Route("/api", func(r chi.Router) {
		s.setupAPIRoutes(r, authMiddleware)
	})

	// Serve static files in production or proxy to Vite in development
	if s.serveStatic {
		s.setupStaticRoutes()
	} else {
		s.setupDevProxyRoutes()
	}
}

// setupDevProxyRoutes sets up routes that proxy to Vite dev server during development
func (s *Server) setupDevProxyRoutes() {
	if s.devProxy == nil {
		s.logger.Warn("Development proxy not configured, but serve_static is false")
		return
	}

	// Proxy all non-API, non-WebSocket requests to Vite dev server
	s.router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		if !s.devProxy.ShouldProxy(r.URL.Path) {
			http.NotFound(w, r)
			return
		}

		s.devProxy.ServeHTTP(w, r)
	})

	s.logger.Info("Development proxy routes configured - non-API requests will be forwarded to Vite dev server")
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

	// Handle public files (favicon.ico, site.webmanifest, etc.)
	s.router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:] // Remove leading slash

		// Only serve files that are likely to be public files (not paths with slashes)
		if path != "" && filepath.Dir(path) != "." {
			return // Skip if it's a nested path, let other handlers deal with it
		}

		// Try to open the file from the embedded filesystem
		file, err := frontend.DistDirFS.Open(path)
		if err != nil {
			return // File doesn't exist, let the next handler try
		}
		defer file.Close()

		// Get file info to check if it's a directory
		info, err := file.Stat()
		if err != nil || info.IsDir() {
			return // Skip directories or files we can't stat
		}

		// Set appropriate content type
		ext := filepath.Ext(path)
		mimeType := mime.TypeByExtension(ext)
		if mimeType != "" {
			w.Header().Set("Content-Type", mimeType)
		}

		// Serve the file
		http.ServeContent(w, r, path, info.ModTime(), file.(io.ReadSeeker))
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
