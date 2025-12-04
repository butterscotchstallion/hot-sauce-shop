package routes

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"hotsauceshop/lib"

	"github.com/gavv/httpexpect/v2"
)

func signInAndGetSessionId(testUsername string, testPassword string, e *httpexpect.Expect) string {
	var signInResponse lib.SignInResponse
	e.POST("/api/v1/user/sign-in").
		WithJSON(lib.LoginRequest{
			Username: testUsername,
			Password: testPassword,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&signInResponse)
	return signInResponse.Results.SessionId
}

func TestGetBoardPosts(t *testing.T) {
	e := httpexpect.Default(t, "http://localhost:8081")
	e.GET("/api/v1/boards").
		Expect().
		Status(http.StatusOK).JSON().Object().
		Value("results").Object().
		Value("boards").Array().Length().Gt(0)
}

func TestCreateBoard(t *testing.T) {
	/**
	 * 1. Add a new board
	 * 2. Check board detail to confirm it exists
	 * 3. Verify response
	 * 4. Delete board
	 */
	e := httpexpect.Default(t, "http://localhost:8081")

	// Log in and get a session cookie
	testUsername := os.Getenv("TEST_USER_NAME")
	testPassword := os.Getenv("TEST_USER_PASSWORD")
	if testUsername == "" || testPassword == "" {
		t.Fatal("TEST_USER_NAME or TEST_USER_PASSWORD not set")
	}
	sessionID := signInAndGetSessionId(testUsername, testPassword, e)
	if len(sessionID) == 0 {
		t.Fatal("Failed to get user session id")
	}

	// Add board
	newBoardPayload := lib.AddBoardRequest{
		DisplayName: "Test Board Name!",
		Description: "Testing testing 1-2-3",
	}
	var addBoardResponse lib.AddBoardResponse
	e.POST("/api/v1/boards").
		WithCookie("sessionId", sessionID).
		WithJSON(newBoardPayload).
		Expect().
		Status(http.StatusCreated).
		JSON().
		Decode(&addBoardResponse)
	if addBoardResponse.Status != "OK" {
		t.Fatal("Error adding board")
	}

	// Verify board exists now
	e.GET(fmt.Sprintf("/api/v1/boards/%v", addBoardResponse.Results.Slug)).
		Expect().
		Status(http.StatusOK).JSON().Object().
		Value("results").Object().
		Value("board").NotNull()
}
