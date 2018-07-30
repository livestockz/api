package feed

import (
	"time"

	"github.com/guregu/null"
	"github.com/satori/go.uuid"
)

const (
	Deleted_True    string = "1"
	Deleted_False   string = "0"
	Deleted_Any     string = ""
	Feed_Incoming   string = "incoming"
	Feed_Outgoing   string = "outgoing"
	Feed_Adjustment string = "adjustment"
)

type FeedType struct {
	ID      uuid.UUID   `json:"id"`
	Name    string      `json:"name"`
	Unit    string      `json:"unit"`
	Status  int32       `json:"status"`
	Deleted bool        `json:"deleted"`
	Created time.Time   `json:"created"`
	Updated null.String `json:"updated"`
}

type FeedIncoming struct {
	ID         uuid.UUID `json:"id"`
	FeedType   FeedType  `json:"feed_type"`
	FeedTypeID uuid.UUID `json:"-"`
	Qty        float64   `json:"qty"`
	Remarks    string    `json:"remarks"`
	Created    time.Time `json:"created"`
}

type FeedAdjustment struct {
	ID         uuid.UUID `json:"id"`
	FeedType   FeedType  `json:"feed_type"`
	FeedTypeID uuid.UUID `json:"-"`
	Qty        float64   `json:"qty"`
	Remarks    string    `json:"remarks"`
	Created    time.Time `json:"created"`
}
