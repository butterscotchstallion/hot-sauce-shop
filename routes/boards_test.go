package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"testing"

	"hotsauceshop/lib"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var config lib.HotSauceShopConfig
var dbPool *pgxpool.Pool

// Used to add post flair to board post
var postFlairIds = []int{1, 2, 3}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	var configReadErr error
	config, configReadErr = lib.ReadConfig(lib.ConfigFilename)
	if configReadErr != nil {
		panic("Could not read config")
	}
	dbPool = lib.InitDB(config.Database.Dsn)
}

func TestGetBoardPosts(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	e.GET("/api/v1/boards").
		Expect().
		Status(http.StatusOK).JSON().Object().
		Value("results").Object().
		Value("boards").Array().Length().Gt(0)
}

func CreateBoardAndVerify(
	t *testing.T, e *httpexpect.Expect, sessionID string, newBoardPayload lib.AddBoardRequest,
) lib.AddBoardResponse {
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

func DeleteBoardAndVerify(t *testing.T, e *httpexpect.Expect, sessionID string, slug string) {
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

func createBoardPost(
	t *testing.T, e *httpexpect.Expect, newPost lib.AddPostRequest, sessionID string, boardSlug string,
	expectedStatusCode int,
) lib.AddPostResponse {
	var addPostResponse lib.AddPostResponse
	e.POST(fmt.Sprintf("/api/v1/boards/%v/posts", boardSlug)).
		WithCookie("sessionId", sessionID).
		WithJSON(newPost).
		Expect().
		Status(expectedStatusCode).
		JSON().
		Decode(&addPostResponse)
	expectedStatus := "OK"
	if expectedStatusCode != http.StatusCreated {
		expectedStatus = "ERROR"
	}
	if addPostResponse.Status != expectedStatus {
		t.Fatalf("Failed to add post: expected status mismatch: %s", addPostResponse.Status)
	}
	if expectedStatusCode == http.StatusCreated {
		if addPostResponse.Results.Post.Title != newPost.Title {
			t.Fatal("New post title mismatch")
		}
		if addPostResponse.Results.Post.PostText != newPost.PostText {
			t.Fatal("New post text mismatch")
		}
	}
	return addPostResponse
}

type VerifyPostDetailRequest struct {
	t              *testing.T
	e              *httpexpect.Expect
	sessionId      string
	expectedStatus int
	post           lib.BoardPost
	boardResponse  lib.AddBoardResponse
}

func verifyPostDetail(verifyPostDetailRequest VerifyPostDetailRequest) {
	var postDetailResponse lib.PostDetailResponse
	verifyPostDetailRequest.
		e.GET(
		fmt.Sprintf(
			"/api/v1/posts/%v/%v",
			verifyPostDetailRequest.boardResponse.Results.Slug,
			verifyPostDetailRequest.post.Slug)).
		Expect().
		Status(verifyPostDetailRequest.expectedStatus).
		JSON().
		Decode(&postDetailResponse)
	if verifyPostDetailRequest.expectedStatus == http.StatusOK {
		if postDetailResponse.Status != "OK" {
			verifyPostDetailRequest.t.Fatal("Failed to get post detail")
		}
		if postDetailResponse.Results.Post.Title != verifyPostDetailRequest.post.Title {
			verifyPostDetailRequest.t.Fatal("New post title mismatch")
		}
		if postDetailResponse.Results.Post.PostText != verifyPostDetailRequest.post.PostText {
			verifyPostDetailRequest.t.Fatal("New post text mismatch")
		}
	}
}

type CreateBoardAndPostAndVerifyRequest struct {
	T               *testing.T
	E               *httpexpect.Expect
	UnprivSessionId string
	AdminSessionId  string
	ExpectedStatus  int
	BoardPayload    lib.AddBoardRequest
}

type CreateBoardAndPostAndVerifyResponse struct {
	AddBoardResponse lib.AddBoardResponse
	AddPostResponse  lib.AddPostResponse
}

func createBoardAndPostAndVerify(request CreateBoardAndPostAndVerifyRequest) CreateBoardAndPostAndVerifyResponse {
	postUUID, postUUIDErr := uuid.NewRandom()
	if postUUIDErr != nil {
		request.T.Fatal("Failed to generate post UUID")
	}
	postName := postUUID.String()
	boardResponse := CreateBoardAndVerify(request.T, request.E, request.AdminSessionId, request.BoardPayload)

	// TODO: figure out how the heck to do images
	newPost := lib.AddPostRequest{
		Title:        postName,
		ParentSlug:   "",
		PostText:     "Follow the white rabbit, Neo.",
		PostFlairIds: postFlairIds,
	}
	addPostResponse := createBoardPost(
		request.T,
		request.E,
		newPost,
		request.UnprivSessionId,
		boardResponse.Results.Slug,
		http.StatusCreated,
	)

	// Verify with post detail
	verifyPostDetail(VerifyPostDetailRequest{
		t:              request.T,
		e:              request.E,
		sessionId:      request.UnprivSessionId,
		expectedStatus: http.StatusOK,
		post:           addPostResponse.Results.Post,
		boardResponse:  boardResponse,
	})

	return CreateBoardAndPostAndVerifyResponse{
		AddBoardResponse: boardResponse,
		AddPostResponse:  addPostResponse,
	}
}

func deleteBoardPostAndVerify(t *testing.T, e *httpexpect.Expect, sessionID string, postSlug string) {
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
	newBoardPayload := lib.AddBoardRequest{
		DisplayName:            GenerateUniqueName(),
		Description:            "Testing testing 1-2-3",
		IsVisible:              true,
		IsPrivate:              false,
		IsOfficial:             false,
		IsPostApprovalRequired: false,
		MinKarmaRequiredToPost: 0,
	}
	addBoardResponse := CreateBoardAndVerify(t, e, sessionID, newBoardPayload)

	log.Printf("Deleting board with slug %v", addBoardResponse.Results.Slug)

	// Clean up
	DeleteBoardAndVerify(t, e, sessionID, addBoardResponse.Results.Slug)
}

func TestCreateBoardPost(t *testing.T) {
	/**
	 * 1. Add a board post with an image
	 * 2. Verify board post detail page
	 * 3. Delete post
	 */
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	sessionID := signInAndGetSessionId(t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword)
	response := createBoardAndPostAndVerify(CreateBoardAndPostAndVerifyRequest{
		T:               t,
		E:               e,
		UnprivSessionId: sessionID,
		AdminSessionId:  sessionID,
		ExpectedStatus:  http.StatusCreated,
		BoardPayload: lib.AddBoardRequest{
			DisplayName:            GenerateUniqueName(),
			ThumbnailFilename:      "mr-brainly.jpg",
			Description:            "",
			IsPostApprovalRequired: false,
			IsPrivate:              false,
			IsOfficial:             false,
			IsVisible:              false,
			MinKarmaRequiredToPost: 0,
		},
	})
	if response.AddPostResponse.Results.Post.Slug == "" {
		t.Fatal("Failed to create board post: slug is blank")
	}
	deleteBoardPostAndVerify(t, e, sessionID, response.AddPostResponse.Results.Post.Slug)
}

func TestCreateBoardPostWithoutSession(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	postUUID, postUUIDErr := uuid.NewRandom()
	if postUUIDErr != nil {
		t.Fatal("Failed to generate post UUID")
	}
	postName := postUUID.String()
	// TODO: figure out how the heck to do images
	newPost := lib.AddPostRequest{
		Title:        postName,
		ParentSlug:   "",
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
	// Hard-coded to add post flair id #1
	response := createBoardAndPostAndVerify(CreateBoardAndPostAndVerifyRequest{
		T:               t,
		E:               e,
		UnprivSessionId: sessionID,
		AdminSessionId:  sessionID,
		ExpectedStatus:  http.StatusCreated,
		BoardPayload: lib.AddBoardRequest{
			DisplayName:            GenerateUniqueName(),
			ThumbnailFilename:      "mr-brainly.jpg",
			Description:            "",
			IsPostApprovalRequired: false,
			IsPrivate:              false,
			IsOfficial:             false,
			IsVisible:              false,
			MinKarmaRequiredToPost: 0,
		},
	})

	var postDetail lib.PostDetailResponse
	e.GET(fmt.Sprintf("/api/v1/posts/sauces/%v", response.AddPostResponse.Results.Post.Slug)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&postDetail)
	if postDetail.Status != "OK" {
		t.Fatal("Failed to get post detail")
	}

	// Verify that post flair id #1 is present
	if len(postDetail.Results.PostFlairs) == 0 {
		t.Fatal("Post flairs is empty!")
	}

	found := false
	for _, postFlair := range postDetail.Results.PostFlairs {
		if slices.Contains(postFlairIds, postFlair.Id) {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Post flair id #1 not found!")
	}

	deleteBoardPostAndVerify(t, e, sessionID, response.AddBoardResponse.Results.Slug)
}

func TestGetPostsFlairsMap(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)

	// All flairs
	var postFlairsResponse lib.PostFlairsResponse
	e.GET("/api/v1/post-flairs").
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&postFlairsResponse)
	if postFlairsResponse.Status != "OK" {
		t.Fatal("Failed to get post flairs")
	}

	// posts/flairs association
	var postsFlairsResponse lib.PostsFlairsResponse
	e.GET("/api/v1/posts-flairs").
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&postsFlairsResponse)

	if postsFlairsResponse.Status != "OK" {
		t.Fatal("Expected 200 OK for postsFlairsResponse")
	}

	if len(postsFlairsResponse.Results.PostsFlairs) == 0 {
		t.Fatal("Empty postsFlairs response from API")
	}

	postFlairIdMap := lib.GetPostFlairIdMap(postFlairsResponse.Results.PostFlairs)

	if len(postFlairIdMap) == 0 {
		t.Fatal("postsFlairsMap is empty!")
	}

	postsFlairsMap := lib.GetPostsFlairsMap(postsFlairsResponse.Results.PostsFlairs, postFlairIdMap)

	if len(postsFlairsMap) == 0 {
		t.Fatal("postsFlairsMap is empty!")
	}

	numFound := 0
	for boardPostId, postFlairs := range postsFlairsMap {
		if len(postFlairs) == 0 {
			t.Fatal("postFlairs length is 0!")
		}
		for _, postFlair := range postFlairs {
			found := slices.Contains(postsFlairsMap[boardPostId], postFlairIdMap[postFlair.Id])
			if found {
				numFound++
			}
		}
	}
	if numFound != len(postsFlairsMap) {
		t.Fatal("Could not find post flair in map")
	}
}

func TestGetPostFlairsQuery(t *testing.T) {
	postId := 42
	query := lib.GetPostFlairQuery(postId, []int{1, 2, 3})
	expectedQuery := fmt.Sprintf(
		`INSERT INTO posts_flairs (board_post_id, post_flair_id) VALUES (%v, 1),(%v, 2),(%v, 3)`,
		postId, postId, postId,
	)
	if query != expectedQuery {
		t.Fatalf("Expected query %v, got %v", expectedQuery, query)
	}
}

type UpdateBoardAndVerifyRequest struct {
	t                      *testing.T
	e                      *httpexpect.Expect
	sessionID              string
	newBoardResponse       lib.AddBoardResponse
	expectedResponseStatus string
	expectedStatus         int
	payload                lib.UpdateBoardRequest
}

func updateBoardAndVerify(updatedBoardAndVerifyRequest UpdateBoardAndVerifyRequest) {
	// Update boardSlug details
	var updateBoardDetailsResponse lib.GenericResponse
	updatedBoardAndVerifyRequest.e.
		PUT(fmt.Sprintf("/api/v1/boards/%v", updatedBoardAndVerifyRequest.newBoardResponse.Results.Slug)).
		WithCookie("sessionId", updatedBoardAndVerifyRequest.sessionID).
		WithJSON(updatedBoardAndVerifyRequest.payload).
		Expect().
		Status(updatedBoardAndVerifyRequest.expectedStatus).
		JSON().
		Decode(&updateBoardDetailsResponse)
	if updateBoardDetailsResponse.Status != updatedBoardAndVerifyRequest.expectedResponseStatus {
		updatedBoardAndVerifyRequest.t.Fatal("Failed to update board details")
	}

	// Verify board details
	var boardDetailResponse lib.BoardDetailResponse
	updatedBoardAndVerifyRequest.e.
		GET(fmt.Sprintf("/api/v1/boards/%v", updatedBoardAndVerifyRequest.newBoardResponse.Results.Slug)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&boardDetailResponse)
	if boardDetailResponse.Status != "OK" {
		updatedBoardAndVerifyRequest.t.Fatal("Failed to get board details")
	}

	// Verify updated board details - only if the expected status is OK
	if updatedBoardAndVerifyRequest.expectedResponseStatus == "OK" {
		if boardDetailResponse.Results.Board.IsPrivate != updatedBoardAndVerifyRequest.payload.IsPrivate {
			updatedBoardAndVerifyRequest.t.Fatal("Updated board is not private")
		}
		if boardDetailResponse.Results.Board.IsOfficial != updatedBoardAndVerifyRequest.payload.IsOfficial {
			updatedBoardAndVerifyRequest.t.Fatal("Updated board is not official")
		}
		if boardDetailResponse.Results.Board.IsPostApprovalRequired !=
			updatedBoardAndVerifyRequest.payload.IsPostApprovalRequired {
			updatedBoardAndVerifyRequest.t.Fatal("Updated board requires post approval")
		}
		if boardDetailResponse.Results.Board.Description != updatedBoardAndVerifyRequest.payload.Description {
			updatedBoardAndVerifyRequest.t.Fatal("Updated board description does not match")
		}
	}
}

/**
 * 1. Create a board
 * 2. Update board details
 * 3. Verify board details
 */
func TestUpdateBoardDetails(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	sessionID := signInAndGetSessionId(t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword)

	// Create a new board
	newBoardPayload := lib.AddBoardRequest{
		DisplayName:            GenerateUniqueName(),
		Description:            "Testing testing 1-2-3",
		IsVisible:              true,
		IsPrivate:              false,
		IsOfficial:             false,
		IsPostApprovalRequired: false,
		MinKarmaRequiredToPost: 0,
	}
	newBoardResponse := CreateBoardAndVerify(t, e, sessionID, newBoardPayload)
	if newBoardResponse.Status != "OK" {
		t.Fatal("Failed to create board")
	}

	updateBoardPayload := lib.UpdateBoardRequest{
		IsVisible:              true,
		IsPrivate:              true,
		IsOfficial:             true,
		IsPostApprovalRequired: true,
		Description:            "Let this corny slice of Americana be your tomb!",
		ThumbnailFilename:      "mr-brainly.jpg",
	}
	updateBoardAndVerify(UpdateBoardAndVerifyRequest{
		t:                      t,
		e:                      e,
		sessionID:              sessionID,
		newBoardResponse:       newBoardResponse,
		expectedResponseStatus: "OK",
		expectedStatus:         http.StatusOK,
		payload:                updateBoardPayload,
	})

	DeleteBoardAndVerify(t, e, sessionID, newBoardResponse.Results.Slug)
}

func TestUpdateBoardDetailsWithUnprivilegedUser(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	unprivSessionID := signInAndGetSessionId(
		t, e, config.TestUsers.UnprivilegedUsername, config.TestUsers.UnprivilegedPassword,
	)
	adminSessionID := signInAndGetSessionId(
		t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword,
	)

	// Create a new boardSlug
	newBoardPayload := lib.AddBoardRequest{
		DisplayName:            GenerateUniqueName(),
		Description:            "Testing testing 1-2-3",
		IsVisible:              true,
		IsPrivate:              false,
		IsOfficial:             false,
		IsPostApprovalRequired: false,
		MinKarmaRequiredToPost: 0,
	}
	newBoardResponse := CreateBoardAndVerify(t, e, adminSessionID, newBoardPayload)
	if newBoardResponse.Status != "OK" {
		t.Fatal("Failed to create boardSlug")
	}

	updateBoardAndVerify(UpdateBoardAndVerifyRequest{
		t:                      t,
		e:                      e,
		sessionID:              unprivSessionID,
		newBoardResponse:       newBoardResponse,
		expectedResponseStatus: "ERROR",
		expectedStatus:         http.StatusForbidden,
		payload: lib.UpdateBoardRequest{
			IsVisible:              true,
			IsPrivate:              true,
			IsOfficial:             true,
			IsPostApprovalRequired: true,
			Description:            "We don't talk about Bruno no no no",
			ThumbnailFilename:      "mr-brainly.jpg",
		},
	})

	DeleteBoardAndVerify(t, e, adminSessionID, newBoardResponse.Results.Slug)
}

func TestGetPostList(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	getPostListAndVerify(GetPostListAndVerifyParams{
		E:                        e,
		T:                        t,
		SessionId:                "",
		ShowUnapproved:           false,
		FilterByUserJoinedBoards: false,
	})
}

/**
 * 1. Add a new board
 * 2. Add a new post
 * 3. Verify the post is not visible to the unprivileged user when viewing the post detail
 */
func TestBoardRequiresPostApprovalWithUnprivilegedUser(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	unprivSessionID := signInAndGetSessionId(
		t, e, config.TestUsers.UnprivilegedUsername, config.TestUsers.UnprivilegedPassword,
	)
	adminSessionID := signInAndGetSessionId(
		t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword,
	)

	newBoardPayload := lib.AddBoardRequest{
		DisplayName: GenerateUniqueName(),
		Description: "IsPostApprovalRequired board test",
		IsVisible:   true,
		IsPrivate:   false,
		IsOfficial:  false,
		// Specifying this setting to make the post require approval
		IsPostApprovalRequired: true,
		MinKarmaRequiredToPost: 0,
	}
	newBoardResponse := CreateBoardAndVerify(t, e, adminSessionID, newBoardPayload)
	postName := GenerateUniqueName()
	newPost := lib.AddPostRequest{
		Title:        postName,
		ParentSlug:   "",
		PostText:     "Testing post approval with unprivileged user.",
		PostFlairIds: postFlairIds,
	}
	newPostResponse := createBoardPost(t, e, newPost, unprivSessionID, newBoardResponse.Results.Slug, http.StatusCreated)

	verifyPostDetail(VerifyPostDetailRequest{
		t:              t,
		e:              e,
		sessionId:      unprivSessionID,
		expectedStatus: http.StatusNotFound,
		post:           newPostResponse.Results.Post,
		boardResponse:  newBoardResponse,
	})
	DeleteBoardAndVerify(t, e, adminSessionID, newBoardResponse.Results.Slug)
	deleteBoardPostAndVerify(t, e, adminSessionID, newPostResponse.Results.Post.Slug)
}

/**
 * 1. Add a new board with isPostApprovalRequired set to true
 * 2. Add a new post with an unprivileged user
 * 3. Get the post list and verify that it is empty for unprivileged and privileged users
 */
func TestBoardPostListApprovedFilterWithPermissionTest(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	unprivSessionID := signInAndGetSessionId(
		t, e, config.TestUsers.UnprivilegedUsername, config.TestUsers.UnprivilegedPassword,
	)
	adminSessionID := signInAndGetSessionId(
		t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword,
	)

	// Need to create board as admin, but create post and verify post as unprivileged user
	boardResponse := CreateBoardAndVerify(t, e, adminSessionID, lib.AddBoardRequest{
		DisplayName:            GenerateUniqueName(),
		Description:            "Testing post list approved filter with unprivileged user.",
		IsVisible:              true,
		IsPrivate:              false,
		IsOfficial:             false,
		IsPostApprovalRequired: true,
		MinKarmaRequiredToPost: 0,
	})

	addPostResponse := createBoardPost(t, e, lib.AddPostRequest{
		Title:        GenerateUniqueName(),
		ParentSlug:   "",
		PostText:     "Meow meow meow meow meow",
		PostImages:   nil,
		Slug:         GenerateUniqueName(),
		PostFlairIds: nil,
	}, unprivSessionID, boardResponse.Results.Slug, http.StatusCreated)

	// To simplify things, unapproved posts will only be visible through the board moderation queue.
	// Expecting that the unapproved post is NOT available in the post detail.
	verifyPostDetail(VerifyPostDetailRequest{
		t:              t,
		e:              e,
		sessionId:      unprivSessionID,
		expectedStatus: http.StatusNotFound,
		post:           addPostResponse.Results.Post,
		boardResponse:  boardResponse,
	})

	// Verify the post list
	var isInList bool
	var postListResponse lib.PostListResponse

	// Check that the post is not visible to the unprivileged user
	postListResponse = getPostListAndVerify(GetPostListAndVerifyParams{
		E:                        e,
		T:                        t,
		SessionId:                unprivSessionID,
		ShowUnapproved:           true,
		BoardName:                boardResponse.Results.Slug,
		VerifyPostListIsNotEmpty: false,
		FilterByUserJoinedBoards: false,
	})
	newBoardSlug := boardResponse.Results.Slug
	newPostSlug := addPostResponse.Results.NewPostSlug

	isInList = isPostSlugInList(addPostResponse.Results.Post.Slug, postListResponse.Results.Posts)
	if isInList {
		t.Fatal("Post was found in the post list, but should not as an unprivileged user")
	}

	// Check that the post is visible to the admin user
	postListResponse = getPostListAndVerify(GetPostListAndVerifyParams{
		E:                        e,
		T:                        t,
		SessionId:                adminSessionID,
		ShowUnapproved:           true,
		BoardName:                boardResponse.Results.Slug,
		VerifyPostListIsNotEmpty: true,
		FilterByUserJoinedBoards: false,
	})

	isInList = isPostSlugInList(newPostSlug, postListResponse.Results.Posts)
	if !isInList {
		t.Fatal("Post was NOT found in the post list as an admin user")
	}

	// Clean up
	DeleteBoardAndVerify(t, e, adminSessionID, newBoardSlug)
	deleteBoardPostAndVerify(t, e, unprivSessionID, newPostSlug)
}

/**
 * 1. Create a board
 * 2. Join the board
 * 3. Create a post on the board
 * 3. Verify that the post list contains the newly created post while filtering on joined boards
 */
func TestGetPostsFilteredByUserJoinedBoards(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	unprivSessionID := signInAndGetSessionId(
		t, e, config.TestUsers.UnprivilegedUsername, config.TestUsers.UnprivilegedPassword,
	)
	adminSessionId := signInAndGetSessionId(
		t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword,
	)

	// Create a new board and post
	boardResponse := CreateBoardAndVerify(t, e, adminSessionId, lib.AddBoardRequest{
		DisplayName:            GenerateUniqueName(),
		Description:            "Testing get posts filtered by user joined boards",
		IsVisible:              true,
		IsPrivate:              false,
		IsOfficial:             false,
		IsPostApprovalRequired: false,
		MinKarmaRequiredToPost: 0,
	})
	postResponse := createBoardPost(t, e, lib.AddPostRequest{
		Title:        GenerateUniqueName(),
		ParentSlug:   "",
		PostText:     "hello",
		PostImages:   nil,
		Slug:         "",
		PostFlairIds: nil,
	}, unprivSessionID, boardResponse.Results.Slug, http.StatusCreated)

	// Join the new board
	joinBoardWithCurrentUser(boardResponse.Results.BoardId, e, t, unprivSessionID)

	// Verify the post list
	postListResponse := getPostListAndVerify(GetPostListAndVerifyParams{
		E:                        e,
		T:                        t,
		SessionId:                unprivSessionID,
		ShowUnapproved:           false,
		FilterByUserJoinedBoards: true,
		BoardName:                boardResponse.Results.Slug,
	})

	foundPost := false
	for _, post := range postListResponse.Results.Posts {
		if post.Slug == postResponse.Results.Post.Slug {
			foundPost = true
			break
		}
	}

	if !foundPost {
		t.Fatal("Post was not found in the post list")
	}

	deleteBoardPostAndVerify(t, e, adminSessionId, postResponse.Results.Post.Slug)
	DeleteBoardAndVerify(t, e, adminSessionId, boardResponse.Results.Slug)
}

/**
 * 1. Create a new board
 * 2. Add a new post for upvoting
 * 3. Set board minKarmaRequired to 1
 * 4. Attempt to post - expected response status code 403
 * 5. Up vote post from step 2
 * 6. Attempt to post - expected response status code 201
 */
func TestMinKarmaRequiredToPost(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	unprivSessionID := signInAndGetSessionId(
		t, e, config.TestUsers.UnprivilegedUsername, config.TestUsers.UnprivilegedPassword,
	)
	adminSessionId := signInAndGetSessionId(
		t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword,
	)

	// 0. Clear posts/votes so the user karma is 0
	const UnprivUserId = 3
	deletePostFlairsErr := lib.DeletePostFlairsByUserId(dbPool, UnprivUserId)
	if deletePostFlairsErr != nil {
		t.Fatalf("Error deleting user post flairs: %v", deletePostFlairsErr)
	}
	votesDeleted, deleteVotesErr := lib.DeleteVotesByPostId(dbPool, UnprivUserId)
	if deleteVotesErr != nil {
		t.Fatalf("Error deleting user votes: %v", deleteVotesErr)
	}
	log.Printf("Deleted %d votes", votesDeleted)
	deletePostsErr := lib.DeletePostsByUserId(dbPool, UnprivUserId)
	if deletePostsErr != nil {
		t.Fatalf("Error deleting user posts: %v", deletePostsErr)
	}

	// 1 + 2. Create a new board and post
	boardAndPostResponse := createBoardAndPostAndVerify(CreateBoardAndPostAndVerifyRequest{
		T:               t,
		E:               e,
		UnprivSessionId: unprivSessionID,
		AdminSessionId:  adminSessionId,
		ExpectedStatus:  http.StatusOK,
		BoardPayload: lib.AddBoardRequest{
			DisplayName:            GenerateUniqueName(),
			ThumbnailFilename:      "",
			Description:            "meow",
			IsPostApprovalRequired: false,
			IsPrivate:              false,
			IsOfficial:             false,
			IsVisible:              false,
			MinKarmaRequiredToPost: 0,
		},
	})

	// 3. Update board minKarmaRequired to 1
	updateBoardAndVerify(UpdateBoardAndVerifyRequest{
		t:                      t,
		e:                      e,
		sessionID:              adminSessionId,
		newBoardResponse:       boardAndPostResponse.AddBoardResponse,
		expectedResponseStatus: "OK",
		expectedStatus:         http.StatusOK,
		payload: lib.UpdateBoardRequest{
			IsVisible:              true,
			IsPrivate:              false,
			IsOfficial:             false,
			IsPostApprovalRequired: false,
			Description:            "meow",
			ThumbnailFilename:      "",
			MinKarmaRequiredToPost: 1,
		},
	})

	// 4. Attempt to post - expected response status code 403
	addPostResponse := createBoardPost(t, e, lib.AddPostRequest{
		Title:        GenerateUniqueName(),
		ParentSlug:   "",
		PostText:     "403 response",
		PostImages:   nil,
		Slug:         "",
		PostFlairIds: nil,
	}, unprivSessionID, boardAndPostResponse.AddBoardResponse.Results.Slug, http.StatusForbidden)
	if addPostResponse.ErrorCode != lib.ErrorCodeInsufficientKarma {
		t.Fatalf("Expected error code to be %s", lib.ErrorCodeInsufficientKarma)
	}

	// 5. upvote post
	var voteResponse lib.VoteResponse
	voteRequestPayload := lib.AddUpdateVoteRequest{
		VoteValue: 1,
	}
	e.POST(fmt.Sprintf("/api/v1/votes/%d", boardAndPostResponse.AddPostResponse.Results.NewPostId)).
		WithCookie("sessionId", unprivSessionID).
		WithJSON(voteRequestPayload).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&voteResponse)
	if voteResponse.Status != "OK" {
		t.Fatal("Failed to upvote post")
	}

	// 6. Attempt to post - expected response status code 201
	createBoardPost(t, e, lib.AddPostRequest{
		Title:        GenerateUniqueName(),
		ParentSlug:   "",
		PostText:     "201 response",
		PostImages:   nil,
		Slug:         "",
		PostFlairIds: nil,
	}, unprivSessionID, boardAndPostResponse.AddBoardResponse.Results.Slug, http.StatusCreated)

	// Delete votes before cleaning up - FK constraint will prevent deletion
	var deletedVotes int64
	deletedVotes, deleteVotesErr = lib.DeleteVotesByPostId(dbPool, UnprivUserId)
	if deleteVotesErr != nil {
		t.Fatalf("Error deleting user votes: %v", deleteVotesErr)
	}
	log.Printf("Deleted %d votes", deletedVotes)

	deleteBoardPostAndVerify(t, e, unprivSessionID, boardAndPostResponse.AddPostResponse.Results.Post.Slug)
	DeleteBoardAndVerify(t, e, adminSessionId, boardAndPostResponse.AddBoardResponse.Results.Slug)
}

/**
 * - Test that board list endpoint returns a list of boards
 * based on the new filter: omitEmpty
 * ------------------------------------------------------------------
 * - Scenario 1: default filter setting is omitEmpty=0, and we should
 * create a board with a post to verify that the board is
 * included in the list
 * -------------------------------------------------------------------
 * - Scenario 2: test omitEmpty=1 by deleting the post and checking
 * that the board is not included in the list
 */
func TestBoardListHasPostsFilter(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	unprivSessionId := signInAndGetSessionId(
		t, e, config.TestUsers.UnprivilegedUsername, config.TestUsers.UnprivilegedPassword,
	)
	adminSessionID := signInAndGetSessionId(
		t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword,
	)

	// Scenario 1: test default omitEmpty=0 setting

	// Create a board
	boardResponse := CreateBoardAndVerify(t, e, adminSessionID, lib.AddBoardRequest{
		DisplayName:            GenerateUniqueName(),
		Description:            "Testing board list has posts filter",
		IsVisible:              true,
		IsPrivate:              false,
		IsOfficial:             false,
		IsPostApprovalRequired: false,
		MinKarmaRequiredToPost: 0,
	})
	// Add a post to that board
	postResponse := createBoardPost(t, e, lib.AddPostRequest{
		Title:        GenerateUniqueName(),
		ParentSlug:   "",
		PostText:     "We don't talk about Bruno no no no",
		PostImages:   nil,
		Slug:         "",
		PostFlairIds: nil,
	}, unprivSessionId, boardResponse.Results.Slug, http.StatusCreated)

	// Verify the board list returns board
	boardListResponse := getBoardList(GetBoardListRequest{
		T:               t,
		E:               e,
		OmitEmptyBoards: false,
	})
	foundBoard := isBoardInList(boardResponse.Results.Slug, boardListResponse.Results.Boards)
	if !foundBoard {
		t.Fatal("Board was not found in the board list")
	}

	// Scenario 2: test omitEmpty=1 setting
	deleteBoardPostAndVerify(t, e, unprivSessionId, postResponse.Results.Post.Slug)

	// Get the board list with omitEmpty set to true
	boardListResponse = getBoardList(GetBoardListRequest{
		T:               t,
		E:               e,
		OmitEmptyBoards: true,
	})
	foundBoard = isBoardInList(boardResponse.Results.Slug, boardListResponse.Results.Boards)
	// Board is now empty - it should not be in the list
	if foundBoard {
		t.Fatal("Board was found in the board list")
	}

	// Clean up
	DeleteBoardAndVerify(t, e, adminSessionID, boardResponse.Results.Slug)
}

func TestGetBoardUsersByRoleName(t *testing.T) {
	const unprivUserId = 3
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	adminSessionID := signInAndGetSessionId(
		t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword,
	)
	addBoardResponse := CreateBoardAndVerify(t, e, adminSessionID, lib.AddBoardRequest{
		DisplayName:            GenerateUniqueName(),
		Description:            "Testing get board users by role name",
		IsVisible:              true,
		IsPrivate:              false,
		IsOfficial:             false,
		IsPostApprovalRequired: false,
		MinKarmaRequiredToPost: 0,
	})
	addAdminErr := lib.AddBoardAdmin(dbPool, unprivUserId, addBoardResponse.Results.BoardId)
	if addAdminErr != nil {
		t.Fatalf("Error adding board admin role: %v", addAdminErr)
	}
	addModErr := lib.AddBoardModerator(dbPool, unprivUserId, addBoardResponse.Results.BoardId)
	if addModErr != nil {
		t.Fatalf("Error adding board moderator role: %v", addModErr)
	}

	// Verify board details
	detailsResponse := getBoardDetailsAndVerify(GetBoardDetailsRequest{
		T:    t,
		E:    e,
		Slug: addBoardResponse.Results.Slug,
	})

	if len(detailsResponse.Results.Admins) != 1 {
		t.Fatalf("Expected 1 admin, got %d", len(detailsResponse.Results.Admins))
	}

	if len(detailsResponse.Results.Moderators) != 1 {
		t.Fatalf("Expected 1 moderator, got %d", len(detailsResponse.Results.Moderators))
	}

	adminInList := isUserIdInUserList(unprivUserId, detailsResponse.Results.Admins)
	if !adminInList {
		t.Fatal("Admin user was not found in the board details")
	}
	modInList := isUserIdInUserList(unprivUserId, detailsResponse.Results.Moderators)
	if !modInList {
		t.Fatal("Moderator user was not found in the board details")
	}

	// Clean up
	DeleteBoardAndVerify(t, e, adminSessionID, addBoardResponse.Results.Slug)
}

func TestCreateBoardMakesCreatorAdmin(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	adminSessionID := signInAndGetSessionId(
		t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword,
	)
	addBoardResponse := CreateBoardAndVerify(t, e, adminSessionID, lib.AddBoardRequest{
		DisplayName:            GenerateUniqueName(),
		Description:            "Testing create board makes creator admin",
		IsVisible:              true,
		IsPrivate:              false,
		IsOfficial:             false,
		IsPostApprovalRequired: false,
		MinKarmaRequiredToPost: 0,
	})

	// Verify board details
	detailsResponse := getBoardDetailsAndVerify(GetBoardDetailsRequest{
		T:    t,
		E:    e,
		Slug: addBoardResponse.Results.Slug,
	})

	if len(detailsResponse.Results.Admins) != 1 {
		t.Fatalf("Expected 1 admin, got %d", len(detailsResponse.Results.Admins))
	}
	adminInList := isUserIdInUserList(unprivUserId, detailsResponse.Results.Admins)
	if !adminInList {
		t.Fatal("Admin user was not found in the board details")
	}
}
