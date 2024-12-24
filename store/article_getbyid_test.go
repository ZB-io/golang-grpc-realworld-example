package store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/DATA-DOG/go-sqlmock"
)

var ErrConnectionFailed = errors.New("connection failed")type ExpectedQuery struct {
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

type Assertions struct {
	t TestingT
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
func TestGetByID(t *testing.T) {
	assert := assert.New(t)

	scenarios := []struct {
		Name        string
		SetupMock   func(mock sqlmock.Sqlmock)
		ID          uint
		ExpectError bool
	}{
		{
			Name: "Valid Article ID",
			ID:   1,
			SetupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "tags", "author", "favorites_count"}).
					AddRow(1, "test title", "test description", "test body", []model.Tag{}, model.User{}, 0)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			ExpectError: false,
		},
		{
			Name: "Invalid Article ID",
			ID:   1,
			SetupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(gorm.ErrRecordNotFound)
			},
			ExpectError: true,
		},
		{
			Name: "Test other scenario where Preload fails",
			ID:   1,
			SetupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(gorm.ErrInvalidSQL)
			},
			ExpectError: true,
		},
		{
			Name: "Test Scenarios where there is an issue with connection to the database",
			ID:   1,
			SetupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(ErrConnectionFailed)
			},
			ExpectError: true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.Name, func(t *testing.T) {
			t.Logf("Running test case: %s", s.Name)

			db, mock, err := sqlmock.New()
			assert.Nil(err)
			gormDB, err := gorm.Open("mysql", db)
			assert.Nil(err)

			store := &ArticleStore{
				db: gormDB,
			}

			s.SetupMock(mock)

			article, err := store.GetByID(s.ID)

			if s.ExpectError {
				assert.NotNil(err, "Expected an error but did not get one")
			} else {
				assert.Nil(err, "Expected no error but got one")
				assert.NotNil(article, "Expected an article but got nil")
			}
		})
	}
}
