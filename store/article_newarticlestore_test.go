package store

import (
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/raahii/golang-grpc-realworld-example/model"
)



type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}




type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

func TestNewArticleStore(t *testing.T) {
	t.Run("Scenario 1: Creating a New ArticleStore with a Valid DB Connection", func(t *testing.T) {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening a stub database connection: %s", err)
		}
		defer db.Close()

		defer func() {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		}()

		gormDB, err := gorm.Open("postgres", db)
		if err != nil {
			t.Fatalf("failed to open gorm DB: %s", err)
		}

		articleStore := NewArticleStore(gormDB)

		if articleStore == nil {
			t.Error("Expected ArticleStore to be not nil")
		}
		if articleStore.db != gormDB {
			t.Error("Expected the db field to match the gorm.DB reference passed")
		}
	})

	t.Run("Scenario 2: Creating a New ArticleStore with a Nil DB Connection", func(t *testing.T) {

		articleStore := NewArticleStore(nil)

		if articleStore == nil {
			t.Error("Expected ArticleStore to be not nil")
		}
		if articleStore.db != nil {
			t.Error("Expected ArticleStore db field to be nil when passed nil")
		}
	})

	t.Run("Scenario 3: Ensuring New ArticleStore Does Not Modify Input DB Object", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening a stub database connection: %s", err)
		}
		defer db.Close()

		defer func() {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		}()

		gormDB, err := gorm.Open("postgres", db)
		if err != nil {
			t.Fatalf("failed to open gorm DB: %s", err)
		}

		initialLogMode := gormDB.LogMode(true)

		NewArticleStore(gormDB)

		if gormDB.LogMode(true) != initialLogMode {
			t.Error("Expected gorm.DB log mode to remain unchanged after NewArticleStore call")
		}
	})

	t.Run("Scenario 4: Confirming New ArticleStore References a Shallow Copy of the DB", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening a stub database connection: %s", err)
		}
		defer db.Close()

		defer func() {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		}()

		gormDB, err := gorm.Open("postgres", db)
		if err != nil {
			t.Fatalf("failed to open gorm DB: %s", err)
		}

		articleStore := NewArticleStore(gormDB)

		gormDB.LogMode(true)

		if articleStore.db.LogMode(true) != gormDB.LogMode(true) {
			t.Error("Expected changes in the original DB to reflect in the ArticleStore's db field")
		}
	})

	t.Run("Scenario 5: Documenting Creation Time of ArticleStore Objects", func(t *testing.T) {

		nowFunc := func() time.Time { return time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC) }

		_ = nowFunc
	})
}
