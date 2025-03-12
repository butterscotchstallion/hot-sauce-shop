package lib

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
)

type InventoryItem struct {
	Id               int        `json:"id" db:"id"`
	Name             string     `json:"name" db:"name"`
	ShortDescription string     `json:"shortDescription" db:"short_description"`
	Description      string     `json:"description" db:"description"`
	Slug             string     `json:"slug" db:"slug"`
	Price            float32    `json:"price" db:"price"`
	SpiceRating      int8       `json:"spiceRating" db:"spice_rating"`
	CreatedAt        time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt        *time.Time `json:"updatedAt" db:"updated_at"`
}

func GetInventoryItemsOrderedByName(c *pgx.Conn, logger *slog.Logger, limit int, offset int) ([]InventoryItem, error) {
	limitClause := fmt.Sprintf("LIMIT %d\n", limit)
	offsetClause := fmt.Sprintf("OFFSET %d\n", offset)

	query := `
		SELECT id,
		       name, 
		       description,
		       short_description,
		       slug,
		       price,
		       created_at,
		       updated_at,
		       spice_rating
		FROM inventories
		ORDER BY name
	` + limitClause + offsetClause
	rows, err := c.Query(context.Background(), query)
	if err != nil {
		logger.Error(fmt.Sprintf("Error running inventory item query: %v", err))
		return nil, err
	}

	logger.Info(fmt.Sprintf("Inventory rows: %v", rows))

	inventoryItems, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[InventoryItem])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}

	return inventoryItems, err
}

func GetInventoryItemBySlug(c *pgx.Conn, slug string) (InventoryItem, error) {
	const query = `
		SELECT name, 
		       description,
		       short_description AS shortDescription,
		       slug,
		       price,
		       created_at AS createdAt,
		       updated_at AS updatedAt, 
		       spice_rating
		FROM inventories
		WHERE slug = @slug
	`
	inventoryItem := InventoryItem{}
	err := c.QueryRow(context.Background(), query, slug).Scan(&inventoryItem)
	if err != nil {
		return inventoryItem, err
	}
	return inventoryItem, nil
}

func InventoryItemExists(c *pgx.Conn, id int) (bool, error) {
	exists := false
	const query = `
		SELECT EXISTS (
			SELECT 1 FROM inventories WHERE id = $1
        )
	`
	err := c.QueryRow(context.Background(), query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetInventoryItemById(c *pgx.Conn, id int) (InventoryItem, error) {
	const query = `
		SELECT name, 
		       description,
		       short_description AS shortDescription,
		       slug,
		       price,
		       created_at AS createdAt,
		       updated_at AS updatedAt, 
		       spice_rating
		FROM inventories
		WHERE id = @id
	`
	inventoryItem := InventoryItem{}
	err := c.QueryRow(context.Background(), query, id).Scan(&inventoryItem)
	if err != nil {
		return inventoryItem, err
	}
	return inventoryItem, nil
}

func GetTotalInventoryItems(c *pgx.Conn) (int32, error) {
	const query = `
		SELECT COUNT(*)
		FROM inventories
	`
	var count int32
	err := c.QueryRow(context.Background(), query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
