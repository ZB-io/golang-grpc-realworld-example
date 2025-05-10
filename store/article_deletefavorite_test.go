package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestDeleteFavorite(t *testing.T) {
	type testCase struct {
		name      string
		setup     func(mock sqlmock.Sqlmock, a *model.Article, u *model.User)
		article   *model.Article
		user      *model.User
		expectErr bool
	}

	tests := []testCase{
		{
			name: "Successful Deletion from Favorited Users",
			setup: func(mock sqlmock.Sqlmock, a *model.Article, u *model.User) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(u.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE "articles" SET "favorites_count" = "favorites_count" - ? WHERE "articles"."deleted_at" IS NULL`).
					WithArgs(1, a.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			article:   &model.Article{Model: gorm.Model{ID: 1}, FavoritesCount: 1, FavoritedUsers: []model.User{{Model: gorm.Model{ID: 1}}}},
			user:      &model.User{Model: gorm.Model{ID: 1}},
			expectErr: false,
		},
		{
			name: "No Change When User Not in Favorited List",
			setup: func(mock sqlmock.Sqlmock, a *model.Article, u *model.User) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(u.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			},
			article:   &model.Article{Model: gorm.Model{ID: 1}, FavoritesCount: 1, FavoritedUsers: []model.User{}},
			user:      &model.User{Model: gorm.Model{ID: 2}},
			expectErr: false,
		},
		{
			name: "Transaction Rollback When Deletion Fails",
			setup: func(mock sqlmock.Sqlmock, a *model.Article, u *model.User) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(u.ID).
					WillReturnError(errors.New("delete error"))
				mock.ExpectRollback()
			},
			article:   &model.Article{Model: gorm.Model{ID: 1}, FavoritesCount: 1, FavoritedUsers: []model.User{{Model: gorm.Model{ID: 1}}}},
			user:      &model.User{Model: gorm.Model{ID: 1}},
			expectErr: true,
		},
		{
			name: "Transaction Rollback When Update Fails",
			setup: func(mock sqlmock.Sqlmock, a *model.Article, u *model.User) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(u.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE "articles" SET "favorites_count" = "favorites_count" - ? WHERE "articles"."deleted_at" IS NULL`).
					WithArgs(1, a.ID).
					WillReturnError(errors.New("update error"))
				mock.ExpectRollback()
			},
			article:   &model.Article{Model: gorm.Model{ID: 1}, FavoritesCount: 1, FavoritedUsers: []model.User{{Model: gorm.Model{ID: 1}}}},
			user:      &model.User{Model: gorm.Model{ID: 1}},
			expectErr: true,
		},
		{
			name: "Handle Nil Arguments Gracefully",
			setup: func(mock sqlmock.Sqlmock, a *model.Article, u *model.User) {

			},
			article:   nil,
			user:      nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gdb, err := gorm.Open("postgres", db)
			assert.NoError(t, err)

			store := ArticleStore{db: gdb}

			if tt.setup != nil {
				tt.setup(mock, tt.article, tt.user)
			}

			err = store.DeleteFavorite(tt.article, tt.user)

			if tt.expectErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Did not expect error but got one")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
