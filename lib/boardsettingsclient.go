package lib

import (
	"context"
	"errors"
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
	UpdatedAt              time.Time
	BoardId                int
}

type BoardSettingsUpdateRequest struct {
	IsOfficial             bool `json:"isOfficial"`
	IsPostApprovalRequired bool `json:"isPostApprovalRequired"`
	BoardId                int  `json:"boardId"`
}

func GetBoardSettings(dbPool *pgxpool.Pool, boardSlug string) (BoardSettings, error) {
	const query = `SELECT 
		bs.is_official AS isOfficial, 
		bs.is_post_approval_required AS isPostApprovalRequired,
		bs.updated_at AS updatedAt,
		bs.board_id AS boardId
		FROM board_settings bs
		JOIN boards b ON b.id = bs.board_id
		WHERE b.slug = $1`
	var boardSettings BoardSettings
	err := dbPool.QueryRow(
		context.Background(), query, boardSlug).
		Scan(
			&boardSettings.IsOfficial,
			&boardSettings.IsPostApprovalRequired,
			&boardSettings.UpdatedAt,
			&boardSettings.BoardId,
		)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return boardSettings, err
	}
	return boardSettings, nil
}

func SetBoardSettings(dbPool *pgxpool.Pool, settings BoardSettingsUpdateRequest) error {
	const query = `INSERT INTO board_settings (is_official, is_post_approval_required, updated_at, board_id)
		VALUES ($1, $2, NOW(), $3)
		ON CONFLICT (board_id) DO UPDATE
		SET is_official = $1, 
		    is_post_approval_required = $2,
		    updated_at = NOW()
		WHERE board_settings.board_id = $3`
	_, err := dbPool.Exec(
		context.Background(), query, settings.IsOfficial, settings.IsPostApprovalRequired, settings.BoardId)
	return err
}
