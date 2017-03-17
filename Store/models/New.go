package models

import (
	"database/sql"
	"github.com/joaquinicolas/Elca/Store"
)

type News struct {
	Id sql.NullInt64
	Text sql.NullString
	AddCh chan *Store.Storer
	ReadCh chan *Store.Storer
	ListCh chan *Store.Storer
}

