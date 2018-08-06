package batch

import (
	"time"

	"github.com/guregu/null"
	"github.com/livestockz/api/domain/feed"
	"github.com/satori/go.uuid"
)

const (
	Deleted_True     string = "1"
	Deleted_False    string = "0"
	Deleted_Any      string = ""
	Pool_Inactive    string = "inactive"
	Pool_Assigned    string = "assigned"
	Pool_Maintenance string = "maintenance"
)

type Batch struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Status  int32     `json:"status"`
	Deleted bool      `json:"deleted"`
	Created time.Time `json:"created"`
	Updated null.Time `json:"updated"`
	//Pool []Pool
}

type Pool struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Status  string    `json:"status"`
	Deleted bool      `json:"deleted"`
	Created time.Time `json:"created"`
	Updated null.Time `json:"updated"`
}

type BatchCycle struct {
	ID      uuid.UUID `json:"id"`
	Batch   Batch     `json:"batch"`
	BatchID uuid.UUID `json:"-"`
	Pool    Pool      `json:"pool"`
	PoolID  uuid.UUID `json:"-"`
	Weight  float64   `json:"weight"`
	Amount  float64   `json:"amount"`
	Feeding []Feeding `json:"feeding"`
	Deaths  []Death   `json:"deaths"`
	Start   time.Time `json:"start"`
	Finish  null.Time `json:"finish"`
	Created time.Time `json:"created"`
	Updated null.Time `json:"updated"`
}

type Death struct {
	ID           uuid.UUID `json:"id"`
	BatchCycleID uuid.UUID `json:"batch_cycle_id"`
	DeathDate    time.Time `json:"death_date"`
	Weight       float64   `json:"weight"`
	Amount       float64   `json:"amount"`
	Remarks      string    `json:"remarks"`
	Created      time.Time `json:"created"`
}

type Feeding struct {
	ID           uuid.UUID     `json:"id"`
	BatchCycleID uuid.UUID     `json:"batch_cycle_id"`
	FeedType     feed.FeedType `json:"feed_type"`
	FeedTypeID   uuid.UUID     `json:"-"`
	FeedingDate  time.Time     `json:"feeding_date"`
	Qty          float64       `json:"qty"`
	Remarks      string        `json:"remarks"`
	Created      time.Time     `json:"created"`
}

type CutOff struct {
	ID           uuid.UUID `json:"id"`
	BatchID      uuid.UUID `json:"batch_id"`
	BatchCycleID uuid.UUID `json:"batch_cycle_id"`
	SummaryDate  time.Time `json:"summary_date"`
	Weight       float64   `json:"weight"`
	Amount       float64   `json:"amount"`
	ADG          float64   `json:"adg"`
	FCR          float64   `json:"fcr"`
	SR           float64   `json:"sr"`
	Created      time.Time `json:"created"`
}
