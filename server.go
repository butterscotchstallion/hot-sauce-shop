package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"hotsauceshop/routes"
)

var dbPool *pgxpool.Pool

func initDB() {
	var err error
	dbUrl := os.Getenv("DATABASE_URL")
	if len(dbUrl) == 0 {
		log.Fatalf("ERROR: Could not get DB url!")
	}

	dbPool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Connected to postgres")
}

func main() {
	initDB()
	defer dbPool.Close()

	r := gin.Default()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	var wsConn *websocket.Conn
	routes.WS(r, wsConn, logger)
	routes.Products(r, dbPool, wsConn, logger)
	routes.Tags(r, dbPool)
	routes.Cart(r, dbPool, logger)
	routes.User(r, dbPool, logger)
	routes.Session(r, dbPool, logger)

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(wsConn)

	err := r.Run("localhost:8081")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func setUpEventBus() {

}
