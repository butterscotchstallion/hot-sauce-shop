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

func getPostIdAsNumberOrError(c *gin.Context) (int, error) {
	postIdSlug := c.Param("postId")
	postId, postIdSlugErr := strconv.Atoi(postIdSlug)
	if postIdSlugErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "ERROR",
			"message": "Post id not found",
		})
		return 0, postIdSlugErr
	}
	return postId, nil
}

//nolint:funlen
func Votes(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger) {
	// User vote map
	r.GET("/api/v1/vote-map", func(c *gin.Context) {
		// Check user
		userId, userSessionErr := GetUserIdFromSessionOrError(c, dbPool, logger)
		if userSessionErr != nil || userId == 0 {
			return
		}

		voteMap, voteMapErr := lib.GetUserVoteMap(dbPool, userId)
		if voteMapErr != nil {
			logger.Error(fmt.Sprintf("Error fetching vote map: %v", voteMapErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": voteMapErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"voteMap": voteMap,
			},
		})
	})

	// Get vote sum for a specific post
	r.GET("/api/v1/votes/:postId", func(c *gin.Context) {
		postId, postIdErr := getPostIdAsNumberOrError(c)
		if postIdErr != nil {
			return
		}

		voteSumMap, voteSumErr := lib.GetVoteSumMapByPostId(dbPool, postId)
		if voteSumErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error fetching vote sum",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"voteSumMap": voteSumMap,
			},
		})
	})

	var addUpdateVoteRequest lib.AddUpdateVoteRequest
	r.POST("/api/v1/votes/:postId", func(c *gin.Context) {
		postId, postIdErr := getPostIdAsNumberOrError(c)
		if postIdErr != nil {
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

		lib.SendWebsocketMessage(lib.WebsocketMessage{
			MessageType: "boardPostUserVoted",
			Data: gin.H{
				"postId": postId,
				"voteId": voteId,
			},
		}, logger)

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"voteId":  voteId,
				"message": "Vote submitted",
			},
		})
	})
}
