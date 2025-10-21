package lib

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddPost(dbPool *pgxpool.Pool, post AddPostRequest, userId int, boardId int) (int, error) {
	lastInsertId := 0
	const query = `
		INSERT INTO board_posts (title, thumbnail_filename, created_by_user_id, board_id, parent_id, slug, post_text) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	insertErr := dbPool.QueryRow(
		context.Background(),
		query,
		post.Title,
		post.ThumbnailFilename,
		userId,
		boardId,
		post.ParentId,
		post.Slug,
		post.PostText,
	).Scan(&lastInsertId)
	if insertErr != nil {
		return 0, insertErr
	}
	return lastInsertId, nil
}

func AddUpdateVote(dbPool *pgxpool.Pool, userId int, postId int, voteValue int) (int, error) {
	if voteValue != -1 && voteValue != 1 {
		return 0, errors.New("invalid vote value")
	}
	lastInsertId := 0
	const query = `
		INSERT INTO votes (user_id, post_id, value) 
		VALUES ($1, $2, $3)
		ON CONFLICT(user_id, post_id)
		    DO UPDATE SET value = $3, updated_at = NOW()
		RETURNING id
	`
	insertErr := dbPool.QueryRow(context.Background(), query, userId, postId, voteValue).Scan(&lastInsertId)
	if insertErr != nil {
		return 0, insertErr
	}
	return lastInsertId, nil
}

// GetUserVoteMap
// Returns a map of post id to vote value
// TODO: add filter for board_id
func GetUserVoteMap(dbPool *pgxpool.Pool, userId int) (map[int]int, error) {
	voteMap := make(map[int]int)
	const query = `
		SELECT post_id, value FROM votes WHERE user_id = $1
	`
	rows, err := dbPool.Query(context.Background(), query, userId)
	if err != nil {
		return voteMap, err
	}
	defer rows.Close()
	return voteMap, nil
}
