package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"testing"
	"time"

	"hotsauceshop/lib"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
)

var config lib.HotSauceShopConfig

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	var configReadErr error
	config, configReadErr = lib.ReadConfig("../config.toml")
	if configReadErr != nil {
		panic("Could not read config")
	}
}

func TestGetBoardPosts(t *testing.T) {
	e := httpexpect.Default(t, "http://localhost:8081")
	e.GET("/api/v1/boards").
		Expect().
		Status(http.StatusOK).JSON().Object().
		Value("results").Object().
		Value("boards").Array().Length().Gt(0)
}

func createBoardAndVerify(t *testing.T, e *httpexpect.Expect, sessionID string) lib.AddBoardResponse {
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
	// checking linter
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
	return addBoardResponse
}

func deleteBoardAndVerify(t *testing.T, e *httpexpect.Expect, sessionID string, slug string) {
	var boardDeleteResponse lib.BoardDeleteResponse
	e.DELETE(fmt.Sprintf("/api/v1/boards/%v", slug)).
		WithCookie("sessionId", sessionID).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&boardDeleteResponse)
	if boardDeleteResponse.Status != "OK" {
		t.Fatal("Failed to delete board")
	}
}

func createBoardPostAndVerify(t *testing.T, e *httpexpect.Expect, sessionID string) string {
	postUUID, postUUIDErr := uuid.NewRandom()
	if postUUIDErr != nil {
		t.Fatal("Failed to generate post UUID")
	}
	postName := postUUID.String()

	// TODO: figure out how the heck to do images
	postFlairIds := []int{1}
	newPost := lib.AddPostRequest{
		Title:        postName,
		ParentId:     0,
		PostText:     "Follow the white rabbit, Neo.",
		Slug:         postName,
		PostFlairIds: postFlairIds,
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
	return postName
}

func deleteBoardPostAndVerify(t *testing.T, e *httpexpect.Expect, sessionID string, postSlug string) {
	// Delete post
	var boardPostDeleteResponse lib.BoardPostDeleteResponse
	e.DELETE(fmt.Sprintf("/api/v1/boards/posts/%v", postSlug)).
		WithCookie("sessionId", sessionID).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&boardPostDeleteResponse)
	if boardPostDeleteResponse.Status != "OK" {
		t.Fatal("Failed to delete board post")
	}
}

func TestCreateBoard(t *testing.T) {
	/**
	 * 1. Add a new board
	 * 2. Check board detail to confirm it exists
	 * 3. Verify response
	 * 4. Delete board
	 */
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	sessionID := signInAndGetSessionId(t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword)

	addBoardResponse := createBoardAndVerify(t, e, sessionID)

	log.Printf("Deleting board with slug %v", addBoardResponse.Results.Slug)

	// Clean up
	deleteBoardAndVerify(t, e, sessionID, addBoardResponse.Results.Slug)
}

func TestCreateBoardPost(t *testing.T) {
	/**
	 * 1. Add a board post with an image
	 * 2. Verify board post detail page
	 * 3. Delete post
	 */
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	sessionID := signInAndGetSessionId(t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword)
	postName := createBoardPostAndVerify(t, e, sessionID)
	deleteBoardPostAndVerify(t, e, sessionID, postName)
}

func TestCreateBoardPostWithoutSession(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	postUUID, postUUIDErr := uuid.NewRandom()
	if postUUIDErr != nil {
		t.Fatal("Failed to generate post UUID")
	}
	postName := postUUID.String()
	// TODO: figure out how the heck to do images
	postFlairIds := []int{1}
	newPost := lib.AddPostRequest{
		Title:        postName,
		ParentId:     0,
		PostText:     "Follow the white rabbit, Neo.",
		Slug:         postName,
		PostFlairIds: postFlairIds,
	}
	var addPostResponse lib.AddPostResponse
	// Probably should create the board here but...
	e.POST("/api/v1/boards/sauces/posts").
		WithJSON(newPost).
		Expect().
		Status(http.StatusUnauthorized).
		JSON().
		Decode(&addPostResponse)
	if addPostResponse.Status != "ERROR" {
		t.Fatal("Adding post without session should have failed")
	}
}

func TestGetPostFlairs(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	var postFlairsResponse lib.PostFlairsResponse
	e.GET("/api/v1/post-flairs").
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&postFlairsResponse)
	if postFlairsResponse.Status != "OK" {
		t.Fatal("Unexpected response for post flairs")
	}
}

func TestGetPostFlairsForPost(t *testing.T) {
	/**
	 * 1. Create a post and attach new flair id
	 * 2. Get flairs for the post and verify the new one is present
	 * NOTE: no plans at this time to allow creation of flair through the UI/API
	 */
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	sessionID := signInAndGetSessionId(t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword)
	postName := createBoardPostAndVerify(t, e, sessionID)

	var postDetail lib.BoardPostResponse
	e.GET(fmt.Sprintf("/api/v1/board/sauces/posts/%v", postName)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&postDetail)
	if postDetail.Status != "OK" {
		t.Fatal("Failed to get post detail")
	}

	deleteBoardPostAndVerify(t, e, sessionID, postName)
}

func TestGetPostsFlairsMap(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)

	var postFlairsResponse lib.PostFlairsResponse
	e.GET("/api/v1/post-flairs").
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&postFlairsResponse)
	if postFlairsResponse.Status != "OK" {
		t.Fatal("Failed to get post flairs")
	}
	postsFlairs := []lib.PostsFlairs{
		{
			Id:          1,
			BoardPostId: 1,
			PostFlairId: 2,
			CreatedAt:   time.Now(),
		},
		{
			Id:          2,
			BoardPostId: 2,
			PostFlairId: 3,
			CreatedAt:   time.Now(),
		},
	}
	postFlairIdMap := lib.GetPostFlairIdMap(postFlairsResponse.Results.PostFlairs)
	postsFlairsMap := lib.GetPostsFlairsMap(postsFlairs, postFlairIdMap)
	found := false
	for boardPostId, postFlairs := range postsFlairsMap {
		if len(postFlairs) == 0 {
			t.Fatal("postFlairs length is 0!")
		}
		if boardPostId == 1 {
			for _, postFlair := range postFlairs {
				found = slices.Contains(postsFlairsMap[boardPostId], postFlairIdMap[postFlair.Id])
				if found {
					break
				}
			}
		}
	}
	if !found {
		t.Fatal("Could not find post flair in map")
	}
}
