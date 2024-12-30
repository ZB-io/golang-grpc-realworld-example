package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestUserStoreGetByID(t *testing.T) {

	tests := []struct {
		name         string
		id           uint
		mockSetup    func(mock sqlmock.Sqlmock)
		expectedUser *model.User
		expectError  bool
		errorMessage string
	}{
		{
			name: "Successfully Retrieve User by ID",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {

				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(1, "testuser", "testuser@example.com", "password", "test bio", "test image")
				mock.ExpectQuery("SELECT * FROM \"users\" WHERE (.*)").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "password",
				Bio:      "test bio",
				Image:    "test image",
			},
			expectError:  false,
			errorMessage: "",
		},
		{
			name: "Fail to Retrieve User Due to Non-Existent ID",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery("SELECT * FROM \"users\" WHERE (.*)").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser: nil,
			expectError:  true,
			errorMessage: "record not found",
		},
		{
			name: "Database Error During User Retrieval",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery("SELECT * FROM \"users\" WHERE (.*)").
					WithArgs(1).
					WillReturnError(errors.New("database connection error"))
			},
			expectedUser: nil,
			expectError:  true,
			errorMessage: "database connection error",
		},
		{
			name: "Invalid ID Handling (Zero ID)",
			id:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {

			},
			expectedUser: nil,
			expectError:  true,
			errorMessage: "invalid parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("failed to open gorm db, %v", err)
			}

			userStore := &UserStore{db: gormDB}

			tt.mockSetup(mock)

			user, err := userStore.GetByID(tt.id)

			if tt.expectError {
				assert.Nil(t, user, "expected user to be nil")
				assert.Error(t, err)
				assert.EqualError(t, err, tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Logf("Test case '%s' executed successfully", tt.name)
		})
	}
}
