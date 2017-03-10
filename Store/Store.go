package Store

import (
	"github.com/mattn/go-sqlite3"
	"database/sql"
	"github.com/joaquinicolas/Elca/libs"
)

var stores map[string] *Store

type Store struct {
	DriverName string
	DataSource string
}

type NewStore func(dsn string)(*Store)


func Register(name string, store Store)  {
	_, ok := stores[name]
	if ok {
		libs.Warning.Println("The dsn alredy exists")
		return
	}
	stores[name] = &store

}
func init()  {
	
}


