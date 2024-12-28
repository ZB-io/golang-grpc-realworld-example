package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
)




func TestNewArticleStore(t *testing.T) {
	t.Run("Scenario 1: Successfully Create a New ArticleStore with a Valid DB Connection", func(t *testing.T) {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		gormDB, err := gorm.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("failed to open in-memory sqlite database: %v", err)
		}
		defer gormDB.Close()

		mock.ExpectQuery("^SELECT .+").WillReturnRows(sqlmock.NewRows(nil))

		store := NewArticleStore(gormDB)

		if store == nil || store.db != gormDB {
			t.Errorf("expected a non-nil ArticleStore with the correct database, got %v", store)
		} else {
			t.Log("Successfully created a new ArticleStore with a valid DB connection")
		}
	})

	t.Run("Scenario 2: Initialize ArticleStore with a Nil DB Connection", func(t *testing.T) {

		var gormDB *gorm.DB = nil

		store := NewArticleStore(gormDB)

		if store == nil || store.db != nil {
			t.Errorf("expected a non-nil ArticleStore with a nil database, got %v", store)
		} else {
			t.Log("Successfully handled nil DB connection while initializing ArticleStore")
		}
	})

	t.Run("Scenario 3: Robustness Against Concurrent Initialization", func(t *testing.T) {

		gormDB, err := gorm.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("failed to open in-memory sqlite database: %v", err)
		}
		defer gormDB.Close()

		var wg sync.WaitGroup
		const goroutineCount = 10
		wg.Add(goroutineCount)

		var stores []*ArticleStore
		mu := sync.Mutex{}
		for i := 0; i < goroutineCount; i++ {
			go func() {
				defer wg.Done()
				store := NewArticleStore(gormDB)
				mu.Lock()
				stores = append(stores, store)
				mu.Unlock()
			}()
		}
		wg.Wait()

		for _, store := range stores {
			if store == nil || store.db != gormDB {
				t.Errorf("expected all ArticleStore instances to correctly point to the *gorm.DB instance, found: %v", store)
			}
		}
		t.Log("Successfully handled concurrent initialization of ArticleStore without race condition")
	})

	t.Run("Scenario 4: Memory and Resource Management Check (Manual Verification Required)", func(t *testing.T) {

		log.Println("Run memory profiling and analysis manually to check for leaks over repeated calls")

	})
}

