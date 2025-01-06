package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreGetTags(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		mockDBFunc  func(db sqlmock.Sqlmock)
		expectedErr error
		expectedRes []model.Tag
	}{
		{
			name: "Normal operation with present tags",
			mockDBFunc: func(db sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "tag1").
					AddRow(2, "tag2")
				db.ExpectQuery("^SELECT (.+) FROM \"tags\"").WillReturnRows(rows)
			},
			expectedErr: nil,
			expectedRes: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "tag1"},
				{Model: gorm.Model{ID: 2}, Name: "tag2"},
			},
		},
		{
			name: "Normal operation with no tags",
			mockDBFunc: func(db sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"})
				db.ExpectQuery("^SELECT (.+) FROM \"tags\"").WillReturnRows(rows)
			},
			expectedErr: nil,
			expectedRes: []model.Tag{},
		},
		{
			name: "Error handling when database query fails",
			mockDBFunc: func(db sqlmock.Sqlmock) {
				db.ExpectQuery("^SELECT (.+) FROM \"tags\"").WillReturnError(errors.New("database error"))
			},
			expectedErr: errors.New("database error"),
			expectedRes: nil,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock the database
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			gormDB, err := gorm.Open("postgres", mockDB)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening gorm database", err)
			}

			// Prepare the mock database for the test case
			tc.mockDBFunc(mock)

			// Create an ArticleStore with the mock database
			store := &ArticleStore{db: gormDB}

			// Call the GetTags function and check the result
			tags, err := store.GetTags()
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, tags)
			}
		})
	}
}
