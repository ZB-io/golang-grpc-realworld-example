package store

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestUserStoreCreate tests the Create method assuming the external dependencies are imported correctly
func TestUserStoreCreate(t *testing.T) {
	tests := []struct {
		name       string
		user       model.User
		setupMock  func(mock sqlmock.Sqlmock)
		expectErr  bool
	}{
		{
			name: "Successful User Creation",
			user: model.User{
				Username: "unique_user",
				Email:    "unique@example.com",
				Password: "password",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name: "User Creation with Duplicate Username",
			user: model.User{
				Username: "duplicate_user",
				Email:    "unique2@example.com",
				Password: "password",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users"`).WillReturnError(sqlmock.ErrDeadlock) // Mocking error for duplication
				mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name: "User Creation with Null Email",
			user: model.User{
				Username: "null_email_user",
				Password: "password",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users"`).WillReturnError(sqlmock.ErrNoRows) // Use a more specific error if possible
				mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name: "User Creation with Database Connection Failure",
			user: model.User{
				Username: "db_fail_user",
				Email:    "db_fail@example.com",
				Password: "password",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(sqlmock.ErrConnDone) // Simulating database connection failure
			},
			expectErr: true,
		},
		{
			name: "User Creation with Maximum Field Lengths",
			user: model.User{
				Username: "maxfieldlengthsuser".Repeat(5), // Hypothetical max length
				Email:    "maxfieldlengthsuser@example.com",
				Password: "longenoughpassword",
				Bio:      "This is a bio with maximum length allowed for the field.",
				Image:    "http://example.com/image_max.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name: "User Creation with Special Characters",
			user: model.User{
				Username: "user@!#$%",
				Email:    "special@character.com",
				Password: "password",
				Bio:      "Test bio with special chars",
				Image:    "http://example.com/image.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error creating sqlmock: %v", err)
			}
			defer db.Close()

			gdb, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm db: %v", err)
			}

			store := &UserStore{db: gdb}
			tt.setupMock(mock)

			err = store.Create(&tt.user)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got none in scenario: %s", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v in scenario: %s", err, tt.name)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
