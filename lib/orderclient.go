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
	CouponTypeName   string    `json:"couponTypeName"`
}

type ShippingOption struct {
	Id                     int     `json:"id"`
	Name                   string  `json:"name"`
	Description            string  `json:"description"`
	Price                  float64 `json:"price"`
	TimeToShipUnitQuantity int     `json:"timeToShipUnitQuantity"`
	TimeToShipUnit         string  `json:"timeToShipUnit"`
}

func GetShippingOptions(dbPool *pgxpool.Pool) ([]ShippingOption, error) {
	const query = `SELECT * FROM shipping_options ORDER BY price DESC`
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	shippingOptions, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[ShippingOption])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}
	return shippingOptions, nil
}

func GetCouponByCode(dbPool *pgxpool.Pool, code string) (CouponCode, error) {
	const query = `
		SELECT code, description, created_at, updated_at, expires_at, reduction_percent,
		       ct.name AS coupon_type_name
		FROM coupons
		JOIN coupon_types ct ON ct.id = coupons.coupon_type_id
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
