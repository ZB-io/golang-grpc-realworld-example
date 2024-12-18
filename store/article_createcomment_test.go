package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

// Testing store.CreateComment with different scenarios
func TestArcticleStoreCreateComment(t *testing.T) {
	// Define test cases
	testCases := []struct {
		desc        string
		comment     model.Comment
		mockError   error
		expectedErr error
	}{
		{
			desc: "Valid comment creation",
			comment: model.Comment{
				Body:      "This is a comment",
				UserID:    1,
				ArticleID: 1,
			},
			mockError:   nil,
			expectedErr: nil,
		},
		{
			desc: "Database error at comment creation",
			comment: model.Comment{
				Body:      "This is a comment",
				UserID:    2,
				ArticleID: 2,
			},
			mockError:   errors.New("Database error"),
			expectedErr: errors.New("Database error"),
		},
		{
			desc: "Invalid comment with missing fields",
			comment: model.Comment{
				Body:      "",
				UserID:    0,
				ArticleID: 0,
			},
			mockError:   nil,
			expectedErr: errors.New("Invalid comment"),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// Create mock SQL db
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Unexpected error in creating mock sql db: %v", err)
			}
			gdb, _ := gorm.Open("mysql", db)
			mock.ExpectBegin()
			mock.ExpectExec("CREATE").WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tC.mockError)
			mock.ExpectCommit()
			// Create new instance of store
			s := &ArticleStore{db: gdb}
			// Create new instance of model.Comment with provided parameters
			c := &tC.comment
			// Function call
			err = s.CreateComment(c)
			// Check for expected errors
			if tC.expectedErr == nil {
				assert.Nil(t, err, tC.desc+": Error should be nil but got: "+err.Error())
			} else {
				assert.NotNil(t, err, tC.desc+": Error should not be nil")
				assert.Equal(t, tC.expectedErr.Error(), err.Error(), tC.desc+": Error should be "+tC.expectedErr.Error()+" but got: "+err.Error())
			}
		})
	}
}
