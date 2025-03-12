package routes

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"hotsauceshop/lib"
)

const USER_ID = 1

func Cart(r *gin.Engine, conn *pgx.Conn, logger *slog.Logger) {
	r.GET("/api/v1/cart", func(c *gin.Context) {
		var res gin.H

		cartItems, err := lib.GetCartItems(conn)
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
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		// Create cart item
		err := lib.UpdateCart(conn, newCart)
		if err != nil {
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
