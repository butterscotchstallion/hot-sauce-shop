package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"hotsauceshop/ent"
)

var client *ent.Client

func initDB() {
	var err error
	client, err = ent.Open("postgres",
		"host=localhost port=5432 sslmode=disable user=hotsauceboss "+
			"dbname=hotsauceshop password='unshaken context deniable shabby jimmy plunging'")
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	log.Println("Connected to postgres")
	/*defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("failed closing client: %v", err)
		}
	}(client)*/
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	log.Println("database schema updated")
}

func main() {
	initDB()

	r := gin.Default()

	r.GET("/api/v1/products", func(c *gin.Context) {
		var res gin.H
		inventory, err := client.Inventory.Query().All(c)
		if err != nil {
			res = gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Error fetching inventory: %v", err),
			}
		} else {
			res = gin.H{
				"status": "OK",
				"results": gin.H{
					"inventory": inventory,
				},
			}
		}
		c.JSON(200, res)
	})

	err := r.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
