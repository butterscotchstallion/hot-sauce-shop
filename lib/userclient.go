package lib

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id             int       `json:"id"`
	Username       string    `json:"username"`
	Password       string    `json:"-"`
	AvatarFilename string    `json:"avatarFilename"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func VerifyUsernameAndPassword(dbPool *pgxpool.Pool, logger *slog.Logger, username string, password string) (bool, error) {
	var dbPassword string
	const query = `SELECT password FROM users WHERE username = $1`
	err := dbPool.QueryRow(context.Background(), query, username).Scan(&dbPassword)
	noRowsReturned := errors.Is(err, pgx.ErrNoRows)

	if noRowsReturned {
		logger.Error(fmt.Sprintf("No user found with username: %v", username))
		return false, nil
	}

	// Includes err no rows
	if err != nil {
		logger.Error(fmt.Sprintf("Error running user query: %v", err))
		return false, err
	}

	passwordMatch := VerifyPassword(password, dbPassword)

	if !passwordMatch {
		logger.Error(fmt.Sprintf("Passwords do not match for username: %v", username))
	}

	return passwordMatch, nil
}

func UserIdExists(dbPool *pgxpool.Pool, id int) (bool, error) {
	exists := false
	const query = `
		SELECT EXISTS (
			SELECT 1 FROM users WHERE id = $1
        )
	`
	err := dbPool.QueryRow(context.Background(), query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func IsValidAccountNameAndPassword(username string, password string) bool {

	return true
}

func GetUserById(dbPool *pgxpool.Pool, id int) (User, error) {
	const query = `
		SELECT 
		id, 
		username,
		password,
		avatar_filename AS avatarFilename,
		created_at AS createdAt,
		updated_at AS updatedAt
		FROM users WHERE id = @id
	`
	row := dbPool.QueryRow(context.Background(), query, id)
	var user User
	err := row.Scan(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

/*
GetUserIdFromSession
Placeholder until user system is implemented
*/
func GetUserIdFromSession() int {
	return 1
}
