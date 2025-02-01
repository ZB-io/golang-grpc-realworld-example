package store

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/raahii/golang-grpc-realworld-example/model"
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
				t.Fatal("NewArticleStore returned nil")
			}

			if !reflect.DeepEqual(got.db, tt.want.db) {
				t.Errorf("NewArticleStore().db = %v, want %v", got.db, tt.want.db)
			}
		})
	}
}

func TestNewArticleStoreDifferentConnections(t *testing.T) {
	mockDB1 := &gorm.DB{}
	mockDB2 := &gorm.DB{}

	store1 := NewArticleStore(mockDB1)
	store2 := NewArticleStore(mockDB2)

	if store1.db == store2.db {
		t.Error("ArticleStore instances should have different DB references")
	}
}

func TestNewArticleStoreImmutability(t *testing.T) {
	mockDB := &gorm.DB{}
	store1 := NewArticleStore(mockDB)
	store2 := NewArticleStore(mockDB)

	if store1 == store2 {
		t.Error("NewArticleStore should return distinct instances")
	}

	if store1.db != store2.db {
		t.Error("ArticleStore instances should have the same DB reference")
	}
}

func TestNewArticleStoreType(t *testing.T) {
	mockDB := &gorm.DB{}
	store := NewArticleStore(mockDB)

	if _, ok := interface{}(store).(*ArticleStore); !ok {
		t.Error("NewArticleStore should return a pointer to ArticleStore")
	}
}


/*
ROOST_METHOD_HASH=IsFavorited_7ef7d3ed9e
ROOST_METHOD_SIG_HASH=IsFavorited_f34d52378f

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) 

*/
func TestArticleStoreIsFavorited(t *testing.T) {
	tests := []struct {
		name            string
		article         *model.Article
		user            *model.User
		mockCountResult int
		mockCountError  error
		expectedResult  bool
		expectedError   error
	}{
		{
			name:            "Article is favorited by the user",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 1,
			mockCountError:  nil,
			expectedResult:  true,
			expectedError:   nil,
		},
		{
			name:            "Article is not favorited by the user",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Nil article parameter",
			article:         nil,
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Nil user parameter",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            nil,
			mockCountResult: 0,
			mockCountError:  nil,
			expectedResult:  false,
			expectedError:   nil,
		},
		{
			name:            "Database error",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 0,
			mockCountError:  errors.New("database error"),
			expectedResult:  false,
			expectedError:   errors.New("database error"),
		},
		{
			name:            "Multiple favorites for the same article and user",
			article:         &model.Article{Model: gorm.Model{ID: 1}},
			user:            &model.User{Model: gorm.Model{ID: 1}},
			mockCountResult: 2,
			mockCountError:  nil,
			expectedResult:  true,
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				countResult: tt.mockCountResult,
				countError:  tt.mockCountError,
			}

			store := &ArticleStore{
				db: &gorm.DB{Value: mockDB},
			}

			result, err := store.IsFavorited(tt.article, tt.user)

			if result != tt.expectedResult {
				t.Errorf("Expected result %v, but got %v", tt.expectedResult, result)
			}

			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("Expected error %v, but got %v", tt.expectedError, err)
			}
		})
	}
}

