package lib

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryItemReview struct {
	Id              int       `json:"id"`
	Title           string    `json:"title"`
	Comment         string    `json:"comment"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Rating          int       `json:"rating"`
	SpiceRating     int       `json:"spiceRating"`
	InventoryItemId int       `json:"inventoryItemId"`
	UserId          int       `json:"userId"`
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
