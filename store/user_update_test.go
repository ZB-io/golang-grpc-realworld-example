package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"golang.org/x/sync/errgroup"
)

func TestUpdate(t *testing.T) {
	tests := []struct {
		name      string
		user      *model.User
		mockSetup func(sqlmock.Sqlmock)
		expectErr bool
		errType   error
	}{
		{
			name: "Successfully update a user in the database",
			user: &model.User{ID: 1, Username: "updated_user"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WithArgs("updated_user", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name: "Fail to update a user due to database error",
			user: &model.User{ID: 1, Username: "user_error"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WithArgs("user_error", 1).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			expectErr: true,
			errType:   errors.New("database error"),
		},
		{
			name: "Attempt to update a non-existent user",
			user: &model.User{ID: 999, Username: "non_existent"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WithArgs("non_existent", 999).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			expectErr: true,
			errType:   gorm.ErrRecordNotFound,
		},
		{
			name: "Attempt to update with invalid user data",
			user: &model.User{ID: 1, Username: ""},
			mockSetup: func(mock sqlmock.Sqlmock) {

			},
			expectErr: true,
			errType:   errors.New("validation error"),
		},
		{
			name: "Concurrent update scenarios",
			user: &model.User{ID: 1, Username: "concurrent_user"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WithArgs("concurrent_user", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Error initializing mock database: %v", err)
			}
			defer db.Close()

			gdb, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Error opening gorm db: %v", err)
			}
			defer gdb.Close()

			store := &UserStore{db: gdb}
			if tt.name == "Concurrent update scenarios" {
				var g errgroup.Group
				for i := 0; i < 5; i++ {
					g.Go(func() error {
						tt.mockSetup(mock)
						return store.Update(tt.user)
					})
				}

				if err := g.Wait(); err != nil && !tt.expectErr {
					t.Errorf("Unexpected error in concurrent test: %v", err)
				}
			} else {
				tt.mockSetup(mock)
				err := store.Update(tt.user)

				if !tt.expectErr && err != nil {
					t.Errorf("Expected no error, but got %v", err)
				}

				if tt.expectErr && (err == nil || !errors.Is(err, tt.errType)) {
					t.Errorf("Expected error type %v, but got %v", tt.errType, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}


