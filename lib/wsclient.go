package lib

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebsocketMessage struct {
	MessageType string `json:"messageType"`
	Data        gin.H  `json:"data"`
}

var clients = make(map[*websocket.Conn]bool)

func HandleWSConnection(c *gin.Context, logger *slog.Logger, wsConn *websocket.Conn) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return c.Request.Header.Get("Origin") == "http://localhost:5173"
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("Error closing WS connection: %v", err.Error()))
		}
	}(conn)
	clients[conn] = true
	logger.Info(fmt.Sprintf("Client connected: %v", conn.RemoteAddr()))

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Error(fmt.Sprintf("WS read error: %v", err))
			delete(clients, conn)
			break
		}
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				logger.Error(fmt.Sprintf("WS write error: %v", err))
				closeErr := client.Close()
				if closeErr != nil {
					logger.Error(fmt.Sprintf("Error closing WS connection: %v", err.Error()))
				}
				delete(clients, client)
			}
		}
	}
}

func SendNotification(message WebsocketMessage, logger *slog.Logger) {
	for client := range clients {
		err := client.WriteJSON(message)
		if err != nil {
			logger.Error(fmt.Sprintf("WS write error: %v", err))
			err := client.Close()
			if err != nil {
				return
			}
			delete(clients, client)
		}
	}
}
