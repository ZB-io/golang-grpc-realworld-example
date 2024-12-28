package db

import (
	"testing"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"sync"
	"os"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestDropTestDB(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *gorm.DB
		validate func(error) error
	}{
		{
			name: "Scenario 1: Successful Closure of an Open Database Connection",
			setup: func() *gorm.DB {
				db, _, err := sqlmock.New()
				if err != nil {
					t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
				}
				gormDB, err := gorm.Open("sqlite3", db)
				if err != nil {
					t.Fatalf("Failed to open gorm db connection: %s", err)
				}
				return gormDB
			},
			validate: func(err error) error {
				if err != nil {
					return fmt.Errorf("expected no error, got %v", err)
				}
				return nil
			},
		},
		{
			name: "Scenario 2: Handling Already Closed Database Connection",
			setup: func() *gorm.DB {
				db, _, err := sqlmock.New()
				if err != nil {
					t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
				}
				gormDB, err := gorm.Open("sqlite3", db)
				if err != nil {
					t.Fatalf("Failed to open gorm db connection: %s", err)
				}
				_ = gormDB.Close()
				return gormDB
			},
			validate: func(err error) error {
				if err != nil {
					return fmt.Errorf("expected no error when trying to close already closed connection, got %v", err)
				}
				return nil
			},
		},
		{
			name: "Scenario 3: Error Handling with Null Database Input",
			setup: func() *gorm.DB {
				return nil
			},
			validate: func(err error) error {

				if err != nil {
					return fmt.Errorf("expected no error with nil input, got %v", err)
				}
				return nil
			},
		},
		{
			name: "Scenario 4: Handling of Concurrent Close Operations",
			setup: func() *gorm.DB {
				db, _, err := sqlmock.New()
				if err != nil {
					t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
				}
				gormDB, err := gorm.Open("sqlite3", db)
				if err != nil {
					t.Fatalf("Failed to open gorm db connection: %s", err)
				}
				return gormDB
			},
			validate: func(err error) error {
				if err != nil {
					return fmt.Errorf("expected no error during concurrent closures, got %v", err)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Starting test case: %s", tt.name)

			db := tt.setup()
			errorChan := make(chan error)
			defer close(errorChan)

			if tt.name == "Scenario 4: Handling of Concurrent Close Operations" {
				var wg sync.WaitGroup
				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						errorChan <- DropTestDB(db)
					}()
				}
				wg.Wait()

				select {
				case err := <-errorChan:
					if validateErr := tt.validate(err); validateErr != nil {
						t.Error(validateErr)
					}
				default:
					t.Log("No errors from concurrent operations.")
				}
			} else {

				stdout := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

				err := DropTestDB(db)

				w.Close()
				out, _ := ioutil.ReadAll(r)
				os.Stdout = stdout

				t.Logf("Captured output: %s", out)

				if validateErr := tt.validate(err); validateErr != nil {
					t.Error(validateErr)
				}
			}

			t.Logf("Finished test case: %s", tt.name)
		})
	}
}
