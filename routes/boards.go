package routes

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"hotsauceshop/lib"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Boards(
	r *gin.Engine,
	dbPool *pgxpool.Pool,
	logger *slog.Logger,
	store *persistence.InMemoryStore,
) {
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

		mods, modsErr := lib.GetBoardModerators(dbPool, boardSlug, 0)
		if modsErr != nil {
			logger.Error(fmt.Sprintf("Error fetching mods: %v", modsErr.Error()))
		}

		if board == (lib.Board{}) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "ERROR",
				"message": "Board not found",
			})
			return
		}

		numBoardMembers, numBoardMembersErr := lib.GetNumBoardMembers(dbPool, boardSlug)
		if numBoardMembersErr != nil {
			logger.Error(fmt.Sprintf("Error fetching num board members: %v", numBoardMembersErr.Error()))
		}

		totalPosts, totalPostsErr := lib.GetTotalPostsByBoardSlug(dbPool, boardSlug)
		if totalPostsErr != nil {
			logger.Error(fmt.Sprintf("Error fetching total posts: %v", totalPostsErr.Error()))
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"board":           board,
				"moderators":      mods,
				"numBoardMembers": numBoardMembers,
				"totalPosts":      totalPosts,
			},
		})
	}))

	// Board total posts
	r.GET("/api/v1/total-posts/:boardSlug", func(c *gin.Context) {
		boardSlug := c.Param("boardSlug")

		if boardSlug == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "ERROR",
				"message": "Board not found",
			})
			return
		}

		totalPosts, totalPostsErr := lib.GetTotalPostsByBoardSlug(dbPool, boardSlug)
		if totalPostsErr != nil {
			logger.Error(fmt.Sprintf("Error fetching total posts: %v", totalPostsErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": totalPostsErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"totalPosts": totalPosts,
			},
		})
	})

	// Total post reply map
	r.GET("/api/v1/total-replies", func(c *gin.Context) {
		boardSlug := c.DefaultQuery("boardSlug", "")
		totalPostReplyMap, replyMapErr := lib.GetTotalPostReplyCountByBoardSlug(dbPool, boardSlug)
		if replyMapErr != nil {
			logger.Error(fmt.Sprintf("Error fetching total post replies: %v", replyMapErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": replyMapErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"totalPostReplyMap": totalPostReplyMap,
			},
		})
	})

	// Post detail
	r.GET("/api/v1/posts/:boardSlug/:postSlug", func(c *gin.Context) {
		boardSlug := c.Param("boardSlug")
		postSlug := c.Param("postSlug")
		post, getPostDetailErr := lib.GetPostDetail(dbPool, boardSlug, postSlug)
		if getPostDetailErr != nil {
			logger.Error(fmt.Sprintf("Error fetching post details: %v", getPostDetailErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": getPostDetailErr.Error(),
			})
			return
		}

		if post == (lib.BoardPost{}) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "ERROR",
				"message": "Post not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"post": post,
			},
		})
	})

	// All posts
	r.GET("/api/v1/posts", cache.CachePage(store, time.Minute*1, func(c *gin.Context) {
		var posts []lib.BoardPost
		var getPostsErr error

		parentIdParam := c.DefaultQuery("parentId", "0")
		parentId, err := strconv.Atoi(parentIdParam)
		if err != nil {
			logger.Error(fmt.Sprintf("Error parsing parentId: %v", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": "Invalid parentId",
			})
			return
		}

		boardSlug := c.DefaultQuery("boardSlug", "")
		postSlug := c.DefaultQuery("postSlug", "")

		logger.Info(
			fmt.Sprintf(
				"GetPosts: fetching posts with parentId: %v, boardSlug: %v, postSlug: %v",
				parentId,
				boardSlug,
				postSlug,
			),
		)

		posts, getPostsErr = lib.GetPosts(dbPool, boardSlug, postSlug, parentId, logger)

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

	// Pin post
	r.POST("/api/v1/boards/pin/:boardSlug/:postSlug", func(c *gin.Context) {
		boardSlug := c.Param("boardSlug")
		postSlug := c.Param("postSlug")

		// Check user
		userId, userSessionErr := GetUserIdFromSessionOrError(c, dbPool, logger)
		if userSessionErr != nil || userId == 0 {
			return
		}

		mods, modsErr := lib.GetBoardModerators(dbPool, boardSlug, userId)
		if modsErr != nil {
			logger.Error(fmt.Sprintf("Error fetching mods: %v", modsErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": modsErr.Error(),
			})
			return
		}

		if len(mods) == 0 {
			logger.Error(fmt.Sprintf("Error pinning post: user %v is not a moderator", userId))
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "ERROR",
				"message": "Permission denied",
			})
			return
		}

		pinPostErr := lib.PinBoardPost(dbPool, postSlug)
		if pinPostErr != nil {
			logger.Error(fmt.Sprintf("Error pinning post: %v", pinPostErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": pinPostErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	// Add post
	r.POST("/api/v1/boards/:slug/posts", func(c *gin.Context) {
		// Check user
		userId, userSessionErr := GetUserIdFromSessionOrError(c, dbPool, logger)
		if userSessionErr != nil || userId == 0 {
			return
		}

		var newPost lib.AddPostRequest
		if err := c.ShouldBind(&newPost); err != nil {
			logger.Error(fmt.Sprintf("AddPost: error binding add post request: %v", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		// Check board
		boardSlug := c.Param("slug")
		if boardSlug == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": "Board slug is required",
			})
			return
		}

		// Check board
		logger.Info("GetBoardBySlug: fetching board by slug: " + boardSlug)
		board, getBoardErr := lib.GetBoardBySlug(dbPool, logger, boardSlug)
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

		// Create post slug
		newPostSlug, err := uuid.NewRandom()
		if err != nil {
			logger.Error(fmt.Sprintf("Error generating new post slug: %v", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}
		newPost.Slug = newPostSlug.String()

		// Add post
		newPostId, addPostErr := lib.AddPost(dbPool, newPost, userId, board.Id)
		if addPostErr != nil {
			logger.Error(fmt.Sprintf("Error adding post: %v", addPostErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": addPostErr.Error(),
			})
			return
		}

		isImagePost := false

		// Add post images
		postImagePath := "ui/src/public/images/posts/"
		var savedPostImageInfo []lib.SavedPostImageInfo
		postImages := newPost.PostImages
		for index, postImage := range postImages {
			originalExtension := filepath.Ext(postImage.Filename)
			postImageFilename := fmt.Sprintf("%s-%v%s", newPost.Slug, index, originalExtension)
			fullImagePath := postImagePath + postImageFilename
			thumbnailFilename := lib.GetThumbnailFilename(postImageFilename)
			thumbnailFullPath := postImagePath + thumbnailFilename

			saveFileErr := c.SaveUploadedFile(postImage, fullImagePath)
			if saveFileErr != nil {
				logger.Error(fmt.Sprintf("Error saving post image: %v", saveFileErr.Error()))
				continue
			}

			// Get mime type
			mimeType, mimeTypeErr := mimetype.DetectFile(fullImagePath)
			if mimeTypeErr != nil {
				logger.Error(fmt.Sprintf("Error detecting file type: %v", mimeTypeErr.Error()))
				continue
			}
			logger.Info(fmt.Sprintf("%v has mime type %v", fullImagePath, mimeType.String()))

			imageWidthHeight, imageWidthHeightErr := lib.GetImageWidthAndHeight(fullImagePath)
			if imageWidthHeightErr != nil {
				logger.Error(fmt.Sprintf("Error getting image width: %v", imageWidthHeightErr.Error()))
				continue
			}

			// Assemble image info for use with AddPostImages/thumbnails
			logger.Info(fmt.Sprintf("Post image saved: %v", postImageFilename))

			savedPostImageInfo = append(savedPostImageInfo, lib.SavedPostImageInfo{
				Filename:          postImageFilename,
				FullImagePath:     fullImagePath,
				ThumbnailFilename: thumbnailFilename,
				ThumbnailFullPath: thumbnailFullPath,
				MimeType:          mimeType.String(),
				ImageWidthHeight:  imageWidthHeight,
			})
		}

		// Iterate successfully saved images, and add to DB/thumbnail
		for _, imageInfo := range savedPostImageInfo {
			// Create thumbnail
			createThumbnailErr := lib.CreateThumbnail(
				imageInfo.FullImagePath,
				imageInfo.ThumbnailFullPath,
				imageInfo.MimeType,
			)
			if createThumbnailErr != nil {
				logger.Error(fmt.Sprintf("Error creating thumbnail: %v", createThumbnailErr.Error()))
			}

			thumbWidthHeight, thumbWidthHeightErr := lib.GetImageWidthAndHeight(imageInfo.ThumbnailFullPath)
			if thumbWidthHeightErr != nil {
				logger.Error(fmt.Sprintf("Error getting thumbnail width: %v", thumbWidthHeightErr.Error()))
				continue
			}

			imageInfo.ThumbnailWidthHeight = thumbWidthHeight

			addPostImagesErr := lib.AddPostImages(dbPool, newPostId, imageInfo)
			if addPostImagesErr != nil {
				logger.Error(fmt.Sprintf("Error adding post image to DB: %v", addPostImagesErr.Error()))
			}
			logger.Info(
				fmt.Sprintf(
					"Post image saved to DB: %v with mime type %v",
					imageInfo.Filename,
					imageInfo.MimeType,
				),
			)
			isImagePost = true
		}

		var updatedExperience float64
		var addImagePostExperienceErr error
		var addPostExperienceErr error
		var experienceUpdated bool
		if isImagePost {
			updatedExperience, addImagePostExperienceErr = lib.AddImagePostExperienceToUser(dbPool, userId)
			if addImagePostExperienceErr != nil {
				logger.Error(fmt.Sprintf("Error adding image experience: %v", addImagePostExperienceErr.Error()))
			} else {
				logger.Info("Added image experience to user")
				experienceUpdated = true
			}
		} else {
			updatedExperience, addPostExperienceErr = lib.AddPostExperienceToUser(dbPool, userId)
			if addPostExperienceErr != nil {
				logger.Error(fmt.Sprintf("Error adding post experience: %v", addPostExperienceErr.Error()))
			} else {
				logger.Info("Added post experience to user")
				experienceUpdated = true
			}
		}

		if experienceUpdated {
			updatedLevel := lib.GetUserLevelByExperience(updatedExperience)
			lib.SendNotification(lib.WebsocketMessage{
				MessageType: "userLevelUpdate",
				Data: gin.H{
					"updatedExperience":         updatedExperience,
					"updatedLevel":              updatedLevel,
					"percentageOfLevelComplete": lib.GetPercentageOfLevelComplete(updatedExperience, updatedLevel),
				},
			}, logger)
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  "OK",
			"message": "Post added",
			"results": gin.H{
				"post":      newPost,
				"newPostId": newPostId,
			},
		})
	})

	// Add Board
	r.POST("/api/v1/boards", func(c *gin.Context) {
		// Check role (this also checks if the user is signed in)
		// NOTE: this already sends a JSON response upon failure
		isMessageBoardAdmin, isMessageBoardAdminErr := lib.IsMessageBoardAdmin(c, dbPool, logger)
		if isMessageBoardAdminErr != nil {
			return
		}
		if !isMessageBoardAdmin {
			logger.Error("Error adding board: user is not message board admin")
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "ERROR",
				"message": "Error adding board: permission denied",
			})
			return
		}

		// Check payload
		var newBoard lib.AddBoardRequest
		if err := c.ShouldBind(&newBoard); err != nil {
			logger.Error(fmt.Sprintf("AddBoard: error binding add board request: %v", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		// This shouldn't be an error at this point, but we need the userId
		userId, getUserIdErr := GetUserIdFromSessionOrError(c, dbPool, logger)
		if getUserIdErr != nil {
			logger.Error("AddBoard: error getting user id from session")
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "ERROR",
				"message": "Permission denied",
			})
			return
		}

		// All checks complete at this point, assemble board info!

		boardSlug := slug.Make(newBoard.DisplayName)
		addBoardErr := lib.AddBoard(
			dbPool, boardSlug, newBoard.DisplayName, newBoard.ThumbnailFilename, userId, newBoard.Description,
		)
		if addBoardErr != nil {
			logger.Error(fmt.Sprintf("AddBoard: error adding board: %v", addBoardErr))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error adding board",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  "OK",
			"message": "Board added",
			"results": gin.H{
				"slug":        boardSlug,
				"displayName": newBoard.DisplayName,
			},
		})
	})

	// Delete Board
	r.DELETE("/api/v1/boards/:boardSlug", func(c *gin.Context) {
		// Check role (this also checks if the user is signed in)
		// NOTE: this already sends a JSON response upon failure
		isMessageBoardAdmin, isMessageBoardAdminErr := lib.IsMessageBoardAdmin(c, dbPool, logger)
		if isMessageBoardAdminErr != nil {
			return
		}
		if !isMessageBoardAdmin {
			logger.Error("Error adding board: user is not message board admin")
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "ERROR",
				"message": "Error adding board: permission denied",
			})
			return
		}

		boardSlug := c.Param("boardSlug")
		boardDeletedErr := lib.DeleteBoard(dbPool, boardSlug)
		if boardDeletedErr != nil {
			logger.Error(fmt.Sprintf("Error deleting board: %v", boardDeletedErr))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error deleting board",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "Board deleted",
		})
	})
}
