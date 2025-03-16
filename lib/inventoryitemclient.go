package lib

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryItem struct {
	Id               int        `json:"id" db:"id"`
	Name             string     `json:"name" db:"name"`
	Description      string     `json:"description" db:"description"`
	ShortDescription string     `json:"shortDescription" db:"short_description"`
	Slug             string     `json:"slug" db:"slug"`
	Price            float32    `json:"price" db:"price"`
	SpiceRating      int        `json:"spiceRating" db:"spice_rating"`
	CreatedAt        time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt        *time.Time `json:"updatedAt" db:"updated_at"`
}

func GetInventoryItemsOrderedBySortKey(dbPool *pgxpool.Pool, logger *slog.Logger, limit int, offset int, sort string) ([]InventoryItem, error) {
	// Sort is validated at endpoint
	sortClause := fmt.Sprintf("ORDER BY %s\n", sort)
	offsetClause := fmt.Sprintf("OFFSET %d\n", offset)
	limitClause := ""
	if limit > 0 {
		limitClause = fmt.Sprintf("LIMIT %d\n", limit)
	}

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
	` + sortClause + limitClause + offsetClause
	rows, err := dbPool.Query(context.Background(), query)
	defer rows.Close()
	if err != nil {
		logger.Error(fmt.Sprintf("Error running inventory item query: %v", err))
		return nil, err
	}

	inventoryItems, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[InventoryItem])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}

	return inventoryItems, err
}

func GetInventoryItemBySlug(dbPool *pgxpool.Pool, slug string) (InventoryItem, error) {
	const query = `
		SELECT
			id,
		   	name, 
		   	description,
		   	short_description,
		   	slug,
		   	price,
		   	created_at,
		   	updated_at,
		   	spice_rating
		FROM inventories
		WHERE slug = $1
	`
	inventoryItem := InventoryItem{}
	rows, err := dbPool.Query(context.Background(), query, slug)
	if err != nil {
		return inventoryItem, err
	}
	inventoryItems, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[InventoryItem])
	if collectRowsErr != nil {
		return inventoryItems[0], collectRowsErr
	}

	return inventoryItems[0], nil
}

func InventoryItemExists(dbPool *pgxpool.Pool, id int) (bool, error) {
	exists := false
	const query = `
		SELECT EXISTS (
			SELECT 1 FROM inventories WHERE id = $1
        )
	`
	err := dbPool.QueryRow(context.Background(), query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetTotalInventoryItems(dbPool *pgxpool.Pool) (int32, error) {
	const query = `
		SELECT COUNT(*)
		FROM inventories
	`
	var count int32
	err := dbPool.QueryRow(context.Background(), query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
