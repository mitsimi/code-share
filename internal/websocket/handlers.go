package ws

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"mitsimi.dev/codeShare/internal/api"
	"mitsimi.dev/codeShare/internal/logger"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin checking
		return true
	},
}

// HandleWebSocket creates the HTTP handler for WebSocket connections
func HandleWebSocket(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		userID := api.GetUserID(r)
		if userID == "" {
			userID = "anonymous"
		}
		log := logger.With(zap.String("request_id", requestID), zap.String("user_id", userID))

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Warn("WebSocket upgrade failed", zap.Error(err))
			return
		}

		client := NewClient(hub, conn, userID)
		client.hub.register <- client

		log.Info("WebSocket connection established", zap.String("user_id", userID), zap.String("request_id", requestID))
		// Start client pumps
		go client.WritePump()
		go client.ReadPump()
	}
}

// HandleStats provides WebSocket statistics endpoint
func HandleStats(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		stats := hub.GetStats()

		// Simple JSON response (you might want to use a proper JSON encoder)
		w.Write([]byte(`{
            "total_clients": ` + string(rune(stats["total_clients"].(int))) + `,
            "user_subscriptions": ` + string(rune(stats["user_subscriptions"].(int))) + `,
            "snippet_subscriptions": ` + string(rune(stats["snippet_subscriptions"].(int))) + `,
            "list_subscriptions": ` + string(rune(stats["list_subscriptions"].(int))) + `
        }`))
	}
}
