package routes

import (
	"fmt"
	"net/http"
	"testing"

	"hotsauceshop/lib"

	"github.com/gavv/httpexpect/v2"
)

func signInAndGetSessionId(t *testing.T, e *httpexpect.Expect, username string, password string) string {
	var signInResponse lib.SignInResponse
	e.POST("/api/v1/user/sign-in").
		WithJSON(lib.LoginRequest{
			Username: username,
			Password: password,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&signInResponse)
	sessionID := signInResponse.Results.SessionId
	if len(sessionID) == 0 {
		t.Fatal("Failed to get user session id")
	}
	return sessionID
}

func TestUserJoinBoard(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	sessionID := signInAndGetSessionId(t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword)
	addBoardResponse := createBoardAndVerify(t, e, sessionID)

	var response GenericResponse
	e.POST(fmt.Sprintf("/api/v1/user/boards/%v", addBoardResponse.Results.BoardId)).
		WithCookie("sessionId", sessionID).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&response)
	if response.Status != "OK" {
		t.Fatal("Failed to join board")
	}

	// Only board admins can delete board
	deleteBoardAndVerify(t, e, sessionID, addBoardResponse.Results.Slug)
}
