package db

import (
	"database/sql"
	"errors"
	"testing"
	"time"
	"sync"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var dsn = func() (string, error) {
	return "original-dsn", nil
}

var originalDSNFunc = dsn




func TestNew(t *testing.T) {
	defer resetDSN()

	tests := []struct {
		name             string
		mockDSNSuccess   bool
		mockDSNError     error
		mockOpenFailures int
		expectedDBNotNil bool
		expectedErrorNil bool
	}{
		{
			name:             "Successful Connection with Default Configuration",
			mockDSNSuccess:   true,
			mockDSNError:     nil,
			mockOpenFailures: 0,
			expectedDBNotNil: true,
			expectedErrorNil: true,
		},
		{
			name:             "Connection Retry Logic with Intermediate Failures",
			mockDSNSuccess:   true,
			mockDSNError:     nil,
			mockOpenFailures: 5,
			expectedDBNotNil: true,
			expectedErrorNil: true,
		},
		{
			name:             "Exhausted Retry Attempts with Constant Failures",
			mockDSNSuccess:   true,
			mockDSNError:     nil,
			mockOpenFailures: 10,
			expectedDBNotNil: false,
			expectedErrorNil: false,
		},
		{
			name:             "Error in Fetching DSN",
			mockDSNSuccess:   false,
			mockDSNError:     errors.New("cannot generate DSN"),
			mockOpenFailures: 0,
			expectedDBNotNil: false,
			expectedErrorNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDSN(tt.mockDSNSuccess, tt.mockDSNError)

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			originalOpen := gorm.Open
			defer func() { gorm.Open = originalOpen }()
			gorm.Open = func(dialect string, settings ...interface{}) (*gorm.DB, error) {
				if tt.mockOpenFailures > 0 {
					tt.mockOpenFailures--
					return nil, errors.New("connection failed")
				}
				return &gorm.DB{DB: db}, nil
			}

			gormDB, err := New()

			if tt.expectedDBNotNil {
				assert.NotNil(t, gormDB, "expected a non-nil database connection")
			} else {
				assert.Nil(t, gormDB, "expected a nil database connection")
			}

			if tt.expectedErrorNil {
				assert.Nil(t, err, "expected no error during database connection")
			} else {
				assert.NotNil(t, err, "expected an error during database connection")
			}

			t.Logf("Test '%s' completed: DB Not Nil: %v, Error Nil: %v", tt.name, tt.expectedDBNotNil, tt.expectedErrorNil)
		})
	}
}

func mockDSN(success bool, err error) {
	if success {
		dsn = func() (string, error) {
			return "valid-dsn", nil
		}
	} else {
		dsn = func() (string, error) {
			return "", err
		}
	}
}
func resetDSN() {
	dsn = originalDSNFunc
}

