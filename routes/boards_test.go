package routes

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"hotsauceshop/lib"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
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
	boardUUID, boardUUIDErr := uuid.NewRandom()
	if boardUUIDErr != nil {
		t.Fatal("Failed to generate board UUID")
	}
	boardName := boardUUID.String()
	newBoardPayload := lib.AddBoardRequest{
		DisplayName: boardName,
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
	if addBoardResponse.Results.DisplayName != newBoardPayload.DisplayName {
		t.Fatal("New board display name mismatch")
	}

	// Verify board exists now
	var boardDetailResponse lib.BoardDetailResponse
	e.GET(fmt.Sprintf("/api/v1/boards/%v", addBoardResponse.Results.Slug)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&boardDetailResponse)
	if boardDetailResponse.Status != "OK" {
		t.Fatal("Failed to get board details of newly created board")
	}

	// Clean up
	var boardDeleteResponse lib.BoardDeleteResponse
	e.DELETE(fmt.Sprintf("/api/v1/boards/%v", addBoardResponse.Results.Slug)).
		WithCookie("sessionId", sessionID).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&boardDeleteResponse)
	if boardDetailResponse.Status != "OK" {
		t.Fatal("Failed to delete board")
	}
}
