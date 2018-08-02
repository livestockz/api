package feed

import (
	"database/sql"
	"fmt"

	"github.com/ncrypthic/dbmapper"
	. "github.com/ncrypthic/dbmapper/dialects/mysql"
	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	//feedtype
	ResolveFeedTypePage(page int32, limit int32, deleted string) (*[]FeedType, int32, int32, int32, error)
	ResolveFeedTypeByID(id uuid.UUID) (*FeedType, error)
	InsertFeedType(feedtype *FeedType) (*FeedType, error)
	UpdateFeedTypeByID(feedtype *FeedType) (*FeedType, error)
	RemoveFeedTypeByID(id uuid.UUID) (*FeedType, error)
	RemoveFeedTypeByIDs(ids []uuid.UUID) (*[]FeedType, error)
	//feed incoming
	ResolveFeedIncomingPage(page int32, limit int32) (*[]FeedIncoming, int32, int32, int32, error)
	ResolveFeedIncomingByID(id uuid.UUID) (*FeedIncoming, error)
	InsertFeedIncoming(feedIncoming *FeedIncoming) (*FeedIncoming, error)
	//feed adjustment
	ResolveFeedAdjustmentPage(page int32, limit int32) (*[]FeedAdjustment, int32, int32, int32, error)
	ResolveFeedAdjustmentByID(id uuid.UUID) (*FeedAdjustment, error)
	InsertFeedAdjustment(feedAdjustment *FeedAdjustment) (*FeedAdjustment, error)
}

const (
	//feedtype
	selectFeedType         = `SELECT id, name, unit, status, deleted, created, updated FROM feed_type`
	selectMultipleFeedType = `SELECT id, name, unit, status, deleted, created, updated FROM feed_type WHERE id IN (:ids)`
	insertFeedType         = `INSERT INTO feed_type(id, name, unit, status, deleted, created) VALUES (:id ,:name, :unit, :status, :deleted, NOW())`
	updateFeedType         = `UPDATE feed_type SET name = :name, unit = :unit, status = :status, deleted = :deleted, updated = NOW() WHERE id = :id`
	deleteFeedType         = `UPDATE feed_type SET deleted = 1, updated = NOW() WHERE id = :id`
	//feed incoming
	selectFeedIncoming = `SELECT id, feed_type_id, qty, remarks, created FROM feed_incoming`
	insertFeedIncoming = `INSERT INTO feed_incoming(id, feed_type_id, qty, remarks, created) VALUES (:id ,:feedtype, :qty, :remarks, NOW())`
	//feed adjustment
	selectFeedAdjustment = `SELECT id, feed_type_id, qty, remarks, created FROM feed_adjustment`
	insertFeedAdjustment = `INSERT INTO feed_adjustment(id, feed_type_id, qty, remarks, created) VALUES (:id ,:feedtype, :qty, :remarks, NOW())`
)

type FeedRepository struct {
	DB *sql.DB `inject:"db"`
}

//feedtype
func (repo *FeedRepository) ResolveFeedTypePage(page int32, limit int32, deleted string) (*[]FeedType, int32, int32, int32, error) {
	var start int32
	var end int32

	start = page * limit
	end = limit

	//get data by given page
	var query dbmapper.QueryMapper
	if deleted == Deleted_Any {
		query = dbmapper.Prepare(selectFeedType+" ORDER BY name ASC LIMIT :start, :end").With(
			dbmapper.Param("start", start),
			dbmapper.Param("end", end),
		)
	} else if deleted == Deleted_True {
		query = dbmapper.Prepare(selectFeedType+" WHERE deleted = 1 ORDER BY name ASC LIMIT :start, :end").With(
			dbmapper.Param("start", start),
			dbmapper.Param("end", end),
		)
	} else if deleted == Deleted_False {
		query = dbmapper.Prepare(selectFeedType+" WHERE deleted = 0 ORDER BY name ASC LIMIT :start, :end").With(
			dbmapper.Param("start", start),
			dbmapper.Param("end", end),
		)
	}
	if err := query.Error(); err != nil {
		return nil, page, limit, 0, err
	}

	feedtypes := make([]FeedType, 0)
	err := Parse(repo.DB.Query(query.SQL(), query.Params()...)).Map(feedtypesMapper(&feedtypes))

	if err != nil {
		return nil, page, limit, 0, err
	}

	//get total feedtype
	var summary dbmapper.QueryMapper
	if deleted == Deleted_Any {
		summary = dbmapper.Prepare("SELECT COUNT(*) AS total FROM feed_type")
	} else if deleted == Deleted_True {
		summary = dbmapper.Prepare("SELECT COUNT(*) AS total FROM feed_type WHERE deleted = 1")
	} else if deleted == Deleted_False {
		summary = dbmapper.Prepare("SELECT COUNT(*) AS total FROM feed_type WHERE deleted = 0")
	}

	if err := summary.Error(); err != nil {
		return nil, page, limit, 0, err
	}

	var feedtypesCount int32
	total := make([]int32, 0)
	err = Parse(repo.DB.Query(summary.SQL())).Map(dbmapper.Int32("total", &total))
	if err != nil {
		return nil, page, limit, 0, err
	} else {
		feedtypesCount = total[0]
	}
	return &feedtypes, page, limit, feedtypesCount, nil

}

func (repo *FeedRepository) ResolveFeedTypeByID(id uuid.UUID) (*FeedType, error) {
	query := dbmapper.Prepare(selectFeedType + " WHERE id = :id").With(
		dbmapper.Param("id", id),
	)
	if err := query.Error(); err != nil {
		return nil, err
	}
	feedtypes := make([]FeedType, 0)
	err := Parse(repo.DB.Query(query.SQL(), query.Params()...)).Map(feedtypesMapper(&feedtypes))

	if err != nil {
		return nil, err
	}
	if len(feedtypes) < 1 {
		return nil, fmt.Errorf("growth feedtype with id %s not found", id)
	}
	return &feedtypes[0], nil
}

func (repo *FeedRepository) InsertFeedType(feedtype *FeedType) (*FeedType, error) {

	//insert
	//prepare query and params
	insert := dbmapper.Prepare(insertFeedType).With(
		dbmapper.Param("id", feedtype.ID),
		dbmapper.Param("name", feedtype.Name),
		dbmapper.Param("unit", feedtype.Unit),
		dbmapper.Param("status", feedtype.Status),
		dbmapper.Param("deleted", feedtype.Deleted),
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
			res, err := repo.ResolveFeedTypeByID(feedtype.ID)
			return res, err
		}
	}
}

func (repo *FeedRepository) UpdateFeedTypeByID(feedtype *FeedType) (*FeedType, error) {
	//find whether if data exist
	_, err := repo.ResolveFeedTypeByID(feedtype.ID)

	if err != nil {
		return nil, err
	} else {
		//update
		updater := dbmapper.Prepare(updateFeedType).With(
			dbmapper.Param("name", feedtype.Name),
			dbmapper.Param("unit", feedtype.Unit),
			dbmapper.Param("status", feedtype.Status),
			dbmapper.Param("deleted", feedtype.Deleted),
			dbmapper.Param("id", feedtype.ID),
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
				res, err := repo.ResolveFeedTypeByID(feedtype.ID)
				return res, err
			}
		}
	}
}

func (repo *FeedRepository) RemoveFeedTypeByID(id uuid.UUID) (*FeedType, error) {
	//find whether if data exist
	if _, err := repo.ResolveFeedTypeByID(id); err != nil {
		return nil, err
	} else {
		remover := dbmapper.Prepare(deleteFeedType).With(
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

func (repo *FeedRepository) RemoveFeedTypeByIDs(ids []uuid.UUID) (*[]FeedType, error) {
	for _, v := range ids {
		if _, err := repo.RemoveFeedTypeByID(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func feedtypeMapper(row *FeedType) *dbmapper.MappedColumns {
	return dbmapper.Columns(
		dbmapper.Column("id").As(&row.ID),
		dbmapper.Column("name").As(&row.Name),
		dbmapper.Column("unit").As(&row.Unit),
		dbmapper.Column("status").As(&row.Status),
		dbmapper.Column("deleted").As(&row.Deleted),
		dbmapper.Column("created").As(&row.Created),
		dbmapper.Column("updated").As(&row.Updated),
	)
}

func feedtypesMapper(rows *[]FeedType) dbmapper.RowMapper {
	return func() *dbmapper.MappedColumns {
		row := FeedType{}
		return feedtypeMapper(&row).Then(func() error {
			*rows = append(*rows, row)
			return nil
		})
	}
}

//feed incoming
func (repo *FeedRepository) ResolveFeedIncomingPage(page int32, limit int32) (*[]FeedIncoming, int32, int32, int32, error) {
	var start int32
	var end int32

	start = page * limit
	end = limit

	//get data by given page
	var query dbmapper.QueryMapper
	query = dbmapper.Prepare(selectFeedIncoming+" ORDER BY created ASC LIMIT :start, :end").With(
		dbmapper.Param("start", start),
		dbmapper.Param("end", end),
	)

	if err := query.Error(); err != nil {
		return nil, page, limit, 0, err
	}

	feedIncomings := make([]FeedIncoming, 0)
	err := Parse(repo.DB.Query(query.SQL(), query.Params()...)).Map(feedIncomingsMapper(&feedIncomings))

	if err != nil {
		return nil, page, limit, 0, err
	}

	var newFeedIncomings []FeedIncoming
	for _, feedIncoming := range feedIncomings {
		if feedType, err := repo.ResolveFeedTypeByID(feedIncoming.FeedTypeID); err != nil {
			return nil, page, limit, 0, err
		} else {
			feedIncoming.FeedType = *feedType
			newFeedIncomings = append(newFeedIncomings, feedIncoming)
		}
	}

	//get total feed
	var summary dbmapper.QueryMapper
	summary = dbmapper.Prepare("SELECT COUNT(*) AS total FROM feed_incoming")

	if err := summary.Error(); err != nil {
		return nil, page, limit, 0, err
	}

	var feedIncomingsCount int32
	total := make([]int32, 0)
	err = Parse(repo.DB.Query(summary.SQL())).Map(dbmapper.Int32("total", &total))
	if err != nil {
		//log.Print(err.Error())
		return nil, page, limit, 0, err
	} else {
		feedIncomingsCount = total[0]
	}
	return &newFeedIncomings, page, limit, feedIncomingsCount, nil

}

func (repo *FeedRepository) ResolveFeedIncomingByID(id uuid.UUID) (*FeedIncoming, error) {
	query := dbmapper.Prepare(selectFeedIncoming + " WHERE id = :id").With(
		dbmapper.Param("id", id),
	)
	if err := query.Error(); err != nil {
		return nil, err
	}
	feedIncomings := make([]FeedIncoming, 0)
	err := Parse(repo.DB.Query(query.SQL(), query.Params()...)).Map(feedIncomingsMapper(&feedIncomings))

	if err != nil {
		return nil, err
	}
	if len(feedIncomings) < 1 {
		return nil, fmt.Errorf("feed incoming with id %s not found", id)
	}

	if feedtype, err := repo.ResolveFeedTypeByID(feedIncomings[0].FeedTypeID); err != nil {
		return nil, err
	} else {
		feedIncomings[0].FeedType = *feedtype
	}

	return &feedIncomings[0], nil
}

func (repo *FeedRepository) InsertFeedIncoming(feedIncoming *FeedIncoming) (*FeedIncoming, error) {

	//prepare query and params
	insert := dbmapper.Prepare(insertFeedIncoming).With(
		dbmapper.Param("id", feedIncoming.ID),
		dbmapper.Param("feedtype", feedIncoming.FeedType.ID),
		dbmapper.Param("qty", feedIncoming.Qty),
		dbmapper.Param("remarks", feedIncoming.Remarks),
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
			res, err := repo.ResolveFeedIncomingByID(feedIncoming.ID)
			return res, err
		}
	}
}

func feedIncomingMapper(row *FeedIncoming) *dbmapper.MappedColumns {
	return dbmapper.Columns(
		dbmapper.Column("id").As(&row.ID),
		dbmapper.Column("feed_type_id").As(&row.FeedTypeID),
		dbmapper.Column("qty").As(&row.Qty),
		dbmapper.Column("remarks").As(&row.Remarks),
		dbmapper.Column("created").As(&row.Created),
	)
}

func feedIncomingsMapper(rows *[]FeedIncoming) dbmapper.RowMapper {
	return func() *dbmapper.MappedColumns {
		row := FeedIncoming{}
		return feedIncomingMapper(&row).Then(func() error {
			*rows = append(*rows, row)
			return nil
		})
	}
}

//feed adjustment
func (repo *FeedRepository) ResolveFeedAdjustmentPage(page int32, limit int32) (*[]FeedAdjustment, int32, int32, int32, error) {
	var start int32
	var end int32

	start = page * limit
	end = limit

	//get data by given page
	var query dbmapper.QueryMapper
	query = dbmapper.Prepare(selectFeedAdjustment+" ORDER BY created ASC LIMIT :start, :end").With(
		dbmapper.Param("start", start),
		dbmapper.Param("end", end),
	)

	if err := query.Error(); err != nil {
		return nil, page, limit, 0, err
	}

	feedAdjustments := make([]FeedAdjustment, 0)
	err := Parse(repo.DB.Query(query.SQL(), query.Params()...)).Map(feedAdjustmentsMapper(&feedAdjustments))

	if err != nil {
		return nil, page, limit, 0, err
	}

	var newFeedAdjustments []FeedAdjustment
	for _, feedAdjustment := range feedAdjustments {
		if feedType, err := repo.ResolveFeedTypeByID(feedAdjustment.FeedTypeID); err != nil {
			return nil, page, limit, 0, err
		} else {
			feedAdjustment.FeedType = *feedType
			newFeedAdjustments = append(newFeedAdjustments, feedAdjustment)
		}
	}

	//get total feed
	var summary dbmapper.QueryMapper
	summary = dbmapper.Prepare("SELECT COUNT(*) AS total FROM feed_adjustment")

	if err := summary.Error(); err != nil {
		return nil, page, limit, 0, err
	}

	var feedsCount int32
	total := make([]int32, 0)
	err = Parse(repo.DB.Query(summary.SQL())).Map(dbmapper.Int32("total", &total))
	if err != nil {
		//log.Print(err.Error())
		return nil, page, limit, 0, err
	} else {
		feedsCount = total[0]
	}
	return &newFeedAdjustments, page, limit, feedsCount, nil

}

func (repo *FeedRepository) ResolveFeedAdjustmentByID(id uuid.UUID) (*FeedAdjustment, error) {
	query := dbmapper.Prepare(selectFeedAdjustment + " WHERE id = :id").With(
		dbmapper.Param("id", id),
	)
	if err := query.Error(); err != nil {
		return nil, err
	}
	feedAdjustments := make([]FeedAdjustment, 0)
	err := Parse(repo.DB.Query(query.SQL(), query.Params()...)).Map(feedAdjustmentsMapper(&feedAdjustments))

	if err != nil {
		return nil, err
	}
	if len(feedAdjustments) < 1 {
		return nil, fmt.Errorf("feed adjustment with id %s not found", id)
	}
	if feedtype, err := repo.ResolveFeedTypeByID(feedAdjustments[0].FeedTypeID); err != nil {
		return nil, err
	} else {
		feedAdjustments[0].FeedType = *feedtype
	}

	return &feedAdjustments[0], nil
}

func (repo *FeedRepository) InsertFeedAdjustment(feedAdjustment *FeedAdjustment) (*FeedAdjustment, error) {

	//prepare query and params
	insert := dbmapper.Prepare(insertFeedAdjustment).With(
		dbmapper.Param("id", feedAdjustment.ID),
		dbmapper.Param("feedtype", feedAdjustment.FeedType.ID),
		dbmapper.Param("qty", feedAdjustment.Qty),
		dbmapper.Param("remarks", feedAdjustment.Remarks),
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
			res, err := repo.ResolveFeedAdjustmentByID(feedAdjustment.ID)
			return res, err
		}
	}
}

func feedAdjustmentMapper(row *FeedAdjustment) *dbmapper.MappedColumns {
	return dbmapper.Columns(
		dbmapper.Column("id").As(&row.ID),
		dbmapper.Column("feed_type_id").As(&row.FeedTypeID),
		dbmapper.Column("qty").As(&row.Qty),
		dbmapper.Column("remarks").As(&row.Remarks),
		dbmapper.Column("created").As(&row.Created),
	)
}

func feedAdjustmentsMapper(rows *[]FeedAdjustment) dbmapper.RowMapper {
	return func() *dbmapper.MappedColumns {
		row := FeedAdjustment{}
		return feedAdjustmentMapper(&row).Then(func() error {
			*rows = append(*rows, row)
			return nil
		})
	}
}
