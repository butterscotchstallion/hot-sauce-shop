package lib

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryItemReview struct {
	Id                 int       `json:"id"`
	Title              string    `json:"title"`
	Comment            string    `json:"comment"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	Rating             int       `json:"rating"`
	SpiceRating        int       `json:"spiceRating"`
	InventoryItemId    int       `json:"inventoryItemId"`
	UserId             int       `json:"userId"`
	Username           string    `json:"username"`
	UserAvatarFilename string    `json:"userAvatarFilename"`
	UsernameSlug       string    `json:"usernameSlug"`
}

type InventoryItemReviewRequest struct {
	Title       string `json:"title" validate:"required,min=10,max=255"`
	Comment     string `json:"comment" validate:"required,min=10,max=1000"`
	Rating      int    `json:"rating" validate:"required,min=1,max=5"`
	SpiceRating int    `json:"spiceRating" validate:"required,min=1,max=5"`
}

func AddInventoryItemReview(dbPool *pgxpool.Pool, inventoryItemId int, userId int, req InventoryItemReviewRequest) (bool, error) {
	const query = `
		INSERT INTO inventory_item_reviews (
			title, comment, rating, spice_rating, inventory_item_id, user_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`
	_, insertErr := dbPool.Exec(
		context.Background(),
		query,
		req.Title,
		req.Comment,
		req.Rating,
		req.SpiceRating,
		inventoryItemId,
		userId,
	)
	if insertErr != nil {
		return false, insertErr
	}
	return true, nil
}

func GetInventoryItemReviewsBySlug(dbPool *pgxpool.Pool, logger *slog.Logger, perPage int, offset int, itemSlug string) ([]InventoryItemReview, error) {
	const query = `
			SELECT
				reviews.id,
				reviews.rating,
				reviews.spice_rating,
				reviews.comment,
				reviews.created_at,
				reviews.title,
				reviews.updated_at,
				reviews.inventory_item_id,
				users.id AS userId,
				users.username,
				users.avatar_filename AS userAvatarFilename,
				users.slug AS usernameSlug
			FROM inventory_item_reviews AS reviews
			JOIN inventories ON reviews.inventory_item_id = inventories.id
			JOIN users ON reviews.user_id = users.id
			WHERE inventories.slug = $1
			ORDER BY reviews.created_at DESC
			LIMIT $2
			OFFSET $3
		`
	rows, rowsErr := dbPool.Query(context.Background(), query, itemSlug, perPage, offset)
	if rowsErr != nil {
		logger.Error(fmt.Sprintf("Error fetching reviews: %v", rowsErr.Error()))
		return nil, rowsErr
	}
	defer rows.Close()
	reviews, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[InventoryItemReview])
	if collectRowsErr != nil {
		logger.Error(fmt.Sprintf("Error collecting review rows: %v", collectRowsErr.Error()))
		return nil, collectRowsErr
	}
	return reviews, nil
}
