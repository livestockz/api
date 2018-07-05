package batch

import (
	"database/sql"
	"fmt"

	"github.com/ncrypthic/dbmapper"
	. "github.com/ncrypthic/dbmapper/dialects/mysql"
)

type BatchRepository interface {
	//RessolvePage(page, limit int32, keyword string) (*utils.Page, error)
	//RessolveBatchByIDs(IDs ...int32) ([]Batch, error)
	RessolveBatchByID(ID int32) (*Batch, error)
	//StoreBatch(batch *Batch) (*Batch, error)
	//RemoveBatchByID(ID int32) (*Batch, error)
	//RemoveBatchByIDs(IDs ...int32) ([]Batch, error)
}

const (
	selectBatch = `SELECT id, name, status, deleted, user_id, created, updated FROM batch`
	insertBatch = `INSERT INTO batch(name, status, deleted, user_id, created) VALUES (:name, :status, :deleted, :user_id, now())`
)

type batchRepository struct {
	db *sql.DB
}

func (repo *batchRepository) RessolveBatchByID(ID int32) (*Batch, error) {
	query := dbmapper.Prepare(selectBatch + " Where id = :id").With(
		Param("id", ID),
	)
	if err := query.Error(); err != nil {
		return nil, err
	}
	batches := make([]Batch, 0)
	err := Parse(repo.db.Query(query.SQL(), query.Params()...)).Map(batchesMapper(batches))
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
		dbmapper.Column("user_id").As(&row.UserID),
		dbmapper.Column("created").As(&row.Created),
		dbmapper.Column("updated").As(&row.Updated),
	)
}

func batchesMapper(rows []Batch) dbmapper.RowMapper {
	return func() *dbmapper.MappedColumns {
		row := Batch{}
		return batchMapper(&row).Then(func() error {
			rows = append(rows, row)
			return nil
		})
	}
}

func NewRepository(db *sql.DB) BatchRepository {
	return &batchRepository{db}
}
