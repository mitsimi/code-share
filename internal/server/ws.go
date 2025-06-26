package server

import (
	"github.com/go-chi/chi/v5"
	ws "mitsimi.dev/codeShare/internal/websocket"
)

// setupWebSocketRoutes configures WebSocket routes
func (s *Server) setupWebSocketRoutes(r chi.Router) {
	r.HandleFunc("/", ws.HandleWebSocket(s.wsHub))
	r.HandleFunc("/stats", ws.HandleStats(s.wsHub))
}
