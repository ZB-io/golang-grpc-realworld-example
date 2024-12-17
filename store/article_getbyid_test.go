package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreGetByID(t *testing.T) {
	// define table driven tests
	tests := []struct {
		name       string
		setupMocks func(mock sqlmock.Sqlmock)
		inputID    interface{}
		output     *model.Article
		err        error
	}{
		{
			name: "Successful Retrieval of an Article by ID",
			setupMocks: func(mock sqlmock.Sqlmock) {
				articleRows := sqlmock.NewRows([]string{"ID", "Author", "Tags"}).
					AddRow(1, "Author1", "Tag1")

				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE \"articles\".\"id\" = $").
					WithArgs(1).
					WillReturnRows(articleRows)
			},
			inputID: 1,
			output: &model.Article{
				ID:     1,
				Author: "Author1",
				Tags:   []string{"Tag1"},
			},
			err: nil,
		},
		{
			name: "Unsuccessful Attempt to Retrieve an Article with a Nonexistent ID",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE \"articles\".\"id\" = $").
					WithArgs(100).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			inputID: 100,
			output:  nil,
			err:     gorm.ErrRecordNotFound,
		},
		{
			name: "Failure to Retrieve an Article due to Database Error",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE \"articles\".\"id\" = $").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			inputID: 1,
			output:  nil,
			err:     errors.New("database error"),
		},
	}

	for _, test := range tests {
		// log current scenario
		t.Log(test.name)

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to setup mock for test %v: %v", test.name, err)
		}

		// defer closing database
		defer db.Close()

		test.setupMocks(mock)

		gormDB, err := gorm.Open("postgres", db)
		if err != nil {
			t.Errorf("Failed to open mock db: %v", err)
		}

		articleStore := ArticleStore{db: gormDB}

		// invoke GetByID function
		article, err := articleStore.GetByID(test.inputID.(uint))

		if err == nil {
			// if no error, assert that the returned article matches the expected result
			assert.Equal(t, test.output.ID, article.ID)
			assert.Equal(t, test.output.Author, article.Author)
			assert.Equal(t, test.output.Tags, article.Tags)
		} else {
			// if an error occurs, assert the error message
			assert.Error(t, test.err)
		}

	}

}
