package lib

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Role struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
	Slug        string    `json:"slug"`
	ColorClass  string    `json:"colorClass"`
	Description string    `json:"description"`
}

type Permission struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Slug      string    `json:"slug"`
}

const UserRoleMessageBoardAdmin = "Message Board Admin"
const UserRoleSuperAdmin = "Super Message Board Admin"
const UserRoleUserAdmin = "User Admin"
const UserRoleMessageBoardModerator = "Message Board Moderator"

// UpdateUserRoles
// - Delete existing user roles
// - Add new roles
func UpdateUserRoles(dbPool *pgxpool.Pool, logger *slog.Logger, userId int, roleIds []int) (bool, error) {
	_, rolesDeletedErr := deleteUserRoles(dbPool, userId)
	if rolesDeletedErr != nil {
		logger.Error(fmt.Sprintf("Error deleting user roles: %v", rolesDeletedErr))
		return false, rolesDeletedErr
	}
	if len(roleIds) == 0 {
		logger.Error("No roles provided")
		return true, nil
	}
	for _, roleId := range roleIds {
		const query = `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`
		_, insertRolesErr := dbPool.Exec(context.Background(), query, userId, roleId)
		if insertRolesErr != nil {
			return false, insertRolesErr
		}
	}
	return true, nil
}

func deleteUserRoles(dbPool *pgxpool.Pool, userId int) (bool, error) {
	const query = `DELETE FROM user_roles WHERE user_id = $1`
	_, err := dbPool.Exec(context.Background(), query, userId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetRoleList(dbPool *pgxpool.Pool, logger *slog.Logger) ([]Role, error) {
	const query = `
		SELECT *
		FROM roles r
		ORDER BY r.name
	`
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting role list: %v", err))
		return nil, err
	}
	roles, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[Role])
	if collectRowsErr != nil {
		logger.Error(fmt.Sprintf("Error collecting roles: %v", collectRowsErr))
		return nil, collectRowsErr
	}
	return roles, nil
}

func GetRolesByUserId(dbPool *pgxpool.Pool, logger *slog.Logger, userId int) ([]Role, error) {
	const query = `
		SELECT r.*
		FROM roles r
		LEFT JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
	`
	rows, err := dbPool.Query(context.Background(), query, userId)
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting roles by user id: %v", err))
		return nil, err
	}
	roles, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[Role])
	if collectRowsErr != nil {
		logger.Error(fmt.Sprintf("Error collecting roles by user id: %v", collectRowsErr))
		return nil, collectRowsErr
	}
	return roles, nil
}

// GetUserRoleByRoleName TODO: replace UserHasRole with this
func GetUserRoleByRoleName(dbPool *pgxpool.Pool, userId int, roleName string) (bool, error) {
	const query = `
		SELECT COUNT(*) AS roleCount
		FROM roles r
		LEFT JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
		AND r.name = $2
	`
	var rowCount = 0
	row := dbPool.QueryRow(context.Background(), query, userId, roleName)
	err := row.Scan(&rowCount)
	if err != nil {
		return false, err
	}
	return rowCount > 0, nil
}

func IsSignedInAndUserExists(c *gin.Context, dbPool *pgxpool.Pool, logger *slog.Logger) (int, error) {
	sessionIdCookieValue, err := c.Cookie("sessionId")
	if err != nil || sessionIdCookieValue == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ERROR",
			"message": "No session ID found",
		})
		return 0, err
	}

	user, getUserErr := GetUserBySessionId(dbPool, logger, sessionIdCookieValue)
	if getUserErr != nil || user == (User{}) {
		logger.Error(fmt.Sprintf("Error fetching user: %v", getUserErr))
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "ERROR",
			"message": "No user found for session ID",
		})
		return 0, getUserErr
	}

	return user.Id, nil
}

func UserHasRole(c *gin.Context, dbPool *pgxpool.Pool, logger *slog.Logger, roleName string) (bool, error) {
	userId, userErr := IsSignedInAndUserExists(c, dbPool, logger)
	if userErr != nil {
		return false, userErr
	}

	roles, rolesErr := GetRolesByUserId(dbPool, logger, userId)
	if rolesErr != nil {
		logger.Error(fmt.Sprintf("Error fetching roles: %v", rolesErr.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "ERROR",
			"message": rolesErr.Error(),
		})
		return false, rolesErr
	}

	for _, role := range roles {
		if role.Name == roleName {
			return true, nil
		}
	}

	return false, nil
}

func IsUserAdmin(c *gin.Context, dbPool *pgxpool.Pool, logger *slog.Logger) (bool, error) {
	return UserHasRole(c, dbPool, logger, UserRoleUserAdmin)
}

// IsSuperMessageBoardAdmin Sends JSON response upon failure
func IsSuperMessageBoardAdmin(c *gin.Context, dbPool *pgxpool.Pool, logger *slog.Logger) (bool, error) {
	return UserHasRole(c, dbPool, logger, UserRoleSuperAdmin)
}

// IsMessageBoardAdmin - checks admin status of a specific board
func IsMessageBoardAdmin(dbPool *pgxpool.Pool, boardId int, userId int) (bool, error) {
	const query = `SELECT COUNT(*) 
		FROM user_roles_boards bu
		JOIN roles r ON r.id = bu.role_id
		WHERE bu.board_id = $1
		AND bu.user_id = $2
		AND r.name = 'Message Board Admin'`
	row := dbPool.QueryRow(context.Background(), query, boardId, userId)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetRoleIdsFromRoles(roles []Role) []int {
	var roleIds []int
	for _, role := range roles {
		roleIds = append(roleIds, role.Id)
	}
	return roleIds
}
