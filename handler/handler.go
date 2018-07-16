package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/livestockz/api/domain/batch"
	"github.com/livestockz/api/utils"
	uuid "github.com/satori/go.uuid"
)

type Handler interface {
	HealthHandler(*gin.Context)
	ResolveGrowthBatchPage(*gin.Context)
	ResolveGrowthBatchByID(*gin.Context)
	StoreGrowthBatch(*gin.Context)
	RemoveGrowthBatchByID(*gin.Context)
}

type BatchHandler struct {
	BatchService batch.Service `inject:"batchService"`
}

func (h *BatchHandler) HealthHandler(c *gin.Context) {
	utils.Ok(c, nil)
}

func (h *BatchHandler) ResolveGrowthBatchPage(c *gin.Context) {
	//capture something like this: http://localhost:9090/growth/batch?page=1&limit=10
	q := c.Request.URL.Query()
	p := q.Get("page")
	l := q.Get("limit")
	page, err := strconv.Atoi(p)
	if err != nil {
		//log.Print("Invalid Page")
		//utils.Error(c, fmt.Errorf("Invalid Page"))
		//return
		page = 0
	}
	limit, err := strconv.Atoi(l)
	if err != nil {
		//log.Print("Invalid Limit Set")
		//utils.Error(c, fmt.Errorf("Invalid Limit Set"))
		//return
		limit = 10
	}
	if batches, p, l, total, err := h.BatchService.ResolveGrowthBatchPage(int32(page), int32(limit)); err != nil {
		//log.Print(err.Error())
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

	var id = c.Params.ByName("id")
	var batch batch.Batch
	c.BindJSON(&batch)

	if id == "" {
		if batch.Name == "" {
			utils.Error(c, fmt.Errorf("Incomplete provided data."))
		} else if result, err := h.BatchService.StoreGrowthBatch(&batch); err != nil {
			utils.Error(c, err)
			return
		} else {
			utils.Created(c, &result)
			return
		}
	} else {
		//convert id to UUID
		//compare uuid to batch
		//save if valid
		var uid, err = uuid.FromString(id)
		if err != nil {
			utils.Error(c, fmt.Errorf("Unable to convert given ID to UUID"))
			//fmt.Print("Unable to convert given ID to UUID")
			return
		} else if batch.ID != uid {
			utils.Error(c, fmt.Errorf("Inconsistent ID."))
			//fmt.Print("Inconsistent ID.")
			return
		} else if batch.Name == "" {
			utils.Error(c, fmt.Errorf("Incomplete provided data."))
			//fmt.Print("Incomplete provided data.")
			return
		} else if result, err := h.BatchService.StoreGrowthBatch(&batch); err != nil {
			utils.Error(c, err)
			return
		} else {
			utils.Ok(c, &result)
		}
	}
}

func (h *BatchHandler) RemoveGrowthBatchByID(c *gin.Context) {
	id := c.Params.ByName("id")
	uid, err := uuid.FromString(id)

	if err != nil {
		utils.Error(c, err)
		return
	}

	if _, err := h.BatchService.RemoveGrowthBatchByID(uid); err != nil {
		utils.Error(c, err)
	} else {
		utils.NoContent(c)
	}
}
