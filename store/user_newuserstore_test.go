package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestNewUserStore(t *testing.T) {
	t.Run("Scenario 1: Creating a New UserStore with a Valid gorm.DB Instance", func(t *testing.T) {

		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("unexpected error when opening a sqlmock database connection: %s", err)
		}
		defer db.Close()

		gormDB, err := gorm.Open("mysql", db)
		if err != nil {
			t.Fatalf("failed to open gorm DB with sqlmock: %s", err)
		}
		defer gormDB.Close()

		userStore := NewUserStore(gormDB)

		if userStore.db != gormDB {
			t.Errorf("expected UserStore to have db instance %v, got %v", gormDB, userStore.db)
		}
		t.Log("Successfully created UserStore with valid gorm.DB instance")
	})

	t.Run("Scenario 2: Creating a UserStore with a Nil gorm.DB", func(t *testing.T) {

		var db *gorm.DB = nil

		userStore := NewUserStore(db)

		if userStore == nil {
			t.Error("expected a non-nil UserStore instance")
		}
		if userStore.db != nil {
			t.Errorf("expected UserStore db field to be nil, got %v", userStore.db)
		}
		t.Log("Successfully handled nil database instance in UserStore creation")
	})

	t.Run("Scenario 3: Creating Multiple UserStore Instances with Different gorm.DB Instances", func(t *testing.T) {

		db1, _, _ := sqlmock.New()
		defer db1.Close()
		gormDB1, _ := gorm.Open("mysql", db1)
		defer gormDB1.Close()

		db2, _, _ := sqlmock.New()
		defer db2.Close()
		gormDB2, _ := gorm.Open("mysql", db2)
		defer gormDB2.Close()

		userStore1 := NewUserStore(gormDB1)
		userStore2 := NewUserStore(gormDB2)

		if userStore1.db == userStore2.db {
			t.Errorf("expected different db instances, but got the same instance for both UserStores")
		}
		t.Log("Successfully created distinct UserStore instances for different gorm.DB inputs")
	})

	t.Run("Scenario 4: UserStore Initialization Timing", func(t *testing.T) {

		db, _, _ := sqlmock.New()
		defer db.Close()
		gormDB, _ := gorm.Open("mysql", db)
		defer gormDB.Close()

		timedSetup := func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				NewUserStore(gormDB)
			}
		}

		result := testing.Benchmark(timedSetup)
		duration := result.NsPerOp()
		if duration > 1e6 {
			t.Errorf("creation of UserStore took too long: %d ns", duration)
		}
		t.Logf("UserStore initialization performance within the acceptable limit: %d ns", duration)
	})
}
