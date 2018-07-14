package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/livestockz/api/domain/batch"
	"github.com/livestockz/api/utils"
	uuid "github.com/satori/go.uuid"
)

type Handler interface {
	HealthHandler(*gin.Context)
	ResolveGrowthBatches(*gin.Context)
	ResolveGrowthBatchByID(*gin.Context)
	StoreGrowthBatch(*gin.Context)
	UpdateGrowthBatchByID(*gin.Context)
	RemoveGrowthBatchByID(*gin.Context)
}

type BatchHandler struct {
	BatchService batch.Service `inject:"batchService"`
}

func (h *BatchHandler) HealthHandler(c *gin.Context) {
	utils.Ok(c, nil)
}

func (h *BatchHandler) ResolveGrowthBatches(c *gin.Context) {
	//capture something like this: http://localhost:9090/growth/batch?page=1&limit=10
	q := c.Request.URL.Query()
	p := q.Get("page")
	l := q.Get("limit")
	page, _ := strconv.Atoi(p)
	limit, _ := strconv.Atoi(l)
	if batches, p, l, total, err := h.BatchService.ResolveGrowthBatches(int32(page), int32(limit)); err != nil {
		utils.Error(c, err)
		return
	} else {
		utils.Page(c, batches, p, l, total)
		return
	}
}

func (h *BatchHandler) ResolveGrowthBatchByID(c *gin.Context) {
	id := c.Params.ByName("id")
	uid, err := uuid.FromString(id)

	if err != nil {
		utils.Error(c, err)
		return
	}

	if batch, err := h.BatchService.ResolveGrowthBatchByID(uid); err != nil {
		utils.Error(c, err)
		return
	} else {
		utils.Ok(c, &batch)
		return
	}
}

func (h *BatchHandler) StoreGrowthBatch(c *gin.Context) {
	if batch, err := h.BatchService.StoreGrowthBatch(c); err != nil {
		utils.Error(c, err)
		return
	} else {
		utils.Created(c, &batch)
		return
	}
}

func (h *BatchHandler) UpdateGrowthBatchByID(c *gin.Context) {
	if _, err := h.BatchService.UpdateGrowthBatchByID(c); err != nil {
		utils.Error(c, err)
		return
	} else {
		utils.NoContent(c)
		return
	}
}

func (h *BatchHandler) RemoveGrowthBatchByID(c *gin.Context) {
	id := c.Params.ByName("id")
	uid, err := uuid.FromString(id)

	if err != nil {
		utils.Error(c, err)
		return
	}

	if batch, err := h.BatchService.RemoveGrowthBatchByID(uid); err != nil {
		utils.Error(c, err)
		return
	} else {
		utils.Ok(c, &batch)
		return
	}
}
