package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"hotsauceshop/lib"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
)

func signInAndGetSessionId(t *testing.T, e *httpexpect.Expect) string {
	testUsername := os.Getenv("TEST_USER_NAME")
	testPassword := os.Getenv("TEST_USER_PASSWORD")
	if testUsername == "" || testPassword == "" {
		t.Fatal("TEST_USER_NAME or TEST_USER_PASSWORD not set")
	}

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
	sessionID := signInResponse.Results.SessionId
	if len(sessionID) == 0 {
		t.Fatal("Failed to get user session id")
	}
	return sessionID
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
	sessionID := signInAndGetSessionId(t, e)

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

	log.Printf("Deleting board with slug %v", addBoardResponse.Results.Slug)

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

func TestCreateBoardPost(t *testing.T) {
	/**
	 * 1. Add a board post with an image
	 * 2. Verify board post detail page
	 * 3. Delete post
	 */
	e := httpexpect.Default(t, "http://localhost:8081")
	sessionID := signInAndGetSessionId(t, e)

	// Add board
	postUUID, postUUIDErr := uuid.NewRandom()
	if postUUIDErr != nil {
		t.Fatal("Failed to generate post UUID")
	}
	postName := postUUID.String()
	// TODO: figure out how the heck to do images
	newPost := lib.AddPostRequest{
		Title:    postName,
		ParentId: 0,
		PostText: "Follow the white rabbit, Neo.",
		Slug:     postName,
	}
	var addPostResponse lib.AddPostResponse
	// Probably should create the board here but...
	e.POST("/api/v1/boards/sauces/posts").
		WithCookie("sessionId", sessionID).
		WithJSON(newPost).
		Expect().
		Status(http.StatusCreated).
		JSON().
		Decode(&addPostResponse)
	if addPostResponse.Status != "OK" {
		t.Fatal("Failed to add post")
	}

	// Delete post
	var boardPostDeleteResponse lib.BoardPostDeleteResponse
	e.DELETE(fmt.Sprintf("/api/v1/boards/posts/%v", postName)).
		WithCookie("sessionId", sessionID).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&boardPostDeleteResponse)
	if boardPostDeleteResponse.Status != "OK" {
		t.Fatal("Failed to delete board post")
	}
}
