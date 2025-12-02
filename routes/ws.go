package routes

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
var clients = make(map[*websocket.Conn]bool)

func WS(r *gin.Engine, wsConn *websocket.Conn, logger *slog.Logger) {
	r.GET("/ws", func(c *gin.Context) {
		var err error
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return c.Request.Header.Get("Origin") == "http://localhost:5173"
		}
		wsConn, err = upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		defer func(wsConn *websocket.Conn) {
			err := wsConn.Close()
			if err != nil {
				logger.Error(fmt.Sprintf("Error closing WS connection: %v", err.Error()))
			}
		}(wsConn)
		clients[wsConn] = true
		logger.Info(fmt.Sprintf("Client connected: %v", wsConn.RemoteAddr()))

		for {
			_, msg, err := wsConn.ReadMessage()
			if err != nil {
				logger.Error(fmt.Sprintf("WS read error: %v", err))
				delete(clients, wsConn)
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
	})
}
