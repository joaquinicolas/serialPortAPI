package Store

import (





	"testing"
	"fmt"

	"github.com/joaquinicolas/Elca/Store/models"
)

var store SQLiteStore

func TestNewSQLiteStore(t *testing.T) {
	store, err := GetStore("sqlite3")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(store)

}

func TestSQLiteStore_StoreNews(t *testing.T) {
	n := &models.News{








		Text:"Testing",
			

	}
	lastId, affectedRows, err := store.StoreNews(n)
	if err := nil {
		t.Error(err)
		

	}

	fmt.Println("Last inserted id:", lastId)
	fmt.Println("Affected rows:", affectedRows)

}
func TestSQLiteStore_ListNews(t *testing.T) {
	news, err := store.ListNews()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(news)
}




