package batch

import (
	"fmt"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type Service interface {
	ResolveBatchByID(uuid.UUID) (*Batch, error)
	CreateBatch(*gin.Context) (*Batch, error)
	//ClosePeriod(*Period) (*Period, error)
	//CreatePeriod(Period) (Period, error)
}

type BatchService struct {
	BatchRepository Repository `inject:"batchRepository"`
}

func (svc *BatchService) ResolveBatchByID(id uuid.UUID) (*Batch, error) {
	if batch, err := svc.BatchRepository.ResolveBatchByID(id); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		return batch, nil
	}
}

func (svc *BatchService) CreateBatch(c *gin.Context) (*Batch, error) {
	var batch Batch
	c.BindJSON(&batch)
	var data = &batch
	if data.Name == "" {
		return nil, fmt.Errorf("Incomplete provided data.")
	} else {
		if result, err := svc.BatchRepository.StoreBatch(&batch); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
}

//func (svc *BatchService) ClosePeriod(p *Period) (*Period, error) {
//
//}
