package store

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

// TestUserStoreGetByEmail tests the GetByEmail function with various scenarios
func TestUserStoreGetByEmail(t *testing.T) {
	t.Parallel()

	// Set up the test cases
	type testCase struct {
		description       string
		email             string
		mockExpectations  func(sqlmock.Sqlmock)
		expectedUser      *model.User
		expectError       bool
		expectedErrorText string
	}

	tests := []testCase{
		{
			description: "Scenario 1: Successfully Retrieve User by Valid Email",
			email:       "valid@example.com",
			mockExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM `users` WHERE (email = ?) ORDER BY `users`.`id` ASC LIMIT 1").
					WithArgs("valid@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
						AddRow(1, "username", "valid@example.com", "hashedpassword", "bio", "image"))
			},
			expectedUser: &model.User{
				Username: "username",
				Email:    "valid@example.com",
				// other fields omitted for brevity
			},
			expectError: false,
		},
		{
			description: "Scenario 2: Return Error for Non-existent Email",
			email:       "nonexistent@example.com",
			mockExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM `users` WHERE (email = ?) ORDER BY `users`.`id` ASC LIMIT 1").
					WithArgs("nonexistent@example.com").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:      nil,
			expectError:       true,
			expectedErrorText: gorm.ErrRecordNotFound.Error(),
		},
		{
			description: "Scenario 3: Handle Database Error Softly",
			email:       "valid@example.com",
			mockExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM `users` WHERE (email = ?) ORDER BY `users`.`id` ASC LIMIT 1").
					WithArgs("valid@example.com").
					WillReturnError(gorm.ErrInvalidSQL)
			},
			expectedUser:      nil,
			expectError:       true,
			expectedErrorText: gorm.ErrInvalidSQL.Error(),
		},
		{
			description: "Scenario 4: Input Is an Empty String",
			email:       "",
			mockExpectations: func(mock sqlmock.Sqlmock) {
				// TODO: Handle this case according to expected application logic (e.g., return specific validation error)
				mock.ExpectQuery("SELECT * FROM `users` WHERE (email = ?) ORDER BY `users`.`id` ASC LIMIT 1").
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:      nil,
			expectError:       true,
			expectedErrorText: gorm.ErrRecordNotFound.Error(),
		},
		{
			description: "Scenario 5: Multiple Users with Same Email (Constraint Violation)",
			email:       "duplicate@example.com",
			mockExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM `users` WHERE (email = ?) ORDER BY `users`.`id` ASC LIMIT 1").
					WithArgs("duplicate@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
						AddRow(1, "username1", "duplicate@example.com", "hashedpassword1", "bio", "image").
						AddRow(2, "username2", "duplicate@example.com", "hashedpassword2", "bio", "image"))
			},
			expectedUser:      nil,
			expectError:       true,
			expectedErrorText: "multiple rows returned", // TODO: Adjust based on actual error handling implementation
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			// Setting up sqlmock
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create sqlmock: %v", err)
			}
			defer db.Close()

			// Gorm DB Init using sqlmock
			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to open gorm DB: %v", err)
			}
			defer gormDB.Close()

			// Mock expectations
			tc.mockExpectations(mock)

			// Create a UserStore instance with the mocked DB
			store := &UserStore{db: gormDB}

			// Invoke the function
			result, err := store.GetByEmail(tc.email)

			// Check the results
			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				} else if err.Error() != tc.expectedErrorText {
					t.Fatalf("Expected error text: %v, got: %v", tc.expectedErrorText, err.Error())
				}
				if result != nil {
					t.Fatalf("Expected nil user but got: %+v", *result)
				}
			} else {
				if err != nil {
					t.Fatalf("Did not expect error but got: %v", err)
				}
				if result == nil || result.Email != tc.expectedUser.Email {
					t.Fatalf("Expected user: %v, got: %v", tc.expectedUser, result)
				}
			}

			// Check that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %v", err)
			}
		})
	}
}
