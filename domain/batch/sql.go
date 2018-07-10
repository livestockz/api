package batch

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ncrypthic/dbmapper"
	. "github.com/ncrypthic/dbmapper/dialects/mysql"
)

type BatchRepository interface {
	//RessolvePage(page, limit int32, keyword string) (*utils.Page, error)
	//RessolveBatchByIDs(IDs ...int32) ([]Batch, error)
	RessolveBatchByID(ID int64) (*Batch, error)
	//StoreBatch(batch *Batch) (*Batch, error)
	//RemoveBatchByID(ID int32) (*Batch, error)
	//RemoveBatchByIDs(IDs ...int32) ([]Batch, error)
}

const (
	selectBatch = `SELECT id, name, status, deleted, created, updated FROM growth_batch`
	insertBatch = `INSERT INTO growth_batch(name, status, deleted, created) VALUES (:name, :status, :deleted, now())`
)

type batchRepository struct {
	db *sql.DB
}

func (repo *batchRepository) RessolveBatchByID(ID int64) (*Batch, error) {
	query := dbmapper.Prepare(selectBatch + " WHERE id = :id").With(
		dbmapper.Param("id", ID),
	)
	if err := query.Error(); err != nil {
		return nil, err
	}
	batches := make([]Batch, 0)
	log.Print("sql:", query.SQL())
	log.Print("sql params:", query.Params())
	err := Parse(repo.db.Query(query.SQL(), query.Params()...)).Map(batchesMapper(&batches))

	if err != nil {
		return nil, err
	}
	if len(batches) < 1 {
		return nil, fmt.Errorf("batch with id %d not found", ID)
	}
	return &batches[0], nil
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

func NewRepository(db *sql.DB) BatchRepository {
	return &batchRepository{db}
}
