package lib

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Role struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Slug      string    `json:"slug"`
}

type Permission struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Slug      string    `json:"slug"`
}

func GetRoleList(dbPool *pgxpool.Pool, logger *slog.Logger) ([]Role, error) {
	const query = `
		SELECT r.id, r.name, r.created_at, r.slug
		FROM roles r
	`
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		logger.Error("Error getting role list: %v", err)
		return nil, err
	}
	roles, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[Role])
	if collectRowsErr != nil {
		logger.Error("Error collecting roles: %v", collectRowsErr)
		return nil, collectRowsErr
	}
	return roles, nil
}

func GetRolesByUserId(dbPool *pgxpool.Pool, logger *slog.Logger, userId int) ([]Role, error) {
	const query = `
		SELECT r.id, r.name, r.created_at, r.slug
		FROM roles r
		LEFT JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
	`
	rows, err := dbPool.Query(context.Background(), query, userId)
	if err != nil {
		logger.Error("Error getting roles by user id: %v", err)
		return nil, err
	}
	roles, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[Role])
	if collectRowsErr != nil {
		logger.Error("Error collecting roles by user id: %v", collectRowsErr)
		return nil, collectRowsErr
	}
	return roles, nil
}
