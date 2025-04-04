package lib

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
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
	ReviewCount      int        `json:"reviewCount" db:"review_count"`
	AverageRating    *float32   `json:"averageRating" db:"average_rating"`
}

type ProductAutocompleteSuggestion struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func AddOrUpdateInventoryItem(dbPool *pgxpool.Pool, logger *slog.Logger, inventoryItem InventoryItem) (int, error) {
	const query = `
		INSERT INTO inventories (
			name,
			description,
			short_description,
			slug,
			price,
			spice_rating,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (name) DO UPDATE SET 
		name = $1, description = $2, short_description = $3,
		slug = $4, price = $5, spice_rating = $6, updated_at = NOW()
		RETURNING id
	`
	var id int
	err := dbPool.QueryRow(
		context.Background(),
		query,
		inventoryItem.Name,
		inventoryItem.Description,
		inventoryItem.ShortDescription,
		inventoryItem.Slug,
		inventoryItem.Price,
		inventoryItem.SpiceRating,
	).Scan(&id)
	if err != nil {
		logger.Error(fmt.Sprintf("AddOrUpdateInventoryItem: Error adding/updating inventory item: %v", err))
		return 0, err
	}
	return id, nil
}

func GetInventoryItemsOrderedBySortKey(
	dbPool *pgxpool.Pool, logger *slog.Logger, limit int, offset int, sort string,
	tagIds []int,
) ([]InventoryItem, error) {
	// Sort is validated at endpoint
	direction := "ASC"
	descSorts := []string{"spice_rating", "review_count", "price"}
	if slices.Contains(descSorts, sort) {
		direction = "DESC"
	}
	sortClause := fmt.Sprintf("ORDER BY %s %s\n", sort, direction)
	offsetClause := fmt.Sprintf("OFFSET %d\n", offset)
	limitClause := ""
	if limit > 0 {
		limitClause = fmt.Sprintf("LIMIT %d\n", limit)
	}

	// Converted to int in controller area first
	tagIdsClause := ""
	tagIdsJoinClause := ""
	if len(tagIds) > 0 {
		tagIdsClause = "AND it.tag_id = ANY($1)\n"
		tagIdsJoinClause = "JOIN inventory_tags it ON i.id = it.inventory_id\n"
	}

	query := `
		SELECT i.id,
		       i.name, 
		       i.description,
		       i.short_description,
		       i.slug,
		       i.price,
		       i.created_at,
		       i.updated_at,
		       i.spice_rating,
		       (SELECT COUNT(*) FROM inventory_item_reviews WHERE inventory_item_id = i.id) AS review_count,
		       (SELECT AVG(rating) FROM inventory_item_reviews WHERE inventory_item_id = i.id) AS average_rating
		FROM inventories i
	` + tagIdsJoinClause + `
		WHERE 1=1
	` + tagIdsClause + sortClause + limitClause + offsetClause
	var rows pgx.Rows
	var err error
	if len(tagIdsClause) > 0 {
		rows, err = dbPool.Query(context.Background(), query, tagIds)
	} else {
		rows, err = dbPool.Query(context.Background(), query)
	}
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

func DeleteInventoryItemTags(dbPool *pgxpool.Pool, logger *slog.Logger, inventoryItemId int, tagIds []int) (bool, error) {
	const query = `DELETE FROM inventory_tags WHERE inventory_id = $1 AND tag_id = ANY($2)`
	_, err := dbPool.Exec(context.Background(), query, inventoryItemId, tagIds)
	if err != nil {
		logger.Error(fmt.Sprintf("Error deleting inventory item tags: %v", err))
		return false, err
	}
	return true, nil
}

func UpdateInventoryItemTags(dbPool *pgxpool.Pool, logger *slog.Logger, inventoryItemId int, tagIds []int) (bool, error) {
	if len(tagIds) == 0 {
		logger.Info("No tags provided for item, not updating tags")
		return true, nil
	}

	_, deleteErr := DeleteInventoryItemTags(dbPool, logger, inventoryItemId, tagIds)
	if deleteErr != nil {
		return false, deleteErr
	}
	for _, tagId := range tagIds {
		const query = `INSERT INTO inventory_tags (inventory_id, tag_id) VALUES ($1, $2)`
		_, insertTagsErr := dbPool.Exec(context.Background(), query, inventoryItemId, tagId)
		if insertTagsErr != nil {
			return false, insertTagsErr
		}
	}

	logger.Info(fmt.Sprintf("Updated inventory item tags: %v", tagIds))

	return true, nil
}

func GetAutocompleteSuggestions(dbPool *pgxpool.Pool, logger *slog.Logger, searchQuery string) ([]ProductAutocompleteSuggestion, error) {
	const query = `
		SELECT name, slug
		FROM inventories
		WHERE name ILIKE '%' || $1 || '%'
		ORDER BY name
		LIMIT 10
	`
	rows, err := dbPool.Query(context.Background(), query, searchQuery)
	defer rows.Close()
	if err != nil {
		logger.Error(fmt.Sprintf("Error running inventory item query: %v", err))
		return nil, err
	}

	suggestions, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[ProductAutocompleteSuggestion])
	if collectRowsErr != nil {
		logger.Error(fmt.Sprintf("Error collecting inventory item suggestions: %v", collectRowsErr))
		return nil, collectRowsErr
	}

	return suggestions, nil
}

func GetInventoryItemTags(dbPool *pgxpool.Pool, logger *slog.Logger, inventoryItemId int) ([]Tag, error) {
	const query = `
		SELECT t.id, t.name, t.description, t.slug, t.created_at, t.updated_at
		FROM tags t
		JOIN inventory_tags it ON it.tag_id = t.id
		WHERE it.inventory_id = $1
	`
	rows, err := dbPool.Query(context.Background(), query, inventoryItemId)
	defer rows.Close()
	if err != nil {
		logger.Error(fmt.Sprintf("Error running inventory item tags query: %v", err))
		return nil, err
	}
	tags, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[Tag])
	if collectRowsErr != nil {
		logger.Error(fmt.Sprintf("Error collecting inventory item tags: %v", collectRowsErr))
		return nil, collectRowsErr
	}
	return tags, nil
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
		   	spice_rating,
		   	(SELECT COUNT(*) FROM inventory_item_reviews WHERE inventory_item_id = i.id) AS review_count
		FROM inventories i
		WHERE slug = $1
	`
	inventoryItem := InventoryItem{}
	rows, err := dbPool.Query(context.Background(), query, slug)
	defer rows.Close()
	if err != nil {
		return inventoryItem, err
	}
	inventoryItems, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[InventoryItem])

	if collectRowsErr != nil {
		return inventoryItem, collectRowsErr
	}

	if len(inventoryItems) == 0 {
		return inventoryItem, fmt.Errorf("inventory item not found")
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
