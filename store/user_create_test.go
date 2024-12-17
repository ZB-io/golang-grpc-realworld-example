package store_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserStoreCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Errorf("failed to open mock sql DB: %s", err)
	}

	gormDB, err := gorm.Open("mysql", db) // TODO: change "mysql" to your actual database
	if err != nil {
		fmt.Errorf("failed to open gorm DB: %s", err)
	}

	userStore := &store.UserStore{DB: gormDB}

	validUser := &model.User{Username: "TestUser", Email: "testuser@domain.com", Password: "secret"}

	t.Run("Successful User Creation", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO `users` (.+) VALUES (.+)").
			WithArgs(validUser.Username, validUser.Email, validUser.Password).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := userStore.Create(validUser)
		assert.NoError(t, err, "expected no error with successful user creation")
	})

	t.Run("Invalid User Creation - Unique Username Violation", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO `users` (.+) VALUES (.+)").
			WithArgs(validUser.Username, validUser.Email, validUser.Password).
			WillReturnError(fmt.Errorf("duplicate entry"))

		err := userStore.Create(validUser)
		assert.Error(t, err, "expected duplicate entry error with duplicate username")
	})

	t.Run("Invalid User Creation - Empty Fields", func(t *testing.T) {
		emptyUser := &model.User{}
		err := userStore.Create(emptyUser)
		assert.Error(t, err, "expected error with empty fields")
	})

	t.Run("DB Connection Error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO `users` (.+) VALUES (.+)").
			WithArgs(validUser.Username, validUser.Email, validUser.Password).
			WillReturnError(fmt.Errorf("database connection error"))

		err := userStore.Create(validUser)
		assert.Error(t, err, "expected database connection error")
	})

	t.Run("Panic Handling", func(t *testing.T) {
		// We use a defer function to recover from the panic and assert that it occurred.
		// If everything goes well, the code inside the defer function will execute after the panic, 
		// catching the panic and preventing it from crashing our program.
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("The code panicked during execution.")
			}
		}()

		err := userStore.Create(nil)
		assert.Error(t, err, "expected panic with nil user")
	})
} 
