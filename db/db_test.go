package db

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"os"
	"errors"
	"fmt"
	"io/ioutil"
	"database/sql"
	"sync"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

/*
ROOST_METHOD_HASH=AutoMigrate_94b22622a5
ROOST_METHOD_SIG_HASH=AutoMigrate_2cd152caa7


 */
func TestAutoMigrate(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
		expectedLog string
	}{
		{
			name: "Scenario 1: Successful Migration",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: false,
			expectedLog: "",
		},
		{
			name: "Scenario 2: Database Connection Error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("CREATE TABLE").WillReturnError(errors.New("connection error"))
			},
			expectError: true,
			expectedLog: "migration error: connection error",
		},
		{
			name: "Scenario 3: Missing Model Definitions",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("CREATE TABLE").WillReturnError(errors.New("model definition error"))
			},
			expectError: true,
			expectedLog: "migration error: model definition error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open("sqlmock", db)
			assert.NoError(t, err)

			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			mutex.Lock()
			err = AutoMigrate(gormDB)
			mutex.Unlock()

			if tt.expectError {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Expected no error")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}

	t.Run("Scenario 4: Simultaneous Migrations", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open("sqlmock", db)
		assert.NoError(t, err)

		mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))

		var wg sync.WaitGroup
		numRoutines := 5
		wg.Add(numRoutines)

		for i := 0; i < numRoutines; i++ {
			go func() {
				defer wg.Done()
				mutex.Lock()
				err = AutoMigrate(gormDB)
				assert.NoError(t, err)
				mutex.Unlock()
			}()
		}
		wg.Wait()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Scenario 5 & 6: Different Dialects & Logger Usage", func(t *testing.T) {

	})
}

/*
ROOST_METHOD_HASH=DropTestDB_4c6b54d5e5
ROOST_METHOD_SIG_HASH=DropTestDB_69b51a825b


 */
func TestDropTestDB(t *testing.T) {

	tests := []struct {
		name     string
		setup    func() *gorm.DB
		expected error
	}{
		{
			name: "Successfully Close Database Connection",
			setup: func() *gorm.DB {
				db, mock, _ := sqlmock.New()
				gdb, _ := gorm.Open("sqlmock", db)
				mock.ExpectClose()
				return gdb
			},
			expected: nil,
		},
		{
			name: "Handle Closing an Already Closed Database Connection",
			setup: func() *gorm.DB {
				db, mock, _ := sqlmock.New()
				gdb, _ := gorm.Open("sqlmock", db)
				gdb.Close()
				mock.ExpectClose()
				return gdb
			},
			expected: nil,
		},
		{
			name: "Handle Nil Database Connection",
			setup: func() *gorm.DB {
				return nil
			},
			expected: nil,
		},
		{
			name: "Error Handling Mechanism",
			setup: func() *gorm.DB {
				db, mock, _ := sqlmock.New()
				gdb, _ := gorm.Open("sqlmock", db)
				mock.ExpectClose().WillReturnError(errors.New("close error"))
				return gdb
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.setup()
			err := DropTestDB(db)

			if err != tt.expected {
				t.Errorf("expected error: %v, got: %v", tt.expected, err)
			}

			if db != nil {
				mock, ok := db.DB().(*sqlmock.Sqlmock)
				if ok {
					if err = mock.ExpectationsWereMet(); err != nil {
						t.Errorf("unfulfilled expectations: %s", err)
					}
				}
			}

			t.Logf("Test %s passed.", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=dsn_e202d1c4f9
ROOST_METHOD_SIG_HASH=dsn_b336e03d64


 */
func Testdsn(t *testing.T) {
	type testCase struct {
		description   string
		setupEnv      func()
		expectedDSN   string
		expectedError string
	}

	testCases := []testCase{
		{
			description: "Scenario 1: Successful DSN generation with all environment variables set.",
			setupEnv: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "testuser")
				os.Setenv("DB_PASSWORD", "testpassword")
				os.Setenv("DB_NAME", "testdb")
				os.Setenv("DB_PORT", "3306")
			},
			expectedDSN:   "testuser:testpassword@(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
			expectedError: "",
		},
		{
			description: "Scenario 2: Missing DB_HOST environment variable.",
			setupEnv: func() {
				os.Unsetenv("DB_HOST")
				os.Setenv("DB_USER", "testuser")
				os.Setenv("DB_PASSWORD", "testpassword")
				os.Setenv("DB_NAME", "testdb")
				os.Setenv("DB_PORT", "3306")
			},
			expectedDSN:   "",
			expectedError: "$DB_HOST is not set",
		},
		{
			description: "Scenario 3: Missing DB_USER environment variable.",
			setupEnv: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Unsetenv("DB_USER")
				os.Setenv("DB_PASSWORD", "testpassword")
				os.Setenv("DB_NAME", "testdb")
				os.Setenv("DB_PORT", "3306")
			},
			expectedDSN:   "",
			expectedError: "$DB_USER is not set",
		},
		{
			description: "Scenario 4: Missing DB_PASSWORD environment variable.",
			setupEnv: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "testuser")
				os.Unsetenv("DB_PASSWORD")
				os.Setenv("DB_NAME", "testdb")
				os.Setenv("DB_PORT", "3306")
			},
			expectedDSN:   "",
			expectedError: "$DB_PASSWORD is not set",
		},
		{
			description: "Scenario 5: Missing DB_NAME environment variable.",
			setupEnv: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "testuser")
				os.Setenv("DB_PASSWORD", "testpassword")
				os.Unsetenv("DB_NAME")
				os.Setenv("DB_PORT", "3306")
			},
			expectedDSN:   "",
			expectedError: "$DB_NAME is not set",
		},
		{
			description: "Scenario 6: Missing DB_PORT environment variable.",
			setupEnv: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "testuser")
				os.Setenv("DB_PASSWORD", "testpassword")
				os.Setenv("DB_NAME", "testdb")
				os.Unsetenv("DB_PORT")
			},
			expectedDSN:   "",
			expectedError: "$DB_PORT is not set",
		},
		{
			description: "Scenario 7: All environment variables are empty.",
			setupEnv: func() {
				os.Setenv("DB_HOST", "")
				os.Setenv("DB_USER", "")
				os.Setenv("DB_PASSWORD", "")
				os.Setenv("DB_NAME", "")
				os.Setenv("DB_PORT", "")
			},
			expectedDSN:   "",
			expectedError: "$DB_HOST is not set",
		},
		{
			description: "Scenario 8: Non-standard options for DSN string.",
			setupEnv: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "testuser")
				os.Setenv("DB_PASSWORD", "testpassword")
				os.Setenv("DB_NAME", "testdb")
				os.Setenv("DB_PORT", "3306")
			},
			expectedDSN:   "testuser:testpassword@(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tc.setupEnv()

			dsn, err := dsn()

			if tc.expectedError != "" {
				if err == nil || err.Error() != tc.expectedError {
					t.Errorf("expected error: %s, got: %v", tc.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error, but got: %v", err)
				}
				if dsn != tc.expectedDSN {
					t.Errorf("expected DSN: %s, got: %s", tc.expectedDSN, dsn)
				}
			}

			t.Logf("Test %s: DSN generated = %s, error = %v", tc.description, dsn, err)
		})
	}
}

/*
ROOST_METHOD_HASH=Seed_5ad31c3a6c
ROOST_METHOD_SIG_HASH=Seed_878933cebc


 */
func TestSeed(t *testing.T) {
	type testCase struct {
		description   string
		setupMockDB   func(sqlmock.Sqlmock)
		expectedError error
		prepareEnv    func()
	}

	tests := []testCase{
		{
			description: "Successful User Seed",
			setupMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"users\"").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
			prepareEnv: func() {
				content := `
[Users]
[[Users]]
Name = "John Doe"
Email = "john@example.com"
`
				_ = ioutil.WriteFile("db/seed/users.toml", []byte(content), 0644)
			},
		},
		{
			description: "Missing TOML File",
			setupMockDB: func(mock sqlmock.Sqlmock) {

			},
			expectedError: errors.New("open db/seed/users.toml: no such file or directory"),
			prepareEnv: func() {
				_ = os.Remove("db/seed/users.toml")
			},
		},
		{
			description: "Malformed TOML File",
			setupMockDB: func(mock sqlmock.Sqlmock) {

			},
			expectedError: errors.New("expected a comma after table array, but got a carriage return instead"),
			prepareEnv: func() {
				content := `
[Users
[[Users]]
Name = "Jane Doe"
Email = "jane@example.com
`
				_ = ioutil.WriteFile("db/seed/users.toml", []byte(content), 0644)
			},
		},
		{
			description: "Database Insertion Error",
			setupMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"users\"").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnError(errors.New("insertion error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("insertion error"),
			prepareEnv: func() {
				content := `
[Users]
[[Users]]
Name = "Error User"
Email = "error@example.com"
`
				_ = ioutil.WriteFile("db/seed/users.toml", []byte(content), 0644)
			},
		},
		{
			description: "Empty Users List",
			setupMockDB: func(mock sqlmock.Sqlmock) {

			},
			expectedError: nil,
			prepareEnv: func() {
				content := `
[Users]
`
				_ = ioutil.WriteFile("db/seed/users.toml", []byte(content), 0644)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {

			tt.prepareEnv()
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mockDialector := gorm.Dialector{
				DriverName: "sqlite3",
				Conn:       db,
			}

			gormDB, err := gorm.Open(mockDialector, &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open mock gorm db: %v", err)
			}

			tt.setupMockDB(mock)

			err = Seed(gormDB)

			if tt.expectedError == nil && err != nil || tt.expectedError != nil && err == nil ||
				tt.expectedError != nil && err != nil && tt.expectedError.Error() != err.Error() {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Log(fmt.Sprintf("Test Case %q passed", tt.description))
		})
	}
}

/*
ROOST_METHOD_HASH=New_1d2840dc39
ROOST_METHOD_SIG_HASH=New_f9cc65f555


 */
func TestNew(t *testing.T) {
	t.Run("Scenario 1: Successful Database Connection", func(t *testing.T) {

		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening a stub database connection: %v", err)
		}
		defer mockDB.Close()

		driverName := "sqlmock"
		gorm.RegisterDialect(driverName, &gorm.MySQLDialect{})
		txdb.Register("txdb", driverName, fmt.Sprintf("user:password@tcp(localhost:3306)/testdb"))

		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_USER", "user")
		os.Setenv("DB_PASSWORD", "password")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_PORT", "3306")

		mock.ExpectPing()

		db, err := New()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if db == nil {
			t.Errorf("expected non-nil DB instance")
		}
		t.Log("Successfully connected to the database.")
	})

	t.Run("Scenario 2: DSN Retrieval Failure", func(t *testing.T) {

		os.Setenv("DB_HOST", "")

		db, err := New()

		if db != nil {
			t.Errorf("expected nil DB, got %v", db)
		}
		if err == nil {
			t.Errorf("expected error due to missing environment variable")
		}
		t.Log("Handled DSN retrieval failure correctly.")
	})

	t.Run("Scenario 3: Database Connection Attempt Exceeded", func(t *testing.T) {

		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_USER", "user")
		os.Setenv("DB_PASSWORD", "password")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_PORT", "3306")

		gorm.RegisterDialect("sqlmockFail", &gorm.MySQLDialect{})
		txdb.Register("txdbFail", "sqlmockFail", "wrong_dsn")

		db, err := New()

		if db != nil {
			t.Errorf("expected nil DB due to connection failures, got %v", db)
		}
		if err == nil {
			t.Errorf("expected error after all retries exhausted")
		}
		t.Log("Correctly handled connection attempt exceeded.")
	})

	t.Run("Scenario 4: Temporary Connection Issue Resolved", func(t *testing.T) {

		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening a stub database connection: %v", err)
		}
		defer mockDB.Close()

		driverName := "sqlmock"
		gorm.RegisterDialect(driverName, &gorm.MySQLDialect{})
		txdb.Register("txdb", driverName, fmt.Sprintf("user:password@tcp(localhost:3306)/testdb"))

		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_USER", "user")
		os.Setenv("DB_PASSWORD", "password")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_PORT", "3306")

		mock.ExpectPing().WillReturnError(errors.New("temporary error"))
		mock.ExpectPing().WillReturnError(errors.New("temporary error"))
		mock.ExpectPing().WillReturnError(errors.New("temporary error"))
		mock.ExpectPing().WillReturnError(nil)

		db, err := New()

		if err != nil {
			t.Errorf("expected no error after retries")
		}
		if db == nil {
			t.Errorf("expected non-nil DB after successful retry")
		}
		t.Log("Successfully recovered from temporary connection failure.")
	})

	t.Run("Scenario 5: Database Configuration After Connection", func(t *testing.T) {

		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening a stub database connection: %v", err)
		}
		defer mockDB.Close()

		driverName := "sqlmock"
		gorm.RegisterDialect(driverName, &gorm.MySQLDialect{})
		txdb.Register("txdb", driverName, fmt.Sprintf("user:password@tcp(localhost:3306)/testdb"))

		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_USER", "user")
		os.Setenv("DB_PASSWORD", "password")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_PORT", "3306")

		mock.ExpectPing()

		db, err := New()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if db == nil {
			t.Errorf("expected non-nil DB instance")
		}

		sqlDB := db.DB()
		maxIdleConns := sqlDB.Stats().Idle
		if maxIdleConns != 3 {
			t.Errorf("expected MaxIdleConns to be 3, got %d", maxIdleConns)
		}

		t.Log("Database settings configured correctly after connection.")
	})
}

/*
ROOST_METHOD_HASH=NewTestDB_7feb2c4a7a
ROOST_METHOD_SIG_HASH=NewTestDB_1b71546d9d


 */
func TestNewTestDB(t *testing.T) {

	err := godotenv.Load("../env/test.env")
	if err != nil {

		t.Fatalf("Failed to load env file: %v", err)
	}
	_ = os.Setenv("DB_HOST", "localhost")

	type testCase struct {
		name          string
		setupEnv      func()
		mockDBFunc    func() (*gorm.DB, error)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successful Database Connection",

			setupEnv: func() {
				_ = os.Setenv("DB_HOST", "localhost")
				_ = os.Setenv("DB_USER", "test")
				_ = os.Setenv("DB_PASSWORD", "password")
				_ = os.Setenv("DB_NAME", "testdb")
				_ = os.Setenv("DB_PORT", "3306")
			},
			mockDBFunc: func() (*gorm.DB, error) {
				db, _, err := sqlmock.New()
				if err != nil {
					return nil, err
				}
				return &gorm.DB{DB: db}, nil
			},
			expectedError: nil,
		},
		{
			name: "Missing Environment File",
			setupEnv: func() {
				os.Remove("../env/test.env")
			},
			mockDBFunc:    nil,
			expectedError: godotenv.Load("../env/test.env"),
		},
		{
			name: "Invalid DSN Configuration",
			setupEnv: func() {
				_ = os.Setenv("DB_HOST", "")
			},
			mockDBFunc: nil,
			expectedError: func() error {
				_, err := dsn()
				return err
			}(),
		},
		{
			name: "Simulated gorm.Open Failure",
			setupEnv: func() {
				_ = os.Setenv("DB_HOST", "localhost")

			},
			mockDBFunc: func() (*gorm.DB, error) {
				return nil, gorm.ErrCantStartTransaction
			},
			expectedError: gorm.ErrCantStartTransaction,
		},
		{
			name: "Concurrent Access to Uninitialized txdb",
			setupEnv: func() {
				_ = os.Setenv("DB_HOST", "localhost")
			},
			mockDBFunc: func() (*gorm.DB, error) {
				return &gorm.DB{}, nil
			},
			expectedError: nil,
		},
		{
			name: "Handling UUID Errors in sql.Open",
			setupEnv: func() {
				uuid.SetRand(None)
				_ = os.Setenv("DB_HOST", "localhost")
			},
			mockDBFunc: func() (*gorm.DB, error) {
				db, _, err := sqlmock.New()
				if err != nil {
					return nil, err
				}
				return &gorm.DB{DB: db}, nil
			},
			expectedError: nil,
		},
	}

	var once sync.Once

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			tc.setupEnv()

			var err error
			if tc.mockDBFunc != nil {
				once.Do(func() {

					_, err := tc.mockDBFunc()
					if err != nil {
						t.Fatalf("Failed to create mock db: %v", err)
					}
				})
			}

			dbInstance, err := NewTestDB()

			if tc.expectedError == nil && err != nil {
				t.Fatalf("Expected no error, but got %v", err)
			}
			if tc.expectedError != nil && err == nil {
				t.Fatalf("Expected error %v, but got none", tc.expectedError)
			}

			if dbInstance == nil && tc.expectedError == nil {
				t.Fatalf("Expected non-nil dbInstance, got nil!")
			}
		})
	}

	t.Cleanup(func() {
		os.Clearenv()
	})
}

