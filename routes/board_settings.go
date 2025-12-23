package routes

import (
	"fmt"
	"log/slog"
	"net/http"

	"hotsauceshop/lib"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BoardSettings(
	r *gin.Engine,
	dbPool *pgxpool.Pool,
	logger *slog.Logger,
) {
	r.GET("/api/v1/board-settings/:boardSlug", func(c *gin.Context) {
		// Translate boardSlug -> boardId
		board, boardErr := lib.GetBoardBySlug(dbPool, c.Param("boardSlug"))
		if boardErr != nil {
			logger.Error(fmt.Sprintf("Board settings error: %v", boardErr))
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: "Error getting board settings",
			})
			return
		}

		// Check access requirements
		accessPermitted := canAccessBoardDetails(c, board.Id, dbPool, logger)
		if !accessPermitted {
			return
		}

		// Finally, get board settings if permitted
		boardSettings, boardSettingsErr := lib.GetBoardSettings(dbPool, c.Param("boardSlug"))
		if boardSettingsErr != nil {
			logger.Error(fmt.Sprintf("Board settings error: %v", boardSettingsErr))
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: "Error getting board settings",
			})
			return
		}
		c.JSON(http.StatusOK, lib.BoardSettingsResponse{
			Status:  "OK",
			Results: boardSettings,
		})
	})

	/**
	 * 1. Check if the user is a board admin
	 * 2. Check if the user is super admin
	 * 3. Update board settings
	 */
	r.PUT("/api/v1/board-settings/:boardSlug", func(c *gin.Context) {
		board, boardErr := lib.GetBoardBySlug(dbPool, c.Param("boardSlug"))
		if boardErr != nil {
			logger.Error(fmt.Sprintf("Board settings error: %v", boardErr))
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: "Error getting board settings",
			})
			return
		}

		accessPermitted := canAccessBoardDetails(c, board.Id, dbPool, logger)
		if !accessPermitted {
			return
		}

		var settings lib.BoardSettingsUpdateRequest
		if err := c.ShouldBind(&settings); err != nil {
			logger.Error(fmt.Sprintf("UpdateBoardSettings: error binding settings request: %v", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		// User is either super message board admin or board admin at this point
		updateErr := lib.SetBoardSettings(dbPool, settings)
		if updateErr != nil {
			logger.Error(fmt.Sprintf("Error updating board settings: %v", updateErr))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error updating board settings",
			})
			return
		}

		c.JSON(http.StatusOK, lib.GenericResponse{
			Status:  "OK",
			Message: "Board settings updated",
		})
	})
}

func canAccessBoardDetails(c *gin.Context, boardId int, dbPool *pgxpool.Pool, logger *slog.Logger) bool {
	userId, userSessionErr := GetUserIdFromSessionOrError(c, dbPool, logger)
	if userSessionErr != nil || userId == 0 {
		return false
	}

	isMessageBoardAdmin, isMessageBoardAdminErr := lib.IsMessageBoardAdmin(dbPool, boardId, userId)
	if isMessageBoardAdminErr != nil {
		logger.Error(fmt.Sprintf("Error checking if user is board admin: %v", isMessageBoardAdminErr))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "ERROR",
			"message": isMessageBoardAdminErr.Error(),
		})
		return false
	}

	if !isMessageBoardAdmin {
		// Check role (this also checks if the user is signed in)
		// NOTE: this already sends a JSON response upon failure
		isSuperMessageBoardAdmin, isSuperMessageBoardAdminErr := lib.IsSuperMessageBoardAdmin(c, dbPool, logger)
		if isSuperMessageBoardAdminErr != nil {
			return false
		}
		if !isSuperMessageBoardAdmin {
			if !isSuperMessageBoardAdmin {
				logger.Error("Error updating board settings: user is not super message board admin")
			}
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "ERROR",
				"message": "Error updating board settings: permission denied",
			})
			return false
		}
	}
	return true
}
