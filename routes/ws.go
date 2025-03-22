package routes

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WS(r *gin.Engine, logger *slog.Logger) {
	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer func(conn *websocket.Conn) {
			err := conn.Close()
			if err != nil {
				logger.Error(err.Error())
			}
		}(conn)
		for {
			err := conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
			if err != nil {
				return
			}
			time.Sleep(time.Second)
		}
	})
}
