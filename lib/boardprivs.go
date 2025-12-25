package lib

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetBoardsByRole
// Returns boards the user has the specified role on
func GetBoardsByRole(dbPool *pgxpool.Pool, userId int, roleName string) ([]Board, error) {
	whereClause := ""

	// not user supplied - no need to bind here
	if len(roleName) > 0 {
		whereClause = fmt.Sprintf("AND r.name = '%s'", roleName)
	}

	query := fmt.Sprintf(`SELECT %s
		b.created_by_user_id,
		u.username AS created_by_username,
		u.slug AS created_by_user_slug
		FROM boards b
		JOIN user_roles_boards urb ON urb.board_id = b.id
		JOIN roles r on r.id = urb.role_id
		JOIN users u on urb.user_id = u.id
		WHERE urb.user_id = $1
		%s
	`, getBoardColumns(), whereClause)
	rows, err := dbPool.Query(context.Background(), query, userId)
	if err != nil {
		return nil, err
	}
	boards, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[Board])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}
	return boards, nil
}

func GetUserAdminBoards(dbPool *pgxpool.Pool, userId int) ([]Board, error) {
	return GetBoardsByRole(dbPool, userId, "Message Board Admin")
}

func AddBoardAdmin(dbPool *pgxpool.Pool, userId int, boardId int) error {
	const boardAdminRoleId = 7
	const query = `
		INSERT INTO user_roles_boards (user_id, board_id, role_id)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`
	_, err := dbPool.Exec(context.Background(), query, userId, boardId, boardAdminRoleId)
	if err != nil {
		return err
	}
	return nil
}
