package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"hotsauceshop/lib"
	"hotsauceshop/routes"

	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func main() {
	gin.SetMode(gin.DebugMode)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	config, configReadErr := lib.ReadConfig("config.toml")
	if configReadErr != nil {
		panic("Could not read config")
	}
	dbPool = lib.InitDB(config.Database.Dsn)
	defer dbPool.Close()
	// lib.SetRuntimeConfig(config)

	err := os.Setenv("TZ", config.Server.TimeZone)
	if err != nil {
		logger.Error(fmt.Sprintf("Error setting timezone: %v", err))
	}

	r := gin.Default()

	store := persistence.NewInMemoryStore(time.Minute * config.Cache.DefaultCacheTime)
	var wsConn *websocket.Conn
	routes.WS(r, wsConn, logger)
	routes.Products(r, dbPool, logger, store)
	routes.Tags(r, dbPool, store)
	routes.Cart(r, dbPool, logger)
	routes.User(r, dbPool, logger)
	routes.Session(r, dbPool, logger)
	routes.Admin(r, dbPool, logger, store)
	routes.Orders(r, dbPool, logger)
	routes.Boards(r, dbPool, logger, store)
	routes.Votes(r, dbPool, logger)

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(wsConn)

	runErr := r.Run(config.Server.Address)
	if runErr != nil {
		log.Fatalf("Error starting server: %v", runErr)
	}
}
