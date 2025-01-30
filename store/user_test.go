package store

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
func BenchmarkNewUserStore(b *testing.B) {
	db := &gorm.DB{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewUserStore(db)
	}
}

func TestNewUserStore(t *testing.T) {
	tests := []struct {
		name     string
		db       *gorm.DB
		wantDB   *gorm.DB
		wantNil  bool
		scenario string
	}{
		{
			name:     "Valid DB Connection",
			db:       &gorm.DB{},
			wantDB:   &gorm.DB{},
			wantNil:  false,
			scenario: "Scenario 1: Create a New UserStore with Valid DB Connection",
		},
		{
			name:     "Nil DB Connection",
			db:       nil,
			wantDB:   nil,
			wantNil:  true,
			scenario: "Scenario 2: Create a New UserStore with Nil DB Connection",
		},
		{
			name:     "Verify DB Field Accessibility",
			db:       &gorm.DB{},
			wantDB:   &gorm.DB{},
			wantNil:  false,
			scenario: "Scenario 3: Verify UserStore DB Field Accessibility",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserStore(tt.db)

			if (got == nil) != tt.wantNil {
				t.Errorf("NewUserStore() returned nil: %v, want nil: %v", got == nil, tt.wantNil)
				return
			}

			if !tt.wantNil && !reflect.DeepEqual(got.db, tt.wantDB) {
				t.Errorf("NewUserStore().db = %v, want %v", got.db, tt.wantDB)
			}
		})
	}
}

func TestNewUserStoreMultipleInstances(t *testing.T) {
	db1 := &gorm.DB{Value: &mockDB{identifier: "db1"}}
	db2 := &gorm.DB{Value: &mockDB{identifier: "db2"}}

	store1 := NewUserStore(db1)
	store2 := NewUserStore(db2)

	if !reflect.DeepEqual(store1.db, db1) {
		t.Errorf("store1.db = %v, want %v", store1.db, db1)
	}

	if !reflect.DeepEqual(store2.db, db2) {
		t.Errorf("store2.db = %v, want %v", store2.db, db2)
	}
}

