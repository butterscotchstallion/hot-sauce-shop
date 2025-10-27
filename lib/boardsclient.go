package lib

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
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
	CreatedByUserSlug string     `json:"createdByUserSlug"`
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
	IsPinned          bool       `json:"isPinned"`
}

type AddPostRequest struct {
	Title      string                  `json:"title" form:"title"`
	ParentId   int                     `json:"parentId" form:"parentId"`
	PostText   string                  `json:"postText" form:"postText" binding:"required"`
	PostImages []*multipart.FileHeader `json:"postImages" form:"postImages"`
	Slug       string                  `json:"slug" form:"slug"`
}

func GetBoards(dbPool *pgxpool.Pool) ([]Board, error) {
	// TODO: filter visible boards, or show everything if privileged
	const query = `
		SELECT b.id, b.display_name, b.created_at, b.updated_at, b.slug, b.visible, 
		CASE WHEN b.thumbnail_filename IS NULL THEN '' ELSE b.thumbnail_filename END AS thumbnail_filename,
		CASE WHEN b.description IS NULL THEN '' ELSE b.description END AS description,
		b.created_by_user_id,
		u.username AS created_by_username,
		u.slug AS created_by_user_slug
		FROM boards b
		JOIN users u on u.id = b.created_by_user_id
		ORDER BY b.display_name
	`
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
		ORDER BY bp.is_pinned DESC, bp.created_at DESC
	`
}

func GetTotalPostReplyCountByBoardSlug(dbPool *pgxpool.Pool, boardSlug string) (map[int]int, error) {
	type totalPostReplyCountResult struct {
		Id                  int
		TotalPostReplyCount int
	}
	var replies []totalPostReplyCountResult
	boardSlugClause := ""
	if len(boardSlug) > 0 {
		boardSlugClause = " AND b.slug = $1 "
	}
	query := `
		SELECT bp.id,
			  COALESCE(COUNT(bp.*), 0) AS total_post_reply_count
		FROM board_posts bp
		JOIN boards b ON b.id = bp.board_id
		WHERE 1=1 ` + boardSlugClause + `  AND bp.parent_id > 0 GROUP BY bp.id
	`
	var rows pgx.Rows
	var err error
	if len(boardSlug) > 0 {
		rows, err = dbPool.Query(context.Background(), query, boardSlug)
	} else {
		rows, err = dbPool.Query(context.Background(), query)
	}
	if err != nil {
		return nil, err
	}
	replies, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[totalPostReplyCountResult])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}
	var totalPostReplyCountMap = make(map[int]int)
	for _, reply := range replies {
		totalPostReplyCountMap[reply.Id] = reply.TotalPostReplyCount
	}
	return totalPostReplyCountMap, nil
}

// GetPosts
// Gets posts, optionally filtered by boardSlug/postSlug
func GetPosts(dbPool *pgxpool.Pool, boardSlug string, postSlug string, parentId int) ([]BoardPost, error) {
	whereClause := ""
	if len(boardSlug) > 0 {
		whereClause += " AND b.slug = $1"
	}
	// If the post slug is here, there will also be a board slug
	if len(postSlug) > 0 {
		whereClause += " AND bp.slug = $2"
	}
	if parentId > 0 {
		whereClause += " AND bp.parent_id = $1"
	} else {
		whereClause += " AND bp.parent_id = 0"
	}
	query := getPostsQuery(whereClause)
	var rows pgx.Rows
	var err error
	if len(boardSlug) > 0 && len(postSlug) == 0 {
		rows, err = dbPool.Query(context.Background(), query, boardSlug)
	} else if len(postSlug) > 0 {
		rows, err = dbPool.Query(context.Background(), query, boardSlug, postSlug)
	} else if parentId > 0 {
		rows, err = dbPool.Query(context.Background(), query, parentId)
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
		       u.username AS created_by_username,
		       u.slug AS created_by_user_slug
		FROM boards b
		JOIN users u on u.id = b.created_by_user_id
		WHERE b.slug = $1
	`
	row, err := dbPool.Query(context.Background(), query, slug)
	if err != nil {
		logger.Error(fmt.Sprintf("Error running GetBoardBySlug query: %v", err))
		return Board{}, err
	}
	defer row.Close()
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
	defer row.Close()
	post, collectRowsErr := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[BoardPost])
	if collectRowsErr != nil {
		return BoardPost{}, collectRowsErr
	}
	return post, nil
}

func GetTotalPostsByBoardSlug(dbPool *pgxpool.Pool, boardSlug string) (int, error) {
	const query = `SELECT 
    	COUNT(bp.*) AS totalPosts
		FROM board_posts bp
		JOIN boards b ON b.id = bp.board_id
		WHERE b.slug = $1
	`
	row, err := dbPool.Query(context.Background(), query, boardSlug)
	if err != nil {
		return 0, err
	}
	type totalPostsResult struct {
		TotalPosts int
	}
	result, collectRowsErr := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[totalPostsResult])
	if collectRowsErr != nil {
		return 0, collectRowsErr
	}
	return result.TotalPosts, err
}

func GetBoardModerators(dbPool *pgxpool.Pool, boardSlug string, userId int) ([]User, error) {
	userFilterClause := ""
	if userId > 0 {
		userFilterClause = " AND urb.user_id = $2"
	}
	query := `SELECT u.*
		FROM users u
        JOIN user_roles_boards urb ON urb.user_id = u.id
		JOIN user_roles ur ON ur.role_id = urb.role_id
        JOIN roles r ON r.id = urb.role_id
        JOIN boards b ON b.id = urb.board_id
		WHERE 1=1
		AND b.slug = $1
		AND r.slug = 'message-board-moderator'
	` + userFilterClause
	var rows pgx.Rows
	var err error
	if userId > 0 {
		rows, err = dbPool.Query(context.Background(), query, boardSlug, userId)
	} else {
		rows, err = dbPool.Query(context.Background(), query, boardSlug)
	}
	if err != nil {
		return []User{}, err
	}
	moderators, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}
	return moderators, err
}

func PinBoardPost(dbPool *pgxpool.Pool, postSlug string) error {
	const query = `UPDATE board_posts SET is_pinned = true WHERE slug = $1`
	_, err := dbPool.Exec(context.Background(), query, postSlug)
	if err != nil {
		return err
	}
	return nil
}

func GetNumBoardMembers(dbPool *pgxpool.Pool, boardSlug string) (int, error) {
	const query = `SELECT COUNT(*) AS num_board_members 
		FROM boards_users bu
		JOIN boards b on b.id = bu.board_id
		WHERE b.slug = $1
	`
	var numBoardMembers int
	scanErr := dbPool.QueryRow(context.Background(), query, boardSlug).Scan(&numBoardMembers)
	if scanErr != nil {
		return 0, scanErr
	}
	return numBoardMembers, nil
}

func AddPost(dbPool *pgxpool.Pool, post AddPostRequest, userId int, boardId int) (int, error) {
	lastInsertId := 0
	const query = `
		INSERT INTO board_posts (title, created_by_user_id, board_id, parent_id, slug, post_text) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	insertErr := dbPool.QueryRow(
		context.Background(),
		query,
		post.Title,
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

func AddPostImages(dbPool *pgxpool.Pool, postId int, filename string, thumbnailFilename string) error {
	const query = `
		INSERT INTO board_posts_images (filename, board_post_id, thumbnail_filename) 
		VALUES ($1, $2, $3)
	`
	_, err := dbPool.Exec(context.Background(), query, filename, postId, thumbnailFilename)
	if err != nil {
		return err
	}
	return nil
}

type PostImages struct {
	filename          string
	thumbnailFilename string
}

func GetPostImages(dbPool *pgxpool.Pool, postId int) ([]PostImages, error) {
	const query = `
		SELECT 
		board_post_id, filename, thumbnail_filename 
		FROM board_posts_images
		WHERE board_post_id = $1
	`
	rows, err := dbPool.Query(context.Background(), query, postId)
	if err != nil {
		return nil, err
	}
	postImagesRows, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[PostImages])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}
	return postImagesRows, nil
}
