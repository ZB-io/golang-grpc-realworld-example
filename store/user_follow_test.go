package store

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"testing"
)

func TestUserStoreFollow(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name          string
		userA         *model.User
		userB         *model.User
		expectedError error
	}{
		{
			name:          "Successful Follow",
			userA:         &model.User{Username: "UserA", Email: "usera@example.com"},
			userB:         &model.User{Username: "UserB", Email: "userb@example.com"},
			expectedError: nil,
		},
		{
			name:          "Follow function fails",
			userA:         &model.User{Username: "UserA", Email: "usera@example.com"},
			userB:         &model.User{Username: "UserB", Email: "userb@example.com"},
			expectedError: errors.New("association error"),
		},
		{
			name:          "User tries to follow themselves",
			userA:         &model.User{Username: "UserA", Email: "usera@example.com"},
			userB:         &model.User{Username: "UserA", Email: "usera@example.com"},
			expectedError: errors.New("a user can't follow themselves"),
		},
		{
			name:          "User tries to follow a user they already follow",
			userA:         &model.User{Username: "UserA", Email: "usera@example.com"},
			userB:         &model.User{Username: "UserB", Email: "userb@example.com"},
			expectedError: errors.New("a user can't follow a user they already follow"),
		},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock database connection
			db, _, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			// Create a new Gorm DB connection
			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a Gorm database connection", err)
			}

			// Initialize UserStore with mocked DB connection
			userStore := &UserStore{gormDB}

			// TODO: Mock the behavior of Association and Append functions based on the test case

			// Call the Follow function
			err = userStore.Follow(tc.userA, tc.userB)

			// Validate the results
			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Fatalf("expected error '%s', but got '%v'", tc.expectedError, err)
				}
			} else if err != nil {
				t.Fatalf("expected no error, but got '%v'", err)
			}
		})
	}
}
