package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"hotsauceshop/ent"
	"hotsauceshop/ent/inventory"
	"hotsauceshop/ent/tag"
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

	r.GET("/api/v1/products/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		var res gin.H
		if len(slug) > 0 {
			product, err := client.Inventory.Query().
				Where(inventory.Slug(slug)).
				All(c)
			if err != nil {
				res = gin.H{
					"status":  "ERROR",
					"message": fmt.Sprintf("Error fetching product: %v", err),
				}
			} else {
				res = gin.H{
					"status": "OK",
					"results": gin.H{
						"product": product,
					},
				}
			}
		}
		c.JSON(200, res)
	})

	r.GET("/api/v1/products", func(c *gin.Context) {
		offset := c.DefaultQuery("offset", "1")
		perPage := c.DefaultQuery("perPage", "10")

		perPageInt, err := strconv.Atoi(perPage)
		if err != nil {
			perPageInt = 10
		}

		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			offsetInt = 1
		}

		total, err := client.Inventory.Query().
			Aggregate(ent.Count()).
			Int(c)
		if err != nil {
			log.Printf("Error getting total inventory items: %v", err)
		}

		var res gin.H
		inventoryResults, err := client.Inventory.
			Query().
			Offset(offsetInt).
			Limit(perPageInt).
			Order(ent.Asc(inventory.FieldName)).
			All(c)
		if err != nil {
			res = gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Error fetching inventory: %v", err),
			}
		} else {
			res = gin.H{
				"status": "OK",
				"results": gin.H{
					"inventory": inventoryResults,
					"total":     total,
				},
			}
		}

		c.JSON(200, res)
	})

	r.GET("/api/v1/tags", func(c *gin.Context) {
		var res gin.H
		tags, err := client.Tag.Query().Order(ent.Asc(tag.FieldName)).All(c)
		if err != nil {
			res = gin.H{
				"status":  "ERROR",
				"message": fmt.Sprintf("Error fetching tags: %v", err),
			}
		} else {
			res = gin.H{
				"status": "OK",
				"results": gin.H{
					"tags": tags,
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
