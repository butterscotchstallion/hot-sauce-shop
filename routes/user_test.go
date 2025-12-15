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
	addBoardResponse := CreateBoardAndVerify(t, e, sessionID)

	var response lib.GenericResponse
	e.POST(fmt.Sprintf("/api/v1/user/boards/%v", addBoardResponse.Results.BoardId)).
		WithCookie("sessionId", sessionID).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&response)
	if response.Status != "OK" {
		t.Fatal("Failed to join board")
	}

	// Only board admins can delete boards
	deleteBoardAndVerify(t, e, sessionID, addBoardResponse.Results.Slug)
}

func TestGetUserListWithUnprivilegedUser(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	sessionID := signInAndGetSessionId(t, e, config.TestUsers.UnprivilegedUsername, config.TestUsers.UnprivilegedPassword)
	var userListResponse lib.UserListResponse
	e.GET("/api/v1/user").
		WithCookie("sessionId", sessionID).
		Expect().
		Status(http.StatusForbidden).
		JSON().
		Decode(&userListResponse)
	if userListResponse.Status != "ERROR" {
		t.Fatal("User list response status should have been ERROR")
	}
}
