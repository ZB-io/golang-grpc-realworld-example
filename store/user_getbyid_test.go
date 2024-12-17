package store

import (
	"testing"
	"errors"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

// TestUserStoreGetByID is a test function for the GetByID method of the UserStore struct
func TestUserStoreGetByID(t *testing.T) {
	db, mock, err := sqlmock.New() // creating instance of mocked DB
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening stub database connection", err)
	}

	gdb, err := gorm.Open("postgres", db) // opening connection using mocked DB
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	// setting up framework for response on query
	rows := sqlmock.NewRows([]string{"id", "Username", "Password", "Email"}).
		AddRow(1, "TestUser", "TestPassword", "test@email.com")

	s := UserStore{
        db: gdb,
    }

	t.Run("Successful user retrieval", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		user, err := s.GetByID(1)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if user.Email != "test@email.com" || user.Username != "TestUser" || user.Password != "TestPassword" {
			t.Log("Expected user data does not match the returned data")
		}
	})

	t.Run("Nonexistent user retrieval", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(gorm.ErrRecordNotFound)
		user, err := s.GetByID(100) // providing random non-existent ID
		if err == nil || err != gorm.ErrRecordNotFound {
			t.Errorf("Expected 'record not found' error, got %s", err)
		}
		if user != nil {
			t.Log("Expected 'nil' user for nonexistent id, got non-nil user")
		}
	})

	t.Run("Database connection failure", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(errors.New("database connection error"))
		user, err := s.GetByID(1)
		if err == nil || err.Error() != "database connection error" {
			t.Errorf("Expected 'database connection error', got %s", err)
		}
		if user != nil {
			t.Log("Expected 'nil' user for failed database connection, got non-nil user")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
