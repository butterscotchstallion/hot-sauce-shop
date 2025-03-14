package lib

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartItem struct {
	Id              int        `json:"id" db:"id"`
	InventoryItemId int        `json:"inventoryItem" db:"inventory_item_id"`
	UserId          int        `json:"userId" db:"user_id"`
	Quantity        int        `json:"quantity" db:"quantity"`
	CreatedAt       *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       *time.Time `json:"updatedAt" db:"updated_at"`
}

type AddCartItemRequest struct {
	InventoryItemId  int  `json:"inventoryItemId"`
	Quantity         int  `json:"quantity"`
	UserId           int  `json:"userId"`
	OverrideQuantity bool `json:"overrideQuantity"`
}

func GetCartItems(dbPool *pgxpool.Pool) ([]CartItem, error) {
	const query = `SELECT * FROM cart_items`
	rows, err := dbPool.Query(context.Background(), query)
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
func validateAddCartItemRequest(dbPool *pgxpool.Pool, req AddCartItemRequest) (bool, error) {
	itemExists, itemExistsErr := InventoryItemExists(dbPool, req.InventoryItemId)
	if itemExistsErr != nil {
		return false, itemExistsErr
	}

	userExists, userExistsErr := UserExists(dbPool, req.UserId)
	if userExistsErr != nil {
		return false, userExistsErr
	}

	if req.Quantity < 1 || req.Quantity > 100 {
		return false, errors.New(fmt.Sprintf("quantity must be between 1 and 100: %v", req.Quantity))
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
func UpdateCart(dbPool *pgxpool.Pool, logger *slog.Logger, req AddCartItemRequest) error {
	isValid, validityErr := validateAddCartItemRequest(dbPool, req)
	if validityErr != nil || !isValid {
		return validityErr
	}

	existingCartItem, err := GetCartItemByInventoryItemIdAndUserId(dbPool, req.InventoryItemId, req.UserId)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	logger.Info(fmt.Sprintf("existingCartItem: %v", existingCartItem))

	quantity := 1
	if existingCartItem != (CartItem{}) {
		quantity = existingCartItem.Quantity + 1
	}

	if req.OverrideQuantity {
		quantity = req.Quantity
	}

	logger.Info(fmt.Sprintf("Adding cart item %v with quantity: %v", req.InventoryItemId, quantity))

	_, addCartErr := addCartItem(dbPool, req.InventoryItemId, req.UserId, quantity)
	if addCartErr != nil {
		return addCartErr
	}

	return nil
}

func addCartItem(dbPool *pgxpool.Pool, inventoryItemId int, userId int, quantity int) (int, error) {
	lastInsertId := 0
	const query = `
		INSERT INTO cart_items (quantity, inventory_item_id, user_id, created_at) 
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT(inventory_item_id, user_id)
		    DO UPDATE SET quantity = cart_items.quantity + 1, updated_at = NOW()
		RETURNING id
	`
	insertErr := dbPool.QueryRow(context.Background(), query, quantity, inventoryItemId, userId).Scan(&lastInsertId)
	if insertErr != nil {
		return 0, insertErr
	}
	return lastInsertId, nil
}

func GetCartItemByInventoryItemIdAndUserId(dbPool *pgxpool.Pool, inventoryItemId int, userId int) (CartItem, error) {
	const query = `
		SELECT * 
		FROM cart_items
		WHERE 1=1
		AND inventory_item_id = $1
		AND user_id = $2
	`
	cartItem := CartItem{}
	rows, err := dbPool.Query(context.Background(), query, inventoryItemId, userId)
	if err != nil {
		return cartItem, err
	}
	defer rows.Close()
	cartItem, collectRowsErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[CartItem])
	if collectRowsErr != nil {
		return cartItem, err
	}
	return cartItem, nil
}
