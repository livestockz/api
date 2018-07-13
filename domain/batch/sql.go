package batch

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ncrypthic/dbmapper"
	. "github.com/ncrypthic/dbmapper/dialects/mysql"
	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	//ResolvePage(page, limit int32, keyword string) (*utils.Page, error)
	ResolveBatchByID(id uuid.UUID) (*Batch, error)
	//ResolveBatchByIDs(IDs ...int32) ([]Batch, error)
	StoreBatch(batch *Batch) (*Batch, error)
	//RemoveBatchByID(ID int32) (*Batch, error)
	//RemoveBatchByIDs(IDs ...int32) ([]Batch, error)
}

const (
	findBatch   = `SELECT id, name, status, deleted, created, updated FROM growth_batch`
	insertBatch = `INSERT INTO growth_batch(id, name, status, deleted, created)`
	updateBatch = `UPDATE growth_batch SET name = :name, status = :status, deleted = :deleted, updated = now() WHERE id = :id`
	deleteBatch = `UPDATE growth_batch SET deleted = 1 WHERE id = :id`
)

type BatchRepository struct {
	DB *sql.DB `inject:"db"`
}

func (repo *BatchRepository) ResolveBatchByID(id uuid.UUID) (*Batch, error) {
	query := dbmapper.Prepare(findBatch + " WHERE id = :id").With(
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
		return nil, fmt.Errorf("batch with id %s not found", id)
	}
	return &batches[0], nil
}

func (repo *BatchRepository) StoreBatch(batch *Batch) (*Batch, error) {
	//generate uuid
	id := uuid.Must(uuid.NewV4())
	//prepare query and params
	query := dbmapper.Prepare(insertBatch+" VALUES (:id ,:name, :status, :deleted, :created)").With(
		dbmapper.Param("id", id),
		dbmapper.Param("name", batch.Name),
		dbmapper.Param("status", batch.Status),
		dbmapper.Param("deleted", batch.Deleted),
		dbmapper.Param("created", time.Now()),
	)
	//validate query
	if err := query.Error(); err != nil {
		return nil, err
	} else {
		//insert to database
		_, err := repo.DB.Query(query.SQL(), query.Params()...)
		if err != nil {
			return nil, err
		} else {
			//find inserted data from database based on generated id
			res, err := repo.ResolveBatchByID(id)
			return res, err
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
