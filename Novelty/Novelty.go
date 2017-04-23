package Novelty

import (
	"database/sql"
)

type Novelty struct {
	Id   sql.NullInt64
	Text string
}
