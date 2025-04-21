package routes

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"hotsauceshop/lib"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

func User(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger) {
	// TODO: limit this to admins when RBAC system is complete
	r.GET("/api/v1/user", func(c *gin.Context) {
		users, err := lib.GetUsers(dbPool, logger)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"users": users,
			},
		})
	})

	r.POST("/api/v1/user/sign-in", func(c *gin.Context) {
		loginRequest := LoginRequest{}
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		verifiedUser, errVerifying := lib.VerifyUsernameAndPasswordAndReturnUser(
			dbPool, logger, loginRequest.Username, loginRequest.Password,
		)
		if errVerifying != nil {
			logger.Error(errVerifying.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": errVerifying.Error(),
			})
			return
		}

		if verifiedUser == (lib.User{}) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ERROR",
				"message": "Invalid username or password",
			})
			return
		}

		sessionId, err := lib.AddUserSessionId(dbPool, verifiedUser.Id)
		if err != nil || len(sessionId) == 0 {
			logger.Error(fmt.Sprintf("Error generating sessionId: %v", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "Sign in successful",
			"results": gin.H{
				"sessionId": sessionId,
				"user":      verifiedUser,
			},
		})
	})
}
