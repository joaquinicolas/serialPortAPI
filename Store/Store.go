package Store

import (
	_"github.com/mattn/go-sqlite3"
	"github.com/joaquinicolas/Elca/libs"
	"golang.org/x/net/websocket"
)

var stores map[string] *Storer

type Storer interface {
	Name() string

}

type NewStore func(dsn string)(Storer)

type SQLiteStore struct {
	DriverName string
	DataSource string
}

func (s *SQLiteStore) Name() string{
	return s.DriverName
}

func NewSQLiteStore(dsn string) (Storer){
	return &SQLiteStore{
		DriverName:"sqlite3",
		DataSource:dsn,
	}
}


func Register(name string, store Storer)  {
	_, ok := stores[name]
	if ok {
		libs.Warning.Println("The dsn alredy exists")
		return
	}
	stores[name] = &store

}
func init()  {
}


