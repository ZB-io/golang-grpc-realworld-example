package db

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var mockDSN func() (string, error)
var mockGormOpen func(dialect string, args ...interface{}) (*gorm.DB, error)
var mockSleep func(d time.Duration)

/*
ROOST_METHOD_HASH=AutoMigrate_94b22622a5
ROOST_METHOD_SIG_HASH=AutoMigrate_2cd152caa7
*/
func TestAutoMigrate(t *testing.T) {
	tests := []struct {
		name    string
		dbSetup func() (*gorm.DB, sqlmock.Sqlmock, error)
		wantErr bool
	}{
		{
			name: "Successful Auto-Migration",
			dbSetup: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				gormDB, err := gorm.Open("mysql", db)
				if err != nil {
					return nil, nil, err
				}
				mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))
				return gormDB, mock, nil
			},
			wantErr: false,
		},
		{
			name: "Database Connection Error",
			dbSetup: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				return nil, nil, errors.New("connection error")
			},
			wantErr: true,
		},
		{
			name: "Partial Migration Failure",
			dbSetup: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				gormDB, err := gorm.Open("mysql", db)
				if err != nil {
					return nil, nil, err
				}
				mock.ExpectExec("CREATE TABLE").WillReturnError(errors.New("migration error"))
				return gormDB, mock, nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := tt.dbSetup()
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("Setup failed: %v", err)
				}
				return
			}
			defer db.Close()

			err = AutoMigrate(db)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if mock != nil {
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("There were unfulfilled expectations: %s", err)
				}
			}
		})
	}

	t.Run("Concurrent Migration Attempts", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		gormDB, err := gorm.Open("mysql", db)
		assert.NoError(t, err)
		defer gormDB.Close()

		mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))

		var wg sync.WaitGroup
		errChan := make(chan error, 3)

		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := AutoMigrate(gormDB)
				errChan <- err
			}()
		}

		wg.Wait()
		close(errChan)

		successCount := 0
		for err := range errChan {
			if err == nil {
				successCount++
			}
		}

		assert.Equal(t, 1, successCount, "Expected exactly one successful migration")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})
}

/*
ROOST_METHOD_HASH=DropTestDB_4c6b54d5e5
ROOST_METHOD_SIG_HASH=DropTestDB_69b51a825b
*/
func TestDropTestDB(t *testing.T) {
	tests := []struct {
		name    string
		db      *gorm.DB
		wantErr bool
	}{
		{
			name:    "Successfully Close Database Connection",
			db:      &gorm.DB{},
			wantErr: false,
		},
		{
			name:    "Handle Nil Database Pointer",
			db:      nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DropTestDB(tt.db)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDropTestDBConcurrent(t *testing.T) {
	db := &gorm.DB{}

	var wg sync.WaitGroup
	concurrentCalls := 10

	for i := 0; i < concurrentCalls; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := DropTestDB(db)
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
}

/*
ROOST_METHOD_HASH=dsn_e202d1c4f9
ROOST_METHOD_SIG_HASH=dsn_b336e03d64
*/
func Testdsn(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
		wantErr  bool
		errMsg   string
	}{
		{
			name: "All Environment Variables Set Correctly",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			expected: "user:password@(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
			wantErr:  false,
		},
		{
			name: "Missing DB_HOST Environment Variable",
			envVars: map[string]string{
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			wantErr: true,
			errMsg:  "$DB_HOST is not set",
		},
		{
			name: "Missing DB_USER Environment Variable",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PASSWORD": "password",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			wantErr: true,
			errMsg:  "$DB_USER is not set",
		},
		{
			name: "Missing DB_PASSWORD Environment Variable",
			envVars: map[string]string{
				"DB_HOST": "localhost",
				"DB_USER": "user",
				"DB_NAME": "testdb",
				"DB_PORT": "3306",
			},
			wantErr: true,
			errMsg:  "$DB_PASSWORD is not set",
		},
		{
			name: "Missing DB_NAME Environment Variable",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_PORT":     "3306",
			},
			wantErr: true,
			errMsg:  "$DB_NAME is not set",
		},
		{
			name: "Missing DB_PORT Environment Variable",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "testdb",
			},
			wantErr: true,
			errMsg:  "$DB_PORT is not set",
		},
		{
			name: "All Environment Variables Set with Empty Values",
			envVars: map[string]string{
				"DB_HOST":     "",
				"DB_USER":     "",
				"DB_PASSWORD": "",
				"DB_NAME":     "",
				"DB_PORT":     "",
			},
			expected: ":@(:)/?charset=utf8mb4&parseTime=True&loc=Local",
			wantErr:  false,
		},
		{
			name: "Special Characters in Environment Variables",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user@123",
				"DB_PASSWORD": "p@ssw0rd!",
				"DB_NAME":     "test_db",
				"DB_PORT":     "3306",
			},
			expected: "user@123:p@ssw0rd!@(localhost:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			got, err := dsn()

			if (err != nil) != tt.wantErr {
				t.Errorf("dsn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("dsn() error message = %v, want %v", err.Error(), tt.errMsg)
				return
			}

			if !tt.wantErr && got != tt.expected {
				t.Errorf("dsn() = %v, want %v", got, tt.expected)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=Seed_5ad31c3a6c
ROOST_METHOD_SIG_HASH=Seed_878933cebc
*/
func TestSeed(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func() (*gorm.DB, error)
		tomlContent   string
		expectedError error
		expectedUsers int
		cleanupFunc   func()
	}{
		{
			name: "Successful Seeding of Users",
			setupFunc: func() (*gorm.DB, error) {
				return &gorm.DB{}, nil
			},
			tomlContent: `
				[[Users]]
				username = "user1"
				email = "user1@example.com"
				password = "password1"

				[[Users]]
				username = "user2"
				email = "user2@example.com"
				password = "password2"
			`,
			expectedError: nil,
			expectedUsers: 2,
			cleanupFunc:   func() {},
		},
		{
			name: "Handling Non-existent TOML File",
			setupFunc: func() (*gorm.DB, error) {
				return &gorm.DB{}, nil
			},
			tomlContent:   "",
			expectedError: errors.New("open db/seed/users.toml: no such file or directory"),
			expectedUsers: 0,
			cleanupFunc:   func() {},
		},
		{
			name: "Handling Malformed TOML File",
			setupFunc: func() (*gorm.DB, error) {
				return &gorm.DB{}, nil
			},
			tomlContent: `
				[[Users]
				username = "user1"
				email = "user1@example.com"
				password = "password1"
			`,
			expectedError: errors.New("toml: line 2: expected '=', '.' or ']' after a key"),
			expectedUsers: 0,
			cleanupFunc:   func() {},
		},
		{
			name: "Database Insertion Failure",
			setupFunc: func() (*gorm.DB, error) {
				db := &gorm.DB{}
				db.Error = errors.New("database insertion error")
				return db, nil
			},
			tomlContent: `
				[[Users]]
				username = "user1"
				email = "user1@example.com"
				password = "password1"
			`,
			expectedError: errors.New("database insertion error"),
			expectedUsers: 0,
			cleanupFunc:   func() {},
		},
		{
			name: "Empty TOML File",
			setupFunc: func() (*gorm.DB, error) {
				return &gorm.DB{}, nil
			},
			tomlContent:   "",
			expectedError: nil,
			expectedUsers: 0,
			cleanupFunc:   func() {},
		},
		{
			name: "Large Number of Users",
			setupFunc: func() (*gorm.DB, error) {
				return &gorm.DB{}, nil
			},
			tomlContent: func() string {
				var content string
				for i := 0; i < 10000; i++ {
					content += fmt.Sprintf(`
						[[Users]]
						username = "user%d"
						email = "user%d@example.com"
						password = "password%d"
					`, i, i, i)
				}
				return content
			}(),
			expectedError: nil,
			expectedUsers: 10000,
			cleanupFunc:   func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := tt.setupFunc()
			if err != nil {
				t.Fatalf("Failed to setup test: %v", err)
			}

			tmpfile, err := ioutil.TempFile("", "users.*.toml")
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.tomlContent)); err != nil {
				t.Fatalf("Failed to write to temporary file: %v", err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatalf("Failed to close temporary file: %v", err)
			}

			originalFile := "db/seed/users.toml"
			os.Rename(originalFile, originalFile+".bak")
			os.Rename(tmpfile.Name(), originalFile)

			err = Seed(db)

			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("Seed() error = %v, expectedError %v", err, tt.expectedError)
			}

			var count int
			db.Model(&model.User{}).Count(&count)
			if count != tt.expectedUsers {
				t.Errorf("Expected %d users, but got %d", tt.expectedUsers, count)
			}

			os.Rename(originalFile+".bak", originalFile)
			tt.cleanupFunc()
		})
	}
}

/*
ROOST_METHOD_HASH=New_1d2840dc39
ROOST_METHOD_SIG_HASH=New_f9cc65f555
*/
func TestNew(t *testing.T) {
	originalDSN := dsn
	originalGormOpen := gorm.Open
	originalSleep := time.Sleep

	defer func() {
		dsn = originalDSN
		gorm.Open = originalGormOpen
		time.Sleep = originalSleep
	}()

	tests := []struct {
		name          string
		dsnFunc       func() (string, error)
		gormOpenFunc  func(string, ...interface{}) (*gorm.DB, error)
		sleepFunc     func(time.Duration)
		expectedDB    *gorm.DB
		expectedError error
		expectedCalls int
		maxIdleConns  int
		logMode       bool
	}{
		{
			name: "Successful Database Connection",
			dsnFunc: func() (string, error) {
				return "valid_dsn", nil
			},
			gormOpenFunc: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				return &gorm.DB{}, nil
			},
			sleepFunc:     func(d time.Duration) {},
			expectedDB:    &gorm.DB{},
			expectedError: nil,
			expectedCalls: 1,
			maxIdleConns:  3,
			logMode:       false,
		},
		{
			name: "Database Connection Failure",
			dsnFunc: func() (string, error) {
				return "valid_dsn", nil
			},
			gormOpenFunc: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				return nil, errors.New("connection failed")
			},
			sleepFunc:     func(d time.Duration) {},
			expectedDB:    nil,
			expectedError: errors.New("connection failed"),
			expectedCalls: 10,
			maxIdleConns:  0,
			logMode:       false,
		},
		{
			name: "DSN Retrieval Failure",
			dsnFunc: func() (string, error) {
				return "", errors.New("dsn error")
			},
			gormOpenFunc:  nil,
			sleepFunc:     nil,
			expectedDB:    nil,
			expectedError: errors.New("dsn error"),
			expectedCalls: 0,
			maxIdleConns:  0,
			logMode:       false,
		},
		{
			name: "Successful Connection After Retries",
			dsnFunc: func() (string, error) {
				return "valid_dsn", nil
			},
			gormOpenFunc: func() func(string, ...interface{}) (*gorm.DB, error) {
				count := 0
				return func(dialect string, args ...interface{}) (*gorm.DB, error) {
					count++
					if count < 3 {
						return nil, errors.New("temporary error")
					}
					return &gorm.DB{}, nil
				}
			}(),
			sleepFunc:     func(d time.Duration) {},
			expectedDB:    &gorm.DB{},
			expectedError: nil,
			expectedCalls: 3,
			maxIdleConns:  3,
			logMode:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDSN = tt.dsnFunc
			dsn = func() (string, error) {
				return mockDSN()
			}

			gormOpenCalls := 0
			mockGormOpen = func(dialect string, args ...interface{}) (*gorm.DB, error) {
				gormOpenCalls++
				return tt.gormOpenFunc(dialect, args...)
			}
			gorm.Open = func(dialect string, args ...interface{}) (*gorm.DB, error) {
				return mockGormOpen(dialect, args...)
			}

			if tt.sleepFunc != nil {
				mockSleep = tt.sleepFunc
				time.Sleep = func(d time.Duration) {
					mockSleep(d)
				}
			}

			db, err := New()

			if (err != nil) != (tt.expectedError != nil) {
				t.Errorf("New() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("New() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if (db != nil) != (tt.expectedDB != nil) {
				t.Errorf("New() db = %v, expectedDB %v", db, tt.expectedDB)
				return
			}
			if gormOpenCalls != tt.expectedCalls {
				t.Errorf("New() gorm.Open calls = %d, expected %d", gormOpenCalls, tt.expectedCalls)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=NewTestDB_7feb2c4a7a
ROOST_METHOD_SIG_HASH=NewTestDB_1b71546d9d
*/
func TestNewTestDB(t *testing.T) {
	txdbInitialized = false

	tests := []struct {
		name    string
		setup   func()
		cleanup func()
		wantErr bool
	}{
		{
			name:    "Successful Database Connection and Initialization",
			setup:   func() {},
			cleanup: func() {},
			wantErr: false,
		},
		{
			name: "Environment File Not Found",
			setup: func() {
				os.Rename("../env/test.env", "../env/test.env.bak")
			},
			cleanup: func() {
				os.Rename("../env/test.env.bak", "../env/test.env")
			},
			wantErr: true,
		},
		{
			name:    "Invalid Database Credentials",
			setup:   func() {},
			cleanup: func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()

			got, err := NewTestDB()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTestDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("NewTestDB() returned nil, want non-nil *gorm.DB")
				} else {
					got.Close()
				}
			}
		})
	}
}

func TestNewTestDBConcurrent(t *testing.T) {
	txdbInitialized = false

	const numGoroutines = 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	errChan := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			db, err := NewTestDB()
			if err != nil {
				errChan <- err
				return
			}
			if db == nil {
				errChan <- errors.New("NewTestDB() returned nil")
				return
			}
			db.Close()
		}()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Errorf("NewTestDB() error in goroutine: %v", err)
	}

	if !txdbInitialized {
		t.Errorf("txdbInitialized not set to true after concurrent calls")
	}
}

func TestNewTestDBMultipleCalls(t *testing.T) {
	txdbInitialized = false

	db1, err := NewTestDB()
	if err != nil {
		t.Fatalf("First call to NewTestDB() failed: %v", err)
	}
	defer db1.Close()

	db2, err := NewTestDB()
	if err != nil {
		t.Fatalf("Second call to NewTestDB() failed: %v", err)
	}
	defer db2.Close()

	if db1 == db2 {
		t.Errorf("NewTestDB() returned the same instance for multiple calls")
	}
}

func TestNewTestDBConnectionPool(t *testing.T) {
	txdbInitialized = false

	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("NewTestDB() failed: %v", err)
	}
	defer db.Close()

	sqlDB := db.DB()
	maxIdleConns := sqlDB.Stats().MaxOpenConnections
	if maxIdleConns != 3 {
		t.Errorf("MaxIdleConns = %d, want 3", maxIdleConns)
	}
}

func TestNewTestDBLogMode(t *testing.T) {
	txdbInitialized = false

	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("NewTestDB() failed: %v", err)
	}
	defer db.Close()

	var logOutput string
	db.SetLogger(gorm.Logger{LogWriter: gorm.LogWriter(func(s string, i ...interface{}) {
		logOutput += s
	})})

	db.Exec("SELECT 1")

	if logOutput != "" {
		t.Errorf("LogMode is not set to false, got log output: %s", logOutput)
	}
}
