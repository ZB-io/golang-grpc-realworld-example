package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)








/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9

FUNCTION_DEF=func NewUserStore(db *gorm.DB) *UserStore 

 */
func TestNewUserStore(t *testing.T) {

	type testCase struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		validate func(*testing.T, *UserStore)
	}

	setupMockDB := func(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}

		gormDB, err := gorm.Open("mysql", db)
		if err != nil {
			t.Fatalf("Failed to create GORM DB: %v", err)
		}

		return gormDB, mock
	}

	tests := []testCase{
		{
			name: "Successful UserStore creation with valid DB",
			db: func() *gorm.DB {
				db, _ := setupMockDB(t)
				return db
			}(),
			wantNil: false,
			validate: func(t *testing.T, us *UserStore) {
				assert.NotNil(t, us)
				assert.NotNil(t, us.db)
				t.Log("Successfully created UserStore with valid DB connection")
			},
		},
		{
			name:    "UserStore creation with nil DB",
			db:      nil,
			wantNil: false,
			validate: func(t *testing.T, us *UserStore) {
				assert.NotNil(t, us)
				assert.Nil(t, us.db)
				t.Log("Created UserStore with nil DB connection")
			},
		},
		{
			name: "Verify DB reference integrity",
			db: func() *gorm.DB {
				db, _ := setupMockDB(t)
				return db
			}(),
			wantNil: false,
			validate: func(t *testing.T, us *UserStore) {
				assert.NotNil(t, us)
				assert.Equal(t, us.db, us.db)
				t.Log("DB reference integrity verified")
			},
		},
		{
			name: "Multiple UserStore instances independence",
			db: func() *gorm.DB {
				db, _ := setupMockDB(t)
				return db
			}(),
			wantNil: false,
			validate: func(t *testing.T, us *UserStore) {

				db2, _ := setupMockDB(t)
				us2 := NewUserStore(db2)

				assert.NotNil(t, us)
				assert.NotNil(t, us2)
				assert.NotEqual(t, us.db, us2.db)
				t.Log("Multiple UserStore instances maintain independent DB references")
			},
		},
		{
			name: "UserStore with configured DB connection",
			db: func() *gorm.DB {
				db, _ := setupMockDB(t)
				db.LogMode(true)
				return db
			}(),
			wantNil: false,
			validate: func(t *testing.T, us *UserStore) {
				assert.NotNil(t, us)
				assert.True(t, us.db.LogMode())
				t.Log("UserStore preserves DB configuration")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			us := NewUserStore(tc.db)

			if tc.wantNil {
				assert.Nil(t, us)
				return
			}

			tc.validate(t, us)
		})
	}
}

