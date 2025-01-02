package store

import (
	"errors"
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/jinzhu/gorm"
	"github.com/DATA-DOG/go-sqlmock"
)








func TestArticleStoreDelete(t *testing.T) {

	testCases := []struct {
		name    string
		article *model.Article
		dbMock  func(mock sqlmock.Sqlmock, article *model.Article)
		wantErr error
	}{
		{
			name:    "Normal operation - Deleting an existing article",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			dbMock: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `articles`").
					WithArgs(article.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name:    "Edge case - Deleting a non-existing article",
			article: &model.Article{Model: gorm.Model{ID: 2}},
			dbMock: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `articles`").
					WithArgs(article.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "Error handling - Deleting an article with a closed database connection",
			article: &model.Article{Model: gorm.Model{ID: 3}},
			dbMock: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `articles`").
					WithArgs(article.ID).
					WillReturnError(errors.New("database connection closed"))
				mock.ExpectRollback()
			},
			wantErr: errors.New("database connection closed"),
		},
		{
			name:    "Edge case - Deleting an article with associated tags, comments and favorited users",
			article: &model.Article{Model: gorm.Model{ID: 4}},
			dbMock: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `articles`").
					WithArgs(article.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tc.dbMock(mock, tc.article)

			gdb, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening gorm database", err)
			}

			store := &ArticleStore{db: gdb}

			err = store.Delete(tc.article)

			if err != nil && err != tc.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tc.wantErr)
			}

			if err == nil && tc.wantErr != nil {
				t.Errorf("Delete() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
