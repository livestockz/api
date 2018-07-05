package batch

import (
	"time"

	"github.com/guregu/null"
)

type Batch struct {
	ID      int32       `json:"id"`
	Name    string      `json:"name"`
	Status  int32       `json:"status"`
	Deleted bool        `json:"deleted"`
	UserID  int32       `json:"user_id"`
	Created time.Time   `json:"created"`
	Updated null.String `json:"updated"`
}
