package lib

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Board struct {
	Id                int        `json:"id"`
	DisplayName       string     `json:"displayName"`
	CreatedAt         *time.Time `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	Slug              string     `json:"slug"`
	Visible           bool       `json:"visible"`
	ThumbnailFilename string     `json:"thumbnailFilename"`
	CreatedByUserId   int        `json:"createdByUserId"`
	CreatedByUsername string     `json:"createdByUsername"`
	Description       string     `json:"description"`
}

type BoardPost struct {
	Id                int        `json:"id"`
	Title             string     `json:"title"`
	CreatedAt         *time.Time `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	Slug              string     `json:"slug"`
	ThumbnailFilename string     `json:"thumbnailFilename"`
	CreatedByUserId   int        `json:"createdByUserId"`
	CreatedByUsername string     `json:"createdByUsername"`
	CreatedByUserSlug string     `json:"createdByUserSlug"`
	BoardId           int        `json:"boardId"`
	BoardSlug         string     `json:"boardSlug"`
	BoardName         string     `json:"boardName"`
	ParentId          int        `json:"parentId"`
	PostText          string     `json:"postText"`
	VoteSum           int        `json:"voteSum"`
}

type AddPostRequest struct {
	Title             string `json:"title"`
	Slug              string `json:"slug"`
	ThumbnailFilename string `json:"thumbnailFilename"`
	ParentId          int    `json:"parentId"`
	PostText          string `json:"postText"`
}

func GetBoards(dbPool *pgxpool.Pool) ([]Board, error) {
	// TODO: filter visible boards, or show everything if privileged
	const query = `
		SELECT b.id, b.display_name, b.created_at, b.updated_at, b.slug, b.visible, 
		CASE WHEN b.thumbnail_filename IS NULL THEN '' ELSE b.thumbnail_filename END AS thumbnail_filename,
		CASE WHEN b.description IS NULL THEN '' ELSE b.description END AS description,
		b.created_by_user_id,
		u.username AS created_by_username
		FROM boards b
		JOIN users u on u.id = b.created_by_user_id
		ORDER BY b.display_name
	`
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	boards, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[Board])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}
	return boards, nil
}

func getPostsQuery(whereClause string) string {
	return `
		SELECT 
		    bp.*,
			u.username AS created_by_username,
			u.slug AS created_by_user_slug,
			b.display_name AS boardName,
			b.slug AS boardSlug,
			COALESCE((SELECT SUM(v.value) FROM votes v WHERE v.post_id = bp.id), 0) AS voteSum
		FROM board_posts bp
		JOIN users u on u.id = bp.created_by_user_id
		JOIN boards b ON b.id = bp.board_id
		WHERE 1=1
		` + whereClause + `
		ORDER BY bp.created_at DESC
	`
}

func GetPostReplies(dbPool *pgxpool.Pool, parentId int) ([]BoardPost, error) {
	whereClause := "WHERE bp.parent_id = $1"
	query := getPostsQuery(whereClause)
	var rows pgx.Rows
	var err error
	rows, err = dbPool.Query(context.Background(), query, parentId)
	if err != nil {
		return nil, err
	}
	posts, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[BoardPost])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}
	return posts, nil
}

// GetPosts
// Gets posts, optionally filtered by boardSlug/postSlug
func GetPosts(dbPool *pgxpool.Pool, boardSlug string, postSlug string, logger *slog.Logger) ([]BoardPost, error) {
	whereClause := ""
	if len(boardSlug) > 0 {
		whereClause += " AND b.slug = $1"
	}
	// If the post slug is here, there will also be a board slug
	if len(postSlug) > 0 {
		whereClause += " AND bp.slug = $2"
	}
	query := getPostsQuery(whereClause)
	logger.Info(fmt.Sprintf("GetPosts query: %s", query))
	var rows pgx.Rows
	var err error
	if len(boardSlug) > 0 && len(postSlug) == 0 {
		rows, err = dbPool.Query(context.Background(), query, boardSlug)
		logger.Info("GetPosts query boardSlug")
	} else if len(postSlug) > 0 {
		rows, err = dbPool.Query(context.Background(), query, boardSlug, postSlug)
		logger.Info("GetPosts query boardSlug and postSlug")
	} else {
		rows, err = dbPool.Query(context.Background(), query)
	}
	if err != nil {
		return nil, err
	}
	posts, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[BoardPost])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}
	return posts, nil
}

func GetNumPostsByUserId(dbPool *pgxpool.Pool, userId int) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM board_posts
		WHERE created_by_user_id = $1
	`
	var count int
	err := dbPool.QueryRow(context.Background(), query, userId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetBoardBySlug(dbPool *pgxpool.Pool, logger *slog.Logger, slug string) (Board, error) {
	const query = `
		SELECT b.*,
		       u.username AS created_by_username
		FROM boards b
		JOIN users u on u.id = b.created_by_user_id
		WHERE b.slug = $1
	`
	row, err := dbPool.Query(context.Background(), query, slug)
	if err != nil {
		logger.Error(fmt.Sprintf("Error running GetBoardBySlug query: %v", err))
		return Board{}, err
	}
	board, collectRowsErr := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[Board])
	if collectRowsErr != nil {
		logger.Error(fmt.Sprintf("GetBoardBySlug: error collecting board: %v", collectRowsErr))
		return Board{}, collectRowsErr
	}
	return board, nil
}

func GetPostDetail(dbPool *pgxpool.Pool, boardSlug string, postSlug string) (BoardPost, error) {
	query := `
		SELECT 
		    bp.*,
			u.username AS created_by_username,
			u.slug AS created_by_user_slug,
			b.display_name AS boardName,
			b.slug AS boardSlug,
			COALESCE((SELECT SUM(v.value) FROM votes v WHERE v.post_id = bp.id), 0) AS voteSum
		FROM board_posts bp
		JOIN users u on u.id = bp.created_by_user_id
		JOIN boards b ON b.id = bp.board_id
		WHERE b.slug = $1
		AND bp.slug = $2
	`
	row, err := dbPool.Query(context.Background(), query, boardSlug, postSlug)
	if err != nil {
		return BoardPost{}, err
	}
	post, collectRowsErr := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[BoardPost])
	if collectRowsErr != nil {
		return BoardPost{}, collectRowsErr
	}
	return post, nil
}
