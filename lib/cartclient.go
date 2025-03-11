package lib

import (
	"time"

	"github.com/jackc/pgx/v5"
)

type CartClient struct {
	GetCartById int32
}

type CartItem struct {
	id            int32
	inventoryItem InventoryItem
	quantity      int8
	createdAt     time.Time
	updatedAt     time.Time
}

func GetCartById(c *pgx.Conn, id int32) {

}
