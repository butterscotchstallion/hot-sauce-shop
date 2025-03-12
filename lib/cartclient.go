package lib

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type CartItem struct {
	Id            int32         `json:"id"`
	InventoryItem InventoryItem `json:"inventoryItem"`
	Quantity      int8          `json:"quantity"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}

type AddCartItemRequest struct {
	InventoryItem    InventoryItem `json:"inventoryItem"`
	Quantity         int8          `json:"quantity"`
	UserId           int           `json:"userId"`
	OverrideQuantity bool          `json:"overrideQuantity"`
}

func GetCartItems(c *pgx.Conn) ([]CartItem, error) {
	const query = `SELECT * FROM cart_items`
	rows, err := c.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	cartItems, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[CartItem])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}
	return cartItems, nil
}

/*
validateAddCartItemRequest
1. Check if inventory item exists
2. Check if user exists
3. Check if quantity is > 0
*/
func validateAddCartItemRequest(conn *pgx.Conn, req AddCartItemRequest) (bool, error) {
	itemExists, itemExistsErr := InventoryItemExists(conn, req.InventoryItem.Id)
	if itemExistsErr != nil {
		return false, itemExistsErr
	}

	userExists, userExistsErr := UserExists(conn, req.UserId)
	if userExistsErr != nil {
		return false, userExistsErr
	}

	if req.Quantity <= 1 || req.Quantity > 100 {
		return false, errors.New("quantity must be between 1 and 100")
	}

	return itemExists && userExists, nil
}

/*
UpdateCart
 1. Check if a cart item with this inventory item and user id exists
 2. Quantity is 1 by default
 3. If cart item exists, add 1 to that
 3. If override quantity, update quantity
 4. Add cart item
*/
func UpdateCart(c *pgx.Conn, req AddCartItemRequest) error {
	isValid, validityErr := validateAddCartItemRequest(c, req)
	if validityErr != nil || !isValid {
		return validityErr
	}

	existingCartItem, err := GetCartItemByInventoryItemIdAndUserId(c, req.InventoryItem.Id, req.UserId)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	quantity := int8(1)
	if existingCartItem != (CartItem{}) {
		quantity = existingCartItem.Quantity + 1
	}

	if req.OverrideQuantity {
		quantity = req.Quantity
	}

	_, addCartErr := addCartItem(c, req.InventoryItem.Id, req.UserId, quantity)

	if addCartErr != nil {
		return addCartErr
	}

	return nil
}

func addCartItem(c *pgx.Conn, inventoryItemId int, userId int, quantity int8) (int32, error) {
	lastInsertId := int32(0)
	const query = `
		INSERT INTO cart_items (quantity, inventory_item_id, user_id) 
		VALUES ($1, $2, $3)
		RETURNING id
	`
	insertErr := c.QueryRow(context.Background(), query, quantity, inventoryItemId, userId).Scan(lastInsertId)
	if insertErr != nil {
		return 0, insertErr
	}
	return lastInsertId, nil
}

func GetCartItemByInventoryItemIdAndUserId(c *pgx.Conn, inventoryItemId int, userId int) (CartItem, error) {
	const query = `
		SELECT * 
		FROM cart_items
		WHERE 1=1
		AND inventory_item_id = $1
		AND user_id = $2
	`
	cartItem := CartItem{}
	err := c.QueryRow(context.Background(), query, inventoryItemId, userId).Scan(&cartItem)
	if err != nil {
		return cartItem, err
	}
	return cartItem, nil
}
