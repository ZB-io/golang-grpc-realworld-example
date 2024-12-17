package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock sql DB: %s", err)
	}
	gormDB, _ := gorm.Open("postgres", db)

	s := &UserStore{db: gormDB}

	defer func() {
		db.Close()
		gormDB.Close()
	}()

	t.Run("Successful User Creation", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		assert.NoError(t, s.Create(&model.User{Username: "test", Email: "test@test.com", Password: "test", Bio: "test", Image: "test"}))
	})

	t.Run("Failed user creation due to non-unique username", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO").WillReturnError(&gorm.DB{Error: errors.New("Unique field duplicate")})
		mock.ExpectRollback()
		assert.Error(t, s.Create(&model.User{Username: "test", Email: "different@test.com", Password: "test", Bio: "test", Image: "test"}))
	})

	t.Run("Failed user creation due to non-unique email", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO").WillReturnError(&gorm.DB{Error: errors.New("Unique field duplicate")})
		mock.ExpectRollback()
		assert.Error(t, s.Create(&model.User{Username: "different", Email: "test@test.com", Password: "test", Bio: "test", Image: "test"}))
	})

	t.Run("Failed user creation due to validation errors in user data", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO").WillReturnError(&gorm.DB{Error: errors.New("Validation error")})
		mock.ExpectRollback()
		assert.Error(t, s.Create(&model.User{Username: "", Email: "", Password: "test", Bio: "test", Image: "test"}))
	})

	t.Run("Failed user creation due to database connection issues", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(&gorm.DB{Error: errors.New("DB connection error")})
		assert.Error(t, s.Create(&model.User{Username: "anotherTest", Email: "anotherTest@test.com", Password: "test", Bio: "test", Image: "test"}))
	})
}
