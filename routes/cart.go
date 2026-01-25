package routes

import (
	"fmt"
	"log/slog"
	"net/http"

	"hotsauceshop/lib"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

//nolint:funlen
func Cart(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger) {
	r.GET("/api/v1/cart", func(c *gin.Context) {
		userId, userSessionErr := GetUserIdFromSessionOrError(c, dbPool, logger)
		if userSessionErr != nil || userId == 0 {
			return
		}

		cartItems, err := lib.GetCartItems(dbPool, userId)
		if err != nil {
			logger.Error(fmt.Sprintf("Error getting cart items: %v", err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Error fetching cart: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"cartItems": cartItems,
			},
		})
	})

	/**
	Add cart item
	1. Verify cart request
	2. Verify referenced user
	3. Verify referenced item
	*/
	var newCart lib.AddCartItemRequest
	r.POST("/api/v1/cart", func(c *gin.Context) {
		if err := c.ShouldBindJSON(&newCart); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		userId, userSessionErr := GetUserIdFromSessionOrError(c, dbPool, logger)
		if userSessionErr != nil || userId == 0 {
			return
		}

		// Create cart item
		err := lib.UpdateCart(dbPool, logger, newCart)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status": "OK",
		})
	})

	var deleteRequest lib.DeleteCartItemRequest
	r.DELETE("/api/v1/cart", func(c *gin.Context) {
		if err := c.ShouldBindJSON(&deleteRequest); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		userId, err := GetUserIdFromSessionOrError(c, dbPool, logger)
		if err != nil || userId == 0 {
			return
		}

		deleteErr := lib.DeleteCartItem(dbPool, deleteRequest.InventoryItemId, userId)
		if deleteErr != nil {
			logger.Error(deleteErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": deleteErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})
}
