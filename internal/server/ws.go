package server

import "mitsimi.dev/codeShare/internal/api/handler"

// setupWebSocketRoutes configures WebSocket routes
func (s *Server) setupWebSocketRoutes() {
	// Import WebSocket handler
	websocketHandler := handler.NewWebSocketHandler()

	// WebSocket route - must be outside API group to avoid middleware conflicts
	s.router.HandleFunc("/ws", websocketHandler.WebSocket)
}
