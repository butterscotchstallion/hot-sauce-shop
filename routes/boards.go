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
		posts, getPostsErr := lib.GetPosts(dbPool, boardSlug)
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
		logger.Info("GetBoardBySlug: fetching board by slug: " + boardSlug)
		board, getBoardErr := lib.GetBoardBySlug(dbPool, logger, boardSlug)
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
		posts, getPostsErr := lib.GetPosts(dbPool, "")
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

	// Add post
	// for a reply: reuse this function for the reply route and add parentId as a parameter
	var newPost lib.AddPostRequest
	r.POST("/api/v1/boards/:slug/posts", func(c *gin.Context) {
		// Check user
		userId, userSessionErr := GetUserIdFromSessionOrError(c, dbPool, logger)
		if userSessionErr != nil || userId == 0 {
			return
		}

		// Check request
		if err := c.ShouldBindJSON(&newPost); err != nil {
			logger.Error(fmt.Sprintf("AddPost: error binding requets JSON: %v", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		// Check board
		logger.Info("GetBoardBySlug: fetching board by slug: " + newPost.Slug)
		board, getBoardErr := lib.GetBoardBySlug(dbPool, logger, newPost.Slug)
		if getBoardErr != nil {
			logger.Error(fmt.Sprintf("Error fetching board: %v", getBoardErr.Error()))
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

		// Add post
		_, addPostErr := lib.AddPost(dbPool, newPost, userId, board.Id)
		if addPostErr != nil {
			logger.Error(fmt.Sprintf("Error adding post: %v", addPostErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": addPostErr.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  "OK",
			"message": "Post added",
			"results": gin.H{
				"post": newPost,
			},
		})
	})
}
