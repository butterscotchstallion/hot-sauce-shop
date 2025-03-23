package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WS(r *gin.Engine, wsConn *websocket.Conn, logger *slog.Logger) {
	r.GET("/ws", func(c *gin.Context) {
		var err error
		upgrader.CheckOrigin = func(r *http.Request) bool {
			if c.Request.Header.Get("Origin") != "http://localhost:5173" {
				return false
			}
			return true
		}
		wsConn, err = upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		err = wsConn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
		if err != nil {
			logger.Error(err.Error())
			return
		}
		/*
			for {
				err := wsConn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
				if err != nil {
					logger.Error(err.Error())
					return
				}
				time.Sleep(time.Second)
			}*/
	})
}
