package routes

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgxpool"
	"hotsauceshop/lib"
)

type InventoryItemUpdateRequest struct {
	Name             string  `json:"name" validate:"required,min=3,max=255"`
	Price            float32 `json:"price" validate:"required,min=0.01,max=999999.99"`
	SpiceRating      int     `json:"spiceRating" validate:"required,min=1,max=5"`
	TagIds           []int   `json:"tags"`
	Description      string  `json:"description" validate:"required,min=3,max=1000000"`
	ShortDescription string  `json:"shortDescription" validate:"required,min=3,max=1000"`
}

func toIntArray(str string) []int {
	chunks := strings.Split(str, ",")
	var res []int
	for _, c := range chunks {
		i, err := strconv.Atoi(c)
		if err != nil {
			continue
		}
		res = append(res, i)
	}
	return res
}

func validateInventoryItemAddOrUpdateRequest(c *gin.Context, logger *slog.Logger, itemUpdateRequest InventoryItemUpdateRequest) (InventoryItemUpdateRequest, error) {
	if err := c.ShouldBindJSON(&itemUpdateRequest); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "ERROR",
			"message": "Malformed request body.",
		})
		return itemUpdateRequest, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(itemUpdateRequest)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "ERROR",
			"message": fmt.Sprintf("Validation failed: %v", err),
		})
		return itemUpdateRequest, err
	}

	return itemUpdateRequest, nil
}

func saveInventoryItem(dbPool *pgxpool.Pool, logger *slog.Logger, itemUpdateRequest InventoryItemUpdateRequest) (int, error) {
	var item lib.InventoryItem
	item.Name = itemUpdateRequest.Name
	item.Price = itemUpdateRequest.Price
	item.SpiceRating = itemUpdateRequest.SpiceRating
	item.Description = itemUpdateRequest.Description
	item.ShortDescription = itemUpdateRequest.ShortDescription
	item.Slug = slug.Make(itemUpdateRequest.Name)

	logger.Info(fmt.Sprintf("Saving inventory item: %+v", item))

	itemId, addUpdateItemErr := lib.AddOrUpdateInventoryItem(dbPool, logger, item)
	if addUpdateItemErr != nil {
		logger.Error(addUpdateItemErr.Error())
		return 0, addUpdateItemErr
	}

	return itemId, nil
}

func Products(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger, store *persistence.InMemoryStore) {
	r.GET("/api/v1/products/:slug", cache.CachePage(store, time.Minute*15, func(c *gin.Context) {
		urlSlug := c.Param("slug")
		var res gin.H
		if len(urlSlug) > 0 {
			product, err := lib.GetInventoryItemBySlug(dbPool, urlSlug)
			if err != nil {
				logger.Error(fmt.Sprintf("Error fetching product: %v", err))
				res = gin.H{
					"status":  "ERROR",
					"message": fmt.Sprintf("Error fetching product: %v", err),
				}
				c.JSON(http.StatusInternalServerError, res)
			} else {
				tags, tagsErr := lib.GetInventoryItemTags(dbPool, logger, product.Id)
				if tagsErr != nil {
					logger.Error(fmt.Sprintf("Error fetching tags: %v", tagsErr))
					c.JSON(http.StatusInternalServerError, gin.H{
						"status":  "ERROR",
						"message": fmt.Sprintf("Error fetching tags: %v", tagsErr),
					})
					return
				}
				res = gin.H{
					"status": "OK",
					"results": gin.H{
						"product": product,
						"tags":    tags,
					},
				}
				c.JSON(http.StatusOK, res)
			}
		}
	}))

	r.GET("/api/v1/products", cache.CachePage(store, time.Minute*15, func(c *gin.Context) {
		offset := c.DefaultQuery("offset", "0")
		perPage := c.DefaultQuery("perPage", "10")
		filterTags := c.DefaultQuery("tags", "")

		tagIds := toIntArray(filterTags)

		// Validate sort
		sort := c.DefaultQuery("sort", "name")
		sorts := []string{"name", "price", "spice_rating", "created_at"}
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
			dbPool, logger, perPageInt, offsetInt, sort, tagIds,
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

	// TODO: add product admin role check here
	r.POST("/api/v1/products", func(c *gin.Context) {
		itemUpdateRequest := InventoryItemUpdateRequest{}
		itemUpdateRequest, validationErr := validateInventoryItemAddOrUpdateRequest(c, logger, itemUpdateRequest)
		// Error responses handled in above func
		if validationErr != nil {
			return
		}

		itemId, saveItemErr := saveInventoryItem(dbPool, logger, itemUpdateRequest)
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

		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": fmt.Sprintf("Inventory item #%v added", itemId),
			"results": gin.H{
				"inventoryItemId": itemId,
			},
		})
	})

	r.PUT("/api/v1/products/:slug", func(c *gin.Context) {
		itemUpdateRequest := InventoryItemUpdateRequest{}
		itemUpdateRequest, validationErr := validateInventoryItemAddOrUpdateRequest(c, logger, itemUpdateRequest)
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

		itemId, saveItemErr := saveInventoryItem(dbPool, logger, itemUpdateRequest)
		if saveItemErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error updating inventory item.",
			})
			return
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
