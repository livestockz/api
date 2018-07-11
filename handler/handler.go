package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/livestockz/api/domain/batch"
	"github.com/livestockz/api/utils"
	uuid "github.com/satori/go.uuid"
)

type Handler interface {
	HealthHandler(*gin.Context)
	ResolveBatchByID(*gin.Context)
}

type BatchHandler struct {
	BatchService batch.Service `inject:"BatchService"`
}

func HealthHandler(c *gin.Context) {
	utils.Ok(c, nil)
}

func (h *BatchHandler) ResolveBatchByID(c *gin.Context) {
	id := c.Params.ByName("id")
	ID, err := uuid.FromString(id)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	if batches, err := h.BatchService.ResolveBatchByID(ID); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"data": batches,
		})
	}
}
