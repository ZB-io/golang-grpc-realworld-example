package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)





func TestCreate(t *testing.T) {
	t.Run("Scenario 1: Successfully Create a New User", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO \"users\"").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		assert.NoError(t, err)

		s := &UserStore{db: gormDB}
		user := &model.User{Username: "testuser", Email: "test@example.com"}

		err = s.Create(user)
		assert.NoError(t, err)
		t.Log("Successfully created user with valid input, no database error.")
	})

	t.Run("Scenario 2: Fail to Create a User Due to Database Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO \"users\"").
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		assert.NoError(t, err)

		s := &UserStore{db: gormDB}
		user := &model.User{Username: "testuser", Email: "test@example.com"}

		err = s.Create(user)
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		t.Logf("Expected database error occurred: %v", err)
	})

	t.Run("Scenario 3: Fail to Create a User with Invalid Data", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		assert.NoError(t, err)

		s := &UserStore{db: gormDB}
		user := &model.User{Username: "", Email: ""}

		err = s.Create(user)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "username or email is empty")
		t.Logf("Invalid data should result in an error: %v", err)
	})

	t.Run("Scenario 4: Handle Nil User Input Gracefully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		assert.NoError(t, err)

		s := &UserStore{db: gormDB}

		err = s.Create(nil)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "user input is nil")
		t.Log("Function should gracefully handle nil user input and return an error.")
	})
}

