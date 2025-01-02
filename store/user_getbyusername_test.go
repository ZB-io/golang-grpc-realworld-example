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
	testCases := []struct {
		name     string
		username string
		mock     func(mock sqlmock.Sqlmock)
		wantErr  bool
	}{
		{
			name:     "User exists in the database",
			username: "testUser",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(1, "testUser", "testUser@example.com", "password", "bio", "image")
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").WithArgs("testUser").WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:     "User does not exist in the database",
			username: "nonExistentUser",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").WithArgs("nonExistentUser").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name:     "Database error",
			username: "anyUser",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").WithArgs("anyUser").
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name:     "Empty username",
			username: "",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").WithArgs("").
					WillReturnError(errors.New("empty username"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			tc.mock(mock)

			gdb, _ := gorm.Open("postgres", db)
			store := UserStore{gdb}

			user, err := store.GetByUsername(tc.username)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
				t.Log("Error is expected: ", err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				userResult, ok := user.(*model.User)
				assert.True(t, ok)
				assert.Equal(t, tc.username, userResult.Username)
				t.Log("No error is expected and user should not be nil: ", user)
			}
		})
	}
}
