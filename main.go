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
	batchHandler := new(handler.BatchHandler)
	//BatchService := new(batch.BatchService)

	//register service
	r := gin.Default()
	sc := gocontainer.NewContainer()
	sc.RegisterService("db", db)
	sc.RegisterService("config", cfg)
	sc.RegisterService("batchHandler", batchHandler)
	sc.RegisterService("batchService", new(batch.BatchService))
	sc.RegisterService("batchRepository", new(batch.BatchRepository))
	sc.HandleGracefulShutdown(3 * time.Second)
	if err := sc.Ready(); err != nil {
		panic("Failed to start service container")
	}
	r.GET("/health", batchHandler.HealthHandler)
	r.GET("/ping/:id", batchHandler.ResolveBatchByID)
	r.Run(":9090")
}
