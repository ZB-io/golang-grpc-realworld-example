package github

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
)









/*
ROOST_METHOD_HASH=NewArticleStore_6be2824012
ROOST_METHOD_SIG_HASH=NewArticleStore_3fe6f79a92

FUNCTION_DEF=func NewArticleStore(db *gorm.DB) *ArticleStore 

 */
func TestArticleStoreMethodAccessibility(t *testing.T) {
	store := NewArticleStore(&gorm.DB{})
	storeType := reflect.TypeOf(store)

	expectedMethods := []string{}

	for _, methodName := range expectedMethods {
		if _, ok := storeType.MethodByName(methodName); !ok {
			t.Errorf("ArticleStore is missing expected method: %s", methodName)
		}
	}
}

func TestNewArticleStore(t *testing.T) {
	tests := []struct {
		name string
		db   *gorm.DB
		want *ArticleStore
	}{
		{
			name: "Create ArticleStore with valid DB",
			db:   &gorm.DB{},
			want: &ArticleStore{db: &gorm.DB{}},
		},
		{
			name: "Create ArticleStore with nil DB",
			db:   nil,
			want: &ArticleStore{db: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewArticleStore(tt.db)
			if got == nil {
				t.Fatal("NewArticleStore returned nil")
			}
			if !reflect.DeepEqual(got.db, tt.want.db) {
				t.Errorf("NewArticleStore().db = %v, want %v", got.db, tt.want.db)
			}
		})
	}
}

func TestNewArticleStoreImmutability(t *testing.T) {
	db := &gorm.DB{}
	store1 := NewArticleStore(db)
	store2 := NewArticleStore(db)

	if store1 == store2 {
		t.Error("NewArticleStore should return different instances")
	}

	if store1.db != store2.db {
		t.Error("ArticleStore instances should share the same db connection")
	}
}

