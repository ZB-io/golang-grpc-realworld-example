package db

import (
	"errors"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestNew(t *testing.T) {
	t.Run("Successful Connection to the Database", func(t *testing.T) {

		originalOpen := gormOpen
		defer func() { gormOpen = originalOpen }()
		gormOpen = func(dialect string, args ...interface{}) (*gorm.DB, error) {
			db, mockDB, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error on mockDB: %v", err)
			}
			defer db.Close()

			mockDB.ExpectPing()

			return &gorm.DB{DB: db}, nil
		}

		db, err := New()

		assert.NotNil(t, db, "DB should not be nil on successful connection")
		assert.Nil(t, err, "Error should be nil on successful connection")
	})

	t.Run("DSN Failure Leading to Error", func(t *testing.T) {

		originalDsn := dsn
		defer func() { dsn = originalDsn }()
		dsn = func() (string, error) {
			return "", errors.New("mock DSN error")
		}

		db, err := New()

		assert.Nil(t, db, "DB should be nil when DSN retrieval fails")
		assert.EqualError(t, err, "mock DSN error", "Expected specific DSN error message")
	})

	t.Run("Maximum Retry Attempt Exhaustion", func(t *testing.T) {

		originalOpen := gormOpen
		defer func() { gormOpen = originalOpen }()
		gormOpen = func(dialect string, args ...interface{}) (*gorm.DB, error) {
			return nil, errors.New("mock open error")
		}

		db, err := New()

		assert.Nil(t, db, "DB should be nil after max retry attempts exhaustion")
		assert.EqualError(t, err, "mock open error", "Expected specific open error after retries")
	})

	t.Run("SetMaxIdleConns Configuration", func(t *testing.T) {

		originalOpen := gormOpen
		defer func() { gormOpen = originalOpen }()
		gormOpen = func(dialect string, args ...interface{}) (*gorm.DB, error) {
			db, mockDB, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error on mockDB: %v", err)
			}
			defer db.Close()

			mockDB.ExpectPing()

			return &gorm.DB{DB: db}, nil
		}

		db, err := New()

		if assert.NoError(t, err) && assert.NotNil(t, db) {
			assert.Equal(t, 3, db.DB().Stats().Idle, "MaxIdleConns should be 3")
		}
	})

	t.Run("LogMode Configuration is Disabled", func(t *testing.T) {

		originalOpen := gormOpen
		defer func() { gormOpen = originalOpen }()
		gormOpen = func(dialect string, args ...interface{}) (*gorm.DB, error) {
			db, mockDB, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error on mockDB: %v", err)
			}
			defer db.Close()

			mockDB.ExpectPing()

			return &gorm.DB{DB: db}, nil
		}

		db, err := New()

		if assert.NoError(t, err) && assert.NotNil(t, db) {
			assert.Equal(t, false, db.LogMode(), "LogMode should be disabled (false)")
		}
	})

	t.Run("Graceful Failure on Incorrect Dialect", func(t *testing.T) {

		originalOpen := gormOpen
		defer func() { gormOpen = originalOpen }()
		gormOpen = func(dialect string, args ...interface{}) (*gorm.DB, error) {
			return nil, errors.New("mock incorrect dialect error")
		}

		db, err := New()

		assert.Nil(t, db, "DB should be nil with incorrect dialect")
		assert.EqualError(t, err, "mock incorrect dialect error", "Expected incorrect dialect error")
	})

	t.Run("Test Concurrency Safety on Sync Ops", func(t *testing.T) {

		originalOpen := gormOpen
		defer func() { gormOpen = originalOpen }()
		gormOpen = func(dialect string, args ...interface{}) (*gorm.DB, error) {
			db, mockDB, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error on mockDB: %v", err)
			}
			defer db.Close()

			mockDB.ExpectPing()

			return &gorm.DB{DB: db}, nil
		}
		concurrentTests := 10
		errChan := make(chan error, concurrentTests)
		var wg sync.WaitGroup
		wg.Add(concurrentTests)

		for i := 0; i < concurrentTests; i++ {
			go func() {
				defer wg.Done()
				if _, err := New(); err != nil {
					errChan <- err
				}
			}()
		}
		wg.Wait()
		close(errChan)

		for err := range errChan {
			assert.NoError(t, err, "Unexpected error during concurrent New() execution")
		}
	})
}
