package routes

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"hotsauceshop/lib"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Boards(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger, store *persistence.InMemoryStore) {
	// Board list
	r.GET("/api/v1/boards", cache.CachePage(store, time.Minute*1, func(c *gin.Context) {
		boards, getBoardsErr := lib.GetBoards(dbPool)
		if getBoardsErr != nil {
			logger.Error(fmt.Sprintf("Error fetching boards: %v", getBoardsErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": getBoardsErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"boards": boards,
			},
		})
	}))

	// Board posts
	r.GET("/api/v1/boards/:slug/posts", cache.CachePage(store, time.Minute*1, func(c *gin.Context) {
		boardSlug := c.Param("slug")
		posts, getPostsErr := lib.GetBoardPostsBySlug(dbPool, boardSlug)
		if getPostsErr != nil {
			logger.Error(fmt.Sprintf("Error fetching posts: %v", getPostsErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": getPostsErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"posts": posts,
			},
		})
	}))

	// Board detail
	r.GET("/api/v1/boards/:slug", cache.CachePage(store, time.Minute*1, func(c *gin.Context) {
		boardSlug := c.Param("slug")
		board, getBoardErr := lib.GetBoardBySlug(dbPool, boardSlug)
		if getBoardErr != nil {
			logger.Error(fmt.Sprintf("Error fetching board details: %v", getBoardErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": getBoardErr.Error(),
			})
			return
		}

		if board == (lib.Board{}) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "ERROR",
				"message": "Board not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"board": board,
			},
		})
	}))

	// All posts
	r.GET("/api/v1/posts", cache.CachePage(store, time.Minute*1, func(c *gin.Context) {
		posts, getPostsErr := lib.GetPosts(dbPool)
		if getPostsErr != nil {
			logger.Error(fmt.Sprintf("Error fetching posts: %v", getPostsErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": getPostsErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"posts": posts,
			},
		})
	}))
}
