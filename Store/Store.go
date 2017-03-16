package Store

import (
	_"github.com/mattn/go-sqlite3"
	"github.com/joaquinicolas/Elca/libs"
	"database/sql"
	"github.com/joaquinicolas/Elca/Store/models"
	"github.com/mattn/go-sqlite3"
	"log"
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

//CreateDBCon creates database object
func (s *SQLiteStore) CreateDBCon() *sql.DB{

	sql.Register(s.DriverName,&sqlite3.SQLiteDriver{})
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


	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS news (id INTEGER PRIMARY KEY," +
			" text VARCHAR(250) NOT NULL )")
	if err != nil {
		log.Fatal(err)
	}



	defer db.Close()
	return db
}

func (s *SQLiteStore) ReadNews(id int) *models.News{

	database := s.CreateDBCon()
	stmt, err := database.Prepare("SELECT id,text FROM news WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()
	news := &models.News{}
	err = stmt.QueryRow(id).Scan(&news.Id,&news.Text)
	if err != nil {
		log.Fatal(err)
	}

	return news
}

func (s *SQLiteStore) ListNews() ([]*models.News){
	database := s.CreateDBCon()
	defer database.Close()
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
	database := s.CreateDBCon()
	defer database.Close()
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
	store := NewSQLiteStore("./elca.db")
	Register(store.Name(),store)
}





