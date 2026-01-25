package routes

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"testing"

	"hotsauceshop/lib"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	T                  *testing.T
	E                  *httpexpect.Expect
	Username           string
	Password           string
	AvatarFilename     string
	SessionId          string
	ExpectedStatusCode int
	ExpectedErrorCode  string
}

func GetUserIdFromSessionOrError(c *gin.Context, dbPool *pgxpool.Pool, logger *slog.Logger) (int, error) {
	userId, err := lib.GetUserIdFromSession(c, dbPool, logger)
	if err != nil || userId == 0 {
		if err != nil {
			logger.Error(fmt.Sprintf("GetUserIdFromSessionOrError: %v", err.Error()))
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "ERROR",
				"message": "User not signed in",
			})
			return 0, err
		}
		return 0, nil
	}
	return userId, nil
}

func CreateUserAndVerify(request CreateUserRequest) lib.UserCreateResponse {
	var userCreateResponse lib.UserCreateResponse
	request.E.POST("/api/v1/user").
		WithCookie("sessionId", request.SessionId).
		WithJSON(request).
		Expect().
		Status(request.ExpectedStatusCode).
		JSON().
		Decode(&userCreateResponse)
	// Only verify this stuff if we're expecting it to work
	if request.ExpectedStatusCode == http.StatusCreated {
		if userCreateResponse.Status != "OK" {
			request.T.Fatal("Failed to create user")
		}
		if userCreateResponse.Results.User.Username != request.Username {
			request.T.Fatal("Created user username mismatch")
		}
		if userCreateResponse.Results.User.AvatarFilename != request.AvatarFilename {
			request.T.Fatal("Created user avatar filename mismatch")
		}
	}
	if len(request.ExpectedErrorCode) > 0 {
		if userCreateResponse.ErrorCode != request.ExpectedErrorCode {
			request.T.Fatalf(
				"Expected error code %s, got %s",
				request.ExpectedErrorCode,
				userCreateResponse.ErrorCode,
			)
		}
	}
	return userCreateResponse
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateUsername(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type CreateRandomUserResponse struct {
	Username string
	Password string
	Response lib.UserCreateResponse
}

func CreateRandomUserAndVerify(
	t *testing.T, e *httpexpect.Expect, sessionId string, expectedStatusCode int,
	expectedErrorCode string,
) CreateRandomUserResponse {
	const UsernameLength = 20
	password := GenerateUniqueName()
	hashedPw, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	username := GenerateUsername(UsernameLength)
	response := CreateUserAndVerify(CreateUserRequest{
		T:                  t,
		E:                  e,
		SessionId:          sessionId,
		Username:           username,
		Password:           hashedPw,
		ExpectedStatusCode: expectedStatusCode,
		ExpectedErrorCode:  expectedErrorCode,
	})
	return CreateRandomUserResponse{
		Username: username,
		Password: password,
		Response: response,
	}
}

func GetUserProfileAndVerify(
	t *testing.T, e *httpexpect.Expect, sessionId string, userSlug string,
) lib.UserProfileResponse {
	var userProfileResponse lib.UserProfileResponse
	e.GET(fmt.Sprintf("/api/v1/user/profile/%s", userSlug)).
		WithCookie("sessionId", sessionId).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&userProfileResponse)
	if userProfileResponse.Status != "OK" {
		t.Fatal("Failed to get user profile")
	}
	if userProfileResponse.Results.User == (lib.User{}) {
		t.Fatal("User profile is nil")
	}
	// ensure user moderated boards doesn't contain duplicates
	userModeratedBoards := make(map[string]bool)
	for _, board := range userProfileResponse.Results.UserModeratedBoards {
		userModeratedBoards[board.Slug] = true
	}
	if len(userModeratedBoards) != len(userProfileResponse.Results.UserModeratedBoards) {
		t.Fatal("User moderated boards contains duplicates")
	}
	return userProfileResponse
}

func DeleteUserAndVerify(
	t *testing.T, e *httpexpect.Expect, expectedStatusCode int, sessionId string, userSlug string,
) {
	var deleteUserResponse lib.GenericResponse
	e.DELETE(fmt.Sprintf("/api/v1/user/%s", userSlug)).
		WithCookie("sessionId", sessionId).
		Expect().
		Status(expectedStatusCode).
		JSON().
		Decode(&deleteUserResponse)
	if expectedStatusCode != http.StatusOK {
		return
	}
	if deleteUserResponse.Status != "OK" {
		t.Fatal("Failed to delete user")
	}
}
