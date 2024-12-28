package store

import (
	"errors"
	"sync"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func TestAddFavorite(t *testing.T) {
	t.Parallel()

	type args struct {
		article *model.Article
		user    *model.User
	}

	tests := []struct {
		name        string
		args        args
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "Successfully add a user to favorites",
			args: args{
				article: &model.Article{ID: 1, FavoritesCount: 0},
				user:    &model.User{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO article_favorited_users").
					WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE articles SET favorites_count").
					WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "User already a favorite",
			args: args{
				article: &model.Article{ID: 1, FavoritesCount: 1},
				user:    &model.User{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO article_favorited_users").
					WithArgs(1, 1).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Append user error triggers rollback",
			args: args{
				article: &model.Article{ID: 1, FavoritesCount: 0},
				user:    &model.User{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO article_favorited_users").
					WithArgs(1, 1).WillReturnError(errors.New("insert error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "Error updating favorites count triggers rollback",
			args: args{
				article: &model.Article{ID: 1, FavoritesCount: 0},
				user:    &model.User{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO article_favorited_users").
					WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE articles SET favorites_count").
					WithArgs(sqlmock.AnyArg()).WillReturnError(errors.New("update error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "Concurrency safety with multiple updates",
			args: args{
				article: &model.Article{ID: 1, FavoritesCount: 0},
				user:    &model.User{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO article_favorited_users").
					WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE articles SET favorites_count").
					WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Nil user input",
			args: args{
				article: &model.Article{ID: 1, FavoritesCount: 0},
				user:    nil,
			},
			setupMock:   func(mock sqlmock.Sqlmock) {},
			expectError: true,
		},
		{
			name: "Nil article input",
			args: args{
				article: nil,
				user:    &model.User{ID: 1},
			},
			setupMock:   func(mock sqlmock.Sqlmock) {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error opening stub database connection: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("failed to open gorm DB: %v", err)
			}

			tt.setupMock(mock)

			store := &ArticleStore{db: gormDB}
			err = store.AddFavorite(tt.args.article, tt.args.user)

			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}

	t.Run("Concurrent update test", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening stub database connection: %v", err)
		}
		defer db.Close()

		gormDB, err := gorm.Open("mysql", db)
		if err != nil {
			t.Fatalf("failed to open gorm DB: %v", err)
		}

		var wg sync.WaitGroup
		article := &model.Article{ID: 1, FavoritesCount: 0}

		store := &ArticleStore{db: gormDB}
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO article_favorited_users").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("UPDATE articles SET favorites_count").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(userID int) {
				defer wg.Done()
				user := &model.User{ID: uint(userID)}
				err := store.AddFavorite(article, user)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}(i)
		}

		wg.Wait()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

