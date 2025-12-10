package lib

import (
	"log/slog"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSendReceiveChatMessage(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	sendErr := SendWebsocketMessage(WebsocketMessage{
		MessageType: "chat",
		Data: gin.H{
			"message": "hello world!",
		},
	}, logger)
	if sendErr != nil {
		t.Fatal(sendErr)
	}
}
