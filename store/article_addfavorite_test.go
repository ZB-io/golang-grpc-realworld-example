package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"database/sql"
)

func TestArticleAddFavorite(t *testing.T) {
	var err error

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("sqlmock", db) // open gorm db

	if err != nil {
		t.Fatal("failed to reopen mock db as gorm db.")
	}

	s := &ArticleStore{db: gormDB}

	testCases := []struct {
		name              string
		article           model.Article
		user              model.User
		addFavoritedError error
		updateCountError  error
		wantError         bool
	}{
		{
			name: "Successful Addition of Favorite Article",
			user: model.User{Model: gorm.Model{ID: 1}, Username: "Alice", Email: "alice@gmail.com"},
			article: model.Article{Model: gorm.Model{ID: 1}, Title: "Article", UserID: 1},
			wantError: false,
		},
		{
			name: "User Already Favorited the Article",
			user: model.User{Model: gorm.Model{ID: 1}, Username: "Alice", Email: "alice@gmail.com"},
			article: model.Article{Model: gorm.Model{ID: 1}, Title: "Article", UserID: 1},
			addFavoritedError: gorm.ErrRecordNotFound,
			wantError:         true,
		},
		{
			name: "Failure to Update FavoritesCount",
			user: model.User{Model: gorm.Model{ID: 1}, Username: "Alice", Email: "alice@gmail.com"},
			article: model.Article{Model: gorm.Model{ID: 1}, Title: "Article", UserID: 1},
			updateCountError: gorm.ErrRecordNotFound,
			wantError:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectBegin()

			mock.ExpectExec("^INSERT INTO favorite_articles").
				WithArgs(tc.user.ID, tc.article.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectExec("^UPDATE articles SET favorites_count = favorites_count + 1 WHERE id = ?").
				WithArgs(tc.article.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()

			err := s.AddFavorite(&tc.article, &tc.user)

			if (err != nil) != tc.wantError {
				t.Errorf("ArticleStore.AddFavorite() error = %v, wantError %v", err, tc.wantError)
				return
			}

			if !tc.wantError && tc.article.FavoritesCount != 1 {
				t.Errorf("ArticleStore.AddFavorite() FavoritesCount = %v, want %v", tc.article.FavoritesCount, 1)
				return
			}
		})
	}
}

