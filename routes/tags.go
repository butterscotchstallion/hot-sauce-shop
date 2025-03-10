package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"hotsauceshop/ent"
	"hotsauceshop/ent/tag"
)

func Tags(r *gin.Engine, client *ent.Client) {
	r.GET("/api/v1/tags", func(c *gin.Context) {
		var res gin.H
		tags, err := client.Tag.Query().Order(ent.Asc(tag.FieldName)).All(c)
		if err != nil {
			res = gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Error fetching tags: %v", err),
			}
		} else {
			res = gin.H{
				"status": "OK",
				"results": gin.H{
					"tags": tags,
				},
			}
		}

		c.JSON(200, res)
	})
}
