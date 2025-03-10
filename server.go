package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"hotsauceshop/ent"
	"hotsauceshop/routes"
	_ "hotsauceshop/routes"
)

var client *ent.Client

func initDB() {
	var err error
	client, err = ent.Open("postgres",
		"host=localhost port=5432 sslmode=disable user=hotsauceboss "+
			"dbname=hotsauceshop password='unshaken context deniable shabby jimmy plunging'")
	if err != nil {
		log.Fatalf("Failed opening connection to postgres: %v", err)
	}
	log.Println("Connected to postgres")
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	log.Println("Database schema updated")
}

func closeClient() {
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("failed closing client: %v", err)
		}
	}(client)
}

func main() {
	initDB()

	r := gin.Default()
	routes.Products(r, client)
	routes.Tags(r, client)

	err := r.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
