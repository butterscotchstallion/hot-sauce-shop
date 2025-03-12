package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"hotsauceshop/lib"
)

func Tags(r *gin.Engine, conn *pgx.Conn) {
	r.GET("/api/v1/tags", func(c *gin.Context) {
		var res gin.H
		tags, err := lib.GetTagsOrderedByName(conn)
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
	})
}
