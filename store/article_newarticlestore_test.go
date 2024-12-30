package store

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
)


type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestNewArticleStore(t *testing.T) {
	tests := []struct {
		name string
		db   *gorm.DB
		want *ArticleStore
	}{
		{
			name: "Create ArticleStore with Valid DB Connection",
			db:   &gorm.DB{},
			want: &ArticleStore{db: &gorm.DB{}},
		},
		{
			name: "Create ArticleStore with Nil DB Connection",
			db:   nil,
			want: &ArticleStore{db: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewArticleStore(tt.db)

			if got == nil {
				t.Error("NewArticleStore returned nil")
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArticleStore() = %v, want %v", got, tt.want)
			}

			if got.db != tt.db {
				t.Errorf("NewArticleStore().db = %v, want %v", got.db, tt.db)
			}

			if _, ok := interface{}(got).(*ArticleStore); !ok {
				t.Errorf("NewArticleStore() did not return *ArticleStore")
			}
		})
	}

	t.Run("Create Multiple ArticleStores with Different DB Connections", func(t *testing.T) {
		db1 := &gorm.DB{Value: 1}
		db2 := &gorm.DB{Value: 2}

		store1 := NewArticleStore(db1)
		store2 := NewArticleStore(db2)

		if store1.db != db1 {
			t.Errorf("First ArticleStore has incorrect DB reference")
		}

		if store2.db != db2 {
			t.Errorf("Second ArticleStore has incorrect DB reference")
		}
	})

	t.Run("Performance Test for NewArticleStore", func(t *testing.T) {
		db := &gorm.DB{}
		iterations := 1000

		for i := 0; i < iterations; i++ {
			_ = NewArticleStore(db)
		}

	})
}
