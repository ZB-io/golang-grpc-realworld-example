package store

import (
	"testing"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/require"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserStoreUpdate(t *testing.T) {

	// Mocked DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gdb, err := gorm.Open("postgres", db)
	require.NoError(t, err)

	// Test Scenarios
	cases := []struct{
		name string
		user *model.User
		mock func()
		wantErr bool
	}{
		{
			"Updating a Valid User",
			&model.User{Username: "John", Email: "john@doe.com"},
			func(){
				mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			false,
		},
		{
			"Updating a User with Non-Unique Username",
			&model.User{Username: "John", Email: "john2@doe.com"},
			func(){
				mock.ExpectExec("UPDATE").WillReturnError(errors.New("pq: duplicate key value violates unique constraint"))
			},
			true,
		},
		{
			"Updating a User with Non-Unique Email",
			&model.User{Username: "John2", Email: "john@doe.com"},
			func(){
				mock.ExpectExec("UPDATE").WillReturnError(errors.New("pq: duplicate key value violates unique constraint"))
			},
			true,
		},
		{
			"Updating a User with No Username",
			&model.User{Email: "john3@doe.com"},
			func(){
				mock.ExpectExec("UPDATE").WillReturnError(errors.New("pq: null value in column violates not-null constraint"))
			},
			true,
		},
		{
			"Updating a User with No Email",
			&model.User{Username: "John3"},
			func(){
				mock.ExpectExec("UPDATE").WillReturnError(errors.New("pq: null value in column violates not-null constraint"))
			},
			true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			us := &UserStore{db: gdb}
			err := us.Update(tt.user)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
