package batch

import (
	"time"

	"github.com/guregu/null"
	"github.com/satori/go.uuid"
)

type Batch struct {
	ID      uuid.UUID   `json:"id"`
	Name    string      `json:"name"`
	Status  int32       `json:"status"`
	Deleted bool        `json:"deleted"`
	Created time.Time   `json:"created"`
	Updated null.String `json:"updated"`
}
