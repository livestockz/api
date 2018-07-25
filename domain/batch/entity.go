package batch

import (
	"time"

	"github.com/guregu/null"
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
	ID      uuid.UUID   `json:"id"`
	Name    string      `json:"name"`
	Status  int32       `json:"status"`
	Deleted bool        `json:"deleted"`
	Created time.Time   `json:"created"`
	Updated null.String `json:"updated"`
	//Pool []Pool
}

type Pool struct {
	ID      uuid.UUID   `json:"id"`
	Name    string      `json:"name"`
	Status  string      `json:"status"`
	Deleted bool        `json:"deleted"`
	Created time.Time   `json:"created"`
	Updated null.String `json:"updated"`
}

type BatchCycle struct {
	ID      uuid.UUID   `json:"id"`
	Batch   Batch       `json:"batch"`
	BatchID uuid.UUID   `json:"-"`
	Pool    Pool        `json:"pool"`
	PoolID  uuid.UUID   `json:"-"`
	Weight  float64     `json:"weight"`
	Amount  float64     `json:"amount"`
	Start   time.Time   `json:"start"`
	Finish  null.String `json:"finish"`
	Created time.Time   `json:"created"`
	Updated null.String `json:"updated"`
}
