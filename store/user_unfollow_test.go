package store

import (
	"testing"
	"errors"
	"reflect"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestUserStoreUnfollow(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(sqlmock.Sqlmock) error
		userA         *model.User
		userB         *model.User
		expectedError bool
	}{
		{
			name: "Unfollowing a User Successfully",
			setupFunc: func(mock sqlmock.Sqlmock) error {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM follows WHERE`).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				return nil
			},
			userA:         &model.User{Model: gorm.Model{ID: 1}, Username: "UserA"},
			userB:         &model.User{Model: gorm.Model{ID: 2}, Username: "UserB"},
			expectedError: false,
		},
		{
			name: "Unfollowing a User Not Followed",
			setupFunc: func(mock sqlmock.Sqlmock) error {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM follows WHERE`).WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
				return nil
			},
			userA:         &model.User{Model: gorm.Model{ID: 3}, Username: "UserA"},
			userB:         &model.User{Model: gorm.Model{ID: 4}, Username: "UserB"},
			expectedError: false,
		},
		{
			name: "Unfollowing with a Nonexistent User Entry",
			setupFunc: func(mock sqlmock.Sqlmock) error {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM follows WHERE`).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
				return nil
			},
			userA:         &model.User{Model: gorm.Model{ID: 5}, Username: "UserA"},
			userB:         &model.User{Model: gorm.Model{ID: 999}, Username: "NonExistentUser"},
			expectedError: true,
		},
		{
			name: "Database Error During Unfollow",
			setupFunc: func(mock sqlmock.Sqlmock) error {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM follows WHERE`).WillReturnError(gorm.ErrInvalidSQL)
				mock.ExpectRollback()
				return nil
			},
			userA:         &model.User{Model: gorm.Model{ID: 6}, Username: "UserA"},
			userB:         &model.User{Model: gorm.Model{ID: 7}, Username: "UserB"},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlmock", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a gorm connection", err)
			}

			store := &UserStore{db: gormDB}

			if err := tt.setupFunc(mock); err != nil {
				t.Fatalf("setupFunc failed: %s", err)
			}

			err = store.Unfollow(tt.userA, tt.userB)

			if (err != nil) != tt.expectedError {
				t.Errorf("unexpected error: got %v, want %v", err != nil, tt.expectedError)
			}

			if tt.expectedError {
				t.Log("expected error occurred during unfollow operation")
			} else {
				t.Log("unfollow procedure executed without errors")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}

}
