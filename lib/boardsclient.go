package lib

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"strings"
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
	ThumbnailFilename string     `json:"thumbnailFilename" db:"thumbnail_filename"`
	ThumbnailWidth    *float32   `json:"thumbnailWidth" db:"thumbnail_width"`
	ThumbnailHeight   *float32   `json:"thumbnailHeight" db:"thumbnail_height"`
}

type AddPostRequest struct {
	Title        string                  `json:"title" form:"title"`
	ParentId     int                     `json:"parentId" form:"parentId"`
	PostText     string                  `json:"postText" form:"postText" binding:"required"`
	PostImages   []*multipart.FileHeader `json:"postImages" form:"postImages"`
	Slug         string                  `json:"slug" form:"slug"`
	PostFlairIds []int                   `json:"postFlairIds" form:"postFlairIds"`
}

type AddBoardRequest struct {
	DisplayName       string `json:"displayName" validate:"required,min=10,max=255"`
	ThumbnailFilename string `json:"thumbnailFilename"`
	Description       string `json:"description"`
}

type AddPostResponseResults struct {
	Post      BoardPost
	NewPostId int
}
type AddPostResponse struct {
	Status  string
	Message string
	Results AddPostResponseResults
}

type BoardPostDeleteResponse struct {
	Status  string
	Message string
}

type SavedPostImageInfo struct {
	Filename             string
	FullImagePath        string
	ThumbnailFilename    string
	ThumbnailFullPath    string
	MimeType             string
	ImageWidthHeight     ImageWidthHeight
	ThumbnailWidthHeight ImageWidthHeight
}

type AddBoardResponseResults struct {
	Slug        string `json:"slug"`
	DisplayName string `json:"displayName"`
	BoardId     int    `json:"boardId"`
}

type AddBoardResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Results AddBoardResponseResults
}

type BoardDetailResponseResults struct {
	Board              Board
	Moderators         []string `json:"moderators"`
	NumBoardModerators int      `json:"numBoardModerators"`
	TotalPosts         int      `json:"totalPosts"`
}
type BoardDetailResponse struct {
	Status  string `json:"status"`
	Results BoardDetailResponseResults
}

type BoardDeleteResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type PostFlair struct {
	Id          int    `json:"id"`
	DisplayName string `json:"displayName"`
	Slug        string `json:"slug"`
}
type PostFlairsResponseResults struct {
	PostFlairs []PostFlair
}
type PostFlairsResponse struct {
	Status  string `json:"status"`
	Results PostFlairsResponseResults
}

// PostsFlairs These are the association between each post and flairs
type PostsFlairs struct {
	Id          int       `json:"id"`
	BoardPostId int       `json:"boardPostId"`
	PostFlairId int       `json:"postFlairId"`
	CreatedAt   time.Time `json:"createdAt"`
}

type PostsFlairsResponseResults struct {
	PostsFlairs []PostsFlairs `json:"postFlairs"`
}

type PostsFlairsResponse struct {
	Status  string                     `json:"status"`
	Results PostsFlairsResponseResults `json:"results"`
}

type BoardPostResponseResults struct {
	Board           Board  `json:"board"`
	Moderators      []User `json:"moderators"`
	NumBoardMembers int    `json:"numBoardMembers"`
	TotalPosts      int    `json:"totalPosts"`
}

type BoardPostResponse struct {
	Status  string `json:"status"`
	Results BoardPostResponseResults
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
			COALESCE((
				SELECT bpi.thumbnail_filename
			  	FROM board_posts_images bpi
			  	WHERE bpi.board_post_id = bp.id
			  	ORDER BY bpi.id DESC
			  	LIMIT 1
			), '') AS thumbnail_filename,
			(SELECT bpi.thumbnail_width
			FROM board_posts_images bpi
			WHERE bpi.board_post_id = bp.id
			ORDER BY bpi.id DESC
			LIMIT 1) AS thumbnail_width,
		    (SELECT bpi.thumbnail_height
			FROM board_posts_images bpi
			WHERE bpi.board_post_id = bp.id
			ORDER BY bpi.id DESC
			LIMIT 1) AS thumbnail_height,
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
		WHERE 1=1 ` + boardSlugClause + `  
		AND bp.parent_id > 0
		GROUP BY bp.id
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
func GetPosts(dbPool *pgxpool.Pool, boardSlug string, postSlug string, parentId int, logger *slog.Logger) ([]BoardPost, error) {
	whereClause := ""
	if len(boardSlug) > 0 {
		whereClause += " AND b.slug = $1"
	}
	// If the slug for the post is here, there will also be a board slug
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
		logger.Info(fmt.Sprintf("GetPosts CollectRows Error: %v", collectRowsErr))

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

func GetBoardBySlug(dbPool *pgxpool.Pool, slug string) (Board, error) {
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
		return Board{}, err
	}
	defer row.Close()
	board, collectRowsErr := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[Board])
	if collectRowsErr != nil {
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
			COALESCE((SELECT SUM(v.value) FROM votes v WHERE v.post_id = bp.id), 0) AS voteSum,
			COALESCE((
				SELECT bpi.thumbnail_filename
			  	FROM board_posts_images bpi
			  	WHERE bpi.board_post_id = bp.id
			  	ORDER BY bpi.id DESC
			  	LIMIT 1
			), '') AS thumbnail_filename,
		    (SELECT bpi.thumbnail_width
			FROM board_posts_images bpi
			WHERE bpi.board_post_id = bp.id
			ORDER BY bpi.id DESC
			LIMIT 1) AS thumbnail_width,
		    (SELECT bpi.thumbnail_height
			FROM board_posts_images bpi
			WHERE bpi.board_post_id = bp.id
			ORDER BY bpi.id DESC
			LIMIT 1) AS thumbnail_height
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
		WHERE b.slug = $1
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

func AddPostImages(dbPool *pgxpool.Pool, postId int, imageInfo SavedPostImageInfo) error {
	const query = `
		INSERT INTO board_posts_images (
			filename,
		    board_post_id,
		    thumbnail_filename,
		    mime_type,
		    orig_width,
		    orig_height,
		    thumbnail_width,
		    thumbnail_height
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := dbPool.Exec(
		context.Background(),
		query,
		imageInfo.Filename,
		postId,
		imageInfo.ThumbnailFilename,
		imageInfo.MimeType,
		imageInfo.ImageWidthHeight.Width,
		imageInfo.ImageWidthHeight.Height,
		imageInfo.ThumbnailWidthHeight.Width,
		imageInfo.ThumbnailWidthHeight.Height,
	)
	if err != nil {
		return err
	}
	return nil
}

func AddBoard(dbPool *pgxpool.Pool, slug string, displayName string, thumbnailFilename string, createdByUserId int, description string) (int, error) {
	const query = `
		INSERT INTO boards (
			slug,
			display_name,
			created_at,
		    thumbnail_filename,
			created_by_user_id,
		    description
		)
		VALUES ($1, $2, NOW(), $3, $4, $5)
		RETURNING id
	`
	var boardId int
	err := dbPool.QueryRow(
		context.Background(),
		query,
		slug,
		displayName,
		thumbnailFilename,
		createdByUserId,
		description,
	).Scan(&boardId)
	if err != nil {
		return 0, err
	}
	return boardId, nil
}

func DeleteBoard(dbPool *pgxpool.Pool, boardSlug string) error {
	/**
	 * Board users must be deleted before the board can be deleted
	 * because of FK constraints. TODO: maybe use cascading here.
	 */
	deleteBoardUsersErr := DeleteBoardUsers(dbPool, boardSlug)
	if deleteBoardUsersErr != nil {
		return deleteBoardUsersErr
	}
	const query = `DELETE FROM boards WHERE slug = $1`
	_, err := dbPool.Exec(
		context.Background(),
		query,
		boardSlug,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteBoardUsers(dbPool *pgxpool.Pool, boardSlug string) error {
	board, boardErr := GetBoardBySlug(dbPool, boardSlug)
	if boardErr != nil {
		return boardErr
	}
	const query = `DELETE 
		FROM boards_users bu
       	WHERE board_id = $1
    `
	_, err := dbPool.Exec(
		context.Background(),
		query,
		board.Id,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteBoardPost(dbPool *pgxpool.Pool, boardPostSlug string) error {
	const query = `DELETE FROM board_posts WHERE slug = $1`
	_, err := dbPool.Exec(
		context.Background(),
		query,
		boardPostSlug,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteBoardPostFlairs(dbPool *pgxpool.Pool, boardPostId int) error {
	const query = `DELETE FROM posts_flairs WHERE board_post_id = $1`
	_, err := dbPool.Exec(
		context.Background(),
		query,
		boardPostId,
	)
	if err != nil {
		return err
	}

	return nil
}

func IsUserBoardPostAuthor(dbPool *pgxpool.Pool, userId int, boardPostSlug string) (bool, error) {
	const query = `
		SELECT COUNT(*) as postCount
		FROM board_posts bp
		WHERE bp.slug = $1
		AND bp.created_by_user_id = $2
	`
	var postCount int
	insertErr := dbPool.QueryRow(
		context.Background(),
		query,
		boardPostSlug,
		userId,
	).Scan(&postCount)
	if insertErr != nil {
		return false, insertErr
	}

	return postCount == 1, nil
}

func GetPostFlairs(dbPool *pgxpool.Pool, postId int) ([]PostFlair, error) {
	postIdFilterClause := "WHERE board_post_id = $1"
	query := fmt.Sprintf(`SELECT * FROM post_flairs %v`, postIdFilterClause)
	var rows pgx.Rows
	var rowsErr error
	if postId > 0 {
		rows, rowsErr = dbPool.Query(context.Background(), query, postId)
	} else {
		rows, rowsErr = dbPool.Query(context.Background(), query)
	}
	if rowsErr != nil {
		return []PostFlair{}, rowsErr
	}
	postFlairs, postFlairsErr := pgx.CollectRows(rows, pgx.RowToStructByName[PostFlair])
	if postFlairsErr != nil {
		return []PostFlair{}, postFlairsErr
	}

	return postFlairs, nil
}

// GetPostsFlairs association between each post and its flairs
func GetPostsFlairs(dbPool *pgxpool.Pool) ([]PostsFlairs, error) {
	const query = `SELECT * FROM posts_flairs`
	rows, rowsErr := dbPool.Query(context.Background(), query)
	if rowsErr != nil {
		return []PostsFlairs{}, rowsErr
	}
	postFlairs, postFlairsErr := pgx.CollectRows(rows, pgx.RowToStructByName[PostsFlairs])
	if postFlairsErr != nil {
		return []PostsFlairs{}, postFlairsErr
	}

	return postFlairs, nil
}

func AddPostFlair(dbPool *pgxpool.Pool, postId int, postFlairIds []int) error {
	postFlairsDeletedErr := DeleteBoardPostFlairs(dbPool, postId)
	if postFlairsDeletedErr != nil {
		return postFlairsDeletedErr
	}
	var values []string
	for flairId := range postFlairIds {
		values = append(values, fmt.Sprintf("(%d, %d)", postId, flairId))
	}
	query := fmt.Sprintf(`
		INSERT INTO posts_flairs (board_post_id, post_flair_id) VALUES %v`, strings.Join(values, ","),
	)
	_, insertErr := dbPool.Query(context.Background(), query)
	if insertErr != nil {
		return insertErr
	}

	return nil
}

func GetPostFlairIdMap(postFlairs []PostFlair) map[int]PostFlair {
	postFlairIdMap := make(map[int]PostFlair)
	for _, postFlair := range postFlairs {
		postFlairIdMap[postFlair.Id] = postFlair
	}

	return postFlairIdMap
}

func GetPostsFlairsMap(postsFlairs []PostsFlairs, postFlairIdMap map[int]PostFlair) map[int][]PostFlair {
	postsFlairsMap := make(map[int][]PostFlair)
	for _, postFlair := range postsFlairs {
		// Initialize slice if we haven't already
		if _, exists := postsFlairsMap[postFlair.BoardPostId]; !exists {
			postsFlairsMap[postFlair.BoardPostId] = make([]PostFlair, 0)
		}
		// Ensure that this flair exists in the flair map, and if so,
		// append it to the slice
		if _, exists := postFlairIdMap[postFlair.PostFlairId]; exists {
			postsFlairsMap[postFlair.BoardPostId] = append(
				postsFlairsMap[postFlair.BoardPostId],
				postFlairIdMap[postFlair.PostFlairId],
			)
		}
	}

	return postsFlairsMap
}
