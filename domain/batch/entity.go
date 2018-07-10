package batch

import (
	"time"

	"github.com/guregu/null"
)

type Batch struct {
	ID      int64       `json:"id"`
	Name    string      `json:"name"`
	Status  int32       `json:"status"`
	Deleted bool        `json:"deleted"`
	Created time.Time   `json:"created"`
	Updated null.String `json:"updated"`
}
