package lib

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserSession struct {
	Id        int    `json:"id"`
	UserId    int    `json:"userId"`
	SessionId string `json:"sessionId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Enabled   bool   `json:"enabled"`
}

func GenerateUserSessionId() (string, error) {
	sessionId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return sessionId.String(), nil
}

func AddUserSessionId(dbPool *pgxpool.Pool, userId int) (string, error) {
	sessionId, sessionErr := GenerateUserSessionId()
	if sessionErr != nil {
		return "", sessionErr
	}
	const query = `
		INSERT INTO user_sessions (user_id, session_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT(user_id)
		    DO UPDATE SET session_id = $2, updated_at = NOW()
	`
	_, err := dbPool.Exec(context.Background(), query, userId, sessionId)
	if err != nil {
		return "", err
	}
	return sessionId, nil
}
