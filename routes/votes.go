package routes

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"hotsauceshop/lib"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddUpdateVoteRequest struct {
	VoteValue int `json:"voteValue" binding:"required,oneof=1 -1"`
}

func Votes(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger) {
	var addUpdateVoteRequest AddUpdateVoteRequest
	r.POST("/api/v1/votes/:postId", func(c *gin.Context) {
		postIdSlug := c.Param("postId")

		if postIdSlug == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "ERROR",
				"message": "Post id not found",
			})
			return
		}

		postId, postIdSlugErr := strconv.Atoi(postIdSlug)
		if postIdSlugErr != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "ERROR",
				"message": "Post id not found",
			})
			return
		}

		// Check user
		userId, userSessionErr := GetUserIdFromSessionOrError(c, dbPool, logger)
		if userSessionErr != nil || userId == 0 {
			return
		}

		// Check request
		if err := c.ShouldBindJSON(&addUpdateVoteRequest); err != nil {
			logger.Error(fmt.Sprintf("AddUpdateVote: error binding requests JSON: %v", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		// FK enforcement prevents votes on posts that don't exist
		// TODO: rate limiting: check if user has voted on another post too soon

		logger.Info(
			fmt.Sprintf(
				"Adding vote for user %d on post %d with value %d",
				userId,
				postId,
				addUpdateVoteRequest.VoteValue,
			),
		)

		voteId, voteErr := lib.AddUpdateVote(
			dbPool, userId, postId, addUpdateVoteRequest.VoteValue,
		)
		if voteErr != nil {
			logger.Error(fmt.Sprintf("AddUpdateVote: error adding vote: %v", voteErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": voteErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"voteId":  voteId,
				"message": "Vote submitted",
			},
		})
	})
}
