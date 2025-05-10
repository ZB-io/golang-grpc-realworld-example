package db

import (
	"database/sql"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var (
	gormOpen = func(dialect string, args ...interface{}) (*gorm.DB, error) {
		return gorm.Open(dialect, args...)
	}

	dsn = func() (string, error) {
		return "", errors.New("Not implemented in test")
	}
)

func TestNew(t *testing.T) {
	t.Run("Successful Connection to the Database", func(t *testing.T) {
		// Mocking gorm.Open to return a valid DB instance
		originalOpen := gormOpen
		defer func() { gormOpen = originalOpen }()
		gormOpen = func(dialect string, args ...interface{}) (*gorm.DB, error) {
			db, mockDB, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error on mockDB: %v", err)
			}
			mockDB.ExpectPing()

			return &gorm.DB{DB: db}, nil
		}

		// Act: Call New
		db, err := New()

		// Assert
		assert.NotNil(t, db, "DB should not be nil on successful connection")
		assert.Nil(t, err, "Error should be nil on successful connection")
	})

	t.Run("DSN Failure Leading to Error", func(t *testing.T) {
		// Mock dsn to return an error
		originalDsn := dsn
		defer func() { dsn = originalDsn }()
		dsn = func() (string, error) {
			return "", errors.New("mock DSN error")
		}

		// Act
		db, err := New()

		// Assert
		assert.Nil(t, db, "DB should be nil when DSN retrieval fails")
		assert.EqualError(t, err, "mock DSN error", "Expected specific DSN error message")
	})

	t.Run("Maximum Retry Attempt Exhaustion", func(t *testing.T) {
		// Mocking gorm.Open to fail continuously
		originalOpen := gormOpen
		defer func() { gormOpen = originalOpen }()
		gormOpen = func(dialect string, args ...interface{}) (*gorm.DB, error) {
			return nil, errors.New("mock open error")
		}

		// Act
		db, err := New()

		// Assert
		assert.Nil(t, db, "DB should be nil after max retry attempts exhaustion")
		assert.EqualError(t, err, "mock open error", "Expected specific open error after retries")
	})

	t.Run("SetMaxIdleConns Configuration", func(t *testing.T) {
		// Mocking gorm.Open to return a valid DB instance
		originalOpen := gormOpen
		defer func() { gormOpen = originalOpen }()
		gormOpen = func(dialect string, args ...interface{}) (*gorm.DB, error) {
			db, mockDB, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error on mockDB: %v", err)
			}
			mockDB.ExpectPing()

			return &gorm.DB{DB: db}, nil
		}

		// Act
		db, err := New()

		// Assert
		if assert.NoError(t, err) && assert.NotNil(t, db) {
			assert.Equal(t, 3, db.DB().Stats().Idle, "MaxIdleConns should be 3")
		}
	})

	t.Run("LogMode Configuration is Disabled", func(t *testing.T) {
		// Mocking gorm.Open to return a valid DB instance
		originalOpen := gormOpen
		defer func() { gormOpen = originalOpen }()
		gormOpen = func(dialect string, args ...interface{}) (*gorm.DB, error) {
			db, mockDB, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error on mockDB: %v", err)
			}
			mockDB.ExpectPing()

			return &gorm.DB{DB: db}, nil
		}

		// Act
		db, err := New()

		// Assert
		if assert.NoError(t, err) && assert.NotNil(t, db) {
			assert.Equal(t, false, db.LogMode(), "LogMode should be disabled (false)")
		}
	})

	t.Run("Graceful Failure on Incorrect Dialect", func(t *testing.T) {
		// Mocking gorm.Open with an incorrect dialect
		originalOpen := gormOpen
		defer func() { gormOpen = originalOpen }()
		gormOpen = func(dialect string, args ...interface{}) (*gorm.DB, error) {
			return nil, errors.New("mock incorrect dialect error")
		}

		// Act
		db, err := New()

		// Assert
		assert.Nil(t, db, "DB should be nil with incorrect dialect")
		assert.EqualError(t, err, "mock incorrect dialect error", "Expected incorrect dialect error")
	})

	t.Run("Test Concurrency Safety on Sync Ops", func(t *testing.T) {
		// Mocking gorm.Open to return a valid DB instance
		originalOpen := gormOpen
		defer func() { gormOpen = originalOpen }()
		gormOpen = func(dialect string, args ...interface{}) (*gorm.DB, error) {
			db, mockDB, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error on mockDB: %v", err)
			}
			mockDB.ExpectPing()

			return &gorm.DB{DB: db}, nil
		}

		concurrentTests := 10
		errChan := make(chan error, concurrentTests)
		var wg sync.WaitGroup
		wg.Add(concurrentTests)

		// Run concurrent goroutines
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
