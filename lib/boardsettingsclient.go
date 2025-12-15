package lib

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BoardSettingsResponse struct {
	Status  string        `json:"status"`
	Results BoardSettings `json:"results"`
}
type BoardSettings struct {
	IsOfficial             bool
	IsPostApprovalRequired bool
	UpdatedAt              *time.Time
}

func GetBoardSettings(dbPool *pgxpool.Pool, boardSlug string) (BoardSettings, error) {
	const query = `SELECT 
		bs.is_official, 
		bs.is_post_approval_required,
		bs.updated_at
		FROM board_settings bs
		JOIN boards b ON b.id = bs.board_id
		WHERE b.slug = $1`
	var boardSettings BoardSettings
	err := dbPool.QueryRow(
		context.Background(), query, boardSlug).
		Scan(
			&boardSettings.IsOfficial,
			&boardSettings.IsPostApprovalRequired,
		)
	if err != nil && err != pgx.ErrNoRows {
		return boardSettings, err
	}
	return boardSettings, nil
}

func SetBoardSettings(dbPool *pgxpool.Pool, boardSlug string, settings BoardSettings) error {
	const query = `INSERT INTO board_settings (is_official, is_post_approval_required)
		VALUES ($1, $2)
		ON CONFLICT (board_id) DO UPDATE
		SET is_official = $1, 
		    is_post_approval_required = $2,
		    updated_at = NOW()
		WHERE board_id = $3`
	_, err := dbPool.Exec(
		context.Background(), query, settings.IsOfficial, settings.IsPostApprovalRequired, boardSlug)
	return err
}
