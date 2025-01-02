package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)








func TestArticleStoreUpdate(t *testing.T) {
	var tests = []struct {
		name     string
		mockFunc func(mock sqlmock.Sqlmock, article *model.Article)
		article  *model.Article
		wantErr  bool
	}{
		{
			name: "Successful Article Update",
			mockFunc: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WithArgs(article.Title, article.Description, article.Body).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			article: &model.Article{
				Title:       "Test Title",
				Description: "Test Description",
				Body:        "Test Body",
			},
			wantErr: false,
		},
		{
			name: "Article Update with Non-Existent Article",
			mockFunc: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WithArgs(article.Title, article.Description, article.Body).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			article: &model.Article{
				Title:       "Test Title",
				Description: "Test Description",
				Body:        "Test Body",
			},
			wantErr: true,
		},
		{
			name: "Database Update Error",
			mockFunc: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WithArgs(article.Title, article.Description, article.Body).WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			article: &model.Article{
				Title:       "Test Title",
				Description: "Test Description",
				Body:        "Test Body",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening gorm database", err)
			}

			tt.mockFunc(mock, tt.article)

			articleStore := &ArticleStore{db: gormDB}
			err = articleStore.Update(tt.article)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

