package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/livestockz/api/domain/batch"
)

func main() {
	//set up database
	//1. Connect to database
	db, err := sql.Open("mysql", "root@/livestock?parseTime=true")
	if err != nil {
		log.Print(err)
		panic("Failed connect to database.")
	}

	defer db.Close()

	repo := batch.NewRepository(db)
	// ids := []int32{1,2,3}
	// RessolveBatchByIDs(ids...)

	//set up http server
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		if batches, err := repo.RessolveBatchByID(1); err != nil {
			ctx.JSON(500, gin.H{
				"error": err,
			})
		} else {
			ctx.JSON(200, gin.H{
				"data": batches,
			})
		}
	})

	r.Run(":9000")
}
