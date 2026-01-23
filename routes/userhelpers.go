package routes

import (
	"fmt"
	"log/slog"
	"net/http"
	"testing"

	"hotsauceshop/lib"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	T              *testing.T
	E              *httpexpect.Expect
	Username       string
	Password       string
	AvatarFilename string
	SessionId      string
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

func CreateUserAndVerify(request CreateUserRequest) {
	var userCreateResponse lib.UserCreateResponse
	request.E.POST("/api/v1/user").
		WithCookie("sessionId", request.SessionId).
		WithJSON(request).
		Expect().
		Status(http.StatusCreated).
		JSON().
		Decode(&userCreateResponse)
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
