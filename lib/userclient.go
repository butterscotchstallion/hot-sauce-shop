package lib

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Id             int       `json:"id"`
	Username       string    `json:"username"`
	Password       string    `json:"-"`
	AvatarFilename string    `json:"avatarFilename"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func GetUserById(c *pgx.Conn, id int) (User, error) {
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
	row := c.QueryRow(context.Background(), query, id)
	var user User
	err := row.Scan(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func UserExists(c *pgx.Conn, id int) (bool, error) {
	exists := false
	const query = `
		SELECT EXISTS (
			SELECT 1 FROM users WHERE id = $1
        )
	`
	err := c.QueryRow(context.Background(), query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return true, nil
}
