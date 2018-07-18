package batch

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type Service interface {
	ResolveGrowthBatchPage(page int32, limit int32) (*[]Batch, int32, int32, int32, error)
	ResolveGrowthBatchByID(uuid.UUID) (*Batch, error)
	StoreGrowthBatch(*Batch) (*Batch, error)
	RemoveGrowthBatchByID(uuid.UUID) (*Batch, error)
	RemoveGrowthBatchByIDs([]uuid.UUID) (*[]Batch, error)
	//ClosePeriod(*Period) (*Period, error)
	//CreatePeriod(Period) (Period, error)
}

type BatchService struct {
	BatchRepository Repository `inject:"batchRepository"`
}

func (svc *BatchService) ResolveGrowthBatchPage(page int32, limit int32) (*[]Batch, int32, int32, int32, error) {
	if batches, page, limit, total, err := svc.BatchRepository.ResolveGrowthBatchPage(page, limit); err != nil {
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

func (svc *BatchService) StoreGrowthBatch(batch *Batch) (*Batch, error) {
	if batch.ID == uuid.Nil {
		batch.ID = uuid.Must(uuid.NewV4())
		if result, err := svc.BatchRepository.InsertGrowthBatch(batch); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	} else {
		//update
		if result, err := svc.BatchRepository.UpdateGrowthBatchByID(batch); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
}

func (svc *BatchService) RemoveGrowthBatchByID(id uuid.UUID) (*Batch, error) {
	if _, err := svc.BatchRepository.RemoveGrowthBatchByID(id); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		return nil, nil
	}
}

func (svc *BatchService) RemoveGrowthBatchByIDs(ids []uuid.UUID) (*[]Batch, error) {
	if _, err := svc.BatchRepository.RemoveGrowthBatchByIDs(ids); err != nil {
		return nil, err
	} else {
		return nil, nil
	}
}

//func (svc *BatchService) ClosePeriod(p *Period) (*Period, error) {
//
//}
