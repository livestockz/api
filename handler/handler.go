package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/livestockz/api/domain/batch"
	"github.com/livestockz/api/domain/feed"
	"github.com/livestockz/api/utils"
	uuid "github.com/satori/go.uuid"
)

type BatchHandler struct {
	BatchService batch.Service `inject:"batchService"`
}

type FeedHandler struct {
	FeedService feed.Service `inject:"feedService"`
}

type UUIDRequestModel struct {
	Data []string `json:"ids"`
}

func (h *BatchHandler) HealthHandler(c *gin.Context) {
	utils.Ok(c, nil)
}

func (h *BatchHandler) ResolveGrowthBatchPage(c *gin.Context) {
	//capture something like this: http://localhost:9090/growth/batch?page=1&limit=10
	q := c.Request.URL.Query()
	p := q.Get("page")
	l := q.Get("limit")
	d := q.Get("deleted")
	page, err := strconv.Atoi(p)
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 10
	}

	if d != batch.Deleted_Any && d != batch.Deleted_False && d != batch.Deleted_True {
		utils.Error(c, fmt.Errorf("Unknown deleted status"))
	} else if batches, p, l, total, err := h.BatchService.ResolveGrowthBatchPage(int32(page), int32(limit), d); err != nil {
		utils.Error(c, err)
	} else {
		utils.Page(c, batches, p, l, total)
	}
	return
}

func (h *BatchHandler) ResolveGrowthBatchByID(c *gin.Context) {
	id := c.Params.ByName("batchId")
	uid, err := uuid.FromString(id)

	if err != nil {
		utils.Error(c, err)
	} else if batch, err := h.BatchService.ResolveGrowthBatchByID(uid); err != nil {
		utils.Error(c, err)
	} else {
		utils.Ok(c, &batch)
	}
	return
}

func (h *BatchHandler) StoreGrowthBatch(c *gin.Context) {

	var id = c.Params.ByName("batchId")
	var batch batch.Batch
	c.BindJSON(&batch)

	if id == "" {
		if batch.Name == "" {
			utils.Error(c, fmt.Errorf("Incomplete provided data."))
		} else if result, err := h.BatchService.StoreGrowthBatch(&batch); err != nil {
			utils.Error(c, err)
		} else {
			utils.Created(c, &result)
		}
		return
	} else {
		//convert id to UUID
		//compare uuid to batch
		//save if valid
		var uid, err = uuid.FromString(id)
		if err != nil {
			utils.Error(c, fmt.Errorf("Unable to convert given ID to UUID"))
		} else if batch.ID != uid {
			utils.Error(c, fmt.Errorf("Inconsistent ID."))
		} else if batch.Name == "" {
			utils.Error(c, fmt.Errorf("Incomplete provided data."))
		} else if result, err := h.BatchService.StoreGrowthBatch(&batch); err != nil {
			utils.Error(c, err)
		} else {
			utils.Ok(c, &result)
		}
		return
	}
}

func (h *BatchHandler) RemoveGrowthBatchByID(c *gin.Context) {
	id := c.Params.ByName("batchId")
	uid, err := uuid.FromString(id)

	if err != nil {
		utils.Error(c, err)
	} else if _, err := h.BatchService.RemoveGrowthBatchByID(uid); err != nil {
		utils.Error(c, err)
	} else {
		utils.NoContent(c)
	}
	return
}

func (h *BatchHandler) RemoveGrowthBatchByIDs(c *gin.Context) {
	//process json like : {"ids":["0b86bef7-0e16-47e6-9463-6a0b583e8d4c","6be6e63c-18f3-48ce-831f-f3210a576945"]}
	var ids []uuid.UUID
	reqBody := new(UUIDRequestModel)
	err := c.Bind(reqBody)
	if err != nil {
		utils.Error(c, err)
	} else if len(reqBody.Data) < 1 {
		utils.Error(c, fmt.Errorf("No Batch to be removed."))
	} else {
		for _, v := range reqBody.Data {
			//convert to UUID
			id, err := uuid.FromString(v)
			if err != nil {
				utils.Error(c, err)
				return
			} else {
				ids = append(ids, id)
			}
		}
		//process to services
		_, err := h.BatchService.RemoveGrowthBatchByIDs(ids)
		if err != nil {
			utils.Error(c, err)
		} else {
			utils.NoContent(c)
			return
		}
	}
	return
}

//pool
func (h *BatchHandler) ResolveGrowthPoolPage(c *gin.Context) {
	//capture something like this: http://localhost:9090/growth/pool?page=1&limit=10
	q := c.Request.URL.Query()
	p := q.Get("page")
	l := q.Get("limit")
	d := q.Get("deleted")
	page, err := strconv.Atoi(p)
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 10
	}

	if d != batch.Deleted_Any && d != batch.Deleted_False && d != batch.Deleted_True {
		utils.Error(c, fmt.Errorf("Unknown deleted status"))
	} else if pools, p, l, total, err := h.BatchService.ResolveGrowthPoolPage(int32(page), int32(limit), d); err != nil {
		utils.Error(c, err)
	} else {
		utils.Page(c, pools, p, l, total)
	}
	return
}

func (h *BatchHandler) ResolveGrowthPoolByID(c *gin.Context) {
	id := c.Params.ByName("poolId")
	uid, err := uuid.FromString(id)

	if err != nil {
		utils.Error(c, err)
	} else if pool, err := h.BatchService.ResolveGrowthPoolByID(uid); err != nil {
		utils.Error(c, err)
	} else {
		utils.Ok(c, &pool)
	}
	return
}

func (h *BatchHandler) StoreGrowthPool(c *gin.Context) {

	var id = c.Params.ByName("poolId")
	var pool batch.Pool
	c.BindJSON(&pool)

	if id == "" {
		if pool.Name == "" {
			utils.Error(c, fmt.Errorf("Incomplete provided data."))
		} else if pool.Status != batch.Pool_Assigned && pool.Status != batch.Pool_Inactive && pool.Status != batch.Pool_Maintenance {
			utils.Error(c, fmt.Errorf("Invalid pool status."))
		} else if result, err := h.BatchService.StoreGrowthPool(&pool); err != nil {
			utils.Error(c, err)
		} else {
			utils.Created(c, &result)
		}
		return
	} else {
		//convert id to UUID
		//compare uuid to pool
		//save if valid
		var uid, err = uuid.FromString(id)
		if err != nil {
			utils.Error(c, fmt.Errorf("Unable to convert given ID to UUID"))
		} else if pool.ID != uid {
			utils.Error(c, fmt.Errorf("Inconsistent ID."))
		} else if pool.Name == "" {
			utils.Error(c, fmt.Errorf("Incomplete provided data."))
		} else if pool.Status != batch.Pool_Assigned && pool.Status != batch.Pool_Inactive && pool.Status != batch.Pool_Maintenance {
			utils.Error(c, fmt.Errorf("Invalid pool status."))
		} else if result, err := h.BatchService.StoreGrowthPool(&pool); err != nil {
			utils.Error(c, err)
		} else {
			utils.Ok(c, &result)
		}
		return
	}
}

func (h *BatchHandler) RemoveGrowthPoolByID(c *gin.Context) {
	id := c.Params.ByName("poolId")
	uid, err := uuid.FromString(id)

	if err != nil {
		utils.Error(c, err)
	} else if _, err := h.BatchService.RemoveGrowthPoolByID(uid); err != nil {
		utils.Error(c, err)
	} else {
		utils.NoContent(c)
	}
	return
}

func (h *BatchHandler) RemoveGrowthPoolByIDs(c *gin.Context) {
	//process json like : {"ids":["0b86bef7-0e16-47e6-9463-6a0b583e8d4c","6be6e63c-18f3-48ce-831f-f3210a576945"]}
	var ids []uuid.UUID
	reqBody := new(UUIDRequestModel)
	err := c.Bind(reqBody)
	if err != nil {
		utils.Error(c, err)
	} else if len(reqBody.Data) < 1 {
		utils.Error(c, fmt.Errorf("No Pool to be removed."))
	} else {
		for _, v := range reqBody.Data {
			//convert to UUID
			id, err := uuid.FromString(v)
			if err != nil {
				utils.Error(c, err)
				return
			} else {
				ids = append(ids, id)
			}
		}
		//process to services
		_, err := h.BatchService.RemoveGrowthPoolByIDs(ids)
		if err != nil {
			utils.Error(c, err)
		} else {
			utils.NoContent(c)
		}
		return
	}
	return
}

//growth batch cycle
func (h *BatchHandler) ResolveGrowthBatchCyclePage(c *gin.Context) {
	//capture something like this: http://localhost:9090/growth/batch-cycle?page=1&limit=10
	q := c.Request.URL.Query()
	p := q.Get("page")
	l := q.Get("limit")
	id := c.Params.ByName("batchId")

	batchId, err := uuid.FromString(id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	page, err := strconv.Atoi(p)
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 10
	}

	if batchCycles, p, l, total, err := h.BatchService.ResolveGrowthBatchCyclePage(batchId, int32(page), int32(limit)); err != nil {
		utils.Error(c, err)
	} else {
		utils.Page(c, batchCycles, p, l, total)
	}
	return
}

func (h *BatchHandler) ResolveGrowthBatchCycleByID(c *gin.Context) {
	bid := c.Params.ByName("batchId")
	cid := c.Params.ByName("cycleId")
	batchId, err := uuid.FromString(bid)
	if err != nil {
		utils.Error(c, err)
		return
	}

	cycleId, err := uuid.FromString(cid)
	if err != nil {
		utils.Error(c, err)
		return
	}

	if batchCycle, err := h.BatchService.ResolveGrowthBatchCycleByID(batchId, cycleId); err != nil {
		utils.Error(c, err)
	} else {
		utils.Ok(c, &batchCycle)
	}
	return
}

func (h *BatchHandler) StoreGrowthBatchCycle(c *gin.Context) {

	var bid = c.Params.ByName("batchId")
	var cid = c.Params.ByName("cycleId")

	var bc batch.BatchCycle
	c.BindJSON(&bc)
	if bid == "" {
		utils.Error(c, fmt.Errorf("Invalid batch id."))
	} else if cid == "" {
		_, err := uuid.FromString(bid)
		if err != nil {
			utils.Error(c, err)
		} else if bc.Weight == 0 || bc.Amount == 0 {
			utils.Error(c, fmt.Errorf("Incomplete provided data."))
		} else if result, err := h.BatchService.StoreGrowthBatchCycle(&bc); err != nil {
			utils.Error(c, err)
		} else {
			utils.Created(c, &result)
		}
		return
	} else {
		//convert id to UUID
		//compare uuid to batch cycle
		//save if valid
		batchId, err := uuid.FromString(bid)
		if err != nil {
			utils.Error(c, err)
			return
		}

		cid := c.Params.ByName("cycleId")
		cycleId, err := uuid.FromString(cid)
		if err != nil {
			utils.Error(c, err)
			return
		}

		if bc.Batch.ID != batchId {
			utils.Error(c, fmt.Errorf("Inconsistent ID."))
		} else if bc.ID != cycleId {
			utils.Error(c, fmt.Errorf("Inconsistent ID."))
		} else if bc.Weight == 0 || bc.Amount == 0 {
			utils.Error(c, fmt.Errorf("Incomplete provided data."))
		} else if result, err := h.BatchService.StoreGrowthBatchCycle(&bc); err != nil {
			utils.Error(c, err)
		} else {
			utils.Ok(c, &result)
		}
		return
	}
}

//growth batch cycle death
func (h *BatchHandler) StoreGrowthDeath(c *gin.Context) {
	var bid = c.Params.ByName("batchId")
	var cid = c.Params.ByName("cycleId")

	var bcd batch.Death
	c.BindJSON(&bcd)

	if bid == "" {
		utils.Error(c, fmt.Errorf("Invalid batch id."))
	} else if cid == "" {
		utils.Error(c, fmt.Errorf("Invalid cycle id."))
	} else if bcd.Amount == 0 || bcd.Weight == 0 {
		utils.Error(c, fmt.Errorf("Incomplete data."))
	} else if _, err := uuid.FromString(bid); err != nil {
		utils.Error(c, err)
	} else if cycleId, err := uuid.FromString(cid); err != nil {
		utils.Error(c, err)
	} else if bcd.BatchCycleID != cycleId {
		utils.Error(c, fmt.Errorf("Inconsistent cycle id."))
	} else if result, err := h.BatchService.StoreGrowthDeath(&bcd); err != nil {
		utils.Error(c, err)
	} else {
		utils.Ok(c, &result)
	}
	return
}

//growth batch cycle feeding
func (h *BatchHandler) StoreGrowthFeeding(c *gin.Context) {
	var bid = c.Params.ByName("batchId")
	var cid = c.Params.ByName("cycleId")

	var feeding batch.Feeding
	c.BindJSON(&feeding)

	if bid == "" {
		utils.Error(c, fmt.Errorf("Invalid batch id."))
	} else if cid == "" {
		utils.Error(c, fmt.Errorf("Invalid cycle id."))
	} else if feeding.Qty == 0 {
		utils.Error(c, fmt.Errorf("Incomplete data."))
	} else if _, err := uuid.FromString(bid); err != nil {
		utils.Error(c, err)
	} else if cycleId, err := uuid.FromString(cid); err != nil {
		utils.Error(c, err)
	} else if feeding.BatchCycleID != cycleId {
		utils.Error(c, fmt.Errorf("Inconsistent cycle id."))
	} else if result, err := h.BatchService.StoreGrowthFeeding(&feeding); err != nil {
		utils.Error(c, err)
	} else {
		utils.Ok(c, &result)
	}
	return
}

//growth batch cycle cut off
func (h *BatchHandler) StoreGrowthCutOff(c *gin.Context) {
	var bid = c.Params.ByName("batchId")
	var cid = c.Params.ByName("cycleId")

	var cutoff batch.CutOff
	c.BindJSON(&cutoff)

	if bid == "" {
		utils.Error(c, fmt.Errorf("Invalid batch id."))
	} else if cid == "" {
		utils.Error(c, fmt.Errorf("Invalid cycle id."))
	} else if cutoff.Weight == 0 || cutoff.Amount == 0 {
		utils.Error(c, fmt.Errorf("Incomplete data."))
	} else if result, err := h.BatchService.StoreGrowthCutOff(&cutoff); err != nil {
		utils.Error(c, err)
	} else {
		utils.Ok(c, &result)
	}
	return
}

//growth sales
func (h *BatchHandler) ResolveGrowthSalesByID(c *gin.Context) {
	var sid = c.Params.ByName("salesId")
	if sid == "" {
		utils.Error(c, fmt.Errorf("Invalid Sales ID"))
	} else if salesId, err := uuid.FromString(sid); err != nil {
		utils.Error(c, err)
	} else if salesId == uuid.Nil {
		utils.Error(c, fmt.Errorf("Sales Id cannot be null"))
	} else if result, err := h.BatchService.ResolveGrowthSalesByID(salesId); err != nil {
		utils.Error(c, err)
	} else {
		utils.Ok(c, result)
	}
	return
}
func (h *BatchHandler) StoreGrowthSales(c *gin.Context) {
	var sid = c.Params.ByName("salesId")

	var sales batch.Sales
	c.BindJSON(&sales)

	if sid == "" && sales.ID != uuid.Nil {
		utils.Error(c, fmt.Errorf("Invalid sales id."))
	} else if sales.Qty == 0 {
		utils.Error(c, fmt.Errorf("Incomplete data."))
	} else if sid == "" && sales.ID == uuid.Nil {
		if result, err := h.BatchService.StoreGrowthSales(&sales); err != nil {
			utils.Error(c, err)
		} else {
			utils.Ok(c, &result)
		}
	} else {
		if salesId, err := uuid.FromString(sid); err != nil {
			utils.Error(c, err)
		} else if salesId != sales.ID {
			utils.Error(c, fmt.Errorf("Mismatch given sales Id."))
		} else if result, err := h.BatchService.StoreGrowthSales(&sales); err != nil {
			utils.Error(c, err)
		} else {
			utils.Ok(c, &result)
		}
	}
	return
}

//growth sales detail
func (h *BatchHandler) StoreGrowthSalesDetail(c *gin.Context) {
	var sid = c.Params.ByName("salesId")

	var sales batch.Sales
	c.BindJSON(&sales)

	if salesId, err := uuid.FromString(sid); err != nil {
		utils.Error(c, err)
	} else if salesId != sales.ID {
		utils.Error(c, fmt.Errorf("Mismatch given sales Id."))
	} else if sales.Qty == 0 {
		utils.Error(c, fmt.Errorf("Qty cannot be empty."))
	} else {
		for _, detail := range sales.Detail {
			if detail.BatchCycleID == uuid.Nil {
				utils.Error(c, fmt.Errorf("Invalid batch cycle id."))
				return
			} else if detail.Amount == 0 || detail.Weight == 0 {
				utils.Error(c, fmt.Errorf("Weight and Amount cannot be empty"))
				return
			}
		}

		if result, err := h.BatchService.StoreGrowthSalesDetail(&sales); err != nil {
			utils.Error(c, err)
		} else {
			utils.Ok(c, result)
		}
	}
	return
}

//feedtype
func (h *FeedHandler) ResolveFeedTypePage(c *gin.Context) {
	//capture something like this: http://localhost:9090/feed/feed-type?page=1&limit=10
	q := c.Request.URL.Query()
	p := q.Get("page")
	l := q.Get("limit")
	d := q.Get("deleted")
	page, err := strconv.Atoi(p)
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 10
	}
	if d != feed.Deleted_Any && d != feed.Deleted_False && d != feed.Deleted_True {
		utils.Error(c, fmt.Errorf("Unknown deleted status"))
	} else if feedtypes, p, l, total, err := h.FeedService.ResolveFeedTypePage(int32(page), int32(limit), d); err != nil {
		utils.Error(c, err)
	} else {
		utils.Page(c, feedtypes, p, l, total)
	}
	return
}

func (h *FeedHandler) ResolveFeedTypeByID(c *gin.Context) {
	id := c.Params.ByName("id")
	uid, err := uuid.FromString(id)

	if err != nil {
		utils.Error(c, err)
	} else if feedtype, err := h.FeedService.ResolveFeedTypeByID(uid); err != nil {
		utils.Error(c, err)
	} else {
		utils.Ok(c, &feedtype)
	}
	return
}

func (h *FeedHandler) StoreFeedType(c *gin.Context) {

	var id = c.Params.ByName("id")
	var feedtype feed.FeedType
	c.BindJSON(&feedtype)

	if id == "" {
		if feedtype.Name == "" || feedtype.Unit == "" {
			utils.Error(c, fmt.Errorf("Incomplete provided data."))
		} else if result, err := h.FeedService.StoreFeedType(&feedtype); err != nil {
			utils.Error(c, err)
		} else {
			utils.Created(c, &result)
		}
		return
	} else {
		//convert id to UUID
		//compare uuid to FeedType.ID
		//save if valid
		var uid, err = uuid.FromString(id)
		if err != nil {
			utils.Error(c, fmt.Errorf("Unable to convert given ID to UUID"))
		} else if feedtype.ID != uid {
			utils.Error(c, fmt.Errorf("Inconsistent ID."))
		} else if feedtype.Name == "" {
			utils.Error(c, fmt.Errorf("Incomplete provided data."))
		} else if result, err := h.FeedService.StoreFeedType(&feedtype); err != nil {
			utils.Error(c, err)
		} else {
			utils.Ok(c, &result)
		}
		return
	}
}

func (h *FeedHandler) RemoveFeedTypeByID(c *gin.Context) {
	id := c.Params.ByName("id")
	uid, err := uuid.FromString(id)

	if err != nil {
		utils.Error(c, err)
	} else if _, err := h.FeedService.RemoveFeedTypeByID(uid); err != nil {
		utils.Error(c, err)
	} else {
		utils.NoContent(c)
	}
	return
}

func (h *FeedHandler) RemoveFeedTypeByIDs(c *gin.Context) {
	//process json like : {"ids":["0b86bef7-0e16-47e6-9463-6a0b583e8d4c","6be6e63c-18f3-48ce-831f-f3210a576945"]}
	var ids []uuid.UUID
	reqBody := new(UUIDRequestModel)
	err := c.Bind(reqBody)
	if err != nil {
		utils.Error(c, err)
	} else if len(reqBody.Data) < 1 {
		utils.Error(c, fmt.Errorf("No Feed Types to be removed."))
	} else {
		for _, v := range reqBody.Data {
			//convert to UUID
			id, err := uuid.FromString(v)
			if err != nil {
				utils.Error(c, err)
				return
			} else {
				ids = append(ids, id)
			}
		}
		//process to services
		_, err := h.FeedService.RemoveFeedTypeByIDs(ids)
		if err != nil {
			utils.Error(c, err)
		} else {
			utils.NoContent(c)
		}
		return
	}
	return
}

//feed incoming
func (h *FeedHandler) ResolveFeedIncomingPage(c *gin.Context) {
	//capture something like this: http://localhost:9090/feed/incoming?page=1&limit=10
	q := c.Request.URL.Query()
	p := q.Get("page")
	l := q.Get("limit")
	page, err := strconv.Atoi(p)
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 10
	}

	if feeds, p, l, total, err := h.FeedService.ResolveFeedIncomingPage(int32(page), int32(limit)); err != nil {
		utils.Error(c, err)
	} else {
		utils.Page(c, feeds, p, l, total)
	}
	return
}

func (h *FeedHandler) ResolveFeedIncomingByID(c *gin.Context) {
	id := c.Params.ByName("id")
	uid, err := uuid.FromString(id)

	if err != nil {
		utils.Error(c, err)
	} else if feed, err := h.FeedService.ResolveFeedIncomingByID(uid); err != nil {
		utils.Error(c, err)
	} else {
		utils.Ok(c, &feed)
	}
	return
}

func (h *FeedHandler) StoreFeedIncoming(c *gin.Context) {

	var f feed.FeedIncoming
	c.BindJSON(&f)
	if f.Qty == 0 {
		utils.Error(c, fmt.Errorf("Qty must smaller or bigger than 0"))
	} else if result, err := h.FeedService.StoreFeedIncoming(&f); err != nil {
		utils.Error(c, err)
	} else {
		utils.Created(c, &result)
	}
	return
}

//feed adjustment
func (h *FeedHandler) ResolveFeedAdjustmentPage(c *gin.Context) {
	//capture something like this: http://localhost:9090/feed/feeding?page=1&limit=10
	q := c.Request.URL.Query()
	p := q.Get("page")
	l := q.Get("limit")
	page, err := strconv.Atoi(p)
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 10
	}

	if feedAdjustments, p, l, total, err := h.FeedService.ResolveFeedAdjustmentPage(int32(page), int32(limit)); err != nil {
		utils.Error(c, err)
	} else {
		utils.Page(c, feedAdjustments, p, l, total)
	}
	return
}

func (h *FeedHandler) ResolveFeedAdjustmentByID(c *gin.Context) {
	id := c.Params.ByName("id")
	uid, err := uuid.FromString(id)

	if err != nil {
		utils.Error(c, err)
	} else if feedAdjustment, err := h.FeedService.ResolveFeedAdjustmentByID(uid); err != nil {
		utils.Error(c, err)
	} else {
		utils.Ok(c, &feedAdjustment)
	}
	return
}

func (h *FeedHandler) StoreFeedAdjustment(c *gin.Context) {

	var f feed.FeedAdjustment
	c.BindJSON(&f)

	if f.Qty == 0 {
		utils.Error(c, fmt.Errorf("Qty must smaller or bigger than 0"))
	} else if result, err := h.FeedService.StoreFeedAdjustment(&f); err != nil {
		utils.Error(c, err)
	} else {
		utils.Created(c, &result)
	}
	return
}
