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
	//feed
	ResolveFeedPage(page int32, limit int32, deleted string) (*[]Feed, int32, int32, int32, error)
	ResolveFeedByID(id uuid.UUID) (*Feed, error)
	InsertFeed(feed *Feed) (*Feed, error)
	UpdateFeedByID(feed *Feed) (*Feed, error)
	RemoveFeedByID(id uuid.UUID) (*Feed, error)
	RemoveFeedByIDs(ids []uuid.UUID) (*[]Feed, error)
}

const (
	//feedtype
	selectFeedType = `SELECT id, name, unit, status, deleted, created, updated FROM feed_type`
	insertFeedType = `INSERT INTO feed_type(id, name, unit, status, deleted, created) VALUES (:id ,:name, :unit, :status, :deleted, NOW())`
	updateFeedType = `UPDATE feed_type SET name = :name, unit = :unit, status = :status, deleted = :deleted, updated = NOW() WHERE id = :id`
	deleteFeedType = `UPDATE feed_type SET deleted = 1, updated = NOW() WHERE id = :id`
	//feed
	selectFeed = `SELECT id, feed_type_id, qty, remarks, reference, deleted, created, updated FROM feed_type`
	insertFeed = `INSERT INTO feed_type(id, feed_type_id, qty, remarks, reference, deleted, created) VALUES (:id ,:feedtype, :qty, :remarks, :reference, :deleted, NOW())`
	updateFeed = `UPDATE feed_type SET feed_type_id = :feedtype, qty = :qty, remarks= :remarks, reference = :reference, deleted = :deleted, updated = NOW() WHERE id = :id`
	deleteFeed = `UPDATE feed_type SET deleted = 1, updated = NOW() WHERE id = :id`
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

//feed
func (repo *FeedRepository) ResolveFeedPage(page int32, limit int32, deleted string) (*[]Feed, int32, int32, int32, error) {
	var start int32
	var end int32

	start = page * limit
	end = limit

	//get data by given page
	var query dbmapper.QueryMapper
	if deleted == Deleted_Any {
		query = dbmapper.Prepare(selectFeed+" ORDER BY name ASC LIMIT :start, :end").With(
			dbmapper.Param("start", start),
			dbmapper.Param("end", end),
		)
	} else if deleted == Deleted_True {
		query = dbmapper.Prepare(selectFeed+" WHERE deleted = 1 ORDER BY name ASC LIMIT :start, :end").With(
			dbmapper.Param("start", start),
			dbmapper.Param("end", end),
		)
	} else if deleted == Deleted_False {
		query = dbmapper.Prepare(selectFeed+" WHERE deleted = 0 ORDER BY name ASC LIMIT :start, :end").With(
			dbmapper.Param("start", start),
			dbmapper.Param("end", end),
		)
	}
	if err := query.Error(); err != nil {
		return nil, page, limit, 0, err
	}

	feeds := make([]Feed, 0)
	err := Parse(repo.DB.Query(query.SQL(), query.Params()...)).Map(feedsMapper(&feeds))

	if err != nil {
		return nil, page, limit, 0, err
	}

	//get total feed
	var summary dbmapper.QueryMapper
	if deleted == Deleted_Any {
		summary = dbmapper.Prepare("SELECT COUNT(*) AS total FROM feed")
	} else if deleted == Deleted_True {
		summary = dbmapper.Prepare("SELECT COUNT(*) AS total FROM feed WHERE deleted = 1")
	} else if deleted == Deleted_False {
		summary = dbmapper.Prepare("SELECT COUNT(*) AS total FROM feed WHERE deleted = 0")
	}

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
	return &feeds, page, limit, feedsCount, nil

}

func (repo *FeedRepository) ResolveFeedByID(id uuid.UUID) (*Feed, error) {
	query := dbmapper.Prepare(selectFeed + " WHERE id = :id").With(
		dbmapper.Param("id", id),
	)
	if err := query.Error(); err != nil {
		return nil, err
	}
	feeds := make([]Feed, 0)
	err := Parse(repo.DB.Query(query.SQL(), query.Params()...)).Map(feedsMapper(&feeds))

	if err != nil {
		return nil, err
	}
	if len(feeds) < 1 {
		return nil, fmt.Errorf("feed with id %s not found", id)
	}
	return &feeds[0], nil
	// feedtypes
	//&feed[0].feed_type = feedtype[0]
}

func (repo *FeedRepository) InsertFeed(feed *Feed) (*Feed, error) {

	//prepare query and params
	insert := dbmapper.Prepare(insertFeed).With(
		dbmapper.Param("id", feed.ID),
		dbmapper.Param("feedtype", feed.FeedType),
		dbmapper.Param("qty", feed.Qty),
		dbmapper.Param("remarks", feed.Remarks),
		dbmapper.Param("reference", feed.Reference),
		dbmapper.Param("deleted", feed.Deleted),
	)
	//validate query
	if err := insert.Error(); err != nil {
		//log.Print(err.Error())
		//fmt.Print("\n")
		return nil, err
	} else {
		//insert to database
		if _, err := repo.DB.Exec(insert.SQL(), insert.Params()...); err != nil {
			return nil, err
		} else {
			//find inserted data from database based on generated id
			res, err := repo.ResolveFeedByID(feed.ID)
			return res, err
		}
	}
}

func (repo *FeedRepository) UpdateFeedByID(feed *Feed) (*Feed, error) {
	_, err := repo.ResolveFeedByID(feed.ID)

	if err != nil {
		return nil, err
	} else {
		//update
		//prepare query and params
		updater := dbmapper.Prepare(updateFeedType).With(
			dbmapper.Param("feedtype", feed.FeedType),
			dbmapper.Param("qty", feed.Qty),
			dbmapper.Param("remarks", feed.Remarks),
			dbmapper.Param("reference", feed.Reference),
			dbmapper.Param("deleted", feed.Deleted),
			dbmapper.Param("id", feed.ID),
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
				res, err := repo.ResolveFeedByID(feed.ID)
				return res, err
			}
		}
	}
}

func (repo *FeedRepository) RemoveFeedByID(id uuid.UUID) (*Feed, error) {
	//find whether if data exist
	if _, err := repo.ResolveFeedByID(id); err != nil {
		return nil, err
	} else {
		remover := dbmapper.Prepare(deleteFeed).With(
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

func (repo *FeedRepository) RemoveFeedByIDs(ids []uuid.UUID) (*[]Feed, error) {
	for _, v := range ids {
		if _, err := repo.RemoveFeedByID(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func feedMapper(row *Feed) *dbmapper.MappedColumns {
	return dbmapper.Columns(
		dbmapper.Column("id").As(&row.ID),
		dbmapper.Column("feed_type").As(&row.FeedType),
		dbmapper.Column("qty").As(&row.Qty),
		dbmapper.Column("remarks").As(&row.Remarks),
		dbmapper.Column("reference").As(&row.Reference),
		dbmapper.Column("deleted").As(&row.Deleted),
		dbmapper.Column("created").As(&row.Created),
		dbmapper.Column("updated").As(&row.Updated),
	)
}

func feedsMapper(rows *[]Feed) dbmapper.RowMapper {
	return func() *dbmapper.MappedColumns {
		row := Feed{}
		return feedMapper(&row).Then(func() error {
			*rows = append(*rows, row)
			return nil
		})
	}
}
