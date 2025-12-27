package routes

import (
	"net/http"
	"testing"

	"hotsauceshop/lib"

	"github.com/gavv/httpexpect/v2"
)

type GetPostListAndVerifyParams struct {
	E              *httpexpect.Expect
	T              *testing.T
	SessionId      string
	ShowUnapproved bool
}

func getPostListAndVerify(params GetPostListAndVerifyParams) lib.PostListResponse {
	var postListResponse lib.PostListResponse
	// showUnapprovedStr := "false"
	// if params.ShowUnapproved {
	// 	showUnapprovedStr = "true"
	// }
	params.E.GET("/api/v1/posts").
		WithCookie("sessionId", params.SessionId).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&postListResponse)
	if postListResponse.Status != "OK" {
		params.T.Fatal("Failed to get post list")
	}
	if len(postListResponse.Results.Posts) == 0 {
		params.T.Fatal("Post list is empty!")
	}
	return postListResponse
}
