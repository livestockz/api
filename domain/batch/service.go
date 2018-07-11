package batch

import (
	uuid "github.com/satori/go.uuid"
)

type Service interface {
	ResolveBatchByID(*uuid.UUID) (*Batch, error)
	//ClosePeriod(*Period) (*Period, error)
	//CreatePeriod(Period) (Period, error)
}

type BatchService struct {
	BatchRepository Repository `inject:"BatchRepository"`
}

func (svc *BatchService) ResolveBatchByID(ID *uuid.UUID) (**Batch, error) {
	batch, error := (*BatchRepository).RessolveBatchByID(*ID)
	return &batch, error
}

//func (svc *BatchService) ClosePeriod(p *Period) (*Period, error) {
//
//}
