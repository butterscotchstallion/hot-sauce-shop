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

var dbpool *pgxpool.Pool

func initDB() {
	var err error
	dbUrl := os.Getenv("DATABASE_URL")
	if len(dbUrl) == 0 {
		log.Fatalf("ERROR: Could not get DB url!")
	}

	dbpool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Connected to postgres")
}

func main() {
	initDB()
	defer dbpool.Close()

	r := gin.Default()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	routes.Products(r, dbpool, logger)
	routes.Tags(r, dbpool)
	routes.Cart(r, dbpool, logger)

	err := r.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
