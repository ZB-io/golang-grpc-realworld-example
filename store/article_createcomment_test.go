package store

import (
	"testing"
	"errors"
	"log"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

type ExpectedBegin struct {
	commonExpectation
	delay time.Duration
}

type ExpectedCommit struct {
	commonExpectation
}

type ExpectedExec struct {
	queryBasedExpectation
	result driver.Result
	delay  time.Duration
}

type ExpectedRollback struct {
	commonExpectation
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
func TestArticleStoreCreateComment(t *testing.T) {
	testCases := []struct {
		name          string
		inputComment  model.Comment
		mockSetup     func(sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name: "Successful Comment Creation",
			inputComment: model.Comment{
				Body:      "This is a valid comment",
				UserID:    1,
				ArticleID: 1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"comments\"").
					WithArgs(sqlmock.AnyArg(), "This is a valid comment", 1, 1, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Comment Creation with Missing Body",
			inputComment: model.Comment{
				UserID:    1,
				ArticleID: 1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"comments\"").
					WithArgs(sqlmock.AnyArg(), "", 1, 1, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("not null violation"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("not null violation"),
		},
		{
			name: "Comment Creation with Invalid User ID",
			inputComment: model.Comment{
				Body:      "This is a comment",
				UserID:    9999,
				ArticleID: 1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"comments\"").
					WithArgs(sqlmock.AnyArg(), "This is a comment", 9999, 1, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("foreign key constraint failed"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("foreign key constraint failed"),
		},
		{
			name: "Comment Creation with Database Connection Failure",
			inputComment: model.Comment{
				Body:      "This is a valid comment",
				UserID:    1,
				ArticleID: 1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("connection error"))
			},
			expectedError: errors.New("connection error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tc.mockSetup(mock)

			gormDB, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("failed to open gorm DB: %v", err)
			}
			defer gormDB.Close()

			store := &ArticleStore{db: gormDB}

			err = store.CreateComment(&tc.inputComment)
			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
