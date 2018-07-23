package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/livestockz/api/config"
	"github.com/livestockz/api/domain/batch"
	"github.com/livestockz/api/domain/feed"
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
	feedHandler := new(handler.FeedHandler)
	//BatchService := new(batch.BatchService)

	//register service
	r := gin.Default()
	sc := gocontainer.NewContainer()
	sc.RegisterService("db", db)
	sc.RegisterService("config", cfg)
	sc.RegisterService("batchHandler", batchHandler)
	sc.RegisterService("feedHandler", feedHandler)
	sc.RegisterService("batchService", new(batch.BatchService))
	sc.RegisterService("feedService", new(feed.FeedService))
	sc.RegisterService("batchRepository", new(batch.BatchRepository))
	sc.RegisterService("feedRepository", new(feed.FeedRepository))
	sc.HandleGracefulShutdown(3 * time.Second)
	if err := sc.Ready(); err != nil {
		panic("Failed to start service container")
	}

	growth := r.Group("/growth")
	{
		//batch
		growth.GET("/batch", batchHandler.ResolveGrowthBatchPage)
		growth.GET("/batch/:id", batchHandler.ResolveGrowthBatchByID)
		growth.POST("/batch", batchHandler.StoreGrowthBatch)
		growth.PUT("/batch/:id", batchHandler.StoreGrowthBatch)
		growth.DELETE("/batch", batchHandler.RemoveGrowthBatchByIDs)
		growth.DELETE("/batch/:id", batchHandler.RemoveGrowthBatchByID)
		//pool
		growth.GET("/pool", batchHandler.ResolveGrowthPoolPage)
		growth.GET("/pool/:id", batchHandler.ResolveGrowthPoolByID)
		growth.POST("/pool", batchHandler.StoreGrowthPool)
		growth.PUT("/pool/:id", batchHandler.StoreGrowthPool)
		growth.DELETE("/pool", batchHandler.RemoveGrowthPoolByIDs)
		growth.DELETE("/pool/:id", batchHandler.RemoveGrowthPoolByID)
	}
	feed := r.Group("/feed")
	{
		//feed type
		feed.GET("/feed-type", feedHandler.ResolveFeedTypePage)
		feed.GET("/feed-type/:id", feedHandler.ResolveFeedTypeByID)
		feed.POST("/feed-type", feedHandler.StoreFeedType)
		feed.PUT("/feed-type/:id", feedHandler.StoreFeedType)
		feed.DELETE("/feed-type", feedHandler.RemoveFeedTypeByIDs)
		feed.DELETE("/feed-type/:id", feedHandler.RemoveFeedTypeByID)
		//feeding
		feed.GET("/feeding", feedHandler.ResolveFeedPage)
		feed.GET("/feeding/:id", feedHandler.ResolveFeedByID)
		feed.POST("/feeding", feedHandler.StoreFeed)
		feed.PUT("/feeding/:id", feedHandler.StoreFeed)
		feed.DELETE("/feeding", feedHandler.RemoveFeedByIDs)
		feed.DELETE("/feeding/:id", feedHandler.RemoveFeedByID)
	}

	r.GET("/health", batchHandler.HealthHandler)
	r.Run(":9090")
}
