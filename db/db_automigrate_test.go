package db

import (
	"errors"
	"sync"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/DATA-DOG/go-txdb"
	_ "github.com/mattn/go-sqlite3"
)


type ExpectedExec struct {
	queryBasedExpectation
	result driver.Result
	delay  time.Duration
}

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
func TestAutoMigrate(t *testing.T) {
	mutex.Lock()
	if !txdbInitialized {
		txdb.Register("txdb_memory", "sqlite3", ":memory:")
		txdbInitialized = true
	}
	mutex.Unlock()

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("error opening a mock database connection: %v\n", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("sqlite3", db)
	if err != nil {
		t.Fatalf("failed to connect to the gorm db: %v\n", err)
	}
	defer gormDB.Close()

	tests := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful AutoMigrate with all Models",
			setupMocks: func() {
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name: "Database Error During Migration",
			setupMocks: func() {
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnError(errors.New("DB error"))
			},
			expectedError: errors.New("DB error"),
		},
		{
			name: "Partial Migration Success with Specific Model Failure",
			setupMocks: func() {
				mock.ExpectExec("CREATE TABLE `users`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE `articles`").WillReturnError(errors.New("Schema conflict"))
			},
			expectedError: errors.New("Schema conflict"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := AutoMigrate(gormDB)

			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			}
		})
	}

	t.Run("Concurrent AutoMigrate Calls", func(t *testing.T) {
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(sqlmock.NewResult(1, 1))

		var wg sync.WaitGroup
		const routines = 5
		errorsChan := make(chan error, routines)

		for i := 0; i < routines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := AutoMigrate(gormDB)
				errorsChan <- err
			}()
		}

		wg.Wait()
		close(errorsChan)

		for err := range errorsChan {
			assert.NoError(t, err)
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}
