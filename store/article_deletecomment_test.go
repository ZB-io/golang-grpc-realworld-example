package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

// TestArticleStoreDeleteComment provides unit tests for function DeleteComment
func TestArticleStoreDeleteComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unable to instantiate mock DB: %v", err)
		return
	}
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("unable to use mock database for MimicDB data store: %v", err)
		return
	}

	articleStore := &ArticleStore{db: gormDB}

	// table-driven tests
	tests := []struct {
		inputComment    *model.Comment
		expectError     bool
		description     string
		setupExpectFunc func(comment *model.Comment)
	}{
		// scenario 1: Successful Deletion of Comment
		{
			inputComment: &model.Comment{Model: gorm.Model{ID: 1}},
			expectError:  false,
			description: "delete an existing comment",
			setupExpectFunc: func(comment *model.Comment) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "comments" WHERE "comments"."id" = ?`).WithArgs(comment.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected
				mock.ExpectCommit()
			},
		},
		// scenario 2: Deletion of Non-Existing Comment
		{
			inputComment: &model.Comment{Model: gorm.Model{ID: 42}},
			expectError:  true,
			description: "delete a non-existing comment",
			setupExpectFunc: func(comment *model.Comment) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "comments" WHERE "comments"."id" = ?`).WithArgs(comment.ID).
					WillReturnError(errors.New("record not found"))
				mock.ExpectRollback()
			},
		},
		// scenario 3: Deletion of Comment With Null Model
		{
			inputComment: nil,
			expectError:  true,
			description: "delete a comment with null model",
			setupExpectFunc: func(comment *model.Comment) {
				// no mock setup needed as nothing happens to the database
			},
		},
	}

	for _, test := range tests {
		// arrange
		test.setupExpectFunc(test.inputComment)

		// act
		err := articleStore.DeleteComment(test.inputComment)

		// assert
		if (err == nil) && test.expectError {
			t.Errorf("%s: Expected an error but got none", test.description)
		} else if (err != nil) && !test.expectError {
			t.Errorf("%s: Unexpected error %v", test.description, err)
		}

		// ensure all mock expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("%s: there were unfulfilled expectations: %v", test.description, err)
		}
	}
}
