package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	go_sql_driver "github.com/DATA-DOG/go-sqlmock" // Imported for the sql.DB type used in sqlmock
)

// TestUserStoreUnfollow tests the Unfollow functionality of a UserStore.
func TestUserStoreUnfollow(t *testing.T) {
	// Define the test table
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
				// Arrange: Set User A follows User B in the mock database
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM follows`).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				return nil
			},
			userA: &model.User{Model: gorm.Model{ID: 1}, Username: "UserA"},
			userB: &model.User{Model: gorm.Model{ID: 2}, Username: "UserB"},
			expectedError: false,
		},
		{
			name: "Unfollowing a User Not Followed",
			setupFunc: func(mock sqlmock.Sqlmock) error {
				// Arrange: Ensure no existing follow relationship in the mock database
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM follows`).WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
				return nil
			},
			userA: &model.User{Model: gorm.Model{ID: 3}, Username: "UserA"},
			userB: &model.User{Model: gorm.Model{ID: 4}, Username: "UserB"},
			expectedError: false,
		},
		{
			name: "Unfollowing with a Nonexistent User Entry",
			setupFunc: func(mock sqlmock.Sqlmock) error {
				// Arrange: Scenario where User B does not exist in the database
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM follows`).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
				return nil
			},
			userA: &model.User{Model: gorm.Model{ID: 5}, Username: "UserA"},
			userB: &model.User{Model: gorm.Model{ID: 999}, Username: "NonExistentUser"},
			expectedError: true,
		},
		{
			name: "Database Error During Unfollow",
			setupFunc: func(mock sqlmock.Sqlmock) error {
				// Arrange: Simulate a database error
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM follows`).WillReturnError(gorm.ErrInvalidSQL)
				mock.ExpectRollback()
				return nil
			},
			userA: &model.User{Model: gorm.Model{ID: 6}, Username: "UserA"},
			userB: &model.User{Model: gorm.Model{ID: 7}, Username: "UserB"},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mock database
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

			// Setup the specific expectations according to the test case
			if err := tt.setupFunc(mock); err != nil {
				t.Fatalf("setupFunc failed: %s", err)
			}

			// Act
			err = store.Unfollow(tt.userA, tt.userB)

			// Assert
			if (err != nil) != tt.expectedError {
				t.Errorf("unexpected error: got %v, want %v", err != nil, tt.expectedError)
			}

			if tt.expectedError {
				t.Log("expected error occurred during unfollow operation")
			} else {
				t.Log("unfollow procedure executed without errors")
			}

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}

	// TODO: Add test case for concurrent unfollow requests; requires goroutine coordination.
}
