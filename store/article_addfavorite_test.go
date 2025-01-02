package store

import (
	"testing"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)











func TestArticleStoreAddFavorite(t *testing.T) {

	testCases := []struct {
		name                   string
		setupMock              func(mock sqlmock.Sqlmock)
		expectedError          error
		expectedFavoritesCount int32
	}{
		{
			name: "Adding a new favorite article successfully",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO favorite_articles").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("^UPDATE articles SET favorites_count").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError:          nil,
			expectedFavoritesCount: 1,
		},
		{
			name: "Adding a favorite article when the article is already in the user's favorites",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO favorite_articles").
					WillReturnError(errors.New("article already favorited"))
				mock.ExpectRollback()
			},
			expectedError:          errors.New("article already favorited"),
			expectedFavoritesCount: 0,
		},
		{
			name: "Adding a favorite article when there is a DB error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO favorite_articles").
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectedError:          errors.New("db error"),
			expectedFavoritesCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, _ := sqlmock.New()
			defer db.Close()
			gdb, _ := gorm.Open("postgres", db)
			tc.setupMock(mock)

			article := &model.Article{Title: "Test Article", FavoritesCount: 0}
			user := &model.User{Username: "Test User"}

			store := &ArticleStore{db: gdb}

			err := store.AddFavorite(article, user)

			if (err != nil || tc.expectedError != nil) && (err == nil || tc.expectedError == nil || err.Error() != tc.expectedError.Error()) {
				t.Errorf("got error %v, want %v", err, tc.expectedError)
			}

			if article.FavoritesCount != tc.expectedFavoritesCount {
				t.Errorf("got favorites count %v, want %v", article.FavoritesCount, tc.expectedFavoritesCount)
			}
		})
	}
}
