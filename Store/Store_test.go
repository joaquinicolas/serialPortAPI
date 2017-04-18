package Store

import (
	"fmt"
	"testing"

	"github.com/joaquinicolas/Elca/Novelty"
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
	n := &Novelty.Novelty{
		Text: "Testing",
	}
	fmt.Printf("Store_test.go driver name: %s", store.DriverName)
	lastID, affectedRows, err := store.StoreNovelty(n)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Last inserted id:", lastID)
	fmt.Println("Affected rows:", affectedRows)

}
func TestSQLiteStore_ListNews(t *testing.T) {
	news, err := store.ListNovelty()
	if err != nil {
		t.Error(err)
	}

	for n := range news {
		fmt.Println(n)
	}

	fmt.Printf("News count: %d", len(news))
}
func TestSQLiteStore_ReadNews(t *testing.T) {
	news := store.ReadNovelty(1)
	if news == nil {
		fmt.Println("news cannot be nil")
		t.Fail()
	}

	fmt.Println(news)
}
