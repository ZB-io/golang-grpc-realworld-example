package store

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"testing"
)

func TestArticleStoreCreateComment(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		comment     *model.Comment
		dbError     error
		expectError bool
	}{
		{
			name: "Successful Creation of a Comment",
			comment: &model.Comment{
				Body:   "test body",
				UserID: 1,
			},
			dbError:     nil,
			expectError: false,
		},
		{
			name: "Failure to Create a Comment due to DB Error",
			comment: &model.Comment{
				Body:   "test body",
				UserID: 1,
			},
			dbError:     errors.New("DB error"),
			expectError: true,
		},
		{
			name: "Creation of a Comment with a Null Body",
			comment: &model.Comment{
				Body:   "",
				UserID: 1,
			},
			dbError:     errors.New("DB error"),
			expectError: true,
		},
		{
			name: "Creation of a Comment with a Null UserID",
			comment: &model.Comment{
				Body:   "test body",
				UserID: 0,
			},
			dbError:     errors.New("DB error"),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			gdb, _ := gorm.Open("mysql", db)
			store := &ArticleStore{db: gdb}

			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO `comments`").
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tc.dbError)
			mock.ExpectCommit()

			err := store.CreateComment(tc.comment)

			if tc.expectError {
				if err == nil {
					t.Errorf("[%s] expected error, got nil", tc.name)
				} else {
					t.Logf("[%s] expected error, got error: %s", tc.name, err)
				}
			} else {
				if err != nil {
					t.Errorf("[%s] expected no error, got error: %s", tc.name, err)
				} else {
					t.Logf("[%s] expected no error, got no error", tc.name)
				}
			}
		})
	}
}
