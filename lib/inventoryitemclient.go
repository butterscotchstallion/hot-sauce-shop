package lib

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type InventoryItem struct {
	id               int32
	name             string
	shortDescription string
	description      string
	slug             string
	price            float32
	spiceRating      int8
	createdAt        time.Time
	updatedAt        time.Time
}

func GetInventoryItemsOrderedByName(c *pgx.Conn, limit int32, offset int32) ([]InventoryItem, error) {
	offsetClause := ""
	limitClause := ""

	if limit > 0 {
		limitClause = fmt.Sprintf("LIMIT %d", limit)
	}

	if offset > 0 {
		offsetClause = fmt.Sprintf("OFFSET %d", offset)
	}

	query := fmt.Sprintf(`
		SELECT name, 
		       description,
		       short_description AS shortDescription,
		       slug,
		       price,
		       created_at AS createdAt,
		       updated_at AS updatedAt, 
		       spice_rating
		FROM inventories
		ORDER BY name
		%v
		%v
	`, limit, offset)
	rows, err := c.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

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
