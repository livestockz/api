package main

import (
	"database/sql"
	"log"
	"strconv"

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
	r.GET("/ping/:id", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		ID, err := strconv.ParseInt(id, 0, 64)
		if err != nil {
			ctx.JSON(404, gin.H{
				"error": err,
			})
		}
		if batches, err := repo.RessolveBatchByID(ID); err != nil {
			ctx.JSON(500, gin.H{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(200, gin.H{
				"data": batches,
			})
		}
	})

	r.Run(":9000")
}
