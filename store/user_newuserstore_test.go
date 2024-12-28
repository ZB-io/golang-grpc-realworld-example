package store

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func TestNewUserStore(t *testing.T) {
	t.Run("Scenario 1: Creating a New UserStore with a Valid gorm.DB Instance", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New() // We create a SQL mock DB
		if err != nil {
			t.Fatalf("unexpected error when opening a sqlmock database connection: %s", err)
		}
		defer db.Close()

		gormDB, err := gorm.Open("mysql", db)
		if err != nil {
			t.Fatalf("failed to open gorm DB with sqlmock: %s", err)
		}
		defer gormDB.Close()

		// Act
		userStore := NewUserStore(gormDB)

		// Assert
		if userStore.db != gormDB {
			t.Errorf("expected UserStore to have db instance %v, got %v", gormDB, userStore.db)
		}
		t.Log("Successfully created UserStore with valid gorm.DB instance")
	})

	t.Run("Scenario 2: Creating a UserStore with a Nil gorm.DB", func(t *testing.T) {
		// Arrange
		var db *gorm.DB = nil

		// Act
		userStore := NewUserStore(db)

		// Assert
		if userStore == nil {
			t.Error("expected a non-nil UserStore instance")
		}
		if userStore.db != nil {
			t.Errorf("expected UserStore db field to be nil, got %v", userStore.db)
		}
		t.Log("Successfully handled nil database instance in UserStore creation")
	})

	t.Run("Scenario 3: Creating Multiple UserStore Instances with Different gorm.DB Instances", func(t *testing.T) {
		// Arrange
		db1, mock1, _ := sqlmock.New()
		defer db1.Close()
		gormDB1, _ := gorm.Open("mysql", db1)
		defer gormDB1.Close()

		db2, mock2, _ := sqlmock.New()
		defer db2.Close()
		gormDB2, _ := gorm.Open("mysql", db2)
		defer gormDB2.Close()

		// Act
		userStore1 := NewUserStore(gormDB1)
		userStore2 := NewUserStore(gormDB2)

		// Assert
		if userStore1.db == userStore2.db {
			t.Errorf("expected different db instances, but got the same instance for both UserStores")
		}
		t.Log("Successfully created distinct UserStore instances for different gorm.DB inputs")
	})

	t.Run("Scenario 4: UserStore Initialization Timing", func(t *testing.T) {
		// Arrange
		db, mock, _ := sqlmock.New()
		defer db.Close()
		gormDB, _ := gorm.Open("mysql", db)
		defer gormDB.Close()

		// Act
		timedSetup := func() {
			NewUserStore(gormDB)
		}

		// Assert
		duration := testing.Benchmark(timedSetup).NsPerOp()
		if duration > 1e6 {
			t.Errorf("creation of UserStore took too long: %d ns", duration)
		}
		t.Logf("UserStore initialization performance within the acceptable limit: %d ns", duration)
	})

	// Clean up expectations to ensure all mocks are used.
	// TODO: User modification might be required to ensure mock expectations are correctly handled.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
