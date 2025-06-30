package ws

import (
	"encoding/json"
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

		type StatsResponse struct {
			TotalClients         int `json:"total_clients"`
			UserSubscriptions    int `json:"user_subscriptions"`
			SnippetSubscriptions int `json:"snippet_subscriptions"`
			ListSubscriptions    int `json:"list_subscriptions"`
		}

		response := StatsResponse{
			TotalClients:         stats["total_clients"].(int),
			UserSubscriptions:    stats["user_subscriptions"].(int),
			SnippetSubscriptions: stats["snippet_subscriptions"].(int),
			ListSubscriptions:    stats["list_subscriptions"].(int),
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode stats", http.StatusInternalServerError)
			return
		}
	}
}
