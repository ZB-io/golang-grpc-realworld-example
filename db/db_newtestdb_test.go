package db

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"sync"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/DATA-DOG/go-txdb"
)


func TestNewTestDB(t *testing.T) {

	mutex = sync.Mutex{}
	txdbInitialized = false

	type test struct {
		name          string
		prepare       func()
		expectedError string
		expectDB      bool
	}

	tests := []test{
		{
			name: "Scenario 1: Successful Database Initialization",
			prepare: func() {
				os.Setenv("DSN", "user:password@tcp(localhost:3306)/dbname?charset=utf8&parseTime=True&loc=Local")

			},
			expectDB:      true,
			expectedError: "",
		},
		{
			name: "Scenario 2: Failure to Load Environment File",
			prepare: func() {
				os.Unsetenv("TEST_ENV_PATH")

			},
			expectDB:      false,
			expectedError: "open ../env/test.env: no such file or directory",
		},
		{
			name: "Scenario 3: Failure During DSN Generation",
			prepare: func() {

				dsn = func() (string, error) {
					return "", errors.New("failed to generate DSN")
				}
			},
			expectDB:      false,
			expectedError: "failed to generate DSN",
		},
		{
			name: "Scenario 4: Error Registering Transactional Database",
			prepare: func() {

				txdb.Register = func(tag, dsn, driverName string) error {
					return errors.New("failed to register txdb")
				}
			},
			expectDB:      false,
			expectedError: "failed to register txdb",
		},
		{
			name: "Scenario 5: Error Opening SQL Connection",
			prepare: func() {

				sqlOpenMock := func(string, ...interface{}) (*sql.DB, error) {
					return nil, errors.New("failed to open SQL connection")
				}
				sql.Open = sqlOpenMock
			},
			expectDB:      false,
			expectedError: "failed to open SQL connection",
		},
		{
			name: "Scenario 6: Successfully Utilize Mutex Locking Mechanism",
			prepare: func() {
				os.Setenv("DSN", "user:password@tcp(localhost:3306)/dbname?charset=utf8&parseTime=True&loc=Local")

			},
			expectDB:      true,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()

			db, err := NewTestDB()

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				log.Println("Test case failed as expected:", err)
			} else {
				require.NoError(t, err)
				if tt.expectDB {
					assert.NotNil(t, db)
					err := db.Close()
					if err != nil {
						t.Fatalf("failed to close the database: %v", err)
					}
					log.Println("Database initialized successfully")
				}
			}
		})
	}
}




