package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreUpdate(t *testing.T) {
	// establish an array of anonymous struct objects that hold the test case scenarios to be tested
	testCases := []struct {
		name     string
		user     *model.User
		dbError  error
		expected error
	}{
		{
			name: "Successful User Update",
			user: &model.User{
				Model:     gorm.Model{ID: 1},
				Username:  "testuser",
				Email:     "testuser@gmail.com",
				Password:  "123456",
				Bio:       "test bio",
				Image:     "test.jpg",
			},
			dbError:  nil,
			expected: nil,
		},
		{
			name:    "Unsuccessful User Update due to DB Error",
			user:    &model.User{},
			dbError: errors.New("db error"),
			expected: errors.New("db error"),
		},
		{
			name:     "Unsuccessful User Update due to Empty User",
			user:     nil,
			dbError:  nil,
			expected: errors.New("User object cannot be nil"),
		},
		{
			name: "Unsuccessful User Update due to Invalid User Field(s)",
			user: &model.User{
				Username: "invaliduser",
			},
			dbError:  errors.New("invalid user"),
			expected: errors.New("invalid user"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was unexpected when opening a stub database connection", err)
			}

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("an error '%s' was unexpected when opening gorm database", err)
			}

			defer gormDB.Close()

			// mock the expectation based on the test case
			mock.ExpectBegin()
			mock.ExpectExec("UPDATE").WillReturnError(tc.dbError)
			mock.ExpectCommit()

			// create the user store with the mock db 
			userStore := &UserStore{db: gormDB}

			// call the Update function with the user object from current test case
			err = userStore.Update(tc.user)

			if tc.user == nil {
				assert.NotNil(t, err, "User object cannot be nil")
			} else if tc.expected == nil {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tc.expected.Error())
			}
		})
	}
}
