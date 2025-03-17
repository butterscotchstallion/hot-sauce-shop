package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
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
	routes.Products(r, dbPool, logger)
	routes.Tags(r, dbPool)
	routes.Cart(r, dbPool, logger)

	err := r.Run("localhost:8081")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
