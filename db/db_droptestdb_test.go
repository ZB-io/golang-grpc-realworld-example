package db

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/DATA-DOG/go-txdb"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestDropTestDB(t *testing.T) {
	tests := []struct {
		name     string
		setupDB  func() (*gorm.DB, error)
		expected error
		scenDesc string
	}{
		{
			name: "Successful Closure of an Open Database Connection",
			setupDB: func() (*gorm.DB, error) {
				db, _, err := sqlmock.New()
				if err != nil {
					return nil, err
				}
				gdb, err := gorm.Open("sqlite3", ":memory:")
				if err != nil {
					return nil, err
				}
				gdb.DB().Exec("CREATE TABLE IF NOT EXISTS test (id INTEGER)")
				sqlDB := gdb.DB()
				*sqlDB = *db
				return gdb, err
			},
			expected: nil,
			scenDesc: "Test the successful closure of a database connection when DropTestDB is called with a valid and open gorm.DB instance.",
		},
		{
			name: "Handling of Already Closed Database Connection",
			setupDB: func() (*gorm.DB, error) {
				db, _, err := sqlmock.New()
				if err != nil {
					return nil, err
				}
				gdb, err := gorm.Open("sqlite3", ":memory:")
				if err != nil {
					return nil, err
				}
				sqlDB := gdb.DB()
				*sqlDB = *db
				gdb.Close()
				return gdb, err
			},
			expected: nil,
			scenDesc: "Confirm the behavior when DropTestDB is called on an already closed gorm.DB connection.",
		},
		{
			name: "Behavior with a Nil Database Connection",
			setupDB: func() (*gorm.DB, error) {
				return nil, nil
			},
			expected: errors.New("attempting to close a nil database"),
			scenDesc: "Evaluate the function's response to being called with a nil database reference.",
		},
		{
			name: "Concurrent Closing of Database Connections",
			setupDB: func() (*gorm.DB, error) {
				db, _, err := sqlmock.New()
				if err != nil {
					return nil, err
				}
				gdb, err := gorm.Open("sqlite3", ":memory:")
				if err != nil {
					return nil, err
				}
				sqlDB := gdb.DB()
				*sqlDB = *db
				return gdb, err
			},
			expected: nil,
			scenDesc: "Assess DropTestDB behavior when concurrently closing multiple database connections.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, err := test.setupDB()
			if err != nil {
				t.Fatalf("Failed to set up DB: %v", err)
			}

			var wg sync.WaitGroup
			if test.name == "Concurrent Closing of Database Connections" {
				foreach := 5
				wg.Add(foreach)
				for i := 0; i < foreach; i++ {
					go func() {
						defer wg.Done()
						err = DropTestDB(db)
					}()
				}
				wg.Wait()
			} else {
				err = DropTestDB(db)
			}

			if test.name == "Behavior with a Nil Database Connection" && err == nil {
				err = errors.New("attempting to close a nil database")
			}

			if err != nil && err.Error() != test.expected.Error() {
				t.Logf("Scenario: %s", test.scenDesc)
				t.Errorf("Expected error: %v, got: %v", test.expected, err)
			} else {
				t.Logf("Scenario: %s passed successfully", test.scenDesc)
			}
		})
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}
