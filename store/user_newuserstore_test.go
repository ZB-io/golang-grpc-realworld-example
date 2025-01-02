package store

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"testing"
)


func TestNewUserStore(t *testing.T) {

	testCases := []struct {
		name     string
		db       *gorm.DB
		expected *UserStore
	}{
		{
			name: "Successful UserStore Creation",
			db:   newMockDB(),
			expected: &UserStore{
				db: newMockDB(),
			},
		},
		{
			name:     "Nil DB Parameter",
			db:       nil,
			expected: &UserStore{db: nil},
		},
		{
			name: "Multiple UserStore Creation",
			db:   newMockDB(),
			expected: &UserStore{
				db: newMockDB(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewUserStore(tc.db)
			if result.db != tc.expected.db {
				t.Errorf("NewUserStore() got = %v, want = %v", result, tc.expected)
			} else {
				t.Logf("NewUserStore() passed for scenario: %s", tc.name)
			}
		})
	}
}
func newMockDB() *gorm.DB {
	db, _, err := sqlmock.New()
	if err != nil {
		return nil
	}
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		return nil
	}
	return gormDB
}
