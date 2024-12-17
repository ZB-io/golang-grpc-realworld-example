package store_test

import (
	"errors"
	"testing"

	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/jinzhu/gorm"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestArticleStoreDeleteComment(t *testing.T) {

	mockError := errors.New("mock Error")

	tests := []struct {
		name          string
		prepare       func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name: "Successful Comment Deletion",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Attempt To Delete Nonexistent Comment",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectRollback()
			},
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WillReturnError(mockError)
				mock.ExpectRollback()
			},
			expectedError: mockError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db, mock, _ := sqlmock.New()
			test.prepare(mock)
			gormDb, _ := gorm.Open("postgres", db)

			store := &ArticleStore{gormDb}

			comment := &model.Comment{
				Model: gorm.Model{
					ID: 1,
				},
			}

			err := store.DeleteComment(comment)

			if test.expectedError != nil {
				if err == nil  || err.Error() != test.expectedError.Error() {
					t.Errorf("Expected error %v, Received %v", test.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, Received %v", err)
				}
			}
		})
	}
}
