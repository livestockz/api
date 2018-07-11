package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/livestockz/api/config"
	"github.com/livestockz/api/domain/batch"
	"github.com/livestockz/api/handler"
	"github.com/ncrypthic/gocontainer"
)

func main() {
	// Setup deps
	// 1. database
	cfg := config.Get()
	db, err := sql.Open("mysql", cfg.DatabaseDSN())
	if err != nil {
		log.Print(err)
		panic("Failed connect to database.")
	}

	defer db.Close()

	//register Handler
	BatchHandler := new(handler.BatchHandler)

	//register service
	r := gin.Default()
	sc := gocontainer.NewContainer()
	sc.RegisterService("db", db)
	sc.RegisterService("config", cfg)
	sc.RegisterService("BatchHandler", new(handler.BatchHandler))
	sc.RegisterService("BatchService", new(batch.BatchService))
	sc.RegisterService("BatchRepository", new(batch.BatchRepository))
	sc.HandleGracefulShutdown(3 * time.Second)
	if err := sc.Ready(); err != nil {
		panic("Failed to start service container")
	}
	r.GET("/health", BatchHandler.HealthHandler)
	r.GET("/ping/:id", BatchHandler.ResolveBatchByID)
	r.Run(":9090")
}
