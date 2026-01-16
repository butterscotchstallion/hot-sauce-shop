package routes

import (
	"fmt"
	"net/http"
	"testing"

	"hotsauceshop/lib"

	"github.com/gavv/httpexpect/v2"
)

type GetPostListAndVerifyParams struct {
	E                        *httpexpect.Expect
	T                        *testing.T
	SessionId                string
	ShowUnapproved           bool
	BoardName                string
	VerifyPostListIsNotEmpty bool
	FilterByUserJoinedBoards bool
}

func boolToNumericStr(param bool) string {
	result := "0"
	if param {
		result = "1"
	}
	return result
}

func getPostListAndVerify(params GetPostListAndVerifyParams) lib.PostListResponse {
	var postListResponse lib.PostListResponse
	showUnapprovedStr := 0
	filterByUserJoinedBoardsStr := 0
	if params.FilterByUserJoinedBoards {
		filterByUserJoinedBoardsStr = 1
	}
	if params.ShowUnapproved {
		showUnapprovedStr = 1
	}
	params.E.GET("/api/v1/posts").
		WithQuery("showUnapproved", showUnapprovedStr).
		WithQuery("filterByUserJoinedBoards", filterByUserJoinedBoardsStr).
		WithQuery("boardSlug", params.BoardName).
		WithCookie("sessionId", params.SessionId).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&postListResponse)
	if postListResponse.Status != "OK" {
		params.T.Fatal("Failed to get post list")
	}
	if params.VerifyPostListIsNotEmpty {
		if len(postListResponse.Results.Posts) == 0 {
			params.T.Fatal("Post list is empty!")
		}
	}
	return postListResponse
}

type GetBoardListRequest struct {
	T               *testing.T
	E               *httpexpect.Expect
	OmitEmptyBoards bool
}

func getBoardList(request GetBoardListRequest) lib.BoardListResponse {
	var boardsResponse lib.BoardListResponse
	request.E.GET("/api/v1/boards").
		WithQuery("omitEmpty", boolToNumericStr(request.OmitEmptyBoards)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&boardsResponse)
	if boardsResponse.Status != "OK" {
		request.T.Fatal("Failed to get board list")
	}
	return boardsResponse
}

func isPostSlugInList(postSlug string, posts []lib.BoardPost) bool {
	for _, post := range posts {
		if post.Slug == postSlug {
			return true
		}
	}
	return false
}

func joinBoardWithCurrentUser(boardId int, e *httpexpect.Expect, t *testing.T, sessionId string) {
	var userJoinedBoardsResponse lib.GenericResponse
	e.POST(fmt.Sprintf("/api/v1/user/boards/%v", boardId)).
		WithCookie("sessionId", sessionId).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&userJoinedBoardsResponse)
	if userJoinedBoardsResponse.Status != "OK" {
		t.Fatal("Failed to join board")
	}
}

func isBoardInList(slug string, boards []lib.Board) bool {
	for _, board := range boards {
		if board.Slug == slug {
			return true
		}
	}
	return false
}
