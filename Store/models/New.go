package models

import (
	"database/sql"
)

type News struct {
	Id sql.NullInt64
	Text sql.NullString
}

