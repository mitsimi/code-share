package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"mitsimi.dev/codeShare/internal/api"
	"mitsimi.dev/codeShare/internal/logger"
)

type WebSocketHandler struct {
	upgrader websocket.Upgrader
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{upgrader: websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}}
}

func (h *WebSocketHandler) WebSocket(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	userID := api.GetUserID(r)
	log := logger.With(zap.String("request_id", requestID), zap.String("user_id", userID))

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("WebSocket upgrade failed", zap.Error(err))
		return
	}
	defer conn.Close()

	log.Info("Client connected")

	// Handle messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Error("Read error", zap.Error(err))
			break
		}

		log.Info("Received message", zap.String("message", string(message)))

		// Echo the message back to client
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Error("Write error", zap.Error(err))
			break
		}
	}
}
