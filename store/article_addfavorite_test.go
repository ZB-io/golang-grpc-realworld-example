package store

import (
	"errors"
	"log"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestAddFavorite(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(sqlmock.Sqlmock, *model.Article, *model.User)
		checkFunc func(*model.Article, error)
	}{
		{
			name: "Successfully Add a Favorite",
			setupMock: func(mock sqlmock.Sqlmock, article *model.Article, user *model.User) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WithArgs(sqlmock.AnyArg(), user.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles` SET `favorites_count` = `favorites_count` + ?").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			checkFunc: func(article *model.Article, err error) {
				if err != nil {
					t.Errorf("expected no error, but got %v", err)
				}
				if article.FavoritesCount != 1 {
					t.Errorf("expected favorites count to be %d, but got %d", 1, article.FavoritesCount)
				}
			},
		},
		{
			name: "Error When User Already Favorited the Article",
			setupMock: func(mock sqlmock.Sqlmock, article *model.Article, user *model.User) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WillReturnError(errors.New("duplicate entry"))
				mock.ExpectRollback()
			},
			checkFunc: func(article *model.Article, err error) {
				if err == nil {
					t.Errorf("expected error, but got none")
				}
				if article.FavoritesCount != 0 {
					t.Errorf("expected favorites count to remain %d, got %d", 0, article.FavoritesCount)
				}
			},
		},
		{
			name: "Error Due to Database Transaction Failure",
			setupMock: func(mock sqlmock.Sqlmock, article *model.Article, user *model.User) {
				mock.ExpectBegin().WillReturnError(errors.New("db error"))
			},
			checkFunc: func(article *model.Article, err error) {
				if err == nil {
					t.Errorf("expected error, but got none")
				}
				if article.FavoritesCount != 0 {
					t.Errorf("expected favorites count to remain %d, got %d", 0, article.FavoritesCount)
				}
			},
		},
		{
			name: "Handle Null Article as Input",
			setupMock: func(mock sqlmock.Sqlmock, article *model.Article, user *model.User) {

			},
			checkFunc: func(article *model.Article, err error) {
				if err == nil || article != nil {
					t.Error("expected error and nil article")
				}
			},
		},
		{
			name: "Handle Null User as Input",
			setupMock: func(mock sqlmock.Sqlmock, article *model.Article, user *model.User) {

			},
			checkFunc: func(article *model.Article, err error) {
				if err == nil {
					t.Errorf("expected error, but got none")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock database: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				log.Fatalf("failed to open gorm DB: %v", err)
			}
			defer gormDB.Close()

			articleStore := &ArticleStore{db: gormDB}
			article := &model.Article{}
			user := &model.User{}

			tt.setupMock(mock, article, user)

			err = articleStore.AddFavorite(article, user)

			tt.checkFunc(article, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}
