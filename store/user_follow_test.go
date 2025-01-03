package store

import (
	"errors"
	"fmt"
	"log"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)




type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

func TestFollow(t *testing.T) {
	tests := []struct {
		name            string
		setupMock       func(mock sqlmock.Sqlmock)
		userA           *model.User
		userB           *model.User
		expectedErr     error
		expectedFollows []model.User
	}{
		{
			name: "Successfully Follow a User",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO follows`).WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			userA:           &model.User{Username: "UserA"},
			userB:           &model.User{Username: "UserB"},
			expectedErr:     nil,
			expectedFollows: []model.User{{Username: "UserB"}},
		},
		{
			name: "Attempt to Follow a Non-Existent User",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO follows`).WithArgs().WillReturnError(errors.New("foreign key constraint fails"))
				mock.ExpectRollback()
			},
			userA:           &model.User{Username: "UserA"},
			userB:           &model.User{Username: "NonExistentUser"},
			expectedErr:     errors.New("foreign key constraint fails"),
			expectedFollows: nil,
		},
		{
			name:            "Attempt to Follow a User Already Being Followed",
			setupMock:       func(mock sqlmock.Sqlmock) {},
			userA:           &model.User{Username: "UserA", Follows: []model.User{{Username: "UserB"}}},
			userB:           &model.User{Username: "UserB"},
			expectedErr:     nil,
			expectedFollows: []model.User{{Username: "UserB"}},
		},
		{
			name: "Handle Database Connection Failure",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("connection error"))
			},
			userA:           &model.User{Username: "UserA"},
			userB:           &model.User{Username: "UserB"},
			expectedErr:     errors.New("connection error"),
			expectedFollows: nil,
		},
		{
			name:            "Self-Follow Attempt",
			setupMock:       func(mock sqlmock.Sqlmock) {},
			userA:           &model.User{Username: "UserA"},
			userB:           &model.User{Username: "UserA"},
			expectedErr:     errors.New("self-follow not allowed"),
			expectedFollows: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatalf("failed to open sqlmock database: %v", err)
			}
			defer db.Close()

			gormDb, err := gorm.Open("mysql", db)
			if err != nil {
				log.Fatalf("failed to open gorm db: %v", err)
			}

			store := &UserStore{db: gormDb}

			tt.setupMock(mock)

			err = store.Follow(tt.userA, tt.userB)
			if (tt.expectedErr == nil && err != nil) || (tt.expectedErr != nil && err == nil) {
				t.Fatalf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			if tt.expectedErr == nil {
				mock.ExpectQuery(`SELECT \* FROM follows`).WithArgs(tt.userA.ID, tt.userB.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
func (s *UserStore) Follow(a *model.User, b *model.User) error {
	if a.ID == b.ID {
		return errors.New("self-follow not allowed")
	}
	for _, followedUser := range a.Follows {
		if followedUser.ID == b.ID {
			return nil
		}
	}
	return s.db.Model(a).Association("Follows").Append(b).Error
}
