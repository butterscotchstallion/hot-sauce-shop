package routes

import (
	"fmt"
	"net/http"
	"testing"

	"hotsauceshop/lib"

	"github.com/gavv/httpexpect/v2"
)

func TestGetBoardSettings(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	sessionID := signInAndGetSessionId(t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword)

	addBoardResponse := CreateBoardAndVerify(t, e, sessionID)
	if addBoardResponse.Status != "OK" {
		t.Fatal("Failed to create board")
	}

	settings := lib.BoardSettings{
		IsOfficial:             true,
		IsPostApprovalRequired: true,
	}
	
	var boardSettingsResponse lib.BoardSettingsResponse
	e.GET(fmt.Sprintf("/api/v1/board-settings/%v", addBoardResponse.Results.Slug)).
		WithCookie("sessionId", sessionID).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&boardSettingsResponse)
	if boardSettingsResponse.Status != "OK" {
		t.Fatal("Failed to get board settings")
	}
}
