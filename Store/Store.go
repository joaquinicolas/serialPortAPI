package Store

import (
	_"github.com/mattn/go-sqlite3"
	"github.com/joaquinicolas/Elca/libs"
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"log"
	"sync"
	"github.com/joaquinicolas/Elca/Store/models"
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

//getInstance creates database object
func (s *SQLiteStore) getInstance() *sql.DB{

	createCon := func() {
		sql.Register(s.DriverName,&sqlite3.SQLiteDriver{})
		db, err := sql.Open(s.DriverName,s.DataSource)
		if err != nil {
			libs.Error.Println(err)
			return
		}
		defer db.Close()
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

func (s *SQLiteStore) ListNews() ([]*models.News){
	database := s.getInstance()
	rows, err := database.Query("SELECT * FROM news")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var result []*models.News
	for rows.Next() {
		var data models.News
		err := rows.Scan(&data.Id,&data.Text)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result,&data)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return result
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
	stores = make(map[string] *Storer)
	store := NewSQLiteStore("./elca.db")
	Register(store.Name(),store)
}





