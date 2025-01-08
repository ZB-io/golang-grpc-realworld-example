package undefined

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
			name: "Create UserStore with valid gorm.DB instance",
			db:   &gorm.DB{},
			want: &UserStore{db: &gorm.DB{}},
		},
		{
			name: "Create UserStore with nil gorm.DB instance",
			db:   nil,
			want: &UserStore{db: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserStore(tt.db)

			if got == nil {
				t.Errorf("NewUserStore() returned nil")
				return
			}

			if reflect.TypeOf(got) != reflect.TypeOf(&UserStore{}) {
				t.Errorf("NewUserStore() returned incorrect type, got %v, want *UserStore", reflect.TypeOf(got))
				return
			}

			if !reflect.DeepEqual(got.db, tt.want.db) {
				t.Errorf("NewUserStore().db = %v, want %v", got.db, tt.want.db)
			}
		})
	}

	t.Run("Multiple calls create distinct UserStore objects", func(t *testing.T) {
		db := &gorm.DB{}
		us1 := NewUserStore(db)
		us2 := NewUserStore(db)

		if us1 == us2 {
			t.Errorf("NewUserStore() returned the same instance for multiple calls")
		}

		if us1.db != us2.db {
			t.Errorf("NewUserStore() created UserStores with different db instances")
		}
	})

	t.Run("NewUserStore doesn't modify gorm.DB instance", func(t *testing.T) {
		db := &gorm.DB{}
		initialDB := *db
		_ = NewUserStore(db)

		if !reflect.DeepEqual(*db, initialDB) {
			t.Errorf("NewUserStore() modified the input gorm.DB instance")
		}
	})
}

