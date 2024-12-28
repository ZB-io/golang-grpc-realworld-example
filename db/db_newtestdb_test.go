package db

import (
	"database/sql"
	"errors"
	"os"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestNewTestDB(t *testing.T) {
	tests := []struct {
		name           string
		prepareEnv     func()
		prepareMockDSN func(sqlmock.Sqlmock)
		expectedError  error
	}{
		{
			name: "Successful Database Initialization",
			prepareEnv: func() {
				_ = os.Setenv("DB_HOST", "localhost")
				_ = os.Setenv("DB_USER", "user")
				_ = os.Setenv("DB_PASSWORD", "password")
				_ = os.Setenv("DB_NAME", "dbname")
				_ = os.Setenv("DB_PORT", "3306")
			},
			prepareMockDSN: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT 1").WillReturnRows(sqlmock.NewRows([]string{"1"}))
			},
			expectedError: nil,
		},
		{
			name: "Error Loading Environment File",
			prepareEnv: func() {
				_ = os.Unsetenv("DB_HOST")
				_ = os.Unsetenv("DB_USER")
				_ = os.Unsetenv("DB_PASSWORD")
				_ = os.Unsetenv("DB_NAME")
				_ = os.Unsetenv("DB_PORT")
			},
			prepareMockDSN: func(_ sqlmock.Sqlmock) {},
			expectedError:  errors.New("$DB_HOST is not set"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareEnv()
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func() {
				if err := db.Close(); err != nil {
					t.Errorf("error closing DB: %v", err)
				}
			}()

			tt.prepareMockDSN(mock)

			actualDB, actualErr := NewTestDB()
			if tt.expectedError != nil {
				assert.Nil(t, actualDB)
				assert.EqualError(t, actualErr, tt.expectedError.Error())
			} else {
				assert.NotNil(t, actualDB)
				assert.NoError(t, actualErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unmet expectations: %s", err)
			}
		})
	}
}
