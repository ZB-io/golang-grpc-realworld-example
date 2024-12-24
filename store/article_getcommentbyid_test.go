package store

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ExpectedQuery struct {
	queryBasedExpectation
	rows             driver.Rows
	delay            time.Duration
	rowsMustBeClosed bool
	rowsWereClosed   bool
}

type Rows struct {
	converter driver.ValueConverter
	cols      []string
	def       []*Column
	rows      [][]driver.Value
	pos       int
	nextErr   map[int]error
	closeErr  error
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
func TestArticleStoreGetCommentByID(t *testing.T) {
	var Scenarios = []struct {
		Scenario      string
		ID            uint
		DBSetup       func(mock sqlmock.Sqlmock) error
		ExpectedError error
		Expected      *model.Comment
	}{
		{
			Scenario: "Successful retrieval of a comment by its ID",
			ID:       1,
			DBSetup: func(mock sqlmock.Sqlmock) error {
				rows := sqlmock.NewRows([]string{"id", "body", "user_id", "article_id"}).AddRow(1, "Test comment", 1, 1)
				mock.ExpectQuery("^SELECT \\* FROM \"comments\" WHERE \"id\" = \\?$").WithArgs(1).WillReturnRows(rows)
				return nil
			},
			Expected: &model.Comment{
				Model:     gorm.Model{ID: 1},
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
		},
		{
			Scenario: "Trying to retrieve a comment with an invalid ID",
			ID:       99,
			DBSetup: func(mock sqlmock.Sqlmock) error {
				mock.ExpectQuery("^SELECT \\* FROM \"comments\" WHERE \"id\" = \\?$").WithArgs(99).WillReturnError(gorm.ErrRecordNotFound)
				return nil
			},
			ExpectedError: gorm.ErrRecordNotFound,
		},
		{
			Scenario: "Database Availability Issues",
			ID:       1,
			DBSetup: func(mock sqlmock.Sqlmock) error {
				mock.ExpectQuery("^SELECT \\* FROM \"comments\" WHERE \"id\" = \\?$").WithArgs(1).WillReturnError(gorm.ErrInvalidSQL)
				return nil
			},
			ExpectedError: gorm.ErrInvalidSQL,
		},
	}

	for _, s := range Scenarios {
		t.Run(s.Scenario, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			gormDB, _ := gorm.Open("postgres", db)
			defer func() {
				_ = db.Close()
				_ = gormDB.Close()
			}()

			store := ArticleStore{
				db: gormDB,
			}

			err := s.DBSetup(mock)
			if err != nil {
				t.Fatalf("an error '%s' occurred when setting up the mock database", err)
			}

			comment, err := store.GetCommentByID(s.ID)

			assert.Equal(t, s.ExpectedError, err, "expected error does not match actual error")
			assert.Equal(t, s.Expected, comment, "expected comment does not match actual comment")
		})
	}
}
