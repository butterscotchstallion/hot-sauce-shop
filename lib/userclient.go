package lib

import (
	"context"
	"time"

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

func UserExists(dbPool *pgxpool.Pool, id int) (bool, error) {
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
	return true, nil
}

/*
GetUserIdFromSession
Placeholder until user system is implemented
*/
func GetUserIdFromSession() int {
	return 1
}
