package batch

import (
	"fmt"

	"github.com/livestockz/api/domain/feed"
	uuid "github.com/satori/go.uuid"
)

type Service interface {
	//batch
	ResolveGrowthBatchPage(page int32, limit int32, deleted string) (*[]Batch, int32, int32, int32, error)
	ResolveGrowthBatchByID(uuid.UUID) (*Batch, error)
	StoreGrowthBatch(*Batch) (*Batch, error)
	RemoveGrowthBatchByID(uuid.UUID) (*Batch, error)
	RemoveGrowthBatchByIDs([]uuid.UUID) (*[]Batch, error)
	//pool
	ResolveGrowthPoolPage(page int32, limit int32, deleted string) (*[]Pool, int32, int32, int32, error)
	ResolveGrowthPoolByID(uuid.UUID) (*Pool, error)
	StoreGrowthPool(*Pool) (*Pool, error)
	RemoveGrowthPoolByID(uuid.UUID) (*Pool, error)
	RemoveGrowthPoolByIDs([]uuid.UUID) (*[]Pool, error)
	//batch cycle
	ResolveGrowthBatchCyclePage(batchId uuid.UUID, page int32, limit int32) (*[]BatchCycle, int32, int32, int32, error)
	ResolveGrowthBatchCycleByID(batchId uuid.UUID, cycleId uuid.UUID) (*BatchCycle, error)
	StoreGrowthBatchCycle(*BatchCycle) (*BatchCycle, error)
	//death
	StoreGrowthDeath(*Death) (*Death, error)
	//death
	StoreGrowthFeeding(*Feeding) (*Feeding, error)
}

type BatchService struct {
	BatchRepository Repository   `inject:"batchRepository"`
	FeedService     feed.Service `inject:"feedService"`
}

func (svc *BatchService) ResolveGrowthBatchPage(page int32, limit int32, deleted string) (*[]Batch, int32, int32, int32, error) {
	if batches, page, limit, total, err := svc.BatchRepository.ResolveGrowthBatchPage(page, limit, deleted); err != nil {
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

//pools
func (svc *BatchService) ResolveGrowthPoolPage(page int32, limit int32, deleted string) (*[]Pool, int32, int32, int32, error) {
	if pools, page, limit, total, err := svc.BatchRepository.ResolveGrowthPoolPage(page, limit, deleted); err != nil {
		return nil, 0, 0, 0, err
	} else {
		return pools, page, limit, total, nil
	}
}

func (svc *BatchService) ResolveGrowthPoolByID(id uuid.UUID) (*Pool, error) {
	if pool, err := svc.BatchRepository.ResolveGrowthPoolByID(id); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		return pool, nil
	}
}

func (svc *BatchService) StoreGrowthPool(pool *Pool) (*Pool, error) {
	if pool.ID == uuid.Nil {
		pool.ID = uuid.Must(uuid.NewV4())
		if result, err := svc.BatchRepository.InsertGrowthPool(pool); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	} else {
		//update
		if result, err := svc.BatchRepository.UpdateGrowthPoolByID(pool); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
}

func (svc *BatchService) RemoveGrowthPoolByID(id uuid.UUID) (*Pool, error) {
	if _, err := svc.BatchRepository.RemoveGrowthPoolByID(id); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		return nil, nil
	}
}

func (svc *BatchService) RemoveGrowthPoolByIDs(ids []uuid.UUID) (*[]Pool, error) {
	if _, err := svc.BatchRepository.RemoveGrowthPoolByIDs(ids); err != nil {
		return nil, err
	} else {
		return nil, nil
	}
}

//batch cycle
func (svc *BatchService) ResolveGrowthBatchCyclePage(batchId uuid.UUID, page int32, limit int32) (*[]BatchCycle, int32, int32, int32, error) {
	if batchCycles, page, limit, total, err := svc.BatchRepository.ResolveGrowthBatchCyclePage(batchId, page, limit); err != nil {
		return nil, 0, 0, 0, err
	} else {
		var newBatchCycles []BatchCycle
		for _, batchCycle := range *batchCycles {
			var newFeeding []Feeding
			for _, feeding := range batchCycle.Feeding {
				if feedtype, err := svc.FeedService.ResolveFeedTypeByID(feeding.FeedTypeID); err != nil {
					return nil, 0, 0, 0, err
				} else {
					feeding.FeedType = *feedtype
					newFeeding = append(newFeeding, feeding)
				}
			}
			batchCycle.Feeding = newFeeding
			newBatchCycles = append(newBatchCycles, batchCycle)
		}
		return &newBatchCycles, page, limit, total, nil
	}
}

func (svc *BatchService) ResolveGrowthBatchCycleByID(batchId uuid.UUID, cycleId uuid.UUID) (*BatchCycle, error) {
	if batchCycle, err := svc.BatchRepository.ResolveGrowthBatchCycleByID(batchId, cycleId); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		var newFeeding []Feeding
		for _, feeding := range batchCycle.Feeding {
			if feedtype, err := svc.FeedService.ResolveFeedTypeByID(feeding.FeedTypeID); err != nil {
				return nil, err
			} else {
				feeding.FeedType = *feedtype
				newFeeding = append(newFeeding, feeding)
			}
		}
		batchCycle.Feeding = newFeeding
		return batchCycle, nil
	}
}

func (svc *BatchService) StoreGrowthBatchCycle(batchCycle *BatchCycle) (*BatchCycle, error) {
	if batchCycle.ID == uuid.Nil {
		batchCycle.ID = uuid.Must(uuid.NewV4())
		if result, err := svc.BatchRepository.InsertGrowthBatchCycle(batchCycle); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	} else {
		//update
		if result, err := svc.BatchRepository.UpdateGrowthBatchCycleByID(batchCycle); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
}

//growth death
func (svc *BatchService) StoreGrowthDeath(death *Death) (*Death, error) {
	death.ID = uuid.Must(uuid.NewV4())
	if result, err := svc.BatchRepository.InsertGrowthDeath(death); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

//growth feeding
func (svc *BatchService) StoreGrowthFeeding(feeding *Feeding) (*Feeding, error) {
	feeding.ID = uuid.Must(uuid.NewV4())
	if result, err := svc.BatchRepository.InsertGrowthFeeding(feeding); err != nil {
		return nil, err
	} else if feedtype, err := svc.FeedService.ResolveFeedTypeByID(result.FeedTypeID); err != nil {
		return nil, err
	} else {
		result.FeedType = *feedtype
		return result, nil
	}
}
