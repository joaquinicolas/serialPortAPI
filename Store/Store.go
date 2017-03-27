package Store

import (
	_"github.com/mattn/go-sqlite3"
	"github.com/joaquinicolas/Elca/libs"
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"log"
	"sync"
	"github.com/joaquinicolas/Elca/Store/models"
	"errors"
	"fmt"
)

var once sync.Once
var stores map[string] *Storer

type Storer interface {
	Name() string

}

type NewStore func(dsn string)(*Storer)

type SQLiteStore struct {
	DriverName string
	DataSource string
	db *sql.DB
}

func (s *SQLiteStore) Name() string{
	return s.DriverName
}

//getInstance returns a instance of database object
func (s *SQLiteStore) getInstance() *sql.DB{

	createCon := func() {

		fmt.Println("Opening connection")
		sql.Register(s.DriverName,&sqlite3.SQLiteDriver{})
		db, err := sql.Open(s.DriverName,s.DataSource)
		if err != nil {
			libs.Error.Println(err)
			return
		}
		err = db.Ping()
		if err != nil {
			libs.Error.Println(err)
			return
		}


		_, err = db.Exec(
			"CREATE TABLE IF NOT EXISTS news (id INTEGER PRIMARY KEY," +
				" text VARCHAR(250) NOT NULL )")
		if err != nil {
			log.Fatal(err)
		}

		s.db = db
	}

	fmt.Println("Before open connection")
	once.Do(createCon)
	return s.db

}

func (s *SQLiteStore) ReadNews(id int) *models.News{

	database := s.getInstance()
	stmt, err := database.Prepare("SELECT id,text FROM news WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}

	news := &models.News{}
	err = stmt.QueryRow(id).Scan(&news.Id,&news.Text)
	if err != nil {
		log.Fatal(err)
	}

	return news
}

func (s *SQLiteStore) ListNews() ([]*models.News, error){
	database := s.getInstance()
	rows, err := database.Query("SELECT * FROM news")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []*models.News
	for rows.Next() {
		var data models.News
		err := rows.Scan(&data.Id,&data.Text)
		if err != nil {
			return nil, err
		}

		result = append(result,&data)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// StoreNews stores news and return lastId, rows affected or an error if exists.
func (s *SQLiteStore) StoreNews(n *models.News) (int64,int64,error) {
	database := s.getInstance()
	stmt, err := database.Prepare("INSERT INTO news(text) VALUES (?)")
	if err != nil {
		return 0,0, err
	}

	res, err := stmt.Exec(n.Text)
	if err != nil {
		return 0,0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0,0, err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return 0, 0, err
	}
	return lastId, rowCnt, nil

}

//newSQLiteStore creates an instance of SQLiteStore
func newSQLiteStore(dsn string) (*SQLiteStore){


	return &SQLiteStore{
		DriverName:"sqlite3",
		DataSource:dsn,
	}
}

func GetStore(dsn string) (*Storer, error) {
	if dsn == ""{
		return nil, errors.New("dsn cannot be empty string")
	}
	store, ok := stores[dsn]
	if ok {
		return store, nil
	}

	return nil, errors.New("Store not exists")

}

//register register a storer
func register(name string, store Storer)  {
	_, ok := stores[name]
	if ok {
		libs.Warning.Println("The dsn alredy exists")
		return
	}
	stores[name] = &store

}

func init()  {
	stores = make(map[string] *Storer)
	store := newSQLiteStore("./elca.db")
	register(store.Name(),store)
}





