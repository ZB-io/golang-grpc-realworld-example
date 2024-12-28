package store

import (
	"errors"
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDeleteFavorite(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, err: %s", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create gorm db, err: %s", err)
	}

	articleStore := &ArticleStore{db: gormDB}

	tests := []struct {
		name            string
		article         *model.Article
		user            *model.User
		setupMocks      func(mock sqlmock.Sqlmock)
		expectError     bool
		expectedFavsCnt int
	}{
		{
			name: "Successful Deletion of a Favorited User",
			article: &model.Article{
				ID:             1,
				FavoritesCount: 3,
			},
			user: &model.User{ID: 1},
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "user_favorites"`).WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE "articles"`).WithArgs(2, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError:     false,
			expectedFavsCnt: 2,
		},
		{
			name: "Attempting to Delete Non-Favorited User",
			article: &model.Article{
				ID:             2,
				FavoritesCount: 1,
			},
			user: &model.User{ID: 2},
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "user_favorites"`).WithArgs(2, 2).WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
			},
			expectError:     false,
			expectedFavsCnt: 1,
		},
		{
			name:    "Database Begin Transaction Failure",
			article: &model.Article{ID: 3},
			user:    &model.User{ID: 3},
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin transaction error"))
			},
			expectError:     true,
			expectedFavsCnt: 0,
		},
		{
			name: "Error During User Association Deletion",
			article: &model.Article{
				ID:             4,
				FavoritesCount: 5,
			},
			user: &model.User{ID: 4},
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "user_favorites"`).WithArgs(4, 4).WillReturnError(errors.New("association deletion error"))
				mock.ExpectRollback()
			},
			expectError:     true,
			expectedFavsCnt: 5,
		},
		{
			name: "Error During Update of Favorites Count",
			article: &model.Article{
				ID:             5,
				FavoritesCount: 10,
			},
			user: &model.User{ID: 5},
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "user_favorites"`).WithArgs(5, 5).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE "articles"`).WithArgs(9, 5).WillReturnError(errors.New("favorites count update error"))
				mock.ExpectRollback()
			},
			expectError:     true,
			expectedFavsCnt: 10,
		},
		{
			name: "No Favorited Users Initially",
			article: &model.Article{
				ID:             6,
				FavoritesCount: 0,
			},
			user: &model.User{ID: 6},
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "user_favorites"`).WithArgs(6, 6).WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
			},
			expectError:     false,
			expectedFavsCnt: 0,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks(mock)

			err := articleStore.DeleteFavorite(tc.article, tc.user)

			if (err != nil) != tc.expectError {
				t.Errorf("Expected error: %v, got: %v", tc.expectError, err != nil)
			}

			if tc.article.FavoritesCount != tc.expectedFavsCnt {
				t.Errorf("Expected favorites count: %d, got: %d", tc.expectedFavsCnt, tc.article.FavoritesCount)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test %s completed", tc.name)
		})
	}
}

