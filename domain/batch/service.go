package batch

import (
	"fmt"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type Service interface {
	ResolveGrowthBatches(page int32, limit int32) (*[]Batch, int32, int32, int32, error)
	ResolveGrowthBatchByID(uuid.UUID) (*Batch, error)
	StoreGrowthBatch(*gin.Context) (*Batch, error)
	UpdateGrowthBatchByID(*gin.Context) (*Batch, error)
	RemoveGrowthBatchByID(uuid.UUID) (*Batch, error)
	//ClosePeriod(*Period) (*Period, error)
	//CreatePeriod(Period) (Period, error)
}

type BatchService struct {
	BatchRepository Repository `inject:"batchRepository"`
}

func (svc *BatchService) ResolveGrowthBatches(page int32, limit int32) (*[]Batch, int32, int32, int32, error) {
	if batches, page, limit, total, err := svc.BatchRepository.ResolveGrowthBatches(page, limit); err != nil {
		return nil, 0, 0, 0, err
	} else {
		return batches, page, limit, total, nil
	}
}

func (svc *BatchService) ResolveGrowthBatchByID(id uuid.UUID) (*Batch, error) {
	if batch, err := svc.BatchRepository.ResolveGrowthBatchByID(id); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		return batch, nil
	}
}

func (svc *BatchService) StoreGrowthBatch(c *gin.Context) (*Batch, error) {
	var batch Batch
	c.BindJSON(&batch)
	var data = &batch

	if data.Name == "" {
		return nil, fmt.Errorf("Incomplete provided data.")
	} else if result, err := svc.BatchRepository.StoreGrowthBatch(&batch); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (svc *BatchService) UpdateGrowthBatchByID(c *gin.Context) (*Batch, error) {
	var batch Batch
	c.BindJSON(&batch)
	var data = &batch
	var id = c.Params.ByName("id")
	var uid, err = uuid.FromString(id)
	if err != nil {
		return nil, fmt.Errorf("Unable to convert given ID to UUID")
	} else if data.ID != uid {
		return nil, fmt.Errorf("Inconsistent ID.")
	} else if data.Name == "" {
		return nil, fmt.Errorf("Incomplete provided data.")
	} else if result, err := svc.BatchRepository.StoreGrowthBatch(&batch); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (svc *BatchService) RemoveGrowthBatchByID(id uuid.UUID) (*Batch, error) {
	if batch, err := svc.BatchRepository.RemoveGrowthBatchByID(id); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		return batch, nil
	}
}

//func (svc *BatchService) ClosePeriod(p *Period) (*Period, error) {
//
//}
