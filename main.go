package main

import (
	"context"
	"log"

	_ "github.com/lib/pq"
	"hotsauceshop/ent"
)

func main() {
	client, err := ent.Open("postgres", "host=localhost port=5432 sslmode=disable user=hotsauceboss dbname=hotsauceshop password='unshaken context deniable shabby jimmy plunging'")
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	log.Println("Connected to postgres")
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("failed closing client: %v", err)
		}
	}(client)
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	log.Println("database schema created")
}
