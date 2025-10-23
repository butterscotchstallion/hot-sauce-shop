package routes

import (
	"fmt"
	"log/slog"
	"net/http"

	"hotsauceshop/lib"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
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

	r.GET("/api/v1/user/profile/:userSlug", func(c *gin.Context) {
		userSlug := c.Param("userSlug")
		if len(userSlug) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "ERROR",
				"message": "User not found",
			})
			return
		}

		user, err := lib.GetUserBySlug(dbPool, logger, userSlug)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		// TODO: maybe restrict this to certain roles?
		roles, rolesErr := lib.GetRolesByUserId(dbPool, logger, user.Id)
		if rolesErr != nil {
			logger.Error(fmt.Sprintf("Error fetching roles: %v", rolesErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": rolesErr.Error(),
			})
			return
		}

		userPostCount, userPostCountErr := lib.GetNumPostsByUserId(dbPool, user.Id)
		if userPostCountErr != nil {
			logger.Error(fmt.Sprintf("Error fetching user post count: %v", userPostCountErr.Error()))
		}

		userPostVoteSum, userPostVoteSumErr := lib.GetUserPostVoteSum(dbPool, user.Id)
		if userPostVoteSumErr != nil {
			logger.Error(fmt.Sprintf("Error fetching user post vote sum: %v", userPostVoteSumErr.Error()))
		}

		logger.Info(fmt.Sprintf("User post vote sum: %v", userPostVoteSum))

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"user":            user,
				"roles":           roles,
				"userPostCount":   userPostCount,
				"userPostVoteSum": userPostVoteSum,
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
