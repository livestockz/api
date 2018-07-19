package batch

import (
	"time"

	"github.com/guregu/null"
	"github.com/satori/go.uuid"
)

const (
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
