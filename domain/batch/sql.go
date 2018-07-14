package batch

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ncrypthic/dbmapper"
	. "github.com/ncrypthic/dbmapper/dialects/mysql"
	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	//ResolvePage(page, limit int32, keyword string) (*utils.Page, error)
	ResolveGrowthBatches(page int32, limit int32) (*[]Batch, int32, int32, int32, error)
	ResolveGrowthBatchByID(id uuid.UUID) (*Batch, error)
	//ResolveBatchByIDs(IDs ...int32) ([]Batch, error)
	StoreGrowthBatch(batch *Batch) (*Batch, error)
	RemoveGrowthBatchByID(ID uuid.UUID) (*Batch, error)
	//RemoveBatchByIDs(IDs ...int32) ([]Batch, error)
}

const (
	selectGrowthBatch = `SELECT id, name, status, deleted, created, updated FROM growth_batch`
	insertGrowthBatch = `INSERT INTO growth_batch(id, name, status, deleted, created) VALUES (:id ,:name, :status, :deleted, NOW())`
	updateGrowthBatch = `UPDATE growth_batch SET name = :name, status = :status, deleted = :deleted, updated = NOW() WHERE id = :id`
	deleteGrowthBatch = `UPDATE growth_batch SET deleted = 1 WHERE id = :id`
)

type BatchRepository struct {
	DB *sql.DB `inject:"db"`
}

func (repo *BatchRepository) ResolveGrowthBatches(page int32, limit int32) (*[]Batch, int32, int32, int32, error) {
	var start int32
	var end int32
	if limit == 0 {
		limit = 10
	}
	start = page * limit
	end = start + limit
	//get data by given page
	query := dbmapper.Prepare(selectGrowthBatch+" WHERE deleted = 0 LIMIT :start, :end").With(
		dbmapper.Param("start", start),
		dbmapper.Param("end", end),
	)
	if err := query.Error(); err != nil {
		return nil, page, limit, 0, err
	}

	batches := make([]Batch, 0)
	err := Parse(repo.DB.Query(query.SQL(), query.Params()...)).Map(batchesMapper(&batches))

	if err != nil {
		return nil, page, limit, 0, err
	}

	//get total batch
	summary := dbmapper.Prepare("SELECT COUNT(*) AS total FROM growth_batch WHERE deleted = 0")
	if err := summary.Error(); err != nil {
		return nil, page, limit, 0, err
	}

	var totalItem int32
	total := make([]int32, 0)
	err = Parse(repo.DB.Query(summary.SQL(), summary.Params())).Map(Int32("total", &total))
	if err != nil {
		return nil, page, limit, 0, err
	} else {
		for _, r := range total {
			fmt.Printf("%+v\n", r)
			totalItem = int32(r)
		}
	}
	//fmt.Println(&total)

	return &batches, page, limit, totalItem, nil
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
}

func (repo *BatchRepository) StoreGrowthBatch(batch *Batch) (*Batch, error) {
	//find whether if data exist
	fmt.Print("\n")
	fmt.Print(batch)
	fmt.Print("\n")
	finder := dbmapper.Prepare(selectGrowthBatch + " WHERE id = :id").With(
		dbmapper.Param("id", batch.ID),
	)
	if err := finder.Error(); err != nil {
		return nil, err
	}
	batches := make([]Batch, 0)
	err := Parse(repo.DB.Query(finder.SQL(), finder.Params()...)).Map(batchesMapper(&batches))

	if err != nil {
		fmt.Print(err)
		fmt.Print("\n")
		return nil, err
	}
	if len(batches) < 1 {
		//insert
		fmt.Println("insert")
		fmt.Print("\n")

		//generate uuid
		id := uuid.Must(uuid.NewV4())

		//prepare query and params
		insert := dbmapper.Prepare(insertGrowthBatch).With(
			dbmapper.Param("id", id),
			dbmapper.Param("name", batch.Name),
			dbmapper.Param("status", batch.Status),
			dbmapper.Param("deleted", batch.Deleted),
		)
		log.Print("sql:", insert.SQL())
		fmt.Print("\n")
		log.Print("sql params:", insert.Params())
		fmt.Print("\n")
		//validate query
		if err := insert.Error(); err != nil {
			log.Print(err.Error())
			fmt.Print("\n")
			return nil, err
		} else {
			//insert to database
			if _, err := repo.DB.Exec(insert.SQL(), insert.Params()...); err != nil {
				log.Print(err.Error())
				fmt.Print("\n")
				return nil, err
			} else {
				//find inserted data from database based on generated id
				res, err := repo.ResolveGrowthBatchByID(id)
				return res, err
			}
		}
	} else {
		//update
		fmt.Println("update")
		fmt.Print("\n")
		//prepare query and params
		updater := dbmapper.Prepare(updateGrowthBatch).With(
			dbmapper.Param("name", batch.Name),
			dbmapper.Param("status", batch.Status),
			dbmapper.Param("deleted", batch.Deleted),
			dbmapper.Param("id", batch.ID),
		)
		fmt.Print("\n")
		log.Print("sql:", updater.SQL())
		log.Print("sql params:", updater.Params())
		fmt.Print("\n")
		//validate query
		if err := updater.Error(); err != nil {
			log.Print(err.Error())
			fmt.Print("\n")
			return nil, err
		} else {
			//update to database
			if _, err := repo.DB.Exec(updater.SQL(), updater.Params()...); err != nil {
				log.Print(err.Error())
				fmt.Print("\n")
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
	fmt.Print("\n")
	fmt.Print(id)
	fmt.Print("\n")
	finder := dbmapper.Prepare(selectGrowthBatch + " WHERE id = :id").With(
		dbmapper.Param("id", id),
	)
	if err := finder.Error(); err != nil {
		return nil, err
	}
	batches := make([]Batch, 0)
	err := Parse(repo.DB.Query(finder.SQL(), finder.Params()...)).Map(batchesMapper(&batches))

	if err != nil {
		fmt.Print(err)
		fmt.Print("\n")
		return nil, err
	}
	if len(batches) < 1 {
		return nil, fmt.Errorf("batch with id %s not found", id)
	} else {
		//update
		fmt.Println("update")
		fmt.Print("\n")
		//prepare query and params
		remover := dbmapper.Prepare(deleteGrowthBatch).With(
			dbmapper.Param("id", id),
		)
		fmt.Print("\n")
		log.Print("sql:", remover.SQL())
		log.Print("sql params:", remover.Params())
		fmt.Print("\n")
		//validate query
		if err := remover.Error(); err != nil {
			log.Print(err.Error())
			fmt.Print("\n")
			return nil, err
		} else {
			//update to database
			if _, err := repo.DB.Exec(remover.SQL(), remover.Params()...); err != nil {
				log.Print(err.Error())
				fmt.Print("\n")
				return nil, err
			} else {
				//find inserted data from database based on given id
				res, err := repo.ResolveGrowthBatchByID(id)
				return res, err
			}
		}
	}
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
