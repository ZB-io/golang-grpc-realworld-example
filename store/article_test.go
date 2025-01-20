package store

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
func TestNewArticleStore(t *testing.T) {
	tests := []struct {
		name string
		db   *gorm.DB
		want *ArticleStore
	}{
		{
			name: "Create ArticleStore with Valid Database Connection",
			db:   &gorm.DB{},
			want: &ArticleStore{db: &gorm.DB{}},
		},
		{
			name: "Create ArticleStore with Nil Database Connection",
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

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArticleStore() = %v, want %v", got, tt.want)
			}

			if got.db != tt.db {
				t.Errorf("NewArticleStore().db = %v, want %v", got.db, tt.db)
			}

			if reflect.TypeOf(got) != reflect.TypeOf(&ArticleStore{}) {
				t.Errorf("NewArticleStore() returned incorrect type: got %T, want *ArticleStore", got)
			}
		})
	}
}

func TestNewArticleStoreMultipleInstances(t *testing.T) {
	db1 := &gorm.DB{}
	db2 := &gorm.DB{}

	store1 := NewArticleStore(db1)
	store2 := NewArticleStore(db2)

	if store1 == store2 {
		t.Error("NewArticleStore created identical instances for different DB connections")
	}

	if store1.db != db1 {
		t.Errorf("store1.db = %v, want %v", store1.db, db1)
	}

	if store2.db != db2 {
		t.Errorf("store2.db = %v, want %v", store2.db, db2)
	}
}
