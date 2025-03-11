package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"hotsauceshop/routes"
)

var conn *pgx.Conn

func initDB() {
	var err error
	dbUrl := os.Getenv("DATABASE_URL")
	if len(dbUrl) == 0 {
		log.Fatalf("ERROR: Could not get DB url!")
	}

	conn, err = pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Connected to postgres")
}

func main() {
	initDB()

	r := gin.Default()
	routes.Products(r, client)
	routes.Tags(r, client)
	routes.Cart(r, client)

	err := r.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
