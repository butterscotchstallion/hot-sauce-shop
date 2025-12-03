package routes

import (
	"log/slog"

	"hotsauceshop/lib"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func WS(r *gin.Engine, wsConn *websocket.Conn, logger *slog.Logger) {
	r.GET("/ws", func(c *gin.Context) {
		lib.HandleWSConnection(c, logger, wsConn)
	})
}
