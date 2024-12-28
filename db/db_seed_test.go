package db

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/require"
	"github.com/raahii/golang-grpc-realworld-example/model"
)





func TestSeed(t *testing.T) {
	var mutex sync.Mutex
	var txdbInitialized bool

	tests := []struct {
		name          string
		setup         func() (*gorm.DB, func())
		expectedErr   error
		validate      func(*testing.T, *gorm.DB, error)
		usersTomlPath string
		concurrent    bool
	}{
		{
			name: "Successful Seeding of Users",
			setup: func() (*gorm.DB, func()) {
				db, mock, cleanup := setupMockDB(t)
				users := mockUsers()
				mock.ExpectBegin()
				for _, user := range users {
					mock.ExpectExec("INSERT INTO \"users\"").WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(1, 1))
				}
				mock.ExpectCommit()
				return db, cleanup
			},
			expectedErr: nil,
			validate: func(t *testing.T, db *gorm.DB, err error) {
				require.NoError(t, err)
				users := mockUsers()
				for _, u := range users {
					var user model.User

					db.Where("id = ?", u.ID).First(&user)
					require.Equal(t, u.ID, user.ID)
				}
			},
			usersTomlPath: "mock_valid_users.toml",
		},
		{
			name: "File Not Found Error",
			setup: func() (*gorm.DB, func()) {
				db, _, cleanup := setupMockDB(t)
				return db, cleanup
			},
			expectedErr: errors.New("open db/seed/users.toml: no such file or directory"),
			validate: func(t *testing.T, db *gorm.DB, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "no such file or directory")
			},
			usersTomlPath: "non_existent.toml",
		},
		{
			name: "Invalid TOML File Format",
			setup: func() (*gorm.DB, func()) {
				db, _, cleanup := setupMockDB(t)
				return db, cleanup
			},
			expectedErr: errors.New("expected error"),
			validate: func(t *testing.T, db *gorm.DB, err error) {
				require.Error(t, err)
			},
			usersTomlPath: "mock_invalid_users.toml",
		},
		{
			name: "Database Creation Error",
			setup: func() (*gorm.DB, func()) {
				db, mock, cleanup := setupMockDB(t)
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"users\"").WillReturnError(errors.New("db connection issue"))
				return db, cleanup
			},
			expectedErr: errors.New("db connection issue"),
			validate: func(t *testing.T, db *gorm.DB, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "db connection issue")
			},
			usersTomlPath: "mock_valid_users.toml",
		},
		{
			name: "Empty TOML File",
			setup: func() (*gorm.DB, func()) {
				db, _, cleanup := setupMockDB(t)
				return db, cleanup
			},
			expectedErr: nil,
			validate: func(t *testing.T, db *gorm.DB, err error) {
				require.NoError(t, err)
			},
			usersTomlPath: "mock_empty_users.toml",
		},
		{
			name: "Concurrent Access for Initialization",
			setup: func() (*gorm.DB, func()) {

				db, mock, cleanup := setupMockDB(t)
				mock.ExpectBegin()
				users := mockUsers()
				for _, user := range users {
					mock.ExpectExec("INSERT INTO \"users\"").WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(1, 1))
				}
				mock.ExpectCommit()
				return db, cleanup
			},
			validate: func(t *testing.T, db *gorm.DB, err error) {
				require.NoError(t, err)
			},
			usersTomlPath: "mock_valid_users.toml",
			concurrent:    true,
		},
	}

	for _, test := range tests {
		db, cleanup := test.setup()
		defer cleanup()

		if test.concurrent {
			mutex.Lock()
			if !txdbInitialized {
				txdbInitialized = true
			}
			mutex.Unlock()
		}

		mockFile := test.usersTomlPath
		originalReadFile := ioutil.ReadFile
		defer func() { ioutil.ReadFile = originalReadFile }()
		ioutil.ReadFile = func(filename string) ([]byte, error) {
			if filename == mockFile {
				return ioutil.ReadFile(filename)
			}
			return nil, os.ErrNotExist
		}

		err := Seed(db)
		test.validate(t, db, err)
		t.Logf("%s: Finished", test.name)
	}
}

func mockUsers() []model.User {
	return []model.User{
		{ID: uuid.New()},
		{ID: uuid.New()},
	}
}
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open("sqlite3", db)
	require.NoError(t, err)

	return gormDB, mock, func() {
		db.Close()
	}
}
