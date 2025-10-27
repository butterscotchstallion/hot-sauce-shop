package lib

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

type VoteMap struct {
	PostId int
	Value  int
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
	votes, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[VoteMap])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}
	defer rows.Close()
	for vote := range votes {
		voteMap[votes[vote].PostId] = votes[vote].Value
	}
	return voteMap, nil
}

func GetVoteSumMapByPostId(dbPool *pgxpool.Pool, postId int) (map[int]int, error) {
	voteSumMap := make(map[int]int)
	const query = `
		SELECT post_id, COALESCE((SELECT SUM(v.value) FROM votes v WHERE v.post_id = $1), 0) AS value
		FROM votes
		WHERE post_id = $1
		GROUP BY post_id
	`
	rows, err := dbPool.Query(context.Background(), query, postId)
	if err != nil {
		return voteSumMap, err
	}
	votes, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[VoteMap])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}
	defer rows.Close()
	for vote := range votes {
		voteSumMap[votes[vote].PostId] = votes[vote].Value
	}
	return voteSumMap, nil
}
