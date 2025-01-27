package github

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
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
			name: "Create UserStore with valid gorm.DB",
			db:   &gorm.DB{},
			want: &UserStore{db: &gorm.DB{}},
		},
		{
			name: "Create UserStore with nil gorm.DB",
			db:   nil,
			want: &UserStore{db: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserStore(tt.db)

			if got == nil {
				t.Fatal("NewUserStore returned nil")
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserStore() = %v, want %v", got, tt.want)
			}

			if got.db != tt.db {
				t.Errorf("NewUserStore().db = %v, want %v", got.db, tt.db)
			}

			if _, ok := interface{}(got).(*UserStore); !ok {
				t.Errorf("NewUserStore() did not return a *UserStore")
			}
		})
	}
}

