package store

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"testing"
)





func TestUserStoreUpdate(t *testing.T) {

	var tests = []struct {
		name     string
		user     *model.User
		mockFunc func(mock sqlmock.Sqlmock, user *model.User)
		wantErr  bool
	}{
		{
			name: "Successful User Update",
			user: &model.User{Username: "test", Email: "test@example.com", Password: "password", Bio: "bio", Image: "image"},
			mockFunc: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectExec("^UPDATE (.+) SET (.+)$").WithArgs(user.Username, user.Email, user.Password, user.Bio, user.Image).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "User Update with Nonexistent User",
			user: &model.User{Username: "nonexistent", Email: "nonexistent@example.com", Password: "password", Bio: "bio", Image: "image"},
			mockFunc: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectExec("^UPDATE (.+) SET (.+)$").WithArgs(user.Username, user.Email, user.Password, user.Bio, user.Image).WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "User Update with Invalid User Data",
			user: &model.User{Username: "", Email: "", Password: "", Bio: "", Image: ""},
			mockFunc: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectExec("^UPDATE (.+) SET (.+)$").WithArgs(user.Username, user.Email, user.Password, user.Bio, user.Image).WillReturnError(errors.New("invalid data"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open mock sql db: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("failed to open gorm db: %v", err)
			}

			tt.mockFunc(mock, tt.user)

			store := &UserStore{db: gormDB}
			err = store.Update(tt.user)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

