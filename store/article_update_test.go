package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

// TestArticleStoreUpdate tests the Update function in store.ArticleStore
func TestArticleStoreUpdate(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock, article *model.Article)
		article       *model.Article
		expectedError error
	}{
		{
			"Successful Article Update",
			func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WithArgs(article.ID, article.Title, article.Description, article.Body, article.UserID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			&model.Article{ID: 1, Title: "Updated Article", Description: "Updated Description", Body: "Updated Body", UserID: 1},
			nil,
		},
		{
			"Article Update with Non-Existent Article",
			func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WithArgs(article.ID, article.Title, article.Description, article.Body, article.UserID).WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectRollback()
			},
			&model.Article{ID: 2, Title: "Non-Existent Article", Description: "Does not exist", Body: "Does not exist", UserID: 2},
			gorm.ErrRecordNotFound,
		},
		{
			"Article Update with Null Fields",
			func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WithArgs(article.ID, article.Title, article.Description, article.Body, article.UserID).WillReturnError(gorm.ErrInvalidSQL)
				mock.ExpectRollback()
			},
			&model.Article{ID: 3, Title: "", Description: "", Body: "", UserID: 3},
			gorm.ErrInvalidSQL,
		},
		{
			"Database Connection Error during Article Update",
			func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WithArgs(article.ID, article.Title, article.Description, article.Body, article.UserID).WillReturnError(errors.New("database connection error"))
				mock.ExpectRollback()
			},
			&model.Article{ID: 4, Title: "Valid Article", Description: "Valid Description", Body: "Valid Body", UserID: 4},
			errors.New("database connection error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock database
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening gorm database", err)
			}

			// Setup mock
			tc.setupMock(mock, tc.article)

			// Update article
			store := &ArticleStore{db: gormDB}
			err = store.Update(tc.article)

			// Assert
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
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
