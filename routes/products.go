package routes

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"slices"
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
				logger.Error(fmt.Sprintf("Error fetching product: %v", err))
				res = gin.H{
					"status":  "ERROR",
					"message": fmt.Sprintf("Error fetching product: %v", err),
				}
				c.JSON(http.StatusInternalServerError, res)
			} else {
				res = gin.H{
					"status": "OK",
					"results": gin.H{
						"product": product,
					},
				}
				c.JSON(http.StatusOK, res)
			}
		}
	})

	r.GET("/api/v1/products", func(c *gin.Context) {
		searchQuery := c.DefaultQuery("q", "")
		offset := c.DefaultQuery("offset", "0")
		perPage := c.DefaultQuery("perPage", "10")

		// Validate search
		if len(searchQuery) > 25 {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"status":  "ERROR",
					"message": "Search query must be less than 25 characters",
				},
			)
			return
		}

		// Validate sort
		sort := c.DefaultQuery("sort", "name")
		sorts := []string{"name", "price", "spice_rating"}
		if !slices.Contains(sorts, sort) {
			sort = "name"
		}

		// Validate page/offset
		perPageInt, perPageErr := strconv.Atoi(perPage)
		if perPageErr != nil || perPageInt < 10 || perPageInt > 30 {
			perPageInt = 10
		}

		offsetInt, offsetErr := strconv.Atoi(offset)
		if offsetErr != nil || offsetInt < 0 || offsetInt > 1000000 {
			offsetInt = 0
		}

		total, totalErr := lib.GetTotalInventoryItems(dbPool)
		if totalErr != nil {
			log.Printf("Error getting total inventory items: %v", totalErr)
		}

		var res gin.H
		inventoryResults, err := lib.GetInventoryItemsOrderedBySortKey(
			dbPool, logger, perPageInt, offsetInt, sort, searchQuery,
		)
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

	r.GET("/api/v1/products/autocomplete", func(c *gin.Context) {
		searchQuery := c.DefaultQuery("q", "")
		if len(searchQuery) == 0 || len(searchQuery) > 25 {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"status":  "ERROR",
					"message": "Search query must be between 1-25 characters",
				},
			)
			return
		}

		suggestions, err := lib.GetAutocompleteSuggestions(dbPool, logger, searchQuery)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"status":  "ERROR",
					"message": fmt.Sprintf("Error fetching autocomplete suggestions: %v", err),
				},
			)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"suggestions": suggestions,
			},
		})
	})

	r.POST("/api/v1/products", func(c *gin.Context) {
		/*
			inventory, err := client.Inventory.Create().
			SetName("yummy hot sauce").
			Save(c)
		*/
	})
}
