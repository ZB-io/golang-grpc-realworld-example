package db

import (
	"errors"
	"sync"
	"testing"

	"github.com/jinzhu/gorm"
)

// Mock implementation of gorm.DB for testing purposes
type mockDB struct {
	closeCallCount int
	closeError     error
	mu             sync.Mutex
}

func (m *mockDB) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closeCallCount++
	return m.closeError
}

func TestDropTestDb(t *testing.T) {
	tests := []struct {
		name           string
		db             *gorm.DB
		expectedError  error
		expectedClosed bool
	}{
		{
			name:           "Successfully close the database connection",
			db:             &gorm.DB{Value: &mockDB{}},
			expectedError:  nil,
			expectedClosed: true,
		},
		{
			name:           "Handle nil database object",
			db:             nil,
			expectedError:  nil,
			expectedClosed: false,
		},
		{
			name:           "Handle database close error",
			db:             &gorm.DB{Value: &mockDB{closeError: errors.New("close error")}},
			expectedError:  nil,
			expectedClosed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DropTestDB(tt.db)

			if err != tt.expectedError {
				t.Errorf("DropTestDB() error = %v, expectedError %v", err, tt.expectedError)
			}

			if tt.db != nil {
				mockDB := tt.db.Value.(*mockDB)
				if tt.expectedClosed && mockDB.closeCallCount != 1 {
					t.Errorf("Expected database to be closed, but Close() was called %d times", mockDB.closeCallCount)
				}
				if !tt.expectedClosed && mockDB.closeCallCount != 0 {
					t.Errorf("Expected database not to be closed, but Close() was called %d times", mockDB.closeCallCount)
				}
			}
		})
	}
}

func TestDropTestDbIdempotency(t *testing.T) {
	mockDB := &mockDB{}
	db := &gorm.DB{Value: mockDB}

	// Call DropTestDB twice
	err1 := DropTestDB(db)
	err2 := DropTestDB(db)

	if err1 != nil || err2 != nil {
		t.Errorf("Expected both calls to return nil, got err1: %v, err2: %v", err1, err2)
	}

	if mockDB.closeCallCount != 1 {
		t.Errorf("Expected Close() to be called only once, but was called %d times", mockDB.closeCallCount)
	}
}

func TestDropTestDbConcurrency(t *testing.T) {
	const numGoroutines = 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			mockDB := &mockDB{}
			db := &gorm.DB{Value: mockDB}
			err := DropTestDB(db)
			if err != nil {
				t.Errorf("Concurrent DropTestDB() returned error: %v", err)
			}
			if mockDB.closeCallCount != 1 {
				t.Errorf("Expected Close() to be called once, but was called %d times", mockDB.closeCallCount)
			}
		}()
	}

	wg.Wait()
}
