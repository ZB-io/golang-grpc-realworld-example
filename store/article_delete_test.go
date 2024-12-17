package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArticleStoreDelete(t *testing.T) {
	// table driven test implementation
	tests := []struct {
		name      string
		article   *model.Article
		setupMock func(mock sqlmock.Sqlmock, article *model.Article)
		wantErr   bool
	}{
		{
			name:    "Successful Deletion",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			setupMock: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectExec("DELETE FROM `articles` WHERE `articles`.`id` = \\?").
					WithArgs(article.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:    "Non-existence Deletion",
			article: &model.Article{Model: gorm.Model{ID: 2}},
			setupMock: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectExec("DELETE FROM `articles` WHERE `articles`.`id` = \\?").
					WithArgs(article.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
		{
			name:    "When Database is Unreachable",
			article: &model.Article{Model: gorm.Model{ID: 3}},
			setupMock: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectExec("DELETE FROM `articles` WHERE `articles`.`id` = \\?").
					WithArgs(article.ID).
					WillReturnError(errors.New("database unreachable"))
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// simulate db and mock interface in place of sql.DB
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			gormDB, gormErr := gorm.Open("mysql", mockDB)
			require.NoError(t, gormErr)

			store := &ArticleStore{db: gormDB}

			// using function to set behaviour for mock db
			test.setupMock(mock, test.article)

			// Invoke method to test
			err = store.Delete(test.article)

			if test.wantErr {
				require.Error(t, err)
				t.Log("Error Here: ", err)
			} else {
				require.NoError(t, err)
			}

			// Check if all expectations were met
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
