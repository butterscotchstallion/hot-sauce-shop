package routes

import (
	"fmt"
	"log/slog"

	"hotsauceshop/lib"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// canAccessBoardDetails - checks if user is board admin or super board admin
// does not send an error
func canAccessBoardDetails(c *gin.Context, boardId int, dbPool *pgxpool.Pool, logger *slog.Logger) (bool, error) {
	userId, userSessionErr := lib.GetUserIdFromSession(c, dbPool, logger)
	if userSessionErr != nil || userId == 0 {
		return false, userSessionErr
	}

	isMessageBoardAdmin, isMessageBoardAdminErr := lib.IsMessageBoardAdmin(dbPool, boardId, userId)
	if isMessageBoardAdminErr != nil {
		logger.Error(fmt.Sprintf("Error checking if user is board admin: %v", isMessageBoardAdminErr))
		return false, &lib.InternalServerError{
			Message: isMessageBoardAdminErr.Error(),
		}
	}

	if !isMessageBoardAdmin {
		// Check role (this also checks if the user is signed in)
		// NOTE: this already sends a JSON response upon failure
		isSuperMessageBoardAdmin, isSuperMessageBoardAdminErr := lib.IsSuperMessageBoardAdmin(c, dbPool, logger)
		if isSuperMessageBoardAdminErr != nil {
			return false, isSuperMessageBoardAdminErr
		}
		if !isSuperMessageBoardAdmin {
			if !isSuperMessageBoardAdmin {
				logger.Error("Error updating board settings: user is not super message board admin")
			}
			return false, &lib.StatusForbiddenError{
				Message: "Error updating board settings: permission denied",
			}
		}
	}
	return true, nil
}

// CanBypassPostApproval - checks if the user can bypass the board setting isPostApprovalRequired
func CanBypassPostApproval(c *gin.Context, dbPool *pgxpool.Pool, board lib.Board, logger *slog.Logger) (bool, error) {
	isBoardMod, isBoardModErr := lib.IsUserBoardModerator(dbPool, board.Slug, c.GetInt("userId"))
	if isBoardModErr != nil {
		return false, isBoardModErr
	}
	if isBoardMod {
		return true, nil
	}

	// Sends error response
	isBoardAdminOrSuperAdmin, isBoardAdminOrSuperAdminErr := canAccessBoardDetails(c, board.Id, dbPool, logger)
	if isBoardAdminOrSuperAdminErr != nil {
		return false, isBoardAdminOrSuperAdminErr
	}

	return isBoardAdminOrSuperAdmin, nil
}
