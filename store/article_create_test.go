package store

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/require"
)

// TestArticleStoreCreate test cases
func TestArticleStoreCreate(t *testing.T) {
	testCases := []struct {
		name          string
		article       *model.Article
		mockDBfunc    func(mock sqlmock.Sqlmock, article *model.Article)
		expectedError error
	}{
		{
			name: "Successful Article Creation",
			article: &model.Article{
				Title:       "Golang Tips",
				Description: "Useful tips for Golang developers",
				Body:        "Test your code with unit tests",
			},
			mockDBfunc: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO articles").
					WithArgs(article.Title, article.Description, article.Body).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Failed Article Creation due to Database Error",
			article: &model.Article{
				Title:       "Golang Tips",
				Description: "Useful tips for Golang developers",
				Body:        "Test your code with unit tests",
			},
			mockDBfunc: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO articles").
					WithArgs(article.Title, article.Description, article.Body).
					WillReturnError(gorm.ErrUnaddressable)
				mock.ExpectRollback()
			},
			expectedError: gorm.ErrUnaddressable,
		},
		{
			name: "Failed Article Creation due to Invalid Article Data",
			article: &model.Article{
				Title: "", // Empty title is invalid
			},
			mockDBfunc: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
			expectedError: gorm.ErrInvalidSQL,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			gormDB, _ := gorm.Open("mysql", db)

			// Initialize ArticleStore
			store := &ArticleStore{
				db: gormDB,
			}

			tt.mockDBfunc(mock, tt.article)

			err := store.Create(tt.article)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
