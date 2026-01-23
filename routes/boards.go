package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"

	"hotsauceshop/lib"

	"github.com/gabriel-vasile/mimetype"
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
	r.GET("/api/v1/boards", func(c *gin.Context) {
		omitEmpty := c.DefaultQuery("omitEmpty", "0") == "1"
		boards, getBoardsErr := lib.GetBoards(dbPool, omitEmpty)
		if getBoardsErr != nil {
			logger.Error(fmt.Sprintf("Error fetching boards: %v", getBoardsErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": getBoardsErr.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, lib.BoardListResponse{
			Status: "OK",
			Results: lib.BoardListResponseResults{
				Boards: boards,
			},
		})
	})

	// Board detail
	r.GET("/api/v1/boards/:slug", func(c *gin.Context) {
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

		admins, adminsErr := lib.GetBoardAdmins(dbPool, boardSlug)
		if adminsErr != nil {
			logger.Error(fmt.Sprintf("Error fetching admins: %v", adminsErr.Error()))
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

		c.JSON(http.StatusOK, lib.BoardDetailResponse{
			Status: "OK",
			Results: lib.BoardDetailResponseResults{
				Board:           board,
				Moderators:      mods,
				Admins:          admins,
				NumBoardMembers: numBoardMembers,
				TotalPosts:      totalPosts,
			},
		})
	})

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
		postSlug := c.Param("postSlug")
		post, getPostDetailErr := lib.GetPostDetail(dbPool, postSlug)
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
	r.GET("/api/v1/posts", func(c *gin.Context) {
		paginationData := getValidPaginationData(c)
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

		filterByUserJoinedBoards := c.DefaultQuery("filterByUserJoinedBoards", "0")
		showUnapprovedParam := c.DefaultQuery("showUnapproved", "0")
		boardSlug := c.DefaultQuery("boardSlug", "")
		postSlug := c.DefaultQuery("postSlug", "")
		showUnapproved := false

		if showUnapprovedParam == "1" {
			logger.Info("showUnapproved param is true")
			if len(boardSlug) > 0 {
				board, boardErr := lib.GetBoardBySlug(dbPool, boardSlug)
				if boardErr != nil {
					logger.Error(fmt.Sprintf("Error fetching board: %v", boardErr.Error()))
				}
				// Don't bother checking unless this board requires post approval
				if board.IsPostApprovalRequired {
					logger.Info(fmt.Sprintf("Post approval required for board '%v'", board.DisplayName))

					canBypass, canBypassErr := CanBypassPostApproval(c, dbPool, board, logger)
					if canBypassErr != nil {
						logger.Error(fmt.Sprintf("Error checking if user can bypass post approval: %v", canBypassErr.Error()))
					}
					if canBypass {
						showUnapproved = true
						logger.Info(fmt.Sprintf("Bypassing post approval for board %v", boardSlug))
					}
				} else {
					logger.Info(fmt.Sprintf("Post approval not required for board '%v'", board.DisplayName))
				}
			}
			// TODO: maybe figure out how to handle this for the unfiltered posts page. We could just check if the user
			// is a super board admin but avoid doing it twice since CanBypassPostApproval already does that, but requires
			// a board
		}

		logger.Info(
			fmt.Sprintf(
				"GetPosts: fetching posts with parentId: %v, boardSlug: %v, postSlug: '%v', showUnapproved: %v",
				parentId,
				boardSlug,
				postSlug,
				showUnapproved,
			),
		)

		posts, getPostsErr = lib.GetPosts(
			dbPool,
			boardSlug,
			postSlug,
			parentId,
			logger,
			paginationData,
			showUnapproved,
		)

		if getPostsErr != nil {
			logger.Error(fmt.Sprintf("Error fetching posts: %v", getPostsErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": getPostsErr.Error(),
			})
			return
		}

		if filterByUserJoinedBoards == "1" {
			userId, userIdErr := lib.GetUserIdFromSession(c, dbPool, logger)
			if userIdErr != nil {
				logger.Error(fmt.Sprintf("Error fetching user id from session: %v", userIdErr.Error()))
			}
			userJoinedBoards, userJoinedBoardsErr := lib.GetJoinedBoardsByUserId(dbPool, userId)
			if userJoinedBoardsErr != nil {
				logger.Error(fmt.Sprintf("Error fetching joined boards: %v", userJoinedBoardsErr.Error()))
			}
			var filteredPosts []lib.BoardPost
			for _, post := range posts {
				for _, joinedBoard := range userJoinedBoards {
					if post.BoardId == joinedBoard.Id {
						filteredPosts = append(filteredPosts, post)
					}
				}
			}
			posts = filteredPosts
		}

		totalPosts, totalPostsErr := lib.GetTotalPosts(dbPool)
		if totalPostsErr != nil {
			logger.Error(fmt.Sprintf("Error fetching total posts: %v", totalPostsErr.Error()))
		}

		c.JSON(http.StatusOK, lib.PostListResponse{
			Status: "OK",
			Results: lib.PostListResponseResults{
				Posts:      posts,
				TotalPosts: totalPosts,
			},
		})
	})

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

		// Get board for this post - used to convert boardSlug slug to board id, which
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

		// Check if the board requires minimum karma to post
		if board.MinKarmaRequiredToPost > 0 {
			karma, karmaErr := lib.GetUserPostVoteSum(dbPool, userId)
			if karmaErr != nil {
				logger.Error(fmt.Sprintf("Error fetching user karma: %v", karmaErr.Error()))
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "ERROR",
					"message": karmaErr.Error(),
				})
				return
			}

			if karma < board.MinKarmaRequiredToPost {
				logger.Error(fmt.Sprintf("Error adding post: user %v does not have enough karma to post", userId))
				c.JSON(http.StatusForbidden, lib.GenericResponseWithErrorCode{
					Status:    "ERROR",
					Message:   "Permission denied: insufficient karma",
					ErrorCode: lib.ErrorCodeInsufficientKarma,
				})
				return
			}
			if karma >= board.MinKarmaRequiredToPost {
				logger.Info(fmt.Sprintf("User %v has enough karma to post: %v", userId, karma))
			}
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

		/**
		 * 1. Posts are approved by default
		 * 2. If the board is set to require approval, then check user permissions
		 * 3. If the user is moderator/board admin/super admin, then the post is approved
		 * 4. If not, then the post is unapproved
		 */
		isPostApproved := true
		boardRequiresPostApproval, boardRequiresPostApprovalErr := lib.IsPostApprovalRequiredForBoard(dbPool, board.Id)
		if boardRequiresPostApprovalErr != nil {
			logger.Error(
				fmt.Sprintf(
					"Error checking if board requires post approval: %v",
					boardRequiresPostApprovalErr.Error()),
			)
		}
		if boardRequiresPostApproval {
			logger.Info(fmt.Sprintf("Post approval required for board '%v'", board.DisplayName))
			/**
			 * If something goes wrong here, it's not critical. We will just leave the post
			 * unapproved. We don't want to abandon the entire post process.
			 */
			canBypass, canBypassError := CanBypassPostApproval(c, dbPool, board, logger)
			if canBypassError != nil {
				logger.Error(fmt.Sprintf("Error checking if user can bypass post approval: %v", canBypassError.Error()))
			}
			if canBypass {
				logger.Info("User can bypass post approval")
				isPostApproved = true
			} else {
				logger.Info("User cannot bypass post approval")
				isPostApproved = false
			}
		} else {
			logger.Info(fmt.Sprintf("Post approval not required for board '%v'", board.DisplayName))
		}

		logger.Info(fmt.Sprintf("Post approved: %v", isPostApproved))

		// Add post
		newPostId, addPostErr := lib.AddPost(dbPool, newPost, userId, board.Id, isPostApproved)
		if addPostErr != nil {
			logger.Error(fmt.Sprintf("Error adding post: %v", addPostErr.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": addPostErr.Error(),
			})
			return
		}

		// Add flair for the post
		if len(newPost.PostFlairIds) > 0 {
			addPostFlairErr := lib.AddPostFlair(dbPool, newPostId, newPost.PostFlairIds)
			if addPostFlairErr != nil {
				logger.Error(fmt.Sprintf("Error adding post flair: %v", addPostFlairErr))
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "ERROR",
					"message": addPostFlairErr.Error(),
				})
				return
			}
		}

		// logger.Info(fmt.Sprintf("Post flair added: %v", newPost.PostFlairIds))

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

			imageWidthHeight, imageWidthHeightErr := lib.GetImageWidthAndHeight(fullImagePath, logger)
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
				logger,
			)
			if createThumbnailErr != nil {
				logger.Error(fmt.Sprintf("Error creating thumbnail: %v", createThumbnailErr.Error()))
			}

			thumbWidthHeight, thumbWidthHeightErr := lib.GetImageWidthAndHeight(imageInfo.ThumbnailFullPath, logger)
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
				// This is the bound empty post variable from the request and won't have most of the
				// fields populated
				"post":        newPost,
				"newPostId":   newPostId,
				"newPostSlug": newPostSlug,
			},
		})
	})

	// Add Board
	r.POST("/api/v1/boards", func(c *gin.Context) {
		// Check role (this also checks if the user is signed in)
		// NOTE: this already sends a JSON response upon failure
		isMessageBoardAdmin, isMessageBoardAdminErr := lib.IsSuperMessageBoardAdmin(c, dbPool, logger)
		if isMessageBoardAdminErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error checking if user is super admin",
			})
			return
		}
		if !isMessageBoardAdmin {
			logger.Error("Error adding board: user is not message board admin")
			// the above check is sending a 200 OK somehow. this sends the correct status response
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
			return
		}

		// All checks complete at this point, assemble board info!

		boardSlug := slug.Make(newBoard.DisplayName)
		boardId, addBoardErr := lib.AddBoard(
			dbPool, boardSlug, newBoard, userId,
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

	/**
	 * Deactivate board - if DeactivatedByUserId is 0, then the board will be reactivated
	 */
	r.PUT("/api/v1/boards/:boardSlug/activation-status", func(c *gin.Context) {
		var updateBoardActivationStatusRequest lib.UpdateBoardActivationStatusRequest
		if err := c.ShouldBind(&updateBoardActivationStatusRequest); err != nil {
			logger.Error(fmt.Sprintf("UpdateBoardActivationStatus: error binding request: %v", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		// Check role (this also checks if the user is signed in)
		// NOTE: this already sends a JSON response upon failure
		isMessageBoardAdmin, isMessageBoardAdminErr := lib.IsSuperMessageBoardAdmin(c, dbPool, logger)
		if isMessageBoardAdminErr != nil {
			return
		}
		if !isMessageBoardAdmin {
			logger.Error("Error deactivating board: user is not message board admin")
			return
		}

		// If the board is activated, then userId (deactivatedByUserId) is 0
		var userId int
		var userIdErr error
		if updateBoardActivationStatusRequest.Activated {
			userId = 0
		} else {
			userId, userIdErr = GetUserIdFromSessionOrError(c, dbPool, logger)
			if userIdErr != nil {
				return
			}
		}

		boardSlug := c.Param("boardSlug")
		boardDeactivatedErr := lib.UpdateBoardActivationStatus(dbPool, boardSlug, userId)
		if boardDeactivatedErr != nil {
			logger.Error(fmt.Sprintf("Error deactivating board: %v", boardDeactivatedErr))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error deactivating board",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "Board deactivated",
		})
	})

	// Delete board (only used in tests currently)
	r.DELETE("/api/v1/boards/:boardSlug", func(c *gin.Context) {
		var updateBoardActivationStatusRequest lib.UpdateBoardActivationStatusRequest
		if err := c.ShouldBind(&updateBoardActivationStatusRequest); err != nil {
			logger.Error(fmt.Sprintf("DeleteBoard: error binding request: %v", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		// Check role (this also checks if the user is signed in)
		// NOTE: this already sends a JSON response upon failure
		isMessageBoardAdmin, isMessageBoardAdminErr := lib.IsSuperMessageBoardAdmin(c, dbPool, logger)
		if isMessageBoardAdminErr != nil {
			return
		}
		if !isMessageBoardAdmin {
			logger.Error("Error deleting board: user is not message board admin")
			return
		}

		// If the board is activated, then userId (deactivatedByUserId) is 0
		var userId int
		var userIdErr error
		if updateBoardActivationStatusRequest.Activated {
			userId = 0
		} else {
			userId, userIdErr = GetUserIdFromSessionOrError(c, dbPool, logger)
			if userIdErr != nil {
				return
			}
		}

		boardSlug := c.Param("boardSlug")
		boardDeactivatedErr := lib.UpdateBoardActivationStatus(dbPool, boardSlug, userId)
		if boardDeactivatedErr != nil {
			logger.Error(fmt.Sprintf("Error deleting board: %v", boardDeactivatedErr))
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

		// TODO: refactor this into a reusable function for use in editing permission check

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
		isMessageBoardAdmin, isMessageBoardAdminErr := lib.IsSuperMessageBoardAdmin(c, dbPool, logger)
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

	// Returns a list of boards of which the user is a board admin
	r.GET("/api/v1/board-admin", func(c *gin.Context) {
		userId, getUserIdErr := GetUserIdFromSessionOrError(c, dbPool, logger)
		if getUserIdErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error getting user id",
			})
			return
		}

		adminBoards, adminBoardsErr := lib.GetUserAdminBoards(dbPool, userId)
		if adminBoardsErr != nil {
			logger.Error(fmt.Sprintf("Error getting board admins: %v", adminBoardsErr))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error getting board admins",
			})
			return
		}

		c.JSON(http.StatusOK, lib.BoardAdminsResponse{
			Status: "OK",
			Results: lib.BoardAdminsResponseResults{
				Boards: adminBoards,
			},
		})
	})

	r.POST("/api/v1/board-admin/:boardId", func(c *gin.Context) {
		isSuperMessageBoardAdmin, isSuperMessageBoardAdminErr := lib.IsSuperMessageBoardAdmin(c, dbPool, logger)
		if isSuperMessageBoardAdminErr != nil {
			return
		}
		if !isSuperMessageBoardAdmin {
			logger.Error("Error getting board admins: user is not super message board admin")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "ERROR",
				"message": "Permission denied",
			})
			return
		}

		userId, getUserIdErr := GetUserIdFromSessionOrError(c, dbPool, logger)
		if getUserIdErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error adding board admin",
			})
			return
		}
		boardId, boardIdErr := strconv.Atoi(c.Param("boardId"))
		if boardIdErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": "Error adding board admin",
			})
			return
		}

		addMessageBoardAdminErr := lib.AddBoardAdmin(dbPool, userId, boardId)
		if addMessageBoardAdminErr != nil {
			logger.Error(fmt.Sprintf("Error adding board admin: %v", addMessageBoardAdminErr))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error adding board admin",
			})
			return
		}

		c.JSON(http.StatusOK, lib.GenericResponse{
			Status:  "OK",
			Message: "Board admin added",
		})
	})

	/**
	 * Update Board
	 * 1. Bind update request
	 * 2. Get board details by slug
	 * 3. Check if the user is a board admin or super board admin
	 * 4. Update and send a response
	 */
	r.PUT("/api/v1/boards/:boardSlug", func(c *gin.Context) {
		var updateBoardRequest lib.UpdateBoardRequest
		if err := c.ShouldBind(&updateBoardRequest); err != nil {
			logger.Error(fmt.Sprintf("Update board: error binding update request: %v", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": err.Error(),
			})
			return
		}

		// TODO: fix this - somehow payload is not validating on the front end
		// validate := validator.New(validator.WithRequiredStructEnabled())
		// err := validate.Struct(updateBoardRequest)
		// if err != nil {
		// 	logger.Error(err.Error())
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"status":  "ERROR",
		// 		"message": fmt.Sprintf("Update board: validation failed: %v", err),
		// 	})
		// 	return
		// }

		board, boardErr := lib.GetBoardBySlug(dbPool, c.Param("boardSlug"))
		if boardErr != nil {
			logger.Error(fmt.Sprintf("Error getting board: %v", boardErr))
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: "Error getting board",
			})
			return
		}

		logger.Info(fmt.Sprintf("UpdateBoard: found board id %v", board.Id))

		// This error is logged within the function, and we don't need to check it in this context
		accessPermitted, _ := canAccessBoardDetails(c, board.Id, dbPool, logger)
		if !accessPermitted {
			c.JSON(http.StatusForbidden, lib.GenericResponse{
				Status:  "ERROR",
				Message: "Access denied",
			})
			return
		}

		logger.Info("UpdateBoard: access permitted")

		updateSuccessful, updateBoardErr := lib.UpdateBoard(dbPool, board.Id, updateBoardRequest, logger)
		if updateBoardErr != nil {
			logger.Error(fmt.Sprintf("Error updating board: %v", updateBoardErr))
			c.JSON(http.StatusInternalServerError, lib.GenericResponse{
				Status:  "ERROR",
				Message: "Error updating board",
			})
			return
		}

		if !updateSuccessful {
			logger.Error("Error updating board: no changes made")
			c.JSON(http.StatusNotFound, lib.GenericResponse{
				Status:  "ERROR",
				Message: "No changes made",
			})
			return
		}

		c.JSON(http.StatusOK, lib.GenericResponse{
			Status:  "OK",
			Message: "Board updated",
		})
	})
}
