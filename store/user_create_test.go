package store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
)








func TestUserStoreCreate(t *testing.T) {
	testCases := []struct {
		name           string
		user           model.User
		mockDbResponse error
		expectedError  error
	}{
		{
			name: "Create a New User with Valid Data",
			user: model.User{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "password",
			},
			mockDbResponse: nil,
			expectedError:  nil,
		},
		{
			name: "Create a New User with Existing Email",
			user: model.User{
				Username: "testuser2",
				Email:    "testuser@example.com",
				Password: "password",
			},
			mockDbResponse: gorm.ErrRecordNotFound,
			expectedError:  errors.New("email already in use"),
		},
		{
			name: "Create a New User with Existing Username",
			user: model.User{
				Username: "testuser",
				Email:    "testuser2@example.com",
				Password: "password",
			},
			mockDbResponse: gorm.ErrRecordNotFound,
			expectedError:  errors.New("username already in use"),
		},
		{
			name: "Create a New User with Missing Required Data",
			user: model.User{
				Username: "",
				Email:    "testuser3@example.com",
				Password: "password",
			},
			mockDbResponse: gorm.ErrRecordNotFound,
			expectedError:  errors.New("required data is missing"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open mock sql db: %v", err)
			}
			gdb, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("failed to open gorm db: %v", err)
			}

			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO \"users\"").WillReturnError(tc.mockDbResponse)
			mock.ExpectCommit()

			us := &UserStore{db: gdb}

			err = us.Create(&tc.user)

			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Errorf("expected error %v, but got %v", tc.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
		})
	}
}
