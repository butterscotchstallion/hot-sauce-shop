package routes

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"hotsauceshop/lib"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GenericResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UserBoardsResponseResults struct {
	Boards []lib.Board `json:"boards"`
}

type UserBoardsResponse struct {
	Status  string                    `json:"status"`
	Results UserBoardsResponseResults `json:"results"`
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
	r.GET("/api/v1/user", func(c *gin.Context) {
		isUserAdmin, isUserAdminErr := lib.IsUserAdmin(c, dbPool, logger)
		if isUserAdminErr != nil {
			return
		}
		if !isUserAdmin {
			c.JSON(http.StatusForbidden, GenericResponse{
				Status:  "ERROR",
				Message: "Permission denied",
			})
			return
		}

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

	// Get user profile by slug
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
		if user == (lib.User{}) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "ERROR",
				"message": "User not found",
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

		userModeratedBoards, userModeratedBoardsErr := lib.GetUserModeratedBoards(dbPool, user.Id)
		if userModeratedBoardsErr != nil {
			logger.Error(fmt.Sprintf("Error fetching moderated boards: %v", userModeratedBoardsErr.Error()))
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
				"user":                user,
				"roles":               roles,
				"userPostCount":       userPostCount,
				"userPostVoteSum":     userPostVoteSum,
				"userModeratedBoards": userModeratedBoards,
			},
		})
	})

	// Sign in
	r.POST("/api/v1/user/sign-in", func(c *gin.Context) {
		loginRequest := lib.LoginRequest{}
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, GenericResponse{
				Status:  "ERROR",
				Message: err.Error(),
			})
			return
		}

		verifiedUser, errVerifying := lib.VerifyUsernameAndPasswordAndReturnUser(
			dbPool, logger, loginRequest.Username, loginRequest.Password,
		)
		if errVerifying != nil {
			logger.Error(errVerifying.Error())
			c.JSON(http.StatusInternalServerError, GenericResponse{
				Status:  "ERROR",
				Message: errVerifying.Error(),
			})
			return
		}

		if verifiedUser == (lib.User{}) {
			c.JSON(http.StatusOK, GenericResponse{
				Status:  "ERROR",
				Message: "Invalid username or password",
			})
			return
		}

		sessionId, err := lib.AddUserSessionId(dbPool, verifiedUser.Id)
		if err != nil || len(sessionId) == 0 {
			logger.Error(fmt.Sprintf("Error generating sessionId: %v", err.Error()))
			c.JSON(http.StatusInternalServerError, GenericResponse{
				Status:  "ERROR",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, lib.SignInResponse{
			Status:  "OK",
			Message: "Sign in successful",
			Results: lib.SignInResponseResults{
				SessionId: sessionId,
				User:      verifiedUser,
			},
		})
	})

	// Get user boards
	r.GET("/api/v1/user/boards", func(c *gin.Context) {
		// Check user
		userId, userSessionErr := GetUserIdFromSessionOrError(c, dbPool, logger)
		if userSessionErr != nil || userId == 0 {
			return
		}

		boards, boardsErr := lib.GetJoinedBoardsByUserId(dbPool, userId)
		if boardsErr != nil {
			logger.Error(fmt.Sprintf("user/boards error: %v", boardsErr.Error()))
			c.JSON(http.StatusInternalServerError, GenericResponse{
				Status:  "ERROR",
				Message: boardsErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, UserBoardsResponse{
			Status: "OK",
			Results: UserBoardsResponseResults{
				Boards: boards,
			},
		})
	})

	// Join board
	r.POST("/api/v1/user/boards/:boardId", func(c *gin.Context) {
		boardIdSlug := c.Param("boardId")
		if len(boardIdSlug) == 0 {
			c.JSON(http.StatusNotFound, GenericResponse{
				Status:  "ERROR",
				Message: "Not found",
			})
			return
		}

		boardId, err := strconv.Atoi(boardIdSlug)
		if err != nil {
			logger.Error(fmt.Sprintf("Error parsing boardIdSlug: %v", err.Error()))
			c.JSON(http.StatusBadRequest, GenericResponse{
				Status:  "ERROR",
				Message: "Invalid boardId",
			})
			return
		}

		// Check user
		userId, userSessionErr := GetUserIdFromSessionOrError(c, dbPool, logger)
		if userSessionErr != nil || userId == 0 {
			return
		}

		boardsErr := lib.AddBoardUser(dbPool, userId, boardId)
		if boardsErr != nil {
			logger.Error(fmt.Sprintf("AddBoardUser error: %v", boardsErr.Error()))
			c.JSON(http.StatusInternalServerError, GenericResponse{
				Status:  "ERROR",
				Message: boardsErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, GenericResponse{
			Status:  "OK",
			Message: "",
		})
	})
}
