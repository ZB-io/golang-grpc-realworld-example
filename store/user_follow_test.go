package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock" // Required to mock database interactions
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

// TestUserStoreFollow tests the Follow function of the UserStore
func TestUserStoreFollow(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(mock sqlmock.Sqlmock)
		userA        *model.User
		userB        *model.User
		expectedErr  error
		expectedFollows []model.User
	}{
		{
			name: "Successfully Follow a User",
			setupMock: func(mock sqlmock.Sqlmock) {
				// Simulating successful database operation
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO follows`).WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			userA: &model.User{Username: "UserA"},
			userB: &model.User{Username: "UserB"},
			expectedErr: nil,
			expectedFollows: []model.User{{Username: "UserB"}},
		},
		{
			name: "Attempt to Follow a Non-Existent User",
			setupMock: func(mock sqlmock.Sqlmock) {
				// Simulating an attempt to follow non-existent user
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO follows`).WithArgs().WillReturnError(errors.New("foreign key constraint fails"))
				mock.ExpectCommit()
			},
			userA: &model.User{Username: "UserA"},
			userB: &model.User{Username: "NonExistentUser"},
			expectedErr: errors.New("foreign key constraint fails"),
			expectedFollows: nil,
		},
		{
			name: "Attempt to Follow a User Already Being Followed",
			setupMock: func(mock sqlmock.Sqlmock) {
				// Simulating the scenario of already followed user
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO follows`).WithArgs().WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
			},
			userA: &model.User{Username: "UserA", Follows: []model.User{{Username: "UserB"}}},
			userB: &model.User{Username: "UserB"},
			expectedErr: nil, // Append should handle idempotency
			expectedFollows: []model.User{{Username: "UserB"}},
		},
		{
			name: "Handle Database Connection Failure",
			setupMock: func(mock sqlmock.Sqlmock) {
				// Simulating database connection failure
				mock.ExpectBegin().WillReturnError(errors.New("connection error"))
			},
			userA: &model.User{Username: "UserA"},
			userB: &model.User{Username: "UserB"},
			expectedErr: errors.New("connection error"),
			expectedFollows: nil,
		},
		{
			name: "Self-Follow Attempt",
			setupMock: func(mock sqlmock.Sqlmock) {
				// Simulating self-follow attempt
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO follows`).WithArgs().WillReturnError(errors.New("self-follow not allowed"))
				mock.ExpectCommit()
			},
			userA: &model.User{Username: "UserA"},
			userB: &model.User{Username: "UserA"},
			expectedErr: errors.New("self-follow not allowed"),
			expectedFollows: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New() // Mock DB setup
			if err != nil {
				t.Fatalf("failed to open sqlmock database: %v", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			gormDb, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("failed to open gorm db: %v", err)
			}

			store := UserStore{db: gormDb}

			// Act
			err = store.Follow(tt.userA, tt.userB)
			if tt.expectedErr != nil && err == nil {
				t.Errorf("expected error: %v, got none", tt.expectedErr)
			} else if tt.expectedErr == nil && err != nil {
				t.Errorf("expected no error, got: %v", err)
			} else if tt.expectedErr != nil && err != nil && tt.expectedErr.Error() != err.Error() {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			// Verify follow list if applicable
			if tt.expectedErr == nil && len(tt.expectedFollows) > 0 {
				// TODO: Add logic to verify the state of the follow list after operation, potentially through a mock user retrieval function
				t.Log("Expected: ", tt.expectedFollows, "Actual follows: ", tt.userA.Follows)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
