package batch

import (
	"fmt"
	"time"

	"github.com/guregu/null"
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
	//cut off
	StoreGrowthCutOff(*CutOff) (*CutOff, error)
	//sales
	ResolveGrowthSalesByID(salesId uuid.UUID) (*Sales, error)
	StoreGrowthSales(sales *Sales) (*Sales, error)
	StoreGrowthSalesDetail(sales *Sales) (*Sales, error)
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
		//prepare ids of feed type id
		var ids []uuid.UUID
		for _, batchCycle := range *batchCycles {
			for _, feeding := range batchCycle.Feeding {
				ids = append(ids, feeding.FeedTypeID)
			}
		}
		//find feed type ids
		feedTypes, err := svc.FeedService.ResolveFeedTypeByIDs(ids)
		if err != nil {
			return nil, 0, 0, 0, err
		}

		//map feed type to batchcycle feeding's feed type
		var newBatchCycles []BatchCycle
		for _, batchCycle := range *batchCycles {
			var newFeeding []Feeding
			for _, feeding := range batchCycle.Feeding {
				for _, feedType := range *feedTypes {
					if feeding.FeedTypeID.String() == feedType.ID.String() {
						feeding.FeedType = feedType
						newFeeding = append(newFeeding, feeding)
					}
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

		//populate feed type id
		var ids []uuid.UUID
		for _, feeding := range batchCycle.Feeding {
			ids = append(ids, feeding.FeedTypeID)
		}
		//find feed type ids
		feedTypes, err := svc.FeedService.ResolveFeedTypeByIDs(ids)
		if err != nil {
			return nil, err
		}
		//replace feed type in feeding on batchCycle
		var newFeeding []Feeding
		for _, feeding := range batchCycle.Feeding {
			for _, feedType := range *feedTypes {
				if feeding.FeedTypeID.String() == feedType.ID.String() {
					feeding.FeedType = feedType
					newFeeding = append(newFeeding, feeding)
				}
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

//growth cut off
func (svc *BatchService) StoreGrowthCutOff(cutoff *CutOff) (*CutOff, error) {
	//validate cutoff existed
	if cutoffs, error := svc.BatchRepository.ResolveGrowthSummaryByBatchCycleID(cutoff.BatchCycleID); error != nil {
		return nil, error
	} else if cutoffs != nil {
		return nil, fmt.Errorf("You cannot cutoff this cycle, cutoff existed.")
	}

	//get batch cycle and feeding data
	if batchCycle, error := svc.BatchRepository.ResolveGrowthBatchCycleByID(cutoff.BatchID, cutoff.BatchCycleID); error != nil {
		return nil, error
	} else if feedings, err := svc.BatchRepository.ResolveGrowthFeedingByBatchCycleID(cutoff.BatchCycleID); err != nil {
		return nil, error
	} else {
		//calculate ADG
		days := cutoff.SummaryDate.Sub(batchCycle.Start).Hours() / 24
		cutoff.ADG = (cutoff.Weight - batchCycle.Weight) / days

		//calculate FCR
		var total float64
		for _, feeding := range *feedings {
			total = total + feeding.Qty
		}
		cutoff.FCR = total / (cutoff.Weight - batchCycle.Weight)

		//calculate SR
		cutoff.SR = (cutoff.Amount / batchCycle.Amount) * 100

		//update cycle finish date on batch cycle then insert growth summary
		batchCycle.Finish = null.TimeFrom(cutoff.SummaryDate)
		cutoff.ID = uuid.Must(uuid.NewV4())
		summary, err := svc.BatchRepository.UpdateGrowthBatchCycleAndInsertGrowthSummaryTransaction(batchCycle, cutoff)
		if err != nil {
			return nil, error
		} else {
			return summary, nil
		}

	}
}

//growth sales
func (svc *BatchService) ResolveGrowthSalesByID(salesId uuid.UUID) (*Sales, error) {
	if result, err := svc.BatchRepository.ResolveGrowthSalesByID(salesId); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
func (svc *BatchService) StoreGrowthSales(sales *Sales) (*Sales, error) {
	if sales.ID == uuid.Nil {
		sales.ID = uuid.Must(uuid.NewV4())
		if sales, err := svc.BatchRepository.InsertGrowthSales(sales); err != nil {
			return nil, err
		} else {
			return sales, nil
		}
	} else {
		if sales, err := svc.BatchRepository.UpdateGrowthSalesByID(sales); err != nil {
			return nil, err
		} else {
			return sales, nil
		}
	}
}

func (svc *BatchService) StoreGrowthSalesDetail(sales *Sales) (*Sales, error) {
	//set sales id and create cutoff
	cutoffs := make([]CutOff, 0)
	batchCycles := make([]BatchCycle, 0)
	salesDetail := make([]SalesDetail, 0)
	for _, detail := range sales.Detail {
		detail.ID = uuid.Must(uuid.NewV4())
		salesDetail = append(salesDetail, detail)
		var cutoff CutOff
		if batchCycle, error := svc.BatchRepository.ResolveGrowthBatchCycleByID(detail.BatchID, detail.BatchCycleID); error != nil {
			return nil, error
		} else if feedings, err := svc.BatchRepository.ResolveGrowthFeedingByBatchCycleID(detail.BatchCycleID); err != nil {
			return nil, error
		} else {
			//calculate ADG
			days := cutoff.SummaryDate.Sub(batchCycle.Start).Hours() / 24
			cutoff.ADG = (cutoff.Weight - batchCycle.Weight) / days

			//calculate FCR
			var total float64
			for _, feeding := range *feedings {
				total = total + feeding.Qty
			}
			cutoff.FCR = total / (cutoff.Weight - batchCycle.Weight)

			//calculate SR
			cutoff.SR = (cutoff.Amount / batchCycle.Amount) * 100

			//set cycle finish date on batch cycle then insert growth summary
			batchCycle.Finish = null.TimeFrom(time.Now())
			batchCycles = append(batchCycles, *batchCycle)
			cutoff.ID = uuid.Must(uuid.NewV4())
			cutoff.BatchCycleID = detail.BatchCycleID
			cutoff.BatchID = detail.BatchID
			cutoff.Weight = detail.Weight
			cutoff.Amount = detail.Amount
			cutoff.SummaryDate = time.Now()
			cutoffs = append(cutoffs, cutoff)
		}
	}
	sales.Detail = salesDetail

	if result, err := svc.BatchRepository.UpdateGrowthBatchCycleInsertGrowthSummaryAndInsertSalesDetail(&batchCycles, &cutoffs, sales); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
