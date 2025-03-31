package routes

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"hotsauceshop/lib"
)

type AdminUpdateUserRequest struct {
	User  lib.User   `json:"user"`
	Roles []lib.Role `json:"roles"`
}

func IsSignedInAndUserExists(c *gin.Context, dbPool *pgxpool.Pool, logger *slog.Logger) (int, error) {
	sessionIdCookieValue, err := c.Cookie("sessionId")
	if err != nil || sessionIdCookieValue == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ERROR",
			"message": "No session ID found",
		})
		return 0, err
	}

	user, getUserErr := lib.GetUserBySessionId(dbPool, logger, sessionIdCookieValue)
	if getUserErr != nil || user == (lib.User{}) {
		logger.Error("Error fetching user: %v", getUserErr)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "ERROR",
			"message": "No user found for session ID",
		})
		return 0, getUserErr
	}

	return user.Id, nil
}

func IsUserAdmin(c *gin.Context, dbPool *pgxpool.Pool, logger *slog.Logger) (bool, error) {
	userId, userErr := IsSignedInAndUserExists(c, dbPool, logger)
	if userErr != nil {
		return false, userErr
	}

	roles, rolesErr := lib.GetRolesByUserId(dbPool, logger, userId)
	if rolesErr != nil {
		logger.Error("Error fetching roles: %v", rolesErr.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "ERROR",
			"message": rolesErr.Error(),
		})
	}

	for _, role := range roles {
		if role.Name == "User Admin" {
			return true, nil
		}
	}

	return false, nil
}

func getRoleIdsFromRoles(roles []lib.Role) []int {
	var roleIds []int
	for _, role := range roles {
		roleIds = append(roleIds, role.Id)
	}
	return roleIds
}

func Admin(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger, store *persistence.InMemoryStore) {
	r.PUT("/api/v1/admin/user/:slug", func(c *gin.Context) {
		userSlug := c.Param("slug")

		if userSlug == "" {
			logger.Error("User slug is required")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": "User slug is required",
			})
			return
		}

		var adminUserUpdateRequest AdminUpdateUserRequest
		if err := c.ShouldBindJSON(&adminUserUpdateRequest); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Update request malformed: %v", err.Error()),
			})
			return
		}

		isUserAdmin, isUserAdminErr := IsUserAdmin(c, dbPool, logger)
		if isUserAdminErr != nil {
			logger.Error("Error checking if user is admin: %v", isUserAdminErr)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": isUserAdminErr.Error(),
			})
			return
		}

		if !isUserAdmin {
			logger.Error("User is not an admin")
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "ERROR",
				"message": "User is not an admin",
			})
			return
		}

		// Update user info
		// Update user roles
		_, updateErr := lib.UpdateUserRoles(
			dbPool,
			logger,
			adminUserUpdateRequest.User.Id,
			getRoleIdsFromRoles(adminUserUpdateRequest.Roles),
		)
		if updateErr != nil {
			logger.Error("Error updating user roles: %v", updateErr)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": updateErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	r.GET("/api/v1/admin/roles", cache.CachePage(store, time.Minute*15, func(c *gin.Context) {
		isUserAdmin, isUserAdminErr := IsUserAdmin(c, dbPool, logger)
		if isUserAdminErr != nil {
			logger.Error("Error checking if user is admin: %v", isUserAdminErr)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": isUserAdminErr.Error(),
			})
			return
		}

		if !isUserAdmin {
			logger.Error("User is not an admin")
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "ERROR",
				"message": "User is not an admin",
			})
			return
		}

		roles, roleErr := lib.GetRoleList(dbPool, logger)
		if roleErr != nil {
			logger.Error("Error fetching roles: %v", roleErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": roleErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"roles":  roles,
		})
	}))

	r.GET("/api/v1/admin/user/:slug", func(c *gin.Context) {
		userSlug := c.Param("slug")
		if userSlug == "" {
			logger.Error("User slug is required")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": "User slug is required",
			})
			return
		}

		// TODO: refactor these checks into a func so we can reuse it
		sessionIdCookieValue, cookieErr := c.Cookie("sessionId")
		if cookieErr != nil || sessionIdCookieValue == "" {
			logger.Error("No session ID found")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": "No session ID found",
			})
			return
		}

		sessionUser, getUserErr := lib.GetUserBySessionId(dbPool, logger, sessionIdCookieValue)
		if getUserErr != nil || sessionUser == (lib.User{}) {
			logger.Error("Error fetching session user: %v", getUserErr)
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "ERROR",
				"message": "No user found for session ID",
			})
			return
		}

		user, err := lib.GetUserBySlug(dbPool, logger, userSlug)
		if err != nil || user == (lib.User{}) {
			logger.Error("Error fetching user with slug %v: %v", userSlug, err)
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "ERROR",
				"message": "User not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"user":   user,
		})
	})
}
