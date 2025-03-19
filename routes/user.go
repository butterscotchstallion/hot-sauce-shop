package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"hotsauceshop/lib"
)

type LoginRequest struct {
	username string
	password string
}

func User(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger) {

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

		usernameAndPasswordVerified, errVerifying := lib.VerifyUsernameAndPassword(
			dbPool, loginRequest.username, loginRequest.password,
		)
		if errVerifying != nil {
			logger.Error(errVerifying.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": errVerifying.Error(),
			})
			return
		}

		if !usernameAndPasswordVerified {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ERROR",
				"message": "Invalid username or password",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "Login successful",
		})
	})
}
