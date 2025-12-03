package routes

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestBoardsGetPosts(t *testing.T) {
	e := httpexpect.New(t, "http://localhost:8081")
	e.GET("/api/v1/boards").
		Expect().
		Status(http.StatusOK).JSON().Object().ContainsKey("results").NotEmpty()
}
