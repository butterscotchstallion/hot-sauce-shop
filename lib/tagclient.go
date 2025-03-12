package lib

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Tag struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Desc      string    `json:"description"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func GetTagsOrderedByName(c *pgx.Conn) ([]Tag, error) {
	const query = `
		SELECT name, description, slug, created_at AS createdAt, updated_at AS updatedAt
		FROM tags
		ORDER BY name
	`
	rows, err := c.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	tags, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[Tag])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}

	return tags, nil
}
