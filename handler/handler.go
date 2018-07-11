package handler

import (
	"fmt"

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
	fmt.Print(h.BatchService)
	if batch, err := h.BatchService.ResolveBatchByID(uid); err != nil {
		utils.Error(c, err)
	} else {
		utils.Ok(c, &batch)
	}
	return
}
