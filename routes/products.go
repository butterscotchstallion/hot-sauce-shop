package routes

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"hotsauceshop/lib"
	_ "hotsauceshop/lib"
)

func Products(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger) {
	r.GET("/api/v1/products/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		var res gin.H
		if len(slug) > 0 {
			product, err := lib.GetInventoryItemBySlug(dbPool, slug)
			if err != nil {
				res = gin.H{
					"status":  "ERROR",
					"message": fmt.Sprintf("Error fetching product: %v", err),
				}
			} else {
				res = gin.H{
					"status": "OK",
					"results": gin.H{
						"product": product,
					},
				}
			}
		}
		c.JSON(http.StatusOK, res)
	})

	r.GET("/api/v1/products", func(c *gin.Context) {
		offset := c.DefaultQuery("offset", "0")
		perPage := c.DefaultQuery("perPage", "10")

		perPageInt, perPageErr := strconv.Atoi(perPage)
		if perPageErr != nil {
			perPageInt = 10
		}

		offsetInt, offsetErr := strconv.Atoi(offset)
		if offsetErr != nil {
			offsetInt = 0
		}

		total, totalErr := lib.GetTotalInventoryItems(dbPool)
		if totalErr != nil {
			log.Printf("Error getting total inventory items: %v", totalErr)
		}

		var res gin.H
		inventoryResults, err := lib.GetInventoryItemsOrderedByName(dbPool, logger, perPageInt, offsetInt)
		if err != nil {
			res = gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Error fetching inventory: %v", err),
			}
			c.JSON(http.StatusInternalServerError, res)
		} else {
			res = gin.H{
				"status": "OK",
				"results": gin.H{
					"inventory": inventoryResults,
					"total":     total,
				},
			}
			c.JSON(http.StatusOK, res)
		}
	})

	r.POST("/api/v1/products", func(c *gin.Context) {
		/*
			inventory, err := client.Inventory.Create().
			SetName("yummy hot sauce").
			Save(c)
		*/
	})
}
