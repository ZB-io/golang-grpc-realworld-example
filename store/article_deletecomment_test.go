package store

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"testing"
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

func TestArticleStoreDeleteComment(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDB, _ := gorm.Open("mysql", db)
	defer func() {
		_ = gormDB.Close()
	}()

	store := &ArticleStore{db: gormDB}

	scenarios := []struct {
		desc      string
		comment   *model.Comment
		mockFunc  func(comment *model.Comment)
		expectErr error
	}{
		{
			desc:    "Successful Deletion",
			comment: &model.Comment{Body: "test comment", UserID: 1, ArticleID: 1},
			mockFunc: func(comment *model.Comment) {
				mock.ExpectBegin()
				mock.ExpectExec("delete").WithArgs(comment.ArticleID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: nil,
		},
		{
			desc:    "Deleting Non-Existing Comment",
			comment: &model.Comment{Body: "missing comment", UserID: 2, ArticleID: 2},
			mockFunc: func(comment *model.Comment) {
				mock.ExpectBegin()
				mock.ExpectExec("delete").WithArgs(comment.ArticleID).WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
			},
			expectErr: errors.New("record not found"),
		},
		{
			desc:      "Null Comment Deletion",
			comment:   nil,
			mockFunc:  func(comment *model.Comment) {},
			expectErr: errors.New("comment cannot be nil"),
		},
		{
			desc:    "Comment Deletion with Failed Database Connection",
			comment: &model.Comment{Body: "test comment", UserID: 3, ArticleID: 3},
			mockFunc: func(comment *model.Comment) {
				mock.ExpectBegin()
				mock.ExpectExec("delete").WithArgs(comment.ArticleID).WillReturnError(errors.New("failed connection"))
				mock.ExpectRollback()
			},
			expectErr: errors.New("failed connection"),
		},
	}

	for _, s := range scenarios {
		t.Run(s.desc, func(t *testing.T) {
			s.mockFunc(s.comment)
			err := store.DeleteComment(s.comment)

			if (err != nil && s.expectErr == nil) || (err == nil && s.expectErr != nil) || (err != nil && s.expectErr != nil && err.Error() != s.expectErr.Error()) {
				t.Errorf("Unexpected error: '%v', was expecting error: '%v'", err, s.expectErr)
			}
		})
	}
}
func (s *ArticleStore) DeleteComment(m *model.Comment) error {
	return s.db.Delete(m).Error
}
