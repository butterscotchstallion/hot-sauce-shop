package lib

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id             int        `json:"id"`
	Slug           string     `json:"slug"`
	Username       string     `json:"username"`
	Password       string     `json:"-"`
	AvatarFilename string     `json:"avatarFilename"`
	CreatedAt      *time.Time `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}

func GetUsers(dbPool *pgxpool.Pool, logger *slog.Logger) ([]User, error) {
	const query = `
		SELECT *
		FROM users u
		ORDER BY u.username
	`
	var users []User
	row, err := dbPool.Query(context.Background(), query)
	if err != nil {
		logger.Error(fmt.Sprintf("Error running GetUsers query: %v", err))
		return users, err
	}
	users, collectRowsErr := pgx.CollectRows(row, pgx.RowToStructByName[User])
	if collectRowsErr != nil {
		logger.Error(fmt.Sprintf("GetUsers: error collecting users: %v", collectRowsErr))
		return users, err
	}
	return users, nil
}

func VerifyUsernameAndPasswordAndReturnUser(dbPool *pgxpool.Pool, logger *slog.Logger, username string, password string) (User, error) {
	const query = `SELECT * FROM users WHERE username = $1`
	row, err := dbPool.Query(context.Background(), query, username)
	if err != nil {
		return User{}, err
	}
	user, collectUserErr := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[User])
	noRowsReturned := errors.Is(err, pgx.ErrNoRows)

	logger.Info(fmt.Sprintf("Verifying username and password for username: %v", username))

	if noRowsReturned {
		logger.Error(fmt.Sprintf("No user found with username: %v", username))
		return user, nil
	}

	if collectUserErr != nil {
		logger.Error(fmt.Sprintf("Error collecting user row: %v", collectUserErr))
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
func GetUserBySessionId(dbPool *pgxpool.Pool, logger *slog.Logger, sessionId string) (User, error) {
	const query = `
		SELECT 
		u.id, 
		u.slug,
		u.username,
		u.password,
		u.avatar_filename AS avatarFilename,
		u.created_at AS createdAt,
		u.updated_at AS updatedAt
		FROM users u
		JOIN user_sessions s ON u.id = s.user_id
		WHERE 1=1
		AND s.enabled = true
		AND s.created_at >= DATE_TRUNC('month', current_date - interval '1' month)
		AND s.session_id = $1
	`
	row, err := dbPool.Query(context.Background(), query, sessionId)
	if err != nil {
		logger.Error(fmt.Sprintf("Error running session query: %v", err))
		return User{}, err
	}
	user, collectRowsErr := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[User])
	if collectRowsErr != nil {
		logger.Error(fmt.Sprintf("GetUserBySessionId: error collecting user: %v", collectRowsErr))
		return User{}, err
	}
	return user, nil
}

func GetUserIdFromSession(c *gin.Context, dbPool *pgxpool.Pool, logger *slog.Logger) (int, error) {
	sessionIdCookieValue, err := c.Cookie("sessionId")

	if err != nil || sessionIdCookieValue == "" {
		return 0, err
	}

	user, getUserErr := GetUserBySessionId(dbPool, logger, sessionIdCookieValue)

	if getUserErr != nil || user == (User{}) {
		return 0, getUserErr
	}

	return user.Id, nil
}

func GetUserBySlug(dbPool *pgxpool.Pool, logger *slog.Logger, slug string) (User, error) {
	const query = `
		SELECT 
		u.id, 
		u.slug,
		u.username,
		u.password,
		u.avatar_filename AS avatarFilename,
		u.created_at AS createdAt,
		u.updated_at AS updatedAt
		FROM users u
		WHERE 1=1
		AND u.slug = $1
	`
	row, err := dbPool.Query(context.Background(), query, slug)
	if err != nil {
		logger.Error(fmt.Sprintf("Error running getUserBySlug query: %v", err))
		return User{}, err
	}
	user, collectRowsErr := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[User])
	if collectRowsErr != nil {
		logger.Error(fmt.Sprintf("getUserBySlug: error collecting user: %v", collectRowsErr))
		return User{}, err
	}
	return user, nil
}
