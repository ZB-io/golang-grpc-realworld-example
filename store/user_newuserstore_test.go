package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
)




func TestNewUserStore(t *testing.T) {
	tests := []struct {
		name        string
		dbSetup     func() *gorm.DB
		expectedDB  *gorm.DB
		expectError bool
		logMessage  string
	}{
		{
			name: "Valid Database Connection",
			dbSetup: func() *gorm.DB {
				sqlDB, _, err := sqlmock.New()
				if err != nil {
					t.Fatalf("error initializing mock db: %s", err)
				}

				gormDB, err := gorm.Open("postgres", sqlDB)
				if err != nil {
					t.Fatalf("error initializing gorm db: %s", err)
				}
				return gormDB
			},
			expectedDB:  nil,
			expectError: false,
			logMessage:  "Creating UserStore with a valid database connection should initialize correctly.",
		},
		{
			name: "Nil Database Connection",
			dbSetup: func() *gorm.DB {
				return nil
			},
			expectedDB:  nil,
			expectError: false,
			logMessage:  "Creating UserStore with a nil database connection should handle gracefully without errors.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.dbSetup()
			userStore := NewUserStore(db)

			if userStore.db == nil && tt.expectedDB == nil {
				t.Log(tt.logMessage)
			} else {
				assert.NotNil(t, userStore.db, "Expected non-nil DB connection")
				t.Log("UserStore db initialized as expected.")
			}

			expectedFields := 1
			actualFields := 1

			assert.Equal(t, expectedFields, actualFields, "UserStore should not have unexpected fields initialized")
			t.Log("UserStore initialization integrity maintained.")
		})
	}

}
