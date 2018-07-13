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
	StoreBatch(*gin.Context)
}

type BatchHandler struct {
	BatchService batch.Service `inject:"batchService"`
}

func (h *BatchHandler) HealthHandler(c *gin.Context) {
	utils.Ok(c, nil)
}

func (h *BatchHandler) ResolveBatchByID(c *gin.Context) {
	id := c.Params.ByName("id")
	uid, err := uuid.FromString(id)

	if err != nil {
		utils.Error(c, err)
		return
	}

	if batch, err := h.BatchService.ResolveBatchByID(uid); err != nil {
		utils.Error(c, err)
		return
	} else {
		utils.Ok(c, &batch)
		return
	}
}

func (h *BatchHandler) StoreBatch(c *gin.Context) {
	if batch, err := h.BatchService.CreateBatch(c); err != nil {
		utils.Error(c, err)
		return
	} else {
		utils.Ok(c, &batch)
	}
}
