package lib

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Role struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"createdAt"`
	Slug       string    `json:"slug"`
	ColorClass string    `json:"colorClass"`
}

type Permission struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Slug      string    `json:"slug"`
}

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
