package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

func TestCreate(t *testing.T) {

	testCases := []struct {
		description   string
		article       model.Article
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
	}{
		{
			description: "Successfully Creating an Article",
			article: model.Article{
				Title:       "Test Article",
				Description: "This is a test description.",
				Body:        "This is the body of the test article.",
				UserID:      1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "articles"`).
					WithArgs(sqlmock.AnyArg(), "Test Article", "This is a test description.", "This is the body of the test article.", 1, 0).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			description: "Failure Due to Missing Required Fields",
			article: model.Article{

				UserID: 1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "articles"`).
					WillReturnError(errors.New("missing required fields"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("missing required fields"),
		},
		{
			description: "Database Connection Error",
			article: model.Article{
				Title:       "Test Article",
				Description: "This is a test description.",
				Body:        "This is the body of the test article.",
				UserID:      1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("connection error"))
			},
			expectedError: errors.New("connection error"),
		},
		{
			description: "Unique Constraint Violation",
			article: model.Article{
				Title:       "Duplicate Title",
				Description: "This is a test description.",
				Body:        "This is the body of the test article.",
				UserID:      1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "articles"`).
					WithArgs(sqlmock.AnyArg(), "Duplicate Title", "This is a test description.", "This is the body of the test article.", 1, 0).
					WillReturnError(errors.New("unique constraint violation"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("unique constraint violation"),
		},
		{
			description: "Large Data Input Handling",
			article: model.Article{
				Title:       string(make([]byte, 1000)),
				Description: string(make([]byte, 2000)),
				Body:        string(make([]byte, 5000)),
				UserID:      1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "articles"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1, 0).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock database: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("failed to open gorm db: %v", err)
			}
			defer gormDB.Close()

			tc.setupMock(mock)

			articleStore := ArticleStore{db: gormDB}

			err = articleStore.Create(&tc.article)

			if tc.expectedError == nil {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			} else {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Errorf("expected error '%v', but got '%v'", tc.expectedError, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
