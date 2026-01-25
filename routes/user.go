package routes

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"hotsauceshop/lib"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

func User(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger) {
	r.GET("/api/v1/user", func(c *gin.Context) {
		isUserAdmin, isUserAdminErr := lib.IsUserAdmin(c, dbPool, logger)
		if isUserAdminErr != nil {
			return
		}
		if !isUserAdmin {
			c.JSON(http.StatusForbidden, lib.GenericResponse{
				Status:  "ERROR",
				Message: "Permission denied",
			})
			return
		}

		users, err := lib.GetUsers(dbPool, logger)
		if err != nil {
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, lib.UserListResponse{
			Status: "OK",
			Results: lib.UserListResponseResults{
				Users: users,
			},
		})
	})

	// Get user profile by slug
	r.GET("/api/v1/user/profile/:userSlug", func(c *gin.Context) {
		userSlug := c.Param("userSlug")
		if len(userSlug) == 0 {
			c.JSON(http.StatusNotFound, lib.GenericResponse{
				Status:  "ERROR",
				Message: "User not found",
			})
			return
		}

		user, err := lib.GetUserBySlug(dbPool, logger, userSlug)
		if err != nil {
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: err.Error(),
			})
			return
		}
		if user == (lib.User{}) {
			c.JSON(http.StatusNotFound, lib.GenericResponse{
				Status:  "ERROR",
				Message: "User not found",
			})
			return
		}

		// TODO: maybe restrict this to certain roles?
		roles, rolesErr := lib.GetRolesByUserId(dbPool, logger, user.Id)
		if rolesErr != nil {
			logger.Error(fmt.Sprintf("Error fetching roles: %v", rolesErr.Error()))
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: rolesErr.Error(),
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

		c.JSON(http.StatusOK, lib.UserProfileResponse{
			Status: "OK",
			Results: lib.UserProfileResponseResults{
				User:                user,
				Roles:               roles,
				UserPostCount:       userPostCount,
				UserPostVoteSum:     userPostVoteSum,
				UserModeratedBoards: userModeratedBoards,
			},
		})
	})

	// Sign in
	r.POST("/api/v1/user/sign-in", func(c *gin.Context) {
		loginRequest := lib.LoginRequest{}
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, lib.GenericResponse{
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
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: errVerifying.Error(),
			})
			return
		}

		if verifiedUser == (lib.User{}) {
			c.JSON(http.StatusOK, lib.GenericResponse{
				Status:  "ERROR",
				Message: "Invalid username or password",
			})
			return
		}

		sessionId, err := lib.AddUserSessionId(dbPool, verifiedUser.Id)
		if err != nil || len(sessionId) == 0 {
			logger.Error(fmt.Sprintf("Error generating sessionId: %v", err.Error()))
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
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
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: boardsErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, lib.UserBoardsResponse{
			Status: "OK",
			Results: lib.UserBoardsResponseResults{
				Boards: boards,
			},
		})
	})

	// Join board
	r.POST("/api/v1/user/boards/:boardId", func(c *gin.Context) {
		boardIdSlug := c.Param("boardId")
		boardId, err := strconv.Atoi(boardIdSlug)
		if err != nil {
			logger.Error(fmt.Sprintf("Error parsing boardIdSlug: %v", err.Error()))
			c.JSON(http.StatusBadRequest, lib.GenericResponse{
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
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: boardsErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, lib.GenericResponse{
			Status:  "OK",
			Message: "Board added",
		})
	})

	// Create user
	r.POST("/api/v1/user", func(c *gin.Context) {
		// Check payload
		var payload lib.UserCreatePayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, lib.GenericResponse{
				Status:  "ERROR",
				Message: err.Error(),
			})
			return
		}

		// Check user role
		isUserAdmin, isUserAdminErr := lib.IsUserAdmin(c, dbPool, logger)
		if isUserAdminErr != nil {
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: isUserAdminErr.Error(),
			})
			return
		}
		if !isUserAdmin {
			logger.Error("Error creating user: user is not admin")
			c.JSON(http.StatusForbidden, lib.GenericResponseWithErrorCode{
				Status:    "ERROR",
				Message:   "Permission denied",
				ErrorCode: lib.ErrorCodePermissionDenied,
			})
			return
		}

		// Validate payload
		validate := validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(payload)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Validation failed: %v", err),
			})
			return
		}

		// Check if username exists
		userExists, userExistsErr := lib.UsernameExists(dbPool, payload.Username)
		if userExistsErr != nil {
			logger.Error(fmt.Sprintf("Error checking if username exists: %v", userExistsErr.Error()))
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: userExistsErr.Error(),
			})
			return
		}

		if userExists {
			logger.Error(fmt.Sprintf("User with username %v already exists", payload.Username))
			c.JSON(http.StatusBadRequest, lib.GenericResponseWithErrorCode{
				Status:    "ERROR",
				Message:   fmt.Sprintf("User with username %v already exists", payload.Username),
				ErrorCode: lib.ErrorCodeUserExists,
			})
			return
		}

		// Create user
		user, createUserErr := lib.CreateUser(dbPool, payload)
		if createUserErr != nil {
			logger.Error(fmt.Sprintf("Error creating user: %v", createUserErr.Error()))
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: createUserErr.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, lib.UserCreateResponse{
			Status:  "OK",
			Results: lib.UserCreateResponseResults{User: user},
		})
	})

	// Delete user
	r.DELETE("/api/v1/user/:slug", func(c *gin.Context) {
		slug := c.Param("slug")

		// Check user role
		isUserAdmin, isUserAdminErr := lib.IsUserAdmin(c, dbPool, logger)
		if isUserAdminErr != nil {
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: isUserAdminErr.Error(),
			})
			return
		}
		if !isUserAdmin {
			logger.Error("Error creating user: user is not admin")
			c.JSON(http.StatusForbidden, lib.GenericResponseWithErrorCode{
				Status:    "ERROR",
				Message:   "Permission denied",
				ErrorCode: lib.ErrorCodePermissionDenied,
			})
			return
		}

		user, err := lib.GetUserBySlug(dbPool, logger, slug)
		if err != nil || user == (lib.User{}) {
			logger.Error("Error fetching user with slug %v: %v", slug, err)
			c.JSON(http.StatusNotFound, lib.GenericResponseWithErrorCode{
				Status:    "ERROR",
				Message:   "User not found",
				ErrorCode: lib.ErrorCodeUserNotFound,
			})
			return
		}

		deleteUserErr := lib.DeleteUser(dbPool, user.Id)
		if deleteUserErr != nil {
			logger.Error(fmt.Sprintf("Error deleting user: %v", deleteUserErr.Error()))
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: deleteUserErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, lib.GenericResponse{
			Status:  "OK",
			Message: "User deleted",
		})
	})
}
