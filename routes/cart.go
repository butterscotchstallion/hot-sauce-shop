package routes

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"hotsauceshop/lib"
)

const USER_ID = 1

func Cart(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger) {
	r.GET("/api/v1/cart", func(c *gin.Context) {
		var res gin.H

		cartItems, err := lib.GetCartItems(dbPool)
		if err != nil {
			res = gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Error fetching cart: %v", err),
			}
			c.JSON(500, res)
		} else {
			res = gin.H{
				"status": "OK",
				"results": gin.H{
					"cartItems": cartItems,
				},
			}
			c.JSON(200, res)
		}
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

		// Create cart item
		err := lib.UpdateCart(dbPool, newCart)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"status": "OK",
			})
		}
	})
}
