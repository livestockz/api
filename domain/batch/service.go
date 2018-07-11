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
	BatchRepository Repository `inject:"batchRepository"`
}

func (svc *BatchService) ResolveBatchByID(id uuid.UUID) (*Batch, error) {
	//fmt.Print(svc)
	if batch, err := svc.BatchRepository.ResolveBatchByID(id); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		return batch, nil
	}
}

//func (svc *BatchService) ClosePeriod(p *Period) (*Period, error) {
//
//}
