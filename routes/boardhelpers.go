package routes

import (
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
}

func getPostListAndVerify(params GetPostListAndVerifyParams) lib.PostListResponse {
	var postListResponse lib.PostListResponse
	showUnapprovedStr := 0
	if params.ShowUnapproved {
		showUnapprovedStr = 1
	}
	params.E.GET("/api/v1/posts").
		WithQuery("showUnapproved", showUnapprovedStr).
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

func isPostSlugInList(postSlug string, posts []lib.BoardPost) bool {
	for _, post := range posts {
		if post.Slug == postSlug {
			return true
		}
	}
	return false
}
