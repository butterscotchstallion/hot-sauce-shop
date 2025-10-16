package routes

import (
	"fmt"
	"net/http"
	"time"

	"hotsauceshop/lib"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Tags(r *gin.Engine, dbPool *pgxpool.Pool, store *persistence.InMemoryStore) {
	r.GET("/api/v1/tags", cache.CachePage(store, time.Minute*60, func(c *gin.Context) {
		var res gin.H
		tags, err := lib.GetTagsOrderedByName(dbPool)
		if err != nil {
			res = gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Error fetching tags: %v", err),
			}
			c.JSON(http.StatusInternalServerError, res)
		} else {
			res = gin.H{
				"status": "OK",
				"results": gin.H{
					"tags": tags,
				},
			}
			c.JSON(http.StatusOK, res)
		}
	}))
}
