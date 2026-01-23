package lib

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponseResults struct {
	SessionId string `json:"sessionId"`
	User      User   `json:"user"`
}

type SignInResponse struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Results SignInResponseResults `json:"results"`
}

type UserListResponseResults struct {
	Users []User `json:"users"`
}

type UserListResponse struct {
	Status  string                  `json:"status"`
	Results UserListResponseResults `json:"results"`
}

type GenericResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type GenericResponseWithErrorCode struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type UserBoardsResponseResults struct {
	Boards []Board `json:"boards"`
}

type UserBoardsResponse struct {
	Status  string                    `json:"status"`
	Results UserBoardsResponseResults `json:"results"`
}

type UserProfileResponseResults struct {
	User                User    `json:"user"`
	Roles               []Role  `json:"roles"`
	UserPostCount       int     `json:"userPostCount"`
	UserPostVoteSum     int     `json:"userPostVoteSum"`
	UserModeratedBoards []Board `json:"userModeratedBoards"`
}

type UserProfileResponse struct {
	Status  string                     `json:"status"`
	Results UserProfileResponseResults `json:"results"`
}

type UserCreateResponseResults struct {
	User User `json:"user"`
}
type UserCreateResponse struct {
	Status  string                    `json:"status"`
	Results UserCreateResponseResults `json:"results"`
}

type UserCreatePayload struct {
	Username       string `json:"username"       validate:"required,min=3,max=20"`
	Password       string `json:"password"       validate:"required,min=18,max=100"`
	AvatarFilename string `json:"avatarFilename" validate:"min=5,max=255"`
}

func GetUserPostVoteSum(dbPool *pgxpool.Pool, userId int) (int, error) {
	const query = `
		SELECT COALESCE(SUM(v.value), 0) AS voteSum
		FROM votes v
		LEFT JOIN board_posts bp on bp.id = v.post_id
		WHERE bp.created_by_user_id = $1
	`
	type UserPostVoteSum struct {
		VoteSum int
	}
	rows, err := dbPool.Query(context.Background(), query, userId)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	userPostVoteSum, collectRowErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[UserPostVoteSum])
	if collectRowErr != nil {
		return 0, err
	}
	return userPostVoteSum.VoteSum, nil
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

func VerifyUsernameAndPasswordAndReturnUser(
	dbPool *pgxpool.Pool, logger *slog.Logger, username string, password string) (User, error) {
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
		AND (
		    s.created_at >= DATE_TRUNC('month', current_date - interval '1' month)
		    OR 
		    s.updated_at >= DATE_TRUNC('month', current_date - interval '1' month)
		)
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

func GetJoinedBoardsByUserId(dbPool *pgxpool.Pool, userId int) ([]Board, error) {
	query := fmt.Sprintf(`
		SELECT %s
		b.created_by_user_id,
		u.username AS created_by_username,
		u.slug AS created_by_user_slug
		FROM boards b
		JOIN users u on u.id = b.created_by_user_id
		JOIN boards_users bu ON bu.board_id = b.id
		WHERE bu.user_id = $1
		ORDER BY b.display_name
	`, getBoardColumns())
	rows, err := dbPool.Query(context.Background(), query, userId)
	if err != nil {
		return nil, err
	}
	boards, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[Board])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}
	return boards, nil
}

func AddBoardUser(dbPool *pgxpool.Pool, userId int, boardId int) error {
	const query = "INSERT INTO boards_users (user_id, board_id) VALUES ($1, $2)"
	_, err := dbPool.Exec(context.Background(), query, userId, boardId)
	if err != nil {
		return err
	}
	return nil
}

// GetUserModeratedBoards
// Returns boards the user is a moderator on
func GetUserModeratedBoards(dbPool *pgxpool.Pool, userId int) ([]Board, error) {
	query := fmt.Sprintf(`SELECT DISTINCT %s
		b.created_by_user_id,
		u.username AS created_by_username,
		u.slug AS created_by_user_slug
		FROM boards b
		JOIN user_roles_boards urb ON urb.board_id = b.id
		JOIN users u on urb.user_id = u.id
		WHERE urb.user_id = $1
		ORDER BY b.display_name`, getBoardColumns())
	rows, err := dbPool.Query(context.Background(), query, userId)
	if err != nil {
		return nil, err
	}
	boards, collectRowsErr := pgx.CollectRows(rows, pgx.RowToStructByName[Board])
	if collectRowsErr != nil {
		return nil, collectRowsErr
	}
	return boards, nil
}

func GetUserLevelInfoByUserId(dbPool *pgxpool.Pool, userId int) (UserLevelInfo, error) {
	var experience float64
	const query = `SELECT experience FROM user_experience WHERE user_id = $1`
	err := dbPool.QueryRow(context.Background(), query, userId).Scan(&experience)
	if err != nil {
		return UserLevelInfo{}, err
	}
	var level = 1
	if experience > 0 {
		level = GetUserLevelByExperience(experience)
	}

	percentage := GetPercentageOfLevelComplete(experience, level)

	return UserLevelInfo{
		Level:                     level,
		Experience:                experience,
		PercentageOfLevelComplete: percentage,
	}, nil
}

func CreateUser(dbPool *pgxpool.Pool, payload UserCreatePayload) (User, error) {
	usernameSlug := slug.Make(payload.Username)
	const query = `INSERT INTO users (username, password, avatar_filename, slug) 
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`
	row, rowErr := dbPool.Query(
		context.Background(),
		query,
		payload.Username,
		payload.Password,
		payload.AvatarFilename,
		usernameSlug,
	)
	if rowErr != nil {
		return User{}, rowErr
	}
	user, collectRowsErr := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[User])
	if collectRowsErr != nil {
		return User{}, collectRowsErr
	}
	return user, nil
}

func DeleteUser(dbPool *pgxpool.Pool, userId int) error {
	const query = `DELETE FROM users WHERE id = $1`
	_, err := dbPool.Exec(context.Background(), query, userId)
	return err
}
