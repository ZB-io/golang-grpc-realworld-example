package store

import (
	"testing"
	"database/sql/driver"
	"errors"
	"sync"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/jinzhu/gorm"
)

// TestUserStoreGetByID tests the GetByID function of UserStore
func TestUserStoreGetByID(t *testing.T) {
	// Create a mock database and the store
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %s", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Error creating GORM DB: %s", err)
	}

	store := &UserStore{db: gormDB}

	// Define test user data
	testUser := model.User{
		Model: gorm.Model{ID: 1},
		Username: "testuser",
		Email: "test@example.com",
		Password: "password",
		Bio: "test bio",
		Image: "https://example.com/image.png",
	}

	// Table-driven tests
	tests := []struct {
		name          string
		userID        uint
		setupMock     func()
		expectedUser  *model.User
		expectedError bool
	}{
		{
			name:   "Scenario 1: Retrieve Valid User by ID",
			userID: 1,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(testUser.ID, testUser.Username, testUser.Email, testUser.Password, testUser.Bio, testUser.Image)
				mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND (("users"."id" = ?))`).WithArgs(testUser.ID).
					WillReturnRows(rows)
			},
			expectedUser:  &testUser,
			expectedError: false,
		},
		{
			name:   "Scenario 2: Handle Non-Existent User ID",
			userID: 999,
			setupMock: func() {
				mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND (("users"."id" = ?))`).WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: true,
		},
		{
			name:   "Scenario 3: Database Connection Error Handling",
			userID: 1,
			setupMock: func() {
				mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND (("users"."id" = ?))`).WithArgs(1).
					WillReturnError(errors.New("connection error"))
			},
			expectedUser:  nil,
			expectedError: true,
		},
		{
			name:   "Scenario 4: Handle Empty or Invalid ID Input",
			userID: 0,
			setupMock: func() {
				mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND (("users"."id" = ?))`).WithArgs(0).
					WillReturnError(errors.New("invalid input"))
			},
			expectedUser:  nil,
			expectedError: true,
		},
		{
			name:   "Scenario 5: Simultaneous Requests and Thread Safety",
			userID: 1,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(testUser.ID, testUser.Username, testUser.Email, testUser.Password, testUser.Bio, testUser.Image)
				mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND (("users"."id" = ?))`).WithArgs(testUser.ID).
					WillReturnRows(rows).WillReturnRows(rows)
			},
			expectedUser:  &testUser,
			expectedError: false,
		},
	}
  
	// Iterate through each test case
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mocks specific to the test
			tc.setupMock()

			// If the scenario requires, simulate concurrent requests 
			if tc.name == "Scenario 5: Simultaneous Requests and Thread Safety" {
				var wg sync.WaitGroup
				successCount := 0
				const numGoroutines = 3

				for i := 0; i < numGoroutines; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						user, err := store.GetByID(tc.userID)
						if err == nil && user.ID == tc.userID {
							successCount++
						}
					}()
				}

				wg.Wait()
				if successCount != numGoroutines {
					t.Errorf("Expected %d successful retrievals, got %d", numGoroutines, successCount)
				}
			} else {
				// Act
				user, err := store.GetByID(tc.userID)

				// Assert
				if tc.expectedError {
					if err == nil {
						t.Errorf("Expected error, got nil")
					} else {
						t.Logf("Expected error received: %v", err)
					}
				} else {
					if err != nil {
						t.Errorf("Unexpected error: %v", err)
					}
					if user.ID != tc.expectedUser.ID {
						t.Errorf("Expected user ID %v, got %v", tc.expectedUser.ID, user.ID)
					}
				}
			}

			// Ensure all expectations are met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}
