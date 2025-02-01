package store

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)








/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9

FUNCTION_DEF=func NewUserStore(db *gorm.DB) *UserStore 

*/
func TestNewUserStore(t *testing.T) {
	tests := []struct {
		name string
		db   *gorm.DB
		want *UserStore
	}{
		{
			name: "Create UserStore with valid DB",
			db:   &gorm.DB{},
			want: &UserStore{db: &gorm.DB{}},
		},
		{
			name: "Create UserStore with nil DB",
			db:   nil,
			want: &UserStore{db: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserStore(tt.db)
			if got == nil {
				t.Error("NewUserStore() returned nil")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUserStore_DBReferenceIntegrity(t *testing.T) {
	db := &gorm.DB{}
	store := NewUserStore(db)

	if store.db != db {
		t.Error("NewUserStore() did not maintain DB reference integrity")
	}
}

func TestNewUserStore_MultipleDatabases(t *testing.T) {
	db1 := &gorm.DB{}
	db2 := &gorm.DB{}

	store1 := NewUserStore(db1)
	store2 := NewUserStore(db2)

	if store1 == store2 {
		t.Error("NewUserStore() returned the same instance for different DB connections")
	}

	if store1.db == store2.db {
		t.Error("NewUserStore() used the same DB reference for different DB connections")
	}
}

func TestNewUserStore_MultipleInstances(t *testing.T) {
	db := &gorm.DB{}
	store1 := NewUserStore(db)
	store2 := NewUserStore(db)

	if store1 == store2 {
		t.Error("NewUserStore() returned the same instance for multiple calls")
	}

	if store1.db != store2.db {
		t.Error("NewUserStore() did not use the same DB reference for multiple calls")
	}
}


/*
ROOST_METHOD_HASH=GetFollowingUserIDs_ba3670aa2c
ROOST_METHOD_SIG_HASH=GetFollowingUserIDs_55ccc2afd7

FUNCTION_DEF=func (s *UserStore) GetFollowingUserIDs(m *model.User) ([]uint, error) 

*/
func TestUserStoreGetFollowingUserIDs(t *testing.T) {
	tests := []struct {
		name          string
		user          *model.User
		mockData      []uint
		mockDBError   error
		mockScanError error
		expectedIDs   []uint
		expectedError error
	}{
		{
			name:          "Successful retrieval of following user IDs",
			user:          &model.User{Model: gorm.Model{ID: 1}},
			mockData:      []uint{2, 3, 4},
			expectedIDs:   []uint{2, 3, 4},
			expectedError: nil,
		},
		{
			name:          "User with no followers",
			user:          &model.User{Model: gorm.Model{ID: 1}},
			mockData:      []uint{},
			expectedIDs:   []uint{},
			expectedError: nil,
		},
		{
			name:          "Database error",
			user:          &model.User{Model: gorm.Model{ID: 1}},
			mockDBError:   errors.New("database error"),
			expectedIDs:   []uint{},
			expectedError: errors.New("database error"),
		},
		{
			name:          "Large number of followers",
			user:          &model.User{Model: gorm.Model{ID: 1}},
			mockData:      []uint{2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			expectedIDs:   []uint{2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			expectedError: nil,
		},
		{
			name:          "Invalid user ID",
			user:          &model.User{Model: gorm.Model{ID: 999}},
			mockData:      []uint{},
			expectedIDs:   []uint{},
			expectedError: nil,
		},
		{
			name:          "Scan error",
			user:          &model.User{Model: gorm.Model{ID: 1}},
			mockData:      []uint{2, 3, 4},
			mockScanError: errors.New("scan error"),
			expectedIDs:   []uint{},
			expectedError: errors.New("scan error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRows := &mockRows{data: tt.mockData}
			if tt.mockScanError != nil {
				mockRows.scanFunc = func(dest ...interface{}) error {
					return tt.mockScanError
				}
			}

			mockDB := &mockDB{
				rows:  mockRows,
				error: tt.mockDBError,
			}

			store := &UserStore{db: &gorm.DB{Value: mockDB}}

			ids, err := store.GetFollowingUserIDs(tt.user)

			assert.Equal(t, tt.expectedIDs, ids)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

