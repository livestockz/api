package batch

import (
	"database/sql"
	"fmt"

	"github.com/ncrypthic/dbmapper"
	. "github.com/ncrypthic/dbmapper/dialects/mysql"
	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	//batch
	ResolveGrowthBatchPage(page int32, limit int32, deleted string) (*[]Batch, int32, int32, int32, error)
	ResolveGrowthBatchByID(id uuid.UUID) (*Batch, error)
	InsertGrowthBatch(batch *Batch) (*Batch, error)
	UpdateGrowthBatchByID(batch *Batch) (*Batch, error)
	RemoveGrowthBatchByID(id uuid.UUID) (*Batch, error)
	RemoveGrowthBatchByIDs(ids []uuid.UUID) (*[]Batch, error)
	//pool
	ResolveGrowthPoolPage(page int32, limit int32, deleted string) (*[]Pool, int32, int32, int32, error)
	ResolveGrowthPoolByID(id uuid.UUID) (*Pool, error)
	InsertGrowthPool(batch *Pool) (*Pool, error)
	UpdateGrowthPoolByID(pool *Pool) (*Pool, error)
	RemoveGrowthPoolByID(id uuid.UUID) (*Pool, error)
	RemoveGrowthPoolByIDs(ids []uuid.UUID) (*[]Pool, error)
	//batch cycle
	ResolveGrowthBatchCyclePage(batchId uuid.UUID, page int32, limit int32) (*[]BatchCycle, int32, int32, int32, error)
	ResolveGrowthBatchCycleByID(batchId uuid.UUID, cycleId uuid.UUID) (*BatchCycle, error)
	InsertGrowthBatchCycle(batchCycle *BatchCycle) (*BatchCycle, error)
	UpdateGrowthBatchCycleByID(batchCycle *BatchCycle) (*BatchCycle, error)
}

const (
	//batch
	selectGrowthBatch = `SELECT id, name, status, deleted, created, updated FROM growth_batch`
	insertGrowthBatch = `INSERT INTO growth_batch(id, name, status, deleted, created) VALUES (:id ,:name, :status, :deleted, NOW())`
	updateGrowthBatch = `UPDATE growth_batch SET name = :name, status = :status, deleted = :deleted, updated = NOW() WHERE id = :id`
	deleteGrowthBatch = `UPDATE growth_batch SET deleted = 1, updated = NOW() WHERE id = :id`
	//pool
	selectGrowthPool = `SELECT id, name, status, deleted, created, updated FROM growth_pool`
	insertGrowthPool = `INSERT INTO growth_pool(id, name, status, deleted, created) VALUES (:id ,:name, :status, :deleted, NOW())`
	updateGrowthPool = `UPDATE growth_pool SET name = :name, status = :status, deleted = :deleted, updated = NOW() WHERE id = :id`
	deleteGrowthPool = `UPDATE growth_pool SET deleted = 1, updated = NOW() WHERE id = :id`
	//batch cycle
	selectGrowthBatchCycle = `SELECT id, growth_batch_id, growth_pool_id, cycle_start, cycle_finish, weight, amount, created, updated FROM growth_batch_cycle`
	insertGrowthBatchCycle = `INSERT INTO growth_batch_cycle(id, growth_batch_id, growth_pool_id, cycle_start, weight, amount, created) VALUES (:id ,:batch, :pool, :start, :weight, :amount, NOW())`
	updateGrowthBatchCycle = `UPDATE growth_batch_cycle SET growth_batch_id = :batch, growth_pool_id = :pool, cycle_start = :start, cycle_finish = :finish, weight = :weight, amount = :amount, updated = NOW() WHERE id = :id`
)

type BatchRepository struct {
	DB *sql.DB `inject:"db"`
}

//batch
func (repo *BatchRepository) ResolveGrowthBatchPage(page int32, limit int32, deleted string) (*[]Batch, int32, int32, int32, error) {
	var start int32
	var end int32

	start = page * limit
	end = limit

	//get data by given page
	var query dbmapper.QueryMapper
	if deleted == Deleted_Any {
		query = dbmapper.Prepare(selectGrowthBatch+" ORDER BY name ASC LIMIT :start, :end").With(
			dbmapper.Param("start", start),
			dbmapper.Param("end", end),
		)
	} else if deleted == Deleted_True {
		query = dbmapper.Prepare(selectGrowthBatch+" WHERE deleted = 1 ORDER BY name ASC LIMIT :start, :end").With(
			dbmapper.Param("start", start),
			dbmapper.Param("end", end),
		)
	} else if deleted == Deleted_False {
		query = dbmapper.Prepare(selectGrowthBatch+" WHERE deleted = 0 ORDER BY name ASC LIMIT :start, :end").With(
			dbmapper.Param("start", start),
			dbmapper.Param("end", end),
		)
	}
	if err := query.Error(); err != nil {
		//log.Print(err.Error())
		return nil, page, limit, 0, err
	}

	batches := make([]Batch, 0)
	err := Parse(repo.DB.Query(query.SQL(), query.Params()...)).Map(batchesMapper(&batches))

	if err != nil {
		//log.Print(err.Error())
		return nil, page, limit, 0, err
	}

	//get total batch
	var summary dbmapper.QueryMapper
	if deleted == Deleted_Any {
		summary = dbmapper.Prepare("SELECT COUNT(*) AS total FROM growth_batch")
	} else if deleted == Deleted_True {
		summary = dbmapper.Prepare("SELECT COUNT(*) AS total FROM growth_batch WHERE deleted = 1")
	} else if deleted == Deleted_False {
		summary = dbmapper.Prepare("SELECT COUNT(*) AS total FROM growth_batch WHERE deleted = 0")
	}

	if err := summary.Error(); err != nil {
		//log.Print(err.Error())
		return nil, page, limit, 0, err
	}

	var batchesCount int32
	total := make([]int32, 0)
	err = Parse(repo.DB.Query(summary.SQL())).Map(dbmapper.Int32("total", &total))
	if err != nil {
		//log.Print(err.Error())
		return nil, page, limit, 0, err
	} else {
		batchesCount = total[0]
		//log.Print(batchesCount)
	}
	//fmt.Println(&total)
	return &batches, page, limit, batchesCount, nil

}

func (repo *BatchRepository) ResolveGrowthBatchByID(id uuid.UUID) (*Batch, error) {
	query := dbmapper.Prepare(selectGrowthBatch + " WHERE id = :id").With(
		dbmapper.Param("id", id),
	)
	if err := query.Error(); err != nil {
		return nil, err
	}
	batches := make([]Batch, 0)
	err := Parse(repo.DB.Query(query.SQL(), query.Params()...)).Map(batchesMapper(&batches))

	if err != nil {
		return nil, err
	}
	if len(batches) < 1 {
		return nil, fmt.Errorf("growth batch with id %s not found", id)
	}
	return &batches[0], nil
	// pool u/
	//&batches[0].Pool = pools
}

func (repo *BatchRepository) InsertGrowthBatch(batch *Batch) (*Batch, error) {

	//insert
	//fmt.Println("insert")
	//fmt.Print("\n")

	//prepare query and params
	insert := dbmapper.Prepare(insertGrowthBatch).With(
		dbmapper.Param("id", batch.ID),
		dbmapper.Param("name", batch.Name),
		dbmapper.Param("status", batch.Status),
		dbmapper.Param("deleted", batch.Deleted),
	)
	//log.Print("sql:", insert.SQL())
	//fmt.Print("\n")
	//log.Print("sql params:", insert.Params())
	//fmt.Print("\n")
	//validate query
	if err := insert.Error(); err != nil {
		//log.Print(err.Error())
		//fmt.Print("\n")
		return nil, err
	} else {
		//insert to database
		if _, err := repo.DB.Exec(insert.SQL(), insert.Params()...); err != nil {
			//log.Print(err.Error())
			//fmt.Print("\n")
			return nil, err
		} else {
			//find inserted data from database based on generated id
			res, err := repo.ResolveGrowthBatchByID(batch.ID)
			return res, err
		}
	}
}

func (repo *BatchRepository) UpdateGrowthBatchByID(batch *Batch) (*Batch, error) {
	//find whether if data exist
	//fmt.Print("\n")
	//fmt.Print(batch)
	//fmt.Print("\n")
	_, err := repo.ResolveGrowthBatchByID(batch.ID)

	if err != nil {
		//fmt.Print(err)
		//fmt.Print("\n")
		return nil, err
	} else {
		//update
		//fmt.Println("update")
		//fmt.Print("\n")
		//prepare query and params
		updater := dbmapper.Prepare(updateGrowthBatch).With(
			dbmapper.Param("name", batch.Name),
			dbmapper.Param("status", batch.Status),
			dbmapper.Param("deleted", batch.Deleted),
			dbmapper.Param("id", batch.ID),
		)
		//fmt.Print("\n")
		//log.Print("sql:", updater.SQL())
		//log.Print("sql params:", updater.Params())
		//fmt.Print("\n")
		//validate query
		if err := updater.Error(); err != nil {
			//log.Print(err.Error())
			//fmt.Print("\n")
			return nil, err
		} else {
			//update to database
			if _, err := repo.DB.Exec(updater.SQL(), updater.Params()...); err != nil {
				//log.Print(err.Error())
				//fmt.Print("\n")
				return nil, err
			} else {
				//find inserted data from database based on generated id
				res, err := repo.ResolveGrowthBatchByID(batch.ID)
				return res, err
			}
		}
	}
}

func (repo *BatchRepository) RemoveGrowthBatchByID(id uuid.UUID) (*Batch, error) {
	//find whether if data exist
	//fmt.Print("\n")
	//fmt.Print(id)
	//fmt.Print("\n")
	if _, err := repo.ResolveGrowthBatchByID(id); err != nil {
		return nil, err
	} else {
		remover := dbmapper.Prepare(deleteGrowthBatch).With(
			dbmapper.Param("id", id),
		)
		//fmt.Print("\n")
		//log.Print("sql:", remover.SQL())
		//log.Print("sql params:", remover.Params())
		//fmt.Print("\n")
		//validate query
		if err := remover.Error(); err != nil {
			//log.Print(err.Error())
			//fmt.Print("\n")
			return nil, err
		} else {
			//update to database
			if _, err := repo.DB.Exec(remover.SQL(), remover.Params()...); err != nil {
				//log.Print(err.Error())
				//fmt.Print("\n")
				return nil, err
			} else {
				return nil, nil
			}
		}
	}
}

func (repo *BatchRepository) RemoveGrowthBatchByIDs(ids []uuid.UUID) (*[]Batch, error) {
	for _, v := range ids {
		if _, err := repo.RemoveGrowthBatchByID(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func batchMapper(row *Batch) *dbmapper.MappedColumns {
	return dbmapper.Columns(
		dbmapper.Column("id").As(&row.ID),
		dbmapper.Column("name").As(&row.Name),
		dbmapper.Column("status").As(&row.Status),
		dbmapper.Column("deleted").As(&row.Deleted),
		dbmapper.Column("created").As(&row.Created),
		dbmapper.Column("updated").As(&row.Updated),
	)
}

func batchesMapper(rows *[]Batch) dbmapper.RowMapper {
	return func() *dbmapper.MappedColumns {
		row := Batch{}
		return batchMapper(&row).Then(func() error {
			*rows = append(*rows, row)
			return nil
		})
	}
}

//pool
func (repo *BatchRepository) ResolveGrowthPoolPage(page int32, limit int32, deleted string) (*[]Pool, int32, int32, int32, error) {
	var start int32
	var end int32

	start = page * limit
	end = limit
	//get data by given page
	var query dbmapper.QueryMapper
	if deleted == Deleted_Any {
		query = dbmapper.Prepare(selectGrowthPool+" ORDER BY name ASC LIMIT :start, :end").With(
			dbmapper.Param("start", start),
			dbmapper.Param("end", end),
		)
	} else if deleted == Deleted_True {
		query = dbmapper.Prepare(selectGrowthPool+" WHERE deleted = 1 ORDER BY name ASC LIMIT :start, :end").With(
			dbmapper.Param("start", start),
			dbmapper.Param("end", end),
		)
	} else if deleted == Deleted_False {
		query = dbmapper.Prepare(selectGrowthPool+" WHERE deleted = 0 ORDER BY name ASC LIMIT :start, :end").With(
			dbmapper.Param("start", start),
			dbmapper.Param("end", end),
		)
	}

	if err := query.Error(); err != nil {
		//log.Print(err.Error())
		return nil, page, limit, 0, err
	}

	pools := make([]Pool, 0)
	err := Parse(repo.DB.Query(query.SQL(), query.Params()...)).Map(poolsMapper(&pools))

	if err != nil {
		//log.Print(err.Error())
		return nil, page, limit, 0, err
	}

	//get total batch
	var summary dbmapper.QueryMapper
	if deleted == Deleted_Any {
		summary = dbmapper.Prepare("SELECT COUNT(*) AS total FROM growth_pool")
	} else if deleted == Deleted_True {
		summary = dbmapper.Prepare("SELECT COUNT(*) AS total FROM growth_pool WHERE deleted = 1")
	} else if deleted == Deleted_False {
		summary = dbmapper.Prepare("SELECT COUNT(*) AS total FROM growth_pool WHERE deleted = 0")
	}
	if err := summary.Error(); err != nil {
		//log.Print(err.Error())
		return nil, page, limit, 0, err
	}

	var poolsCount int32
	total := make([]int32, 0)
	err = Parse(repo.DB.Query(summary.SQL())).Map(dbmapper.Int32("total", &total))
	if err != nil {
		return nil, page, limit, 0, err
	} else {
		poolsCount = total[0]
	}
	return &pools, page, limit, poolsCount, nil

}

func (repo *BatchRepository) ResolveGrowthPoolByID(id uuid.UUID) (*Pool, error) {
	query := dbmapper.Prepare(selectGrowthPool + " WHERE id = :id").With(
		dbmapper.Param("id", id),
	)
	if err := query.Error(); err != nil {
		return nil, err
	}
	pools := make([]Pool, 0)
	err := Parse(repo.DB.Query(query.SQL(), query.Params()...)).Map(poolsMapper(&pools))

	if err != nil {
		return nil, err
	}
	if len(pools) < 1 {
		return nil, fmt.Errorf("growth pool with id %s not found", id)
	}
	return &pools[0], nil
}

func (repo *BatchRepository) InsertGrowthPool(pool *Pool) (*Pool, error) {

	//prepare query and params
	insert := dbmapper.Prepare(insertGrowthPool).With(
		dbmapper.Param("id", pool.ID),
		dbmapper.Param("name", pool.Name),
		dbmapper.Param("status", pool.Status),
		dbmapper.Param("deleted", pool.Deleted),
	)
	//validate query
	if err := insert.Error(); err != nil {
		return nil, err
	} else {
		//insert to database
		if _, err := repo.DB.Exec(insert.SQL(), insert.Params()...); err != nil {
			return nil, err
		} else {
			//find inserted data from database based on generated id
			res, err := repo.ResolveGrowthPoolByID(pool.ID)
			return res, err
		}
	}
}

func (repo *BatchRepository) UpdateGrowthPoolByID(pool *Pool) (*Pool, error) {
	//find whether if data exist
	_, err := repo.ResolveGrowthPoolByID(pool.ID)

	if err != nil {
		return nil, err
	} else {
		//update
		//prepare query and params
		updater := dbmapper.Prepare(updateGrowthPool).With(
			dbmapper.Param("name", pool.Name),
			dbmapper.Param("status", pool.Status),
			dbmapper.Param("deleted", pool.Deleted),
			dbmapper.Param("id", pool.ID),
		)
		//validate query
		if err := updater.Error(); err != nil {
			return nil, err
		} else {
			//update to database
			if _, err := repo.DB.Exec(updater.SQL(), updater.Params()...); err != nil {
				return nil, err
			} else {
				//find inserted data from database based on generated id
				res, err := repo.ResolveGrowthPoolByID(pool.ID)
				return res, err
			}
		}
	}
}

func (repo *BatchRepository) RemoveGrowthPoolByID(id uuid.UUID) (*Pool, error) {
	//find whether if data exist
	if _, err := repo.ResolveGrowthPoolByID(id); err != nil {
		return nil, err
	} else {
		remover := dbmapper.Prepare(deleteGrowthPool).With(
			dbmapper.Param("id", id),
		)
		//validate query
		if err := remover.Error(); err != nil {
			return nil, err
		} else {
			//update to database
			if _, err := repo.DB.Exec(remover.SQL(), remover.Params()...); err != nil {
				return nil, err
			} else {
				return nil, nil
			}
		}
	}
}

func (repo *BatchRepository) RemoveGrowthPoolByIDs(ids []uuid.UUID) (*[]Pool, error) {
	for _, v := range ids {
		if _, err := repo.RemoveGrowthPoolByID(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func poolMapper(row *Pool) *dbmapper.MappedColumns {
	return dbmapper.Columns(
		dbmapper.Column("id").As(&row.ID),
		dbmapper.Column("name").As(&row.Name),
		dbmapper.Column("status").As(&row.Status),
		dbmapper.Column("deleted").As(&row.Deleted),
		dbmapper.Column("created").As(&row.Created),
		dbmapper.Column("updated").As(&row.Updated),
	)
}

func poolsMapper(rows *[]Pool) dbmapper.RowMapper {
	return func() *dbmapper.MappedColumns {
		row := Pool{}
		return poolMapper(&row).Then(func() error {
			*rows = append(*rows, row)
			return nil
		})
	}
}

//batch cycle
func (repo *BatchRepository) ResolveGrowthBatchCyclePage(batchId uuid.UUID, page int32, limit int32) (*[]BatchCycle, int32, int32, int32, error) {
	var start int32
	var end int32

	start = page * limit
	end = limit

	//get data by given page
	var query dbmapper.QueryMapper
	query = dbmapper.Prepare(selectGrowthBatchCycle+" WHERE growth_batch_id=:batchId ORDER BY created ASC LIMIT :start, :end").With(
		dbmapper.Param("batchId", batchId),
		dbmapper.Param("start", start),
		dbmapper.Param("end", end),
	)

	if err := query.Error(); err != nil {
		//log.Print(err.Error())
		return nil, page, limit, 0, err
	}

	batchCycles := make([]BatchCycle, 0)
	err := Parse(repo.DB.Query(query.SQL(), query.Params()...)).Map(batchCyclesMapper(&batchCycles))

	if err != nil {
		return nil, page, limit, 0, err
	}

	var newBatchCycles []BatchCycle
	for _, batchCycle := range batchCycles {
		if batch, err := repo.ResolveGrowthBatchByID(batchCycle.BatchID); err != nil {
			return nil, page, limit, 0, err
		} else {
			batchCycle.Batch = *batch
		}
		pool, err := repo.ResolveGrowthPoolByID(batchCycle.PoolID)
		if err != nil {
			return nil, page, limit, 0, err
		} else {
			batchCycle.Pool = *pool
		}
		newBatchCycles = append(newBatchCycles, batchCycle)
	}

	//get total batch cycle
	var summary dbmapper.QueryMapper
	summary = dbmapper.Prepare("SELECT COUNT(*) AS total FROM growth_batch_cycle WHERE growth_batch_id=:batchId").With(
		dbmapper.Param("batchId", batchId),
	)

	if err := summary.Error(); err != nil {
		//log.Print(err.Error())
		return nil, page, limit, 0, err
	}

	var batchCyclesCount int32
	total := make([]int32, 0)
	err = Parse(repo.DB.Query(summary.SQL(), summary.Params()...)).Map(dbmapper.Int32("total", &total))
	if err != nil {
		return nil, page, limit, 0, err
	} else {
		batchCyclesCount = total[0]
	}
	return &newBatchCycles, page, limit, batchCyclesCount, nil

}

func (repo *BatchRepository) ResolveGrowthBatchCycleByID(batchId uuid.UUID, cycleId uuid.UUID) (*BatchCycle, error) {
	query := dbmapper.Prepare(selectGrowthBatchCycle+" WHERE id = :cycleId and growth_batch_id = :batchId").With(
		dbmapper.Param("cycleId", cycleId),
		dbmapper.Param("batchId", batchId),
	)
	if err := query.Error(); err != nil {
		return nil, err
	}
	batchCycles := make([]BatchCycle, 0)
	err := Parse(repo.DB.Query(query.SQL(), query.Params()...)).Map(batchCyclesMapper(&batchCycles))

	if err != nil {
		return nil, err
	}

	if len(batchCycles) < 1 {
		return nil, fmt.Errorf("growth batch cycle with batchId %s and cycleId %s not found", batchId, cycleId)
	} else {
		if batch, err := repo.ResolveGrowthBatchByID(batchCycles[0].BatchID); err != nil {
			return nil, err
		} else {
			batchCycles[0].Batch = *batch
		}

		if pool, err := repo.ResolveGrowthPoolByID(batchCycles[0].PoolID); err != nil {
			return nil, err
		} else {
			batchCycles[0].Pool = *pool
		}

		return &batchCycles[0], nil
	}
}

func (repo *BatchRepository) InsertGrowthBatchCycle(batchCycle *BatchCycle) (*BatchCycle, error) {

	//prepare query and params
	insert := dbmapper.Prepare(insertGrowthBatchCycle).With(
		dbmapper.Param("id", batchCycle.ID),
		dbmapper.Param("batch", batchCycle.Batch.ID),
		dbmapper.Param("pool", batchCycle.Pool.ID),
		dbmapper.Param("start", batchCycle.Start),
		dbmapper.Param("weight", batchCycle.Weight),
		dbmapper.Param("amount", batchCycle.Amount),
	)
	//validate query
	if err := insert.Error(); err != nil {
		return nil, err
	} else {
		//insert to database
		if _, err := repo.DB.Exec(insert.SQL(), insert.Params()...); err != nil {
			return nil, err
		} else {
			//find inserted data from database based on generated id
			res, err := repo.ResolveGrowthBatchCycleByID(batchCycle.Batch.ID, batchCycle.ID)
			return res, err
		}
	}
}

func (repo *BatchRepository) UpdateGrowthBatchCycleByID(batchCycle *BatchCycle) (*BatchCycle, error) {
	//find whether if data exist
	_, err := repo.ResolveGrowthBatchCycleByID(batchCycle.Batch.ID, batchCycle.ID)

	if err != nil {
		return nil, err
	} else {
		//update
		//prepare query and params
		updater := dbmapper.Prepare(updateGrowthBatchCycle).With(
			dbmapper.Param("batch", batchCycle.Batch.ID),
			dbmapper.Param("pool", batchCycle.Pool.ID),
			dbmapper.Param("start", batchCycle.Start),
			dbmapper.Param("finish", batchCycle.Finish),
			dbmapper.Param("weight", batchCycle.Weight),
			dbmapper.Param("amount", batchCycle.Amount),
			dbmapper.Param("id", batchCycle.ID),
		)
		//validate query
		if err := updater.Error(); err != nil {
			return nil, err
		} else {
			//update to database
			if _, err := repo.DB.Exec(updater.SQL(), updater.Params()...); err != nil {
				return nil, err
			} else {
				//find inserted data from database based on generated id
				res, err := repo.ResolveGrowthBatchCycleByID(batchCycle.Batch.ID, batchCycle.ID)
				return res, err
			}
		}
	}
}

func batchCycleMapper(row *BatchCycle) *dbmapper.MappedColumns {
	return dbmapper.Columns(
		dbmapper.Column("id").As(&row.ID),
		dbmapper.Column("growth_batch_id").As(&row.BatchID),
		dbmapper.Column("growth_pool_id").As(&row.PoolID),
		dbmapper.Column("cycle_start").As(&row.Start),
		dbmapper.Column("cycle_finish").As(&row.Finish),
		dbmapper.Column("weight").As(&row.Weight),
		dbmapper.Column("amount").As(&row.Amount),
		dbmapper.Column("created").As(&row.Created),
		dbmapper.Column("updated").As(&row.Updated),
	)
}

func batchCyclesMapper(rows *[]BatchCycle) dbmapper.RowMapper {
	return func() *dbmapper.MappedColumns {
		row := BatchCycle{}
		return batchCycleMapper(&row).Then(func() error {
			*rows = append(*rows, row)
			return nil
		})
	}
}
