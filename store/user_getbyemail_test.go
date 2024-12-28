package store

import (
	"errors"
	"log"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestUserStoreGetByEmail(t *testing.T) {

	type TestData struct {
		email         string
		expectedUser  *model.User
		expectedError error
		setupMocks    func(sqlmock.Sqlmock)
		scenario      string
	}

	userID := 1
	user := model.User{
		ID:    uint(userID),
		Email: "test@example.com",
		Name:  "Test User",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		log.Fatalf("Failed to open the gorm.DB: %v", err)
	}
	defer gormDB.Close()

	s := &UserStore{db: gormDB}

	tests := []TestData{
		{
			email:         "test@example.com",
			expectedUser:  &user,
			expectedError: nil,
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "name"}).
					AddRow(user.ID, user.Email, user.Name)
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE (.+)`).
					WithArgs(user.Email).
					WillReturnRows(rows)
			},
			scenario: "Successfully Retrieve a User by Email",
		},
		{
			email:         "notfound@example.com",
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE (.+)`).
					WithArgs("notfound@example.com").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			scenario: "User Not Found",
		},
		{
			email:         "test@example.com",
			expectedUser:  nil,
			expectedError: errors.New("db connection error"),
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE (.+)`).
					WithArgs(user.Email).
					WillReturnError(errors.New("db connection error"))
			},
			scenario: "Database Connection Error",
		},
		{
			email:         "",
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE (.+)`).
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			scenario: "Empty Email Input",
		},
		{
			email:         "invalid-email-format",
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE (.+)`).
					WithArgs("invalid-email-format").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			scenario: "Invalid Email Format",
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.setupMocks(mock)

			result, err := s.GetByEmail(test.email)

			if test.expectedError != nil {
				if err == nil || err.Error() != test.expectedError.Error() {
					t.Errorf("Expected error %v, but got %v", test.expectedError, err)
				}
				t.Logf("Scenario: %s succeeded: Correct error returned.", test.scenario)
			} else if err != nil {
				t.Errorf("Unexpected error occurred: %v", err)
			}

			if test.expectedUser != nil {
				if *result != *test.expectedUser {
					t.Errorf("Expected user %v, but got %v", *test.expectedUser, *result)
				}
				t.Logf("Scenario: %s succeeded: Correct user returned.", test.scenario)
			} else if result != nil {
				t.Errorf("Expected user to be nil, but got %v", result)
			}
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %s", err)
	}
}


