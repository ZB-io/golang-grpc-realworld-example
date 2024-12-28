package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
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
func TestArticleStoreDeleteComment(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, got error: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Failed to open gorm connection, got error: %v", err)
	}
	defer gormDB.Close()

	articleStore := &ArticleStore{db: gormDB}

	tests := []struct {
		name        string
		comment     *model.Comment
		prepare     func()
		expectedErr error
	}{
		{
			name: "Successfully Delete a Comment",
			comment: &model.Comment{
				Body:      "A valid comment",
				UserID:    1,
				ArticleID: 1,
			},
			prepare: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"comments\" WHERE .*").
					WithArgs(sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name: "Attempt to Delete a Non-Existent Comment",
			comment: &model.Comment{
				Body:      "A comment that does not exist",
				UserID:    2,
				ArticleID: 2,
			},
			prepare: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"comments\" WHERE .*").
					WithArgs(sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
			},
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Attempt to Delete a Comment with a Non-Existent ArticleID",
			comment: &model.Comment{
				Body:      "Valid comment with non-existent ArticleID",
				UserID:    1,
				ArticleID: 999,
			},
			prepare: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"comments\" WHERE .*").
					WithArgs(sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
			},
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Error due to Database Connection Issue",
			comment: &model.Comment{
				Body:      "Any comment",
				UserID:    1,
				ArticleID: 1,
			},
			prepare: func() {
				mock.ExpectBegin().WillReturnError(gorm.ErrInvalidTransaction)
			},
			expectedErr: gorm.ErrInvalidTransaction,
		},
		{
			name: "Deleting a Comment with Constraints",
			comment: &model.Comment{
				Body:      "Comment with constraints",
				UserID:    1,
				ArticleID: 1,
			},
			prepare: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"comments\" WHERE .*").
					WithArgs(sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			err := articleStore.DeleteComment(tt.comment)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				t.Logf("Expected error: %v", err)
			} else {
				assert.NoError(t, err)
				t.Log("Expected comment deletion successful")
			}
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
