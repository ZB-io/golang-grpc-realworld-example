package store

import (
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/DATA-DOG/go-sqlmock"
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
func TestArticleStoreDelete(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("could not open gorm DB: %v", err)
	}

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name         string
		setup        func()
		input        *model.Article
		expectedErr  bool
		expectedDesc string
	}{
		{
			name: "Successfully delete an existing article",
			setup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "articles"`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			input:        &model.Article{Model: gorm.Model{ID: 1}, Title: "Test Title"},
			expectedErr:  false,
			expectedDesc: "This test aims to verify that the Delete function correctly removes an existing article from the database without errors.",
		},
		{
			name: "Attempt to delete a non-existing article",
			setup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "articles"`).
					WithArgs(2).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			input:        &model.Article{Model: gorm.Model{ID: 2}, Title: "Non-existent"},
			expectedErr:  false,
			expectedDesc: "This test ensures that the Delete function handles the case of a non-existing article gracefully.",
		},

		{
			name: "Delete an article when the database connection is down",
			setup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "articles"`).
					WithArgs(3).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			input:        &model.Article{Model: gorm.Model{ID: 3}, Title: "DB Down"},
			expectedErr:  true,
			expectedDesc: "Test the Delete function's behavior and error handling when the database connection is unavailable.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := store.Delete(tt.input)

			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Logf(tt.expectedDesc)
		})
	}
}
