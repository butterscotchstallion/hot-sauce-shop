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
	addBoardResponse := CreateBoardAndVerify(t, e, sessionID, lib.AddBoardRequest{
		DisplayName:            GenerateUniqueName(),
		Description:            "Testing testing 1-2-3",
		IsVisible:              true,
		IsPrivate:              false,
		IsOfficial:             false,
		IsPostApprovalRequired: false,
	})

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
	DeleteBoardAndVerify(t, e, sessionID, addBoardResponse.Results.Slug)
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

func TestGetUserAdminBoards(t *testing.T) {
	// Get board admin session
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	sessionID := signInAndGetSessionId(t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword)

	// Create boardSlug
	addBoardResponse := CreateBoardAndVerify(t, e, sessionID, lib.AddBoardRequest{
		DisplayName:            GenerateUniqueName(),
		Description:            "Testing testing 1-2-3",
		IsVisible:              true,
		IsPrivate:              false,
		IsOfficial:             false,
		IsPostApprovalRequired: false,
	})

	// Add board admin
	var addBoardAdminResponse lib.GenericResponse
	e.POST(fmt.Sprintf("/api/v1/board-admin/%d", addBoardResponse.Results.BoardId)).
		WithCookie("sessionId", sessionID).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&addBoardAdminResponse)
	if addBoardAdminResponse.Status != "OK" {
		t.Fatal("Failed to add board admin")
	}

	// Get board admins and verify
	var userBoardsResponse lib.UserBoardsResponse
	e.GET("/api/v1/board-admin").
		WithCookie("sessionId", sessionID).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&userBoardsResponse)
	if userBoardsResponse.Status != "OK" {
		t.Fatal("Failed to get board admins")
	}

	boardFound := false
	for _, board := range userBoardsResponse.Results.Boards {
		if board.Id == addBoardResponse.Results.BoardId {
			boardFound = true
			break
		}
	}
	if !boardFound {
		t.Fatal("Board admin not found in board admins list")
	}
}

func TestGetUserProfile(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	sessionID := signInAndGetSessionId(t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword)
	GetUserProfileAndVerify(t, e, sessionID, config.TestUsers.BoardAdminUsername)
}

func TestGetKnownUserProfile(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	sessionID := signInAndGetSessionId(t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword)
	profile := GetUserProfileAndVerify(t, e, sessionID, config.TestUsers.BoardAdminUsername)
	if len(profile.Results.Roles) == 0 {
		t.Fatal("User roles are empty")
	}
	if profile.Results.UserPostCount == 0 {
		t.Fatal("User post count is 0")
	}
	if profile.Results.UserPostVoteSum == 0 {
		t.Fatal("User post vote sum is 0")
	}
	if len(profile.Results.UserModeratedBoards) == 0 {
		t.Fatal("User moderated boards is empty")
	}
}

/**
 * Create user tests
 * - Covers user creation
 * - Tests permissions
 * - Tests user creation with valid and invalid data
 * - Tests user sign in
 * - Tests user profile
 * - Tests user deletion
 */
func TestCreateUser(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	adminSessionId := signInAndGetSessionId(
		t,
		e,
		config.TestUsers.AdminUsername,
		config.TestUsers.AdminPassword,
	)
	unprivSessionId := signInAndGetSessionId(
		t,
		e,
		config.TestUsers.UnprivilegedUsername,
		config.TestUsers.UnprivilegedPassword,
	)
	// Only admins can create users - test permissions as well
	CreateRandomUserAndVerify(t, e, unprivSessionId, http.StatusForbidden, lib.ErrorCodePermissionDenied)

	// Admin attempt expected to succeed
	newUserInfo := CreateRandomUserAndVerify(t, e, adminSessionId, http.StatusCreated, "")

	// Test sign in
	newUserSessionId := signInAndGetSessionId(
		t,
		e,
		newUserInfo.Username,
		newUserInfo.Password,
	)

	// Test user profile
	GetUserProfileAndVerify(t, e, newUserSessionId, newUserInfo.Response.Results.User.Slug)

	// Test creating users that exist already
	CreateUserAndVerify(CreateUserRequest{
		T:                  t,
		E:                  e,
		Username:           newUserInfo.Username,
		Password:           newUserInfo.Password,
		AvatarFilename:     "",
		SessionId:          adminSessionId,
		ExpectedStatusCode: http.StatusBadRequest,
		ExpectedErrorCode:  lib.ErrorCodeUserExists,
	})

	// Clean up: delete user, test permissions

	// Unprivileged user should not be able to delete the user
	DeleteUserAndVerify(DeleteUserRequest{
		T:                  t,
		E:                  e,
		UserSlug:           newUserInfo.Response.Results.User.Slug,
		SessionId:          unprivSessionId,
		ExpectedStatusCode: http.StatusForbidden,
		ExpectedErrorCode:  lib.ErrorCodePermissionDenied,
	})
	
	// Admin should be able to delete the user
	DeleteUserAndVerify(DeleteUserRequest{
		T:                  t,
		E:                  e,
		UserSlug:           newUserInfo.Response.Results.User.Slug,
		SessionId:          adminSessionId,
		ExpectedStatusCode: http.StatusOK,
		ExpectedErrorCode:  "",
	})
}
