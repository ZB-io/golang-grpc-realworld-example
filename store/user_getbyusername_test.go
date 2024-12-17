package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreGetByUsername(t *testing.T) {
	// Create table of test cases
	testCases := []struct {
		name        string
		username    string
		mock        func(mock sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			name:     "Test GetByUsername with valid username",
			username: "testuser",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows([]string{"id", "username"}).
					AddRow(1, "testuser")
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username = (.+)").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			expectedErr: nil,
		},
		{
			name:     "Test GetByUsername with non-existing username",
			username: "nonexistent",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username = (.+)").
					WithArgs("nonexistent").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name:        "Test GetByUsername with empty username",
			username:    "",
			mock:        func(mock sqlmock.Sqlmock) {},
			expectedErr: errors.New("username cannot be empty"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)

			// Create a mock DB
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			// Mock the database behavior
			tc.mock(mock)

			gdb, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening gorm database", err)
			}

			// TODO: Use correct package name and struct instantiation
			userStore := UserStore{db: gdb}

			// Run our function to test
			user, err := userStore.GetByUsername(tc.username)

			// Validate results
			if tc.expectedErr != nil {
				assert.Errorf(t, err, tc.expectedErr.Error())
			} else {
				assert.Nil(t, err)
			}

			if len(tc.username) > 0 {
				assert.NotNil(t, user)
				assert.Equal(t, tc.username, user.Username)
			} else {
				assert.Nil(t, user)
			}
		})
	}
}
