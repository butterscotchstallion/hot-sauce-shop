package routes

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"hotsauceshop/lib"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

const CacheTimeProductPage = 15 * time.Minute

//nolint:funlen
func Products(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger, store *persistence.InMemoryStore) {
	r.GET("/api/v1/products/:slug", cache.CachePage(store, CacheTimeProductPage, func(c *gin.Context) {
		urlSlug := c.Param("slug")
		var res gin.H
		product, err := lib.GetInventoryItemBySlug(dbPool, urlSlug)
		if err != nil {
			logger.Error(fmt.Sprintf("Error fetching product: %v", err))
			res = gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Error fetching product: %v", err),
			}
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		tags, tagsErr := lib.GetInventoryItemTags(dbPool, logger, product.Id)
		if tagsErr != nil {
			logger.Error(fmt.Sprintf("Error fetching tags: %v", tagsErr))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Error fetching tags: %v", tagsErr),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"product": product,
				"tags":    tags,
			},
		})
	}))

	r.GET("/api/v1/products", cache.CachePage(store, time.Minute*15, func(c *gin.Context) {
		paginationData := lib.GetValidPaginationData(c)
		filterTags := c.DefaultQuery("tags", "")

		tagIds := lib.ToIntArray(filterTags)

		// Validate sort
		sort := c.DefaultQuery("sort", "name")
		sorts := []string{
			"name", "price", "spice_rating", "created_at",
			"review_count", "average_rating", "average_spice_rating"}
		if !slices.Contains(sorts, sort) {
			sort = "name"
		}

		total, totalErr := lib.GetTotalInventoryItems(dbPool)
		if totalErr != nil {
			log.Printf("Error getting total inventory items: %v", totalErr)
		}

		var res gin.H
		inventoryResults, err := lib.GetInventoryItemsOrderedBySortKey(
			dbPool, logger, paginationData.PerPage, paginationData.Offset, sort, tagIds,
		)
		if err != nil {
			res = gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Error fetching inventory: %v", err),
			}
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"inventory": inventoryResults,
				"total":     total,
			},
		})
	}))

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

	/*
		- Attempt to parse request JSON
		- Validate request
		- Check if the item exists
		- Get user from sessionId
		- Add review
	*/
	r.POST("/api/v1/products/:slug/reviews", func(c *gin.Context) {
		// Check request JSON
		var inventoryItemReviewRequest lib.InventoryItemReviewRequest
		if err := c.ShouldBindJSON(&inventoryItemReviewRequest); err != nil {
			logger.Error(fmt.Sprintf("Malformed review request: %v", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": "Malformed request body.",
			})
			return
		}

		// Validate data
		validate := validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(inventoryItemReviewRequest)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Validation failed: %v", err),
			})
			return
		}

		// Check if the user signed in
		signedInUserId, userSessionErr := GetUserIdFromSessionOrError(c, dbPool, logger)
		if userSessionErr != nil || signedInUserId == 0 {
			return
		}

		// Check if item exists
		item, itemErr := lib.GetInventoryItemBySlug(dbPool, c.Param("slug"))
		if itemErr != nil || item == (lib.InventoryItem{}) {
			logger.Error(fmt.Sprintf("Error fetching inventory item: %v", itemErr.Error()))
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "ERROR",
				"message": "Error fetching inventory item.",
			})
			return
		}

		// Add review
		_, reviewErr := lib.AddInventoryItemReview(dbPool, item.Id, signedInUserId, inventoryItemReviewRequest)
		if reviewErr != nil {
			logger.Error(fmt.Sprintf("Error adding review: %v", reviewErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error adding review.",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  "OK",
			"message": "Review added.",
		})
	})

	r.GET("/api/v1/products/:slug/reviews", func(c *gin.Context) {
		paginationData := lib.GetValidPaginationData(c)
		itemSlug := c.Param("slug")
		reviews, reviewsErr := lib.GetInventoryItemReviewsBySlug(
			dbPool,
			logger,
			paginationData.PerPage,
			paginationData.Offset,
			itemSlug,
		)
		if reviewsErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Error fetching reviews: %v", reviewsErr.Error()),
			})
			return
		}

		ratingDistribution, ratingErr := lib.GetInventoryItemReviewRatingDistributionBySlug(dbPool, itemSlug)
		if ratingErr != nil {
			logger.Error(fmt.Sprintf("Error fetching rating distribution: %v", ratingErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Error fetching rating distribution: %v", ratingErr.Error()),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"reviews":            reviews,
				"ratingDistribution": ratingDistribution,
			},
		})
	})

	// TODO: add product admin role check here
	r.POST("/api/v1/products", func(c *gin.Context) {
		itemUpdateRequest := lib.InventoryItemUpdateRequest{}
		itemUpdateRequest, validationErr := lib.ValidateInventoryItemAddOrUpdateRequest(c, logger, itemUpdateRequest)
		// Error responses handled in the above func
		if validationErr != nil {
			return
		}

		itemId, saveItemErr := lib.SaveInventoryItem(dbPool, logger, itemUpdateRequest)
		if saveItemErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error updating inventory item.",
			})
			return
		}

		_, tagUpdateErr := lib.UpdateInventoryItemTags(dbPool, logger, itemId, itemUpdateRequest.TagIds)
		if tagUpdateErr != nil {
			logger.Error(fmt.Sprintf("Error updating product tags: %v", tagUpdateErr.Error()))
		}

		logger.Info(fmt.Sprintf("Inventory item tags updated: %v", itemUpdateRequest.TagIds))

		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": fmt.Sprintf("Inventory item #%v added", itemId),
			"results": gin.H{
				"inventoryItemId": itemId,
			},
		})
	})

	r.PUT("/api/v1/products/:slug", func(c *gin.Context) {
		itemUpdateRequest := lib.InventoryItemUpdateRequest{}
		itemUpdateRequest, validationErr := lib.ValidateInventoryItemAddOrUpdateRequest(c, logger, itemUpdateRequest)
		if validationErr != nil {
			return
		}

		urlSlug := c.Param("slug")
		item, itemErr := lib.GetInventoryItemBySlug(dbPool, urlSlug)
		if itemErr != nil {
			logger.Error(itemErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error fetching inventory item.",
			})
			return
		}

		if item == (lib.InventoryItem{}) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "ERROR",
				"message": "Inventory item not found.",
			})
			return
		}

		itemId, saveItemErr := lib.SaveInventoryItem(dbPool, logger, itemUpdateRequest)
		if saveItemErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error updating inventory item.",
			})
			return
		}

		logger.Info(fmt.Sprintf("Update item req: %v", itemUpdateRequest))

		_, tagUpdateErr := lib.UpdateInventoryItemTags(dbPool, logger, itemId, itemUpdateRequest.TagIds)
		if tagUpdateErr != nil {
			logger.Error(fmt.Sprintf("Error updating product tags: %v", tagUpdateErr.Error()))
		}

		// TODO: add websocket event for item updates here

		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "Inventory item updated.",
			"results": gin.H{
				"inventoryItemId": itemId,
			},
		})
	})
}
