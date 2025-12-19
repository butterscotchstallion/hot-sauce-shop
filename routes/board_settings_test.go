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

	// Create a new board
	addBoardResponse := CreateBoardAndVerify(t, e, sessionID)
	if addBoardResponse.Status != "OK" {
		t.Fatal("Failed to create board")
	}

	// Update board settings
	settings := lib.BoardSettingsUpdateRequest{
		IsOfficial:             true,
		IsPostApprovalRequired: true,
		BoardId:                addBoardResponse.Results.BoardId,
	}
	var updateSettingsResponse lib.GenericResponse
	e.PUT(fmt.Sprintf("/api/v1/board-settings/%v", addBoardResponse.Results.Slug)).
		WithCookie("sessionId", sessionID).
		WithJSON(settings).
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&updateSettingsResponse)
	if updateSettingsResponse.Status != "OK" {
		t.Fatal("Failed to update board settings")
	}

	// Verify board settings
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
	if boardSettingsResponse.Results.BoardId != settings.BoardId {
		t.Fatal("boardId mismatch")
	}
	if boardSettingsResponse.Results.IsOfficial != settings.IsOfficial {
		t.Fatal("isOfficial mismatch")
	}
	if boardSettingsResponse.Results.IsPostApprovalRequired != settings.IsPostApprovalRequired {
		t.Fatal("IsPostApprovalRequired mismatch")
	}

	DeleteBoardAndVerify(t, e, sessionID, addBoardResponse.Results.Slug)
}

func TestGetBoardSettingsWithUnprivilegedUser(t *testing.T) {
	e := httpexpect.Default(t, config.Server.AddressWithProtocol)
	unprivSessionID := signInAndGetSessionId(
		t, e, config.TestUsers.UnprivilegedUsername, config.TestUsers.UnprivilegedPassword,
	)

	adminSessionID := signInAndGetSessionId(
		t, e, config.TestUsers.BoardAdminUsername, config.TestUsers.BoardAdminPassword,
	)

	// Create a new board
	addBoardResponse := CreateBoardAndVerify(t, e, adminSessionID)
	if addBoardResponse.Status != "OK" {
		t.Fatal("Failed to create board")
	}

	// Get settings (should be 403)
	var forbiddenSettingsResponse lib.GenericResponse
	e.GET(fmt.Sprintf("/api/v1/board-settings/%v", addBoardResponse.Results.Slug)).
		WithCookie("sessionId", unprivSessionID).
		Expect().
		Status(http.StatusForbidden).
		JSON().
		Decode(&forbiddenSettingsResponse)
	if forbiddenSettingsResponse.Status != "ERROR" {
		t.Fatal("Expected 403 Forbidden response")
	}
}
