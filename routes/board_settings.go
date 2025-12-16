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
			Status: "OK",
			Results: lib.BoardSettings{
				IsOfficial:             boardSettings.IsOfficial,
				IsPostApprovalRequired: boardSettings.IsPostApprovalRequired,
				UpdatedAt:              boardSettings.UpdatedAt,
			},
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

		userId, userSessionErr := GetUserIdFromSessionOrError(c, dbPool, logger)
		if userSessionErr != nil || userId == 0 {
			return
		}

		isMessageBoardAdmin, isMessageBoardAdminErr := lib.IsMessageBoardAdmin(dbPool, board.Id, userId)
		if isMessageBoardAdminErr != nil {
			logger.Error(fmt.Sprintf("Error checking if user is board admin: %v", isMessageBoardAdminErr))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": isMessageBoardAdminErr.Error(),
			})
			return
		}

		if !isMessageBoardAdmin {
			// Check role (this also checks if the user is signed in)
			// NOTE: this already sends a JSON response upon failure
			isSuperMessageBoardAdmin, isSuperMessageBoardAdminErr := lib.IsSuperMessageBoardAdmin(c, dbPool, logger)
			if isSuperMessageBoardAdminErr != nil {
				return
			}
			if !isSuperMessageBoardAdmin {
				if !isSuperMessageBoardAdmin {
					logger.Error("Error updating board settings: user is not super message board admin")
				}
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  "ERROR",
					"message": "Error updating board settings: permission denied",
				})
				return
			}
		}
	})
}
