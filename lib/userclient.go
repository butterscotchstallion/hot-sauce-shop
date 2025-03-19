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

func VerifyUsernameAndPassword(dbPool *pgxpool.Pool, logger *slog.Logger, username string, password string) (User, error) {
	var user User
	const query = `SELECT id, password FROM users WHERE username = $1`
	row := dbPool.QueryRow(context.Background(), query, username)
	err := row.Scan(&user.Id, &user.Password)
	noRowsReturned := errors.Is(err, pgx.ErrNoRows)

	logger.Info(fmt.Sprintf("Verifying username and password for username: %v", username))

	if noRowsReturned {
		logger.Error(fmt.Sprintf("No user found with username: %v", username))
		return user, nil
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Error running user query: %v", err))
		return user, err
	}

	passwordMatch := VerifyPassword(password, user.Password)
	if !passwordMatch {
		logger.Error(fmt.Sprintf("Passwords do not match for username: %v", username))
	}

	return user, nil
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

/*
GetUserBySessionId
- Filter non-expired sessions too
*/
func GetUserBySessionId(dbPool *pgxpool.Pool, sessionId string) (User, error) {
	const query = `
		SELECT 
		u.id, 
		u.username,
		u.password,
		u.avatar_filename AS avatarFilename,
		u.created_at AS createdAt,
		u.updated_at AS updatedAt
		FROM users u
		JOIN user_sessions s ON u.id = s.user_id
		WHERE 1=1
		AND s.enabled = true
		AND s.created_at < NOW() - INTERVAL '30 day'
		AND s.session_id = $1
	`
	row := dbPool.QueryRow(context.Background(), query, sessionId)
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
