package lib

import (
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
	Title       string `json:"title"`
	Comment     string `json:"comment"`
	Rating      int    `json:"rating"`
	SpiceRating int    `json:"spiceRating"`
}

func AddInventoryItemReview(dbPool *pgxpool.Pool, req InventoryItemReviewRequest) (bool, error) {
	return true, nil
}
