package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

type ExpectedExec struct {
	queryBasedExpectation
	result driver.Result
	delay  time.Duration
}

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
func TestArticleStoreDelete(t *testing.T) {
	scenarios := []struct {
		desc      string
		mockFunc  func(mock sqlmock.Sqlmock, article *model.Article)
		article   *model.Article
		wantError error
	}{
		{
			desc: "Deleting an existing article",
			mockFunc: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectExec("^DELETE FROM `articles` WHERE .*").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			article:   &model.Article{Model: gorm.Model{ID: 1}},
			wantError: nil,
		},
		{
			desc: "Deleting an article that does not exist",
			mockFunc: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectExec("^DELETE FROM `articles` WHERE .*").WillReturnResult(sqlmock.NewResult(0, 0))
			},
			article:   &model.Article{Model: gorm.Model{ID: 100}},
			wantError: gorm.ErrRecordNotFound,
		},
		{
			desc: "Deleting an article with related data",
			mockFunc: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectExec("^DELETE FROM `articles` WHERE .*").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("^DELETE FROM .* WHERE .*").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			article:   &model.Article{Model: gorm.Model{ID: 1}},
			wantError: nil,
		},
		{
			desc: "Passing a nil to the Delete function",
			mockFunc: func(mock sqlmock.Sqlmock, article *model.Article) {
			},
			article:   nil,
			wantError: errors.New("nil article"),
		},
	}

	for _, s := range scenarios {
		t.Run(s.desc, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening gorm database", err)
			}

			defer gormDB.Close()

			articleStore := &ArticleStore{db: gormDB}

			s.mockFunc(mock, s.article)

			err = articleStore.Delete(s.article)

			if s.wantError != nil && !errors.Is(err, s.wantError) {
				t.Errorf("%s: expected %v error, but got %v", s.desc, s.wantError, err)
			}

			if s.wantError == nil && err != nil {
				t.Errorf("%s: expected no error, but got %v", s.desc, err)
			}
		})
	}
}
