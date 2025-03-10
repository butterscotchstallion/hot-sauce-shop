package products

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"hotsauceshop/ent"
	"hotsauceshop/ent/inventory"
)

func Routes(r *gin.Engine, client *ent.Client) {
	r.GET("/api/v1/products/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		var res gin.H
		if len(slug) > 0 {
			product, err := client.Inventory.Query().
				Where(inventory.Slug(slug)).
				All(c)
			if err != nil || len(product) == 0 {
				res = gin.H{
					"status":  "ERROR",
					"message": fmt.Sprintf("Error fetching product: %v", err),
				}
			} else {
				res = gin.H{
					"status": "OK",
					"results": gin.H{
						"product": product[0],
					},
				}
			}
		}
		c.JSON(200, res)
	})

	r.GET("/api/v1/products", func(c *gin.Context) {
		offset := c.DefaultQuery("offset", "1")
		perPage := c.DefaultQuery("perPage", "10")

		perPageInt, perPageErr := strconv.Atoi(perPage)
		if perPageErr != nil {
			perPageInt = 10
		}

		offsetInt, offsetErr := strconv.Atoi(offset)
		if offsetErr != nil {
			offsetInt = 1
		}

		total, totalErr := client.Inventory.Query().
			Aggregate(ent.Count()).
			Int(c)
		if totalErr != nil {
			log.Printf("Error getting total inventory items: %v", totalErr)
		}

		var res gin.H
		inventoryResults, err := client.Inventory.
			Query().
			Offset(offsetInt).
			Limit(perPageInt).
			Order(ent.Asc(inventory.FieldName)).
			All(c)
		if err != nil {
			res = gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Error fetching inventory: %v", err),
			}
		} else {
			res = gin.H{
				"status": "OK",
				"results": gin.H{
					"inventory": inventoryResults,
					"total":     total,
				},
			}
		}

		c.JSON(200, res)
	})
}
