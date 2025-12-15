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
}
