package lib

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CouponCode struct {
	Code             string    `json:"code"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	ExpiresAt        time.Time `json:"expiresAt"`
	ReductionPercent int       `json:"reductionPercent"`
}

func GetCouponByCode(dbPool *pgxpool.Pool, code string) (CouponCode, error) {
	const query = `
		SELECT code, description, created_at, updated_at, expires_at, reduction_percent
		FROM coupons
		WHERE UPPER(code) = UPPER($1)
		AND expires_at > NOW()
	`
	row, err := dbPool.Query(context.Background(), query, code)
	if err != nil {
		return CouponCode{}, err
	}
	defer row.Close()
	coupon, collectRowsErr := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[CouponCode])
	if collectRowsErr != nil {
		return CouponCode{}, collectRowsErr
	}
	return coupon, nil
}
