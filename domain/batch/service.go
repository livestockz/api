package batch

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type Service interface {
	ResolveBatchByID(uuid.UUID) (*Batch, error)
	//ClosePeriod(*Period) (*Period, error)
	//CreatePeriod(Period) (Period, error)
}

type BatchService struct {
	BatchRepository Repository `inject:"BatchRepository"`
}

func (svc *BatchService) ResolveBatchByID(ID uuid.UUID) (**Batch, error) {
	if batch, err := svc.BatchRepository.ResolveBatchByID(ID); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		return &batch, nil
	}
}

//func (svc *BatchService) ClosePeriod(p *Period) (*Period, error) {
//
//}
