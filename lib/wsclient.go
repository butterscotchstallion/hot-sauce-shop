package lib

import (
	"fmt"
	"log"
	"log/slog"

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

func SendWSMessage(c *gin.Context, message WebsocketMessage, logger *slog.Logger) {
	w, r := c.Writer, c.Request
	upgradedConnection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer func() { _ = upgradedConnection.Close() }()
	err = upgradedConnection.WriteJSON(message)
	if err != nil {
		logger.Error(fmt.Sprintf("WS write error: %v", err))
	}

	logger.Info(fmt.Sprintf("WS write: %v", message))
}
