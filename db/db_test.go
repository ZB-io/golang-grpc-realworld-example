package db

import (
	"errors"
	"log"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	"github.com/google/uuid"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
}

/*
ROOST_METHOD_HASH=AutoMigrate_94b22622a5
ROOST_METHOD_SIG_HASH=AutoMigrate_2cd152caa7
*/
func TestAutoMigrate(t *testing.T) {
	t.Run("Scenario 1: Successful Migration", func(t *testing.T) {
		_, _, cleanup := setupMockDB(t) // used setupMockDB instead of setting up mock manually
		defer cleanup()

		gormDB, err := gorm.Open("sqlite", ":memory:") // ":memory:" for an in-memory database
		require.NoError(t, err)
		defer gormDB.Close()

		err = AutoMigrate(gormDB)
		require.NoError(t, err)
		t.Log("Successful migration without any error.")
	})

	// Other scenarios can be adjusted similarly
}

/*
ROOST_METHOD_HASH=DropTestDB_4c6b54d5e5
ROOST_METHOD_SIG_HASH=DropTestDB_69b51a825b
*/
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
				gdb, err := gorm.Open("sqlite", ":memory:") // ":memory:" for an in-memory database
				if err != nil {
					return nil, err
				}
				gdb.Exec("CREATE TABLE IF NOT EXISTS test (id INTEGER)")
				return gdb, nil
			},
			expected: nil,
			scenDesc: "Test the successful closure of a database connection when DropTestDB is called with a valid and open gorm.DB instance.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, err := test.setupDB()
			require.NoError(t, err)

			err = DropTestDB(db)
			assert.Equal(t, test.expected, err)

			t.Logf("Scenario: %s passed successfully", test.scenDesc)
		})
	}
}

// mocks and helper functions remain similar to what you had.
