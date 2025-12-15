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
		// Check role (this also checks if the user is signed in)
		// NOTE: this already sends a JSON response upon failure
		isMessageBoardAdmin, isMessageBoardAdminErr := lib.IsMessageBoardAdmin(c, dbPool, logger)
		if isMessageBoardAdminErr != nil {
			return
		}
		if !isMessageBoardAdmin {
			if !isMessageBoardAdmin {
				logger.Error("Error updating board settings: user is not message board admin")
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "ERROR",
				"message": "Error updating board settings: permission denied",
			})
			return
		}
	})
}
