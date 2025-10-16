package routes

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"hotsauceshop/lib"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Boards(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger, store *persistence.InMemoryStore) {
	r.GET("/api/v1/boards", cache.CachePage(store, time.Minute*300, func(c *gin.Context) {
		boards, getBoardsErr := lib.GetBoards(dbPool)
		if getBoardsErr != nil {
			logger.Error(fmt.Sprintf("Error fetching boards: %v", getBoardsErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": getBoardsErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"boards": boards,
		})
	}))
}
