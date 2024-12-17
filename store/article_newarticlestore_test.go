package store

import (
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

// TestNewArticleStore tests the NewArticleStore function
func TestNewArticleStore(t *testing.T) {
	// define test cases; true indicates that we expect a non-nil *gorm.DB for the returned ArticleStore 
	testCases := []struct {
		name   string
		expect bool
	}{
		{"Test Successful Initialization of Article Store", true},
		{"Test Function Call with Nil DB Reference", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			var mockDB *gorm.DB = nil
			if tc.expect { // Create DB mock only when it is expected
				var mock sqlmock.Sqlmock
				var err error
				mockDB, mock, err = sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was unexpected when opening a stub database connection", err)
				}
				defer mockDB.Close()
				// TODO: add any mock expectations as required here.
			}

			articleStore := NewArticleStore(mockDB) // invocation of the function with mock DB
			assert.Equal(t, tc.expect, articleStore.db != nil, "Expected condition: %v; got: %v", tc.expect, articleStore.db != nil)

			if tc.expect && articleStore.db != nil {
				if err := mock.ExpectationsWereMet(); err != nil { // this checks if all sqlmock calls are expected as per stub
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			}
		})
	}

	// Concurrent ArticleStore creations
	t.Run("Test Concurrent Calls to NewArticleStore", func(t *testing.T) {
		var wg sync.WaitGroup
		var mockDB *gorm.DB
		var mock sqlmock.Sqlmock
		var err error
		
		mockDB, mock, err = sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was unexpected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				articleStore := NewArticleStore(mockDB)
				// assert not nil and same as supplied db
				if articleStore == nil || articleStore.db != mockDB {
					t.Errorf("Mismatch in expected and actual db instances. Expected: %v, Actual: %v", mockDB, articleStore.db)
				}
			}()
		}
		wg.Wait()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
