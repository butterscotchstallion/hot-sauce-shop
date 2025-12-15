package routes

import (
	"database/sql"
	"errors"
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

//nolint:funlen
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
		board, getBoardErr := lib.GetBoardBySlug(dbPool, boardSlug)
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

		c.JSON(http.StatusOK, lib.BoardPostResponse{
			Status: "OK",
			Results: lib.BoardPostResponseResults{
				Board:           board,
				Moderators:      mods,
				NumBoardMembers: numBoardMembers,
				TotalPosts:      totalPosts,
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

	// Total post/reply map
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
		if getPostDetailErr != nil && !errors.Is(getPostDetailErr, sql.ErrNoRows) {
			logger.Error(fmt.Sprintf("Error fetching post details: %v", getPostDetailErr.Error()))
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: getPostDetailErr.Error(),
			})
			return
		}

		if post == (lib.BoardPost{}) || errors.Is(getPostDetailErr, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, lib.GenericResponse{
				Status:  "ERROR",
				Message: "Post not found",
			})
			return
		}

		postFlairs, postFlairsErr := lib.GetPostFlairsForPostId(dbPool, post.Id)
		if postFlairsErr != nil {
			logger.Error(fmt.Sprintf("Error fetching post flairs: %v", postFlairsErr.Error()))
		}

		c.JSON(http.StatusOK, lib.PostDetailResponse{
			Status: "OK",
			Results: lib.PostDetailResponseResults{
				Post:       post,
				PostFlairs: postFlairs,
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

		// Get board for this post - used to convert board slug to board id, which
		// is used below when adding the post
		boardSlug := c.Param("slug")
		logger.Info("GetBoardBySlug: fetching board by slug: " + boardSlug)
		board, getBoardErr := lib.GetBoardBySlug(dbPool, boardSlug)
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

		// Create a slug for the post
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

		// Add flair for the post
		addPostFlairErr := lib.AddPostFlair(dbPool, newPostId, newPost.PostFlairIds)
		if addPostFlairErr != nil {
			logger.Error(fmt.Sprintf("Error adding post flair: %v", addPostFlairErr))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": addPostFlairErr.Error(),
			})
			return
		}

		logger.Info(fmt.Sprintf("Post flair added: %v", newPost.PostFlairIds))

		isImagePost := false

		// Add images for the post
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

		// Iterate successfully saved images and add them to DB/thumbnail
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

		/**
		 * Add experience to user based on the activity type
		 */
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

		// Send WS message if experience updated
		if experienceUpdated {
			updatedLevel := lib.GetUserLevelByExperience(updatedExperience)
			sendErr := lib.SendWebsocketMessage(lib.WebsocketMessage{
				MessageType: "userLevelUpdate",
				Data: gin.H{
					"updatedExperience":         updatedExperience,
					"updatedLevel":              updatedLevel,
					"percentageOfLevelComplete": lib.GetPercentageOfLevelComplete(updatedExperience, updatedLevel),
				},
			}, logger)
			if sendErr != nil {
				logger.Error(fmt.Sprintf("Error sending websocket message: %v", sendErr.Error()))
			}
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
		boardId, addBoardErr := lib.AddBoard(
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

		c.JSON(http.StatusCreated, lib.AddBoardResponse{
			Status:  "OK",
			Message: "Board added",
			Results: lib.AddBoardResponseResults{
				Slug:        boardSlug,
				DisplayName: newBoard.DisplayName,
				BoardId:     boardId,
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
			logger.Error("Error deleting board: user is not message board admin")
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "ERROR",
				"message": "Error deleting board: permission denied",
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

	// Delete Post
	r.DELETE("/api/v1/boards/posts/:postSlug", func(c *gin.Context) {
		postSlug := c.Param("postSlug")

		userId, userSessionErr := GetUserIdFromSessionOrError(c, dbPool, logger)
		if userSessionErr != nil || userId == 0 {
			return
		}

		// Check if the user is post author here
		isBoardPostAuthor, isBoardPostAuthorErr := lib.IsUserBoardPostAuthor(dbPool, userId, postSlug)
		if isBoardPostAuthorErr != nil {
			logger.Error(fmt.Sprintf("Delete Post: error checking if user is post author: %v", isBoardPostAuthorErr))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error checking post author",
			})
			return
		}

		// Check role (this also checks if the user is signed in)
		// NOTE: this already sends a JSON response upon failure
		isMessageBoardAdmin, isMessageBoardAdminErr := lib.IsMessageBoardAdmin(c, dbPool, logger)
		if isMessageBoardAdminErr != nil {
			return
		}
		if !isMessageBoardAdmin && !isBoardPostAuthor {
			if !isMessageBoardAdmin {
				logger.Error("Error deleting post: user is not message board admin")
			}
			if !isBoardPostAuthor {
				logger.Error("Error deleting post: user is not post author")
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "ERROR",
				"message": "Error deleting post: permission denied",
			})
			return
		}

		postDeletedErr := lib.DeleteBoardPost(dbPool, postSlug)
		if postDeletedErr != nil {
			logger.Error(fmt.Sprintf("Error deleting post: %v", postDeletedErr))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error deleting post",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "Post deleted",
		})
	})

	// All available post flair listing
	r.GET("/api/v1/post-flairs", func(c *gin.Context) {
		postFlairs, postFlairsErr := lib.GetPostFlairs(dbPool)
		if postFlairsErr != nil {
			logger.Error(fmt.Sprintf("Error getting post flairs: %v", postFlairsErr))
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: "Error getting post flairs",
			})
			return
		}
		c.JSON(http.StatusOK, lib.PostFlairsResponse{
			Status: "OK",
			Results: lib.PostFlairsResponseResults{
				PostFlairs: postFlairs,
			},
		})
	})

	/**
	 * Flairs for each post
	 * This is used to create a map of post-id -> post flairs, which
	 * is then used for the post list page. Individual posts will use a
	 * single query to get post flairs
	 */
	r.GET("/api/v1/posts-flairs", func(c *gin.Context) {
		postFlairs, postFlairsErr := lib.GetPostsFlairs(dbPool)
		if postFlairsErr != nil {
			logger.Error(fmt.Sprintf("Error getting post flairs: %v", postFlairsErr))
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: "Error getting post flairs",
			})
			return
		}
		c.JSON(http.StatusOK, lib.PostsFlairsResponse{
			Status: "OK",
			Results: lib.PostsFlairsResponseResults{
				PostsFlairs: postFlairs,
			},
		})
	})

	// r.GET("/api/v1/posts-flairs/post/:postId", func(c *gin.Context) {
	// 	postIdParam := c.Param("postId")
	// 	postId, postIdErr := strconv.Atoi(postIdParam)
	// 	if postIdErr != nil {
	// 		logger.Error(fmt.Sprintf("Error parsing post id: %v", postIdErr))
	// 		c.JSON(http.StatusInternalServerError, lib.GenericResponse{
	// 			Status:  "ERROR",
	// 			Message: "Error parsing post id",
	// 		})
	// 		return
	// 	}
	// 	postFlairs, postFlairsErr := lib.GetPostFlairsForPostId(dbPool, postId)
	// 	if postFlairsErr != nil {
	// 		logger.Error(fmt.Sprintf("Error getting post flairs: %v", postFlairsErr))
	// 		c.JSON(http.StatusInternalServerError, lib.GenericResponse{
	// 			Status:  "ERROR",
	// 			Message: "Error getting post flairs",
	// 		})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, lib.PostFlairsResponse{
	// 		Status: "OK",
	// 		Results: lib.PostFlairsResponseResults{
	// 			PostFlairs: postFlairs,
	// 		},
	// 	})
	// })
}
