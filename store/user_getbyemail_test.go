package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestUserStoreGetByEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		dbMock   func(mock sqlmock.Sqlmock)
		wantErr  bool
	}{
		{
			name: "Valid Email Test",
			email: "test@example.com",
			dbMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"username", "email", "password", "bio", "image"}).
					AddRow("testUser", "test@example.com", "password", "bio", "image")
				mock.ExpectQuery("^SELECT (.+) FROM (.+) WHERE (.+)$").WithArgs("test@example.com").WillReturnRows(rows)
			},
			wantErr: false,
		},
		// ... other test cases ...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			tt.dbMock(mock)

			gormDB, _ := gorm.Open("postgres", db)
			userStore := &UserStore{db: gormDB}

			user, err := userStore.GetByEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.GetByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if user != nil {
				// Now we are using the model.User structure
				if user.Email != tt.email {
					t.Errorf("Expected email %s, but got %s", tt.email, user.Email)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}
