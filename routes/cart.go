package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"hotsauceshop/ent"
	"hotsauceshop/ent/cartitems"
)

const USER_ID = 1

func Cart(r *gin.Engine, conn *pgx.Conn) {
	r.GET("/api/v1/cart", func(c *gin.Context) {
		var res gin.H

		cartItems, err := client.CartItems.Query().
			Order(ent.Asc(cartitems.FieldCreatedAt)).
			All(c)
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
	type Cart struct {
		inventoryItemId  int32
		overrideQuantity bool
		quantity         int8
	}

	var newCart Cart
	r.POST("/api/v1/cart", func(c *gin.Context) {
		if err := c.ShouldBindJSON(&newCart); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		// Verify user
		user, userErr := client.User.Get(c, USER_ID)
		if userErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": "Could not find user",
			})
			return
		}

		// Verify item
		item, itemErr := client.Inventory.Get(c, int(newCart.inventoryItemId))
		if itemErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": "Could not find inventory item",
			})
			return
		}

		/*allCartItems, allCartItemsErr := client.CartItems.Query().All(c)

		if allCartItemsErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error fetching cart items",
			})
			return
		}

		// Does this cart item exist?
		existingCartItem, cartItemErr := client.CartItems.
			QueryInventory(allCartItems).
			Where(inventory.ID(int(newCart.inventoryItemId))).
			All(c)
		*/

		// Create cart item
		cart, err := client.CartItems.Create().
			SetQuantity(newCart.quantity).
			Save(c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
		} else {
			updatedCart, cartErr := cart.Update().
				AddUser(user).
				AddInventory(item).
				Save(c)
			if cartErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "ERROR",
					"message": cartErr.Error(),
				})
			}

			c.JSON(http.StatusCreated, gin.H{
				"status": "OK",
				"results": gin.H{
					"cart": updatedCart,
				},
			})
		}
	})
}
