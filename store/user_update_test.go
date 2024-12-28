package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestUserStoreUpdate(t *testing.T) {
	// Initialize test scenarios
	tests := []struct {
		name          string
		setup         func(mock sqlmock.Sqlmock, user *model.User)
		user          *model.User
		expectError   bool
		expectedError error
	}{
		{
			name: "Successful User Update",
			setup: func(mock sqlmock.Sqlmock, user *model.User) {
				// Setup for successful update scenario
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE \"users\"").WithArgs(user.Username, user.Email, user.Password, user.Bio, user.Image, user.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			user: &model.User{ // Sample user for a successful update
				Model:    gorm.Model{ID: 1},
				Username: "updated_username",
				Email:    "updated_email@example.com",
				Password: "updated_password",
				Bio:      "updated_bio",
				Image:    "updated_image.jpg",
			},
			expectError: false,
		},
		{
			name: "Update Non-Existent User",
			setup: func(mock sqlmock.Sqlmock, user *model.User) {
				// Setup for updating a non-existent user
				mock.ExpectExec("UPDATE \"users\"").WithArgs(user.Username, user.Email, user.Password, user.Bio, user.Image, user.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			user: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "nonexistent_user",
			},
			expectError:   true,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Update with Partial Data",
			setup: func(mock sqlmock.Sqlmock, user *model.User) {
				// Setup for update with partial data
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE \"users\"").WithArgs(user.Username, user.Bio, user.Image, user.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "user1",
				Email:    "email@example.com",
				Bio:      "new_bio",
				Image:    "new_image.jpg",
			},
			expectError: false,
		},
		{
			name: "Update with Invalid Data",
			setup: func(mock sqlmock.Sqlmock, user *model.User) {
				// Setup for updating a user with data violating constraints
				mock.ExpectExec("UPDATE \"users\"").WithArgs(user.Email, user.ID).
					WillReturnError(errors.New("unique constraint violation"))
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
				Email: "duplicate_email@example.com",
			},
			expectError:   true,
			expectedError: errors.New("unique constraint violation"),
		},
		{
			name: "Update with Empty Fields",
			setup: func(mock sqlmock.Sqlmock, user *model.User) {
				// Setup for update where fields are set to empty strings
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE \"users\"").WithArgs(user.Username, user.Bio, user.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
				Username: "",
				Bio:      "",
			},
			expectError: false, // Based on business logic this might need to be adjusted
		},
		{
			name: "Database Connection Failure",
			setup: func(mock sqlmock.Sqlmock, user *model.User) {
				// Simulate a database connection error during update
				mock.ExpectExec("UPDATE \"users\"").WithArgs(user.ID).
					WillReturnError(errors.New("database connection failure"))
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			expectError:   true,
			expectedError: errors.New("database connection failure"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open a stub database connection: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("failed to initialize gorm DB: %s", err)
			}
			defer gormDB.Close()

			store := &UserStore{db: gormDB} // Assuming UserStore is correctly implemented

			tt.setup(mock, tt.user)

			err = store.Update(tt.user)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v but got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
