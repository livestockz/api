package feed

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type Service interface {
	ResolveFeedTypePage(page int32, limit int32, deleted string) (*[]FeedType, int32, int32, int32, error)
	ResolveFeedTypeByID(uuid.UUID) (*FeedType, error)
	StoreFeedType(*FeedType) (*FeedType, error)
	RemoveFeedTypeByID(uuid.UUID) (*FeedType, error)
	RemoveFeedTypeByIDs([]uuid.UUID) (*[]FeedType, error)

	ResolveFeedPage(page int32, limit int32, deleted string) (*[]Feed, int32, int32, int32, error)
	ResolveFeedByID(uuid.UUID) (*Feed, error)
	StoreFeed(*Feed) (*Feed, error)
	RemoveFeedByID(uuid.UUID) (*Feed, error)
	RemoveFeedByIDs([]uuid.UUID) (*[]Feed, error)
}

type FeedService struct {
	FeedRepository Repository `inject:"feedRepository"`
}

//feed type
func (svc *FeedService) ResolveFeedTypePage(page int32, limit int32, deleted string) (*[]FeedType, int32, int32, int32, error) {
	if feedtypes, page, limit, total, err := svc.FeedRepository.ResolveFeedTypePage(page, limit, deleted); err != nil {
		return nil, 0, 0, 0, err
	} else {
		return feedtypes, page, limit, total, nil
	}
}

func (svc *FeedService) ResolveFeedTypeByID(id uuid.UUID) (*FeedType, error) {
	if feedtype, err := svc.FeedRepository.ResolveFeedTypeByID(id); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		return feedtype, nil
	}
}

func (svc *FeedService) StoreFeedType(feedtype *FeedType) (*FeedType, error) {
	if feedtype.ID == uuid.Nil {
		feedtype.ID = uuid.Must(uuid.NewV4())
		if result, err := svc.FeedRepository.InsertFeedType(feedtype); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	} else {
		//update
		if result, err := svc.FeedRepository.UpdateFeedTypeByID(feedtype); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
}

func (svc *FeedService) RemoveFeedTypeByID(id uuid.UUID) (*FeedType, error) {
	if _, err := svc.FeedRepository.RemoveFeedTypeByID(id); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		return nil, nil
	}
}

func (svc *FeedService) RemoveFeedTypeByIDs(ids []uuid.UUID) (*[]FeedType, error) {
	if _, err := svc.FeedRepository.RemoveFeedTypeByIDs(ids); err != nil {
		return nil, err
	} else {
		return nil, nil
	}
}

//feed
func (svc *FeedService) ResolveFeedPage(page int32, limit int32, deleted string) (*[]Feed, int32, int32, int32, error) {
	if feeds, page, limit, total, err := svc.FeedRepository.ResolveFeedPage(page, limit, deleted); err != nil {
		return nil, 0, 0, 0, err
	} else {
		return feeds, page, limit, total, nil
	}
}

func (svc *FeedService) ResolveFeedByID(id uuid.UUID) (*Feed, error) {
	if feed, err := svc.FeedRepository.ResolveFeedByID(id); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		return feed, nil
	}
}

func (svc *FeedService) StoreFeed(feed *Feed) (*Feed, error) {
	if feed.ID == uuid.Nil {
		feed.ID = uuid.Must(uuid.NewV4())
		if result, err := svc.FeedRepository.InsertFeed(feed); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	} else {
		//update
		if result, err := svc.FeedRepository.UpdateFeedByID(feed); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
}

func (svc *FeedService) RemoveFeedByID(id uuid.UUID) (*Feed, error) {
	if _, err := svc.FeedRepository.RemoveFeedByID(id); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		return nil, nil
	}
}

func (svc *FeedService) RemoveFeedByIDs(ids []uuid.UUID) (*[]Feed, error) {
	if _, err := svc.FeedRepository.RemoveFeedByIDs(ids); err != nil {
		return nil, err
	} else {
		return nil, nil
	}
}
