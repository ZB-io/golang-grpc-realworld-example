package store

import (
	"errors"
	"testing"
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

func TestArticleStoreCreate(t *testing.T) {
	testCases := []struct {
		name        string
		mockArticle *model.Article
		wantedError error
	}{
		{
			name:        "Scenario 1: Article creation with valid data",
			mockArticle: &model.Article{UserID: 1, Body: "Body", Title: "Title", Description: "Description"},
			wantedError: nil,
		},
		{
			name:        "Scenario 2: Article creation with null data fields",
			mockArticle: &model.Article{},
			wantedError: errors.New("record not found"),
		},
		{
			name:        "Scenario 3: Article creation with non-existent user",
			mockArticle: &model.Article{UserID: 2},
			wantedError: errors.New("record not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}

			gdb, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening gorm database", err)
			}

			AS := &ArticleStore{db: gdb}
			mock.ExpectBegin()

			if tc.wantedError == nil {
				mock.ExpectExec("^INSERT INTO \"articles\" .*").WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec("^INSERT INTO \"articles\" .*").WillReturnError(tc.wantedError)
			}

			mock.ExpectCommit()

			err = AS.Create(tc.mockArticle)

			if (err != nil && tc.wantedError == nil) || (err == nil && tc.wantedError != nil) || (err != nil && err.Error() != tc.wantedError.Error()) {
				t.Fatalf("Errors do not match: got error '%s', wanted error '%s'", err, tc.wantedError)
			}

			mock.ExpectationsWereMet()
		})
	}
}
