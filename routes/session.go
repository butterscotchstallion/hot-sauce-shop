package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"hotsauceshop/lib"
)

func Session(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger) {
	r.GET("/api/v1/session", func(c *gin.Context) {
		sessionIdCookieValue, err := c.Cookie("sessionId")

		if err != nil || sessionIdCookieValue == "" {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ERROR",
				"message": "No session ID found",
			})
			return
		}

		user, getUserErr := lib.GetUserBySessionId(dbPool, logger, sessionIdCookieValue)

		if getUserErr != nil || user == (lib.User{}) {
			logger.Error("Error fetching user: %v", getUserErr)
			c.JSON(http.StatusOK, gin.H{
				"status":  "ERROR",
				"message": "No user found for session ID",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"user": user,
			},
		})
	})
}
