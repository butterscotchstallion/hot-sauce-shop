package routes

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"hotsauceshop/lib"

	"github.com/gavv/httpexpect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func SetupSuite(tb testing.TB) {
	log.Println("Setting up test suite")
	dbPool = lib.InitDB()
}

func TestGetBoardPosts(t *testing.T) {
	e := httpexpect.New(t, "http://localhost:8081")
	e.GET("/api/v1/boards").
		Expect().
		Status(http.StatusOK).JSON().Object().
		Value("results").Object().
		Value("boards").Array().Length().Gt(0)
}

func TestCreateBoard(t *testing.T) {
	/**
	 * 1. Create a new board
	 * 2. Check board detail to confirm it exists
	 * 3. Verify response
	 * 4. Delete board
	 */
	boardUUID, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("Error getting board name")
	}
	testBoardName := boardUUID.String()
	errAddingBoard := lib.AddBoard(dbPool, testBoardName, testBoardName, "", 1, "test board")
	if errAddingBoard != nil {
		t.Fatalf("Error adding board: %v", errAddingBoard)
	}

	e := httpexpect.New(t, "http://localhost:8081")
	e.GET(fmt.Sprintf("/api/v1/boards/%v", testBoardName)).
		Expect().
		Status(http.StatusOK).JSON().Object().
		Value("results").Object().
		Value("board").NotNull()
}
