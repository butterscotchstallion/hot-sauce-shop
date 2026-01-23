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

func IsUserBoardModerator(dbPool *pgxpool.Pool, boardSlug string, userId int) (bool, error) {
	mods, err := GetBoardModerators(dbPool, boardSlug, userId)
	return len(mods) > 0, err
}

func GetBoardUsersByRole(dbPool *pgxpool.Pool, boardSlug string, userId int, roleName string) ([]User, error) {
	userFilterClause := ""
	if userId > 0 {
		userFilterClause = " AND urb.user_id = $3"
	}
	query := `SELECT u.*
		FROM users u
        JOIN user_roles_boards urb ON urb.user_id = u.id
		JOIN user_roles ur ON ur.role_id = urb.role_id
        JOIN roles r ON r.id = urb.role_id
        JOIN boards b ON b.id = urb.board_id
		WHERE b.slug = $1
		AND r.name = $2
	` + userFilterClause
	var rows pgx.Rows
	var err error
	if userId > 0 {
		rows, err = dbPool.Query(context.Background(), query, boardSlug, roleName, userId)
	} else {
		rows, err = dbPool.Query(context.Background(), query, boardSlug, roleName)
	}
	if err != nil {
		return []User{}, err
	}
	moderators, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}
	return moderators, err
}

func GetBoardModerators(dbPool *pgxpool.Pool, boardSlug string, userId int) ([]User, error) {
	return GetBoardUsersByRole(dbPool, boardSlug, userId, UserRoleMessageBoardModerator)
}

func GetBoardAdmins(dbPool *pgxpool.Pool, boardSlug string) ([]User, error) {
	return GetBoardUsersByRole(dbPool, boardSlug, 0, UserRoleMessageBoardAdmin)
}
