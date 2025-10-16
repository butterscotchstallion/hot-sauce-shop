package lib

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Board struct {
	Id                int    `json:"id"`
	DisplayName       string `json:"display_name"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	Slug              string `json:"slug"`
	Visible           bool   `json:"visible"`
	ThumbnailFilename string `json:"thumbnail_filename"`
}

func GetBoards(dbPool *pgxpool.Pool) ([]Board, error) {
	// TODO: filter visible boards, or show everything if privileged
	const query = `
		SELECT *
		FROM boards b
		ORDER BY b.display_name `
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	boards, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[Board])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}
	return boards, nil
}
