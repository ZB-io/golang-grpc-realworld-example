package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

// TestUserStoreCreate performs unit testing on the Create method of UserStore 
func TestUserStoreCreate(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		mock    func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Successful User Creation",
			user: &model.User{Username: "testuser", Email: "testuser@mail.com", Password: "password"},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users").WithArgs("testuser", "testuser@mail.com", "password").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "User Creation with Existing Username",
			user: &model.User{Username: "testuser", Email: "otheruser@mail.com", Password: "password"},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users").WithArgs("testuser", "otheruser@mail.com", "password").WillReturnError(errors.New("username already exists"))
			},
			wantErr: true,
		},
		{
			name: "User Creation with Existing Email",
			user: &model.User{Username: "otheruser", Email: "testuser@mail.com", Password: "password"},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users").WithArgs("otheruser", "testuser@mail.com", "password").WillReturnError(errors.New("email already exists"))
			},
			wantErr: true,
		},
		{
			name: "User Creation with Invalid Data",
			user: &model.User{Username: "", Email: "", Password: ""},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users").WithArgs("", "", "").WillReturnError(errors.New("invalid data"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			tt.mock(mock)
			gdb, _ := gorm.Open("postgres", db)

			userStore := &UserStore{db: gdb}
			err := userStore.Create(tt.user)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				t.Log("Error received as expected")
			} else {
				if err != nil {
					t.Errorf("Expected no error but got %v", err)
				}
				t.Log("No error received as expected")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
