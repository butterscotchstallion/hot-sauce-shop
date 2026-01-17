package lib

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"strings"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Board struct {
	Id                     int        `json:"id"`
	DisplayName            string     `json:"displayName"`
	CreatedAt              *time.Time `json:"createdAt"`
	UpdatedAt              *time.Time `json:"updatedAt"`
	Slug                   string     `json:"slug"`
	ThumbnailFilename      string     `json:"thumbnailFilename"`
	CreatedByUserId        int        `json:"createdByUserId"`
	CreatedByUsername      string     `json:"createdByUsername"`
	CreatedByUserSlug      string     `json:"createdByUserSlug"`
	Description            string     `json:"description"`
	IsVisible              bool       `json:"isVisible"`
	IsPrivate              bool       `json:"isPrivate"`
	IsOfficial             bool       `json:"isOfficial"`
	IsPostApprovalRequired bool       `json:"isPostApprovalRequired"`
	MinKarmaRequiredToPost int        `json:"minKarmaRequiredToPost"`
}

type BoardListResponseResults struct {
	Boards []Board `json:"boards"`
}
type BoardListResponse struct {
	Status  string                   `json:"status"`
	Results BoardListResponseResults `json:"results"`
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
	BoardIsOfficial   bool       `json:"boardIsOfficial"`
	ParentId          int        `json:"parentId"`
	PostText          string     `json:"postText"`
	VoteSum           int        `json:"voteSum"`
	IsPinned          bool       `json:"isPinned"`
	IsApproved        bool       `json:"isApproved"`
	ThumbnailFilename string     `db:"thumbnail_filename"      json:"thumbnailFilename"`
	ThumbnailWidth    *float32   `db:"thumbnail_width"         json:"thumbnailWidth"`
	ThumbnailHeight   *float32   `db:"thumbnail_height"        json:"thumbnailHeight"`
}

type AddPostRequest struct {
	Title        string                  `json:"title" form:"title"`
	ParentSlug   string                  `json:"parentSlug" form:"parentSlug"`
	PostText     string                  `json:"postText" form:"postText" binding:"required"`
	PostImages   []*multipart.FileHeader `json:"postImages" form:"postImages"`
	Slug         string                  `json:"slug" form:"slug"`
	PostFlairIds []int                   `form:"postFlairIds" json:"postFlairIds" `
}

type AddBoardRequest struct {
	DisplayName            string `json:"displayName"            validate:"required,min=10,max=255"`
	ThumbnailFilename      string `json:"thumbnailFilename"`
	Description            string `json:"description"`
	IsPostApprovalRequired bool   `json:"isPostApprovalRequired"`
	IsPrivate              bool   `json:"isPrivate"`
	IsOfficial             bool   `json:"isOfficial"`
	IsVisible              bool   `json:"isVisible"`
	MinKarmaRequiredToPost int    `json:"minKarmaRequiredToPost"`
}

type AddPostResponseResults struct {
	Post        BoardPost `json:"post"`
	NewPostId   int       `json:"newPostId"`
	NewPostSlug string    `json:"newPostSlug"`
}
type AddPostResponse struct {
	Status    string                 `json:"status"`
	Message   string                 `json:"message"`
	ErrorCode string                 `json:"errorCode"`
	Results   AddPostResponseResults `json:"results"`
}

type BoardPostDeleteResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
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
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Results AddBoardResponseResults `json:"results"`
}

type BoardDetailResponseResults struct {
	Board              Board  `json:"board"`
	Moderators         []User `json:"moderators"`
	NumBoardModerators int    `json:"numBoardModerators"`
	TotalPosts         int    `json:"totalPosts"`
	NumBoardMembers    int    `json:"numBoardMembers"`
}
type BoardDetailResponse struct {
	Status  string                     `json:"status"`
	Results BoardDetailResponseResults `json:"results"`
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
	PostFlairs []PostFlair `json:"postFlairs"`
}
type PostFlairsResponse struct {
	Status  string                    `json:"status"`
	Results PostFlairsResponseResults `json:"results"`
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
	Status  string                   `json:"status"`
	Results BoardPostResponseResults `json:"results"`
}

type PostDetailResponseResults struct {
	Post       BoardPost   `json:"post"`
	PostFlairs []PostFlair `json:"postFlairs"`
}

type PostDetailResponse struct {
	Status  string                    `json:"status"`
	Results PostDetailResponseResults `json:"results"`
}

type BoardAdminsResponseResults struct {
	Boards []Board `json:"boards"`
}
type BoardAdminsResponse struct {
	Status  string                     `json:"status"`
	Results BoardAdminsResponseResults `json:"results"`
}

type UpdateBoardRequest struct {
	IsVisible              bool   `json:"isVisible"              validate:"required,boolean"`
	IsPrivate              bool   `json:"isPrivate"              validate:"required,boolean"`
	IsOfficial             bool   `json:"isOfficial"             validate:"required,boolean"`
	IsPostApprovalRequired bool   `json:"isPostApprovalRequired" validate:"required,boolean"`
	MinKarmaRequiredToPost int    `json:"minKarmaRequiredToPost" validate:"required,min=0,max=50000"`
	Description            string `json:"description"            validate:"required,min=10,max=1000000"`
	ThumbnailFilename      string `json:"thumbnailFilename"      validate:"required"`
}

type PostListResponseResults struct {
	Posts      []BoardPost `json:"posts"`
	TotalPosts int         `json:"totalPosts"`
}

type PostListResponse struct {
	Status  string                  `json:"status"`
	Results PostListResponseResults `json:"results"`
}

func GetBoards(dbPool *pgxpool.Pool, omitEmpty bool) ([]Board, error) {
	havingClause := ""
	if omitEmpty {
		havingClause = `GROUP BY b.id, u.username, u.slug
			HAVING COUNT(bp.*) > 0
		`
	}
	// TODO: filter visible boards, or show everything if privileged
	// NOTE: inner join used here - boards without posts will not be included in results
	query := fmt.Sprintf(`
		SELECT %s
		b.created_by_user_id,
		u.username AS created_by_username,
		u.slug AS created_by_user_slug
		FROM boards b
		JOIN users u on u.id = b.created_by_user_id
		JOIN board_posts bp ON b.id = bp.board_id
		%s
		ORDER BY b.display_name
	`, getBoardColumns(), havingClause)
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

func GetTotalPosts(dbPool *pgxpool.Pool) (int, error) {
	const query = `SELECT COUNT(*) AS totalPosts FROM board_posts`
	row, err := dbPool.Query(context.Background(), query)
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
	return result.TotalPosts, nil
}

func getPostsQuery(whereClause string, paginationData PaginationData) string {
	limitClause := ""
	offsetClause := ""
	if paginationData.PerPage > 0 {
		limitClause = fmt.Sprintf("LIMIT %d", paginationData.PerPage)
	}
	if paginationData.Offset > 0 {
		offsetClause = fmt.Sprintf("OFFSET %d", paginationData.Offset)
	}
	return fmt.Sprintf(`
		SELECT 
		    bp.*,
			u.username AS created_by_username,
			u.slug AS created_by_user_slug,
			b.display_name AS boardName,
			b.slug AS boardSlug,
			b.is_official AS boardIsOfficial,
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
		`+whereClause+`
		ORDER BY bp.is_pinned DESC, bp.created_at DESC
		%s
		%s
	`, limitClause, offsetClause)
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
func GetPosts(
	dbPool *pgxpool.Pool, boardSlug string, postSlug string, parentId int,
	logger *slog.Logger, paginationData PaginationData, showUnapproved bool,
) ([]BoardPost, error,
) {
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
	}
	/*
		// Results in 0 rows when viewing a board post...but we actually
		// want this because we only want top level posts in this scenario...
		// Results show incorrect parent_id due to joins
		else {
			whereClause += " AND bp.parent_id = 0"
		}
	*/
	if !showUnapproved {
		whereClause += " AND bp.is_approved = true"
	}
	query := getPostsQuery(whereClause, paginationData)
	var params []string
	var rows pgx.Rows
	var err error
	if len(boardSlug) > 0 && len(postSlug) == 0 {
		rows, err = dbPool.Query(context.Background(), query, boardSlug)
		params = append(params, boardSlug)
	} else if len(postSlug) > 0 {
		rows, err = dbPool.Query(context.Background(), query, boardSlug, postSlug)
		params = append(params, boardSlug, postSlug)
	} else if parentId > 0 {
		rows, err = dbPool.Query(context.Background(), query, parentId)
		params = append(params, string(rune(parentId)))
	} else {
		rows, err = dbPool.Query(context.Background(), query)
	}
	logger.Info(
		fmt.Sprintf("GetPosts query: %v", debugQuery(query, params)),
	)
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

func GetPostDetail(dbPool *pgxpool.Pool, postSlug string) (BoardPost, error) {
	query := `
		SELECT 
		    bp.*,
			u.username AS created_by_username,
			u.slug AS created_by_user_slug,
			b.display_name AS boardName,
			b.slug AS boardSlug,
			b.is_official AS boardIsOfficial,
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
		WHERE bp.slug = $1
		AND bp.is_approved = true
	`
	row, err := dbPool.Query(context.Background(), query, postSlug)
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

func getPostParentIdBySlug(dbPool *pgxpool.Pool, postSlug string) (int, error) {
	const query = `SELECT id FROM board_posts WHERE slug = $1`
	var parentId int
	scanErr := dbPool.QueryRow(context.Background(), query, postSlug).Scan(&parentId)
	if scanErr != nil {
		return 0, scanErr
	}
	return parentId, nil
}

// AddPost - isApproved is not part of AddPostRequest because it would require the user to send this property.
// We don't care if the user requests it. This value is derived by checking permissions and board
// settings in the route.
func AddPost(dbPool *pgxpool.Pool, post AddPostRequest, userId int, boardId int, isApproved bool) (int, error) {
	parentPostId, parentPostErr := getPostParentIdBySlug(dbPool, post.ParentSlug)
	if parentPostErr != nil {
		return 0, parentPostErr
	}
	lastInsertId := 0
	const query = `
		INSERT INTO board_posts (title, created_by_user_id, board_id, parent_id, slug, post_text, is_approved) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	insertErr := dbPool.QueryRow(
		context.Background(),
		query,
		post.Title,
		userId,
		boardId,
		parentPostId,
		post.Slug,
		post.PostText,
		isApproved,
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

func AddBoard(
	dbPool *pgxpool.Pool, slug string, board AddBoardRequest, createdByUserId int,
) (int, error) {
	const query = `
		INSERT INTO boards (
			slug,
			display_name,
			created_at,
		    thumbnail_filename,
			created_by_user_id,
		    description,
		    is_post_approval_required
		)
		VALUES ($1, $2, NOW(), $3, $4, $5, $6)
		RETURNING id
	`
	var boardId int
	err := dbPool.QueryRow(
		context.Background(),
		query,
		slug,
		board.DisplayName,
		board.ThumbnailFilename,
		createdByUserId,
		board.Description,
		board.IsPostApprovalRequired,
	).Scan(&boardId)
	if err != nil {
		return 0, err
	}

	return boardId, nil
}

func DeleteBoard(dbPool *pgxpool.Pool, boardSlug string) error {
	// TODO: maybe use cascading here.

	/**
	 * Board users must be deleted before the board can be deleted
	 * because of FK constraints.
	 */
	deleteBoardUsersErr := DeleteBoardUsers(dbPool, boardSlug)
	if deleteBoardUsersErr != nil {
		return deleteBoardUsersErr
	}

	// Delete board posts?

	// Finally, delete the board
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
	post, postErr := GetPostDetail(dbPool, boardPostSlug)
	if postErr != nil {
		return postErr
	}
	deleteFlairErr := DeleteBoardPostFlairs(dbPool, post.Id)
	if deleteFlairErr != nil {
		return deleteFlairErr
	}

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

func GetPostFlairs(dbPool *pgxpool.Pool) ([]PostFlair, error) {
	query := `SELECT * FROM post_flairs`
	var rows pgx.Rows
	var rowsErr error
	rows, rowsErr = dbPool.Query(context.Background(), query)
	if rowsErr != nil {
		return []PostFlair{}, rowsErr
	}
	postFlairs, postFlairsErr := pgx.CollectRows(rows, pgx.RowToStructByName[PostFlair])
	if postFlairsErr != nil {
		return []PostFlair{}, postFlairsErr
	}

	return postFlairs, nil
}

// GetPostFlairsForPostId Post flairs for an individual post
func GetPostFlairsForPostId(dbPool *pgxpool.Pool, boardPostId int) ([]PostFlair, error) {
	const query = `
		SELECT * 
		FROM post_flairs
		LEFT JOIN posts_flairs ON post_flairs.id = posts_flairs.post_flair_id
		WHERE board_post_id = $1
	`
	rows, rowsErr := dbPool.Query(context.Background(), query, boardPostId)
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
// NOTE: this is used to create a map, so there is no post id filter
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

func GetPostFlairQuery(postId int, postFlairIds []int) string {
	var values []string
	for _, flairId := range postFlairIds {
		values = append(values, fmt.Sprintf("(%d, %d)", postId, flairId))
	}
	valuePairs := strings.Join(values, ",")
	query := fmt.Sprintf(`INSERT INTO posts_flairs (board_post_id, post_flair_id) VALUES %v`, valuePairs)

	return query
}

func AddPostFlair(dbPool *pgxpool.Pool, postId int, postFlairIds []int) error {
	postFlairsDeletedErr := DeleteBoardPostFlairs(dbPool, postId)
	if postFlairsDeletedErr != nil {
		logger.Error("AddPostFlair: error deleting existing post flairs")
		return postFlairsDeletedErr
	}
	query := GetPostFlairQuery(postId, postFlairIds)
	// logger.Info(fmt.Sprintf("PostFlairsQuery: %v", query))
	_, insertErr := dbPool.Query(context.Background(), query)
	if insertErr != nil {
		logger.Error("AddPostFlair: error inserting post flairs")
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

func UpdateBoard(
	dbPool *pgxpool.Pool, boardId int, updateBoardRequest UpdateBoardRequest, logger *slog.Logger) (bool, error) {
	const query = `
		UPDATE boards 
		SET description = $1, 
		    is_visible = $2, 
		    is_private = $3,
		    is_official = $4,
		    is_post_approval_required = $5,
		    min_karma_required_to_post = $6,
		    updated_at = NOW()
		WHERE id = $7
	`
	logger.Info(fmt.Sprintf("UpdateBoard: query: %v", query))
	result, err := dbPool.Exec(
		context.Background(),
		query,
		updateBoardRequest.Description,
		updateBoardRequest.IsVisible,
		updateBoardRequest.IsPrivate,
		updateBoardRequest.IsOfficial,
		updateBoardRequest.IsPostApprovalRequired,
		updateBoardRequest.MinKarmaRequiredToPost,
		boardId)
	if err != nil {
		return false, err
	}

	logger.Info(fmt.Sprintf("UpdateBoard: updated board #%v result: %v", boardId, result.RowsAffected()))

	return result.RowsAffected() > 0, nil
}

func getBoardColumns() string {
	return `
		b.id,
		b.display_name,
		b.created_at,
		b.updated_at,
		b.slug,
		CASE WHEN b.thumbnail_filename IS NULL THEN '' ELSE b.thumbnail_filename END AS thumbnail_filename,
		CASE WHEN b.description IS NULL THEN '' ELSE b.description END AS description,
		b.is_visible, 
		b.is_private,
		b.is_official,
		b.is_post_approval_required,
		b.min_karma_required_to_post,
	`
}

func IsPostApprovalRequiredForBoard(dbPool *pgxpool.Pool, boardId int) (bool, error) {
	const query = `SELECT is_post_approval_required FROM boards WHERE id = $1`
	var isPostApprovalRequired bool
	err := dbPool.QueryRow(context.Background(), query, boardId).Scan(&isPostApprovalRequired)
	if err != nil {
		return false, err
	}
	return isPostApprovalRequired, nil
}

func DeletePostsByUserId(dbPool *pgxpool.Pool, userId int) error {
	const query = `DELETE FROM board_posts WHERE board_posts.created_by_user_id = $1`
	_, err := dbPool.Exec(context.Background(), query, userId)
	return err
}

func DeletePostFlairsByUserId(dbPool *pgxpool.Pool, userId int) error {
	const query = `
		DELETE FROM public.posts_flairs
		WHERE id = ANY(
			SELECT pf.id
			FROM posts_flairs pf
			LEFT JOIN board_posts bp ON bp.id = pf.board_post_id
			JOIN users u ON u.id = bp.created_by_user_id
			WHERE bp.created_by_user_id = $1
		)
	`
	_, err := dbPool.Exec(context.Background(), query, userId)
	return err
}
