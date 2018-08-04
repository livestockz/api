package feed

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type Service interface {
	ResolveFeedTypePage(page int32, limit int32, deleted string) (*[]FeedType, int32, int32, int32, error)
	ResolveFeedTypeByIDs([]uuid.UUID) (*[]FeedType, error)
	ResolveFeedTypeByID(uuid.UUID) (*FeedType, error)
	StoreFeedType(*FeedType) (*FeedType, error)
	RemoveFeedTypeByID(uuid.UUID) (*FeedType, error)
	RemoveFeedTypeByIDs([]uuid.UUID) (*[]FeedType, error)

	ResolveFeedIncomingPage(page int32, limit int32) (*[]FeedIncoming, int32, int32, int32, error)
	ResolveFeedIncomingByID(uuid.UUID) (*FeedIncoming, error)
	StoreFeedIncoming(*FeedIncoming) (*FeedIncoming, error)

	ResolveFeedAdjustmentPage(page int32, limit int32) (*[]FeedAdjustment, int32, int32, int32, error)
	ResolveFeedAdjustmentByID(uuid.UUID) (*FeedAdjustment, error)
	StoreFeedAdjustment(*FeedAdjustment) (*FeedAdjustment, error)
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

func (svc *FeedService) ResolveFeedTypeByIDs(ids []uuid.UUID) (*[]FeedType, error) {
	if feedtypes, err := svc.FeedRepository.ResolveFeedTypeByIDs(ids); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		return feedtypes, nil
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

//feed incoming
func (svc *FeedService) ResolveFeedIncomingPage(page int32, limit int32) (*[]FeedIncoming, int32, int32, int32, error) {
	if feedIncomings, page, limit, total, err := svc.FeedRepository.ResolveFeedIncomingPage(page, limit); err != nil {
		return nil, 0, 0, 0, err
	} else {
		return feedIncomings, page, limit, total, nil
	}
}

func (svc *FeedService) ResolveFeedIncomingByID(id uuid.UUID) (*FeedIncoming, error) {
	if feedIncoming, err := svc.FeedRepository.ResolveFeedIncomingByID(id); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		return feedIncoming, nil
	}
}

func (svc *FeedService) StoreFeedIncoming(feedIncoming *FeedIncoming) (*FeedIncoming, error) {
	feedIncoming.ID = uuid.Must(uuid.NewV4())
	if result, err := svc.FeedRepository.InsertFeedIncoming(feedIncoming); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

//feed adjustment
func (svc *FeedService) ResolveFeedAdjustmentPage(page int32, limit int32) (*[]FeedAdjustment, int32, int32, int32, error) {
	if feedAdjustments, page, limit, total, err := svc.FeedRepository.ResolveFeedAdjustmentPage(page, limit); err != nil {
		return nil, 0, 0, 0, err
	} else {
		return feedAdjustments, page, limit, total, nil
	}
}

func (svc *FeedService) ResolveFeedAdjustmentByID(id uuid.UUID) (*FeedAdjustment, error) {
	if feedAdjustment, err := svc.FeedRepository.ResolveFeedAdjustmentByID(id); err != nil {
		return nil, fmt.Errorf("found an error: %s", err.Error())
	} else {
		return feedAdjustment, nil
	}
}

func (svc *FeedService) StoreFeedAdjustment(feedAdjustment *FeedAdjustment) (*FeedAdjustment, error) {
	feedAdjustment.ID = uuid.Must(uuid.NewV4())
	if result, err := svc.FeedRepository.InsertFeedAdjustment(feedAdjustment); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
