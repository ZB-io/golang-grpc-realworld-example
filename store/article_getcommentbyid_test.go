package store

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestArticleStoreGetCommentByID tests the GetCommentByID function of the ArticleStore struct
func TestArticleStoreGetCommentByID(t *testing.T) {
	t.Parallel()

	// Define test cases
	testCases := []struct {
		scenario       string
		dbMockCallback func(db sqlmock.Sqlmock)
		commentID      uint
		expectedError  bool
	}{
		{
			scenario: "Successful retrieval of comment by ID",
			dbMockCallback: func(db sqlmock.Sqlmock) {
				// TODO: Set up DB mock for successful retrieval of comment
			},
			commentID:     1,
			expectedError: false,
		},
		{
			scenario: "Incorrect ID provided",
			dbMockCallback: func(db sqlmock.Sqlmock) {
				// TODO: Set up DB mock for unsuccessful retrieval of comment
			},
			commentID:     999,
			expectedError: true,
		},
		{
			scenario: "Error in accessing the data",
			dbMockCallback: func(db sqlmock.Sqlmock) {
				// TODO: Set up DB mock to return an error
			},
			commentID:     1,
			expectedError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.scenario, func(t *testing.T) {
			// Mock DB
			mockDB, mock, _ := sqlmock.New()
			defer mockDB.Close()
			gormDB, _ := gorm.Open("sqlmock", mockDB)
			test.dbMockCallback(mock)

			// Create a store instance
			store := ArticleStore{db: gormDB}

			// Invoke the GetCommentByID function
			comment, err := store.GetCommentByID(test.commentID)

			// Assert function behaviour as per the scenario
			if test.expectedError {
				assert.Error(t, err)
				assert.Nil(t, comment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comment)
				assert.Equal(t, test.commentID, comment.ID)
			}
		})
	}
}
