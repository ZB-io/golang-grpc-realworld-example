package store

import (
	"testing"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
)




type mockDB struct {
	gorm.DB
	mock sqlmock.Sqlmock
}
func (db *mockDB) Association(column string) *gorm.Association {
	association := &gorm.Association{}
	association.Error = db.DB.Error
	return association
}
func (db *mockDB) Begin() *gorm.DB {
	return &db.DB
}
func (db *mockDB) Commit() *gorm.DB {
	return &db.DB
}
func (db *mockDB) Model(value interface{}) *gorm.DB {
	return &db.DB
}
func (db *mockDB) Rollback() *gorm.DB {
	return &db.DB
}
func TestArticleStoreDeleteFavorite(t *testing.T) {

	testCases := []struct {
		name           string
		article        *model.Article
		user           *model.User
		mockFunc       func() (*mockDB, error)
		expectedError  error
		expectedResult bool
	}{
		{
			name:    "Successful Deletion of Favorite",
			article: &model.Article{FavoritesCount: 1, FavoritedUsers: []model.User{{Username: "test_user"}}},
			user:    &model.User{Username: "test_user"},
			mockFunc: func() (*mockDB, error) {
				db, _, _ := sqlmock.New()
				gormDB, err := gorm.Open("mysql", db)
				mockDB := &mockDB{DB: *gormDB}
				return mockDB, err
			},
			expectedError:  nil,
			expectedResult: true,
		},
		{
			name:    "User Not in FavoritedUsers",
			article: &model.Article{FavoritesCount: 1, FavoritedUsers: []model.User{{Username: "another_user"}}},
			user:    &model.User{Username: "test_user"},
			mockFunc: func() (*mockDB, error) {
				db, _, _ := sqlmock.New()
				gormDB, err := gorm.Open("mysql", db)
				mockDB := &mockDB{DB: *gormDB}
				mockDB.DB.Error = errors.New("user not in favorited users")
				return mockDB, err
			},
			expectedError:  errors.New("user not in favorited users"),
			expectedResult: false,
		},
		{
			name:    "Error During Deletion",
			article: &model.Article{FavoritesCount: 1, FavoritedUsers: []model.User{{Username: "test_user"}}},
			user:    &model.User{Username: "test_user"},
			mockFunc: func() (*mockDB, error) {
				db, _, _ := sqlmock.New()
				gormDB, err := gorm.Open("mysql", db)
				mockDB := &mockDB{DB: *gormDB}
				mockDB.DB.Error = errors.New("error during deletion")
				return mockDB, err
			},
			expectedError:  errors.New("error during deletion"),
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDB, err := tc.mockFunc()
			if err != nil {
				t.Fatalf("error setting up mock DB: %v", err)
			}

			store := &ArticleStore{db: &mockDB.DB}

			err = store.DeleteFavorite(tc.article, tc.user)

			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}

			if tc.expectedResult && tc.article.FavoritesCount != 0 {
				t.Errorf("expected favorites count to be 0, got %v", tc.article.FavoritesCount)
			}

			if !tc.expectedResult && tc.article.FavoritesCount != 1 {
				t.Errorf("expected favorites count to be 1, got %v", tc.article.FavoritesCount)
			}
		})
	}
}
func (db *mockDB) Update(attrs ...interface{}) *gorm.DB {
	return &db.DB
}
