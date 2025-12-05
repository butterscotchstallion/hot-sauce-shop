package lib

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(dsn string) *pgxpool.Pool {
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Connected to postgres")
	return dbPool
}
