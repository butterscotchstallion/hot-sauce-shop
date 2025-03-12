package lib

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Tag struct {
	Id        int        `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Desc      string     `json:"description" db:"description"`
	Slug      string     `json:"slug" db:"slug"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}

func GetTagsOrderedByName(dbPool *pgxpool.Pool) ([]Tag, error) {
	const query = `
		SELECT id, name, description, slug, created_at, updated_at
		FROM tags
		ORDER BY name
	`
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[Tag])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}

	return tags, nil
}
