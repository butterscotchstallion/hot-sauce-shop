package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"hotsauceshop/lib"
)

func Admin(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger) {
	// TODO: implement RBAC checks for all routes here
	r.GET("/admin/users/:slug", func(c *gin.Context) {
		userSlug := c.Param("slug")

		if userSlug == "" {
			c.JSON(400, gin.H{
				"status":  "ERROR",
				"message": "User slug is required",
			})
			return
		}

		// TODO: refactor these checks into a func so we can reuse it
		sessionIdCookieValue, cookieErr := c.Cookie("sessionId")
		if cookieErr != nil || sessionIdCookieValue == "" {
			c.JSON(http.StatusOK, gin.H{
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
			c.JSON(404, gin.H{
				"status":  "ERROR",
				"message": "User not found",
			})
			return
		}

		c.JSON(200, gin.H{
			"status": "OK",
			"user":   user,
		})
	})
}
