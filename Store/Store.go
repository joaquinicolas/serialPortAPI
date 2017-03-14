package Store

import (
	_"github.com/mattn/go-sqlite3"
	"github.com/joaquinicolas/Elca/libs"
	"golang.org/x/net/websocket"
	"database/sql"
	"reflect"
)

var stores map[string] *Storer

type Storer interface {
	Name() string

}

type NewStore func(dsn string)(*Storer)

type SQLiteStore struct {
	DriverName string
	DataSource string
}

func (s *SQLiteStore) Name() string{
	return s.DriverName
}

//CreateDB creates database object
func (s *SQLiteStore)  CreateDB() *sql.DB{

	db, err := sql.Open(s.DriverName,s.DataSource)
	if err != nil {
		libs.Error.Println(err)
		return nil
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		libs.Error.Println(err)
		return nil
	}

	defer db.Close()
	return db
}

//NewSQLiteStore creates an instance of SQLiteStore
func NewSQLiteStore(dsn string) (*SQLiteStore){
	return &SQLiteStore{
		DriverName:"sqlite3",
		DataSource:dsn,
	}
}


//Register register a storer
func Register(name string, store Storer)  {
	_, ok := stores[name]
	if ok {
		libs.Warning.Println("The dsn alredy exists")
		return
	}
	stores[name] = &store

}
func init()  {
	store := NewSQLiteStore("./elca.db")
	Register(store.Name(),store)
}


type Creator interface {
	Store(db *sql.DB,data interface{})
}

type Reader interface {
	Read(db *sql.DB) (interface{})
	List(db *sql.DB) ([]interface{})
}

