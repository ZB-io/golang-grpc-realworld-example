package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)



func TestUserStoreGetByEmail(t *testing.T) {

	tests := []struct {
		name     string
		email    string
		mock     func(mock sqlmock.Sqlmock)
		wantUser *model.User
		wantErr  bool
	}{
		{
			name:  "User Found with Given Email",
			email: "test@example.com",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email"}).
					AddRow(1, "test@example.com")
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").WithArgs("test@example.com").WillReturnRows(rows)
			},
			wantUser: &model.User{Model: gorm.Model{ID: 1}, Email: "test@example.com"},
			wantErr:  false,
		},
		{
			name:  "User Not Found with Given Email",
			email: "nonexistent@example.com",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").WithArgs("nonexistent@example.com").WillReturnError(gorm.ErrRecordNotFound)
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name:  "Database Error",
			email: "error@example.com",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").WithArgs("error@example.com").WillReturnError(errors.New("database error"))
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name:  "Empty Email",
			email: "",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").WithArgs("").WillReturnError(errors.New("invalid input"))
			},
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			gormDB, _ := gorm.Open("postgres", db)
			defer func() {
				gormDB.Close()
				db.Close()
			}()

			tt.mock(mock)

			us := UserStore{db: gormDB}
			user, err := us.GetByEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantUser, user)
		})
	}
}
