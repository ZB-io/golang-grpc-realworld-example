package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

// TestArticleStoreGetCommentByID tests the GetCommentByID function of the ArticleStore struct
func TestArticleStoreGetCommentByID(t *testing.T) {
	// Create a slice of test cases
	testCases := []struct {
		name         string
		mock         func(mock sqlmock.Sqlmock)
		id           uint
		expectedErr  error
		expectedResp *model.Comment
	}{
		{
			name: "Successful retrieval of a comment by ID",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "body", "user_id", "article_id"}).
					AddRow(1, "Test Comment", 1, 1)
				mock.ExpectQuery("^SELECT (.+) FROM \"comments\" WHERE \"comments\".\"deleted_at\" IS NULL AND ((id = $1))$").
					WithArgs(1).
					WillReturnRows(rows)
			},
			id:           1,
			expectedErr:  nil,
			expectedResp: &model.Comment{Model: gorm.Model{ID: 1}, Body: "Test Comment", UserID: 1, ArticleID: 1},
		},
		{
			name: "Attempt to retrieve a comment using an ID that does not exist in the DB",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"comments\" WHERE \"comments\".\"deleted_at\" IS NULL AND ((id = $1))$").
					WithArgs(2).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			id:           2,
			expectedErr:  gorm.ErrRecordNotFound,
			expectedResp: nil,
		},
		{
			name: "Database error during comment retrieval",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"comments\" WHERE \"comments\".\"deleted_at\" IS NULL AND ((id = $1))$").
					WithArgs(3).
					WillReturnError(errors.New("database error"))
			},
			id:           3,
			expectedErr:  errors.New("database error"),
			expectedResp: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			gdb, err := gorm.Open("postgres", db) // TODO: Replace "postgres" with your actual database
			assert.NoError(t, err)

			tc.mock(mock)

			store := ArticleStore{db: gdb}
			resp, err := store.GetCommentByID(tc.id)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedResp, resp)
		})
	}
}
