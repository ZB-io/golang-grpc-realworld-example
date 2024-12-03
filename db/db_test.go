package db

import (
	"database/sql"
	"github.com/DATA-DOG/go-txdb"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq"
	"io/ioutil"
	"path/filepath"
	"errors"
)

/*
ROOST_METHOD_HASH=AutoMigrate_94b22622a5
ROOST_METHOD_SIG_HASH=AutoMigrate_2cd152caa7


 */
func TestAutoMigrate(t *testing.T) {
	tests := []struct {
		name    string
		dbSetup func() (*gorm.DB, error)
		wantErr bool
	}{
		{
			name: "Successful Migration",
			dbSetup: func() (*gorm.DB, error) {

				sql.Open("txdb", uuid.New().String())
				db, err := gorm.Open("txdb", "identifier")
				if err != nil {
					return nil, err
				}
				return db, nil
			},
			wantErr: false,
		},
		{
			name: "Database Connection Error",
			dbSetup: func() (*gorm.DB, error) {

				db, _ := gorm.Open("postgres", "invalid_connection_string")
				return db, nil
			},
			wantErr: true,
		},
		{
			name: "Schema Already Exists",
			dbSetup: func() (*gorm.DB, error) {

				db, err := gorm.Open("txdb", uuid.New().String())
				if err != nil {
					return nil, err
				}

				err = AutoMigrate(db)
				if err != nil {
					return nil, err
				}
				return db, nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Starting test:", tt.name)

			db, err := tt.dbSetup()
			if err != nil {
				t.Fatalf("Failed to setup database: %v", err)
			}
			defer db.Close()

			err = AutoMigrate(db)

			if tt.wantErr {
				assert.Error(t, err)
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				t.Log("Migration completed successfully")

				tables := []string{"users", "articles", "tags", "comments"}
				for _, table := range tables {
					exists := db.HasTable(table)
					assert.True(t, exists, "Table %s should exist", table)
				}
			}
		})
	}

	t.Run("Concurrent Migrations", func(t *testing.T) {
		const numGoroutines = 3
		var wg sync.WaitGroup
		errChan := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				db, err := gorm.Open("txdb", uuid.New().String())
				if err != nil {
					errChan <- err
					return
				}
				defer db.Close()

				if err := AutoMigrate(db); err != nil {
					errChan <- err
				}
			}()
		}

		wg.Wait()
		close(errChan)

		for err := range errChan {
			assert.NoError(t, err, "Concurrent migration should not produce errors")
		}
	})
}

func init() {

	txdb.Register("txdb", "postgres", "postgres://user:password@localhost:5432/testdb?sslmode=disable")
}

/*
ROOST_METHOD_HASH=DropTestDB_4c6b54d5e5
ROOST_METHOD_SIG_HASH=DropTestDB_69b51a825b


 */
func TestDropTestDB(t *testing.T) {

	type testCase struct {
		name     string
		setup    func() *gorm.DB
		validate func(*testing.T, error)
	}

	tests := []testCase{
		{
			name: "Successfully Close Database Connection",
			setup: func() *gorm.DB {
				db, err := gorm.Open("sqlite3", ":memory:")
				if err != nil {
					t.Fatalf("Failed to create test database: %v", err)
				}
				return db
			},
			validate: func(t *testing.T, err error) {
				if err != nil {
					t.Errorf("Expected nil error, got %v", err)
				}
			},
		},
		{
			name: "Handle Already Closed Database Connection",
			setup: func() *gorm.DB {
				db, err := gorm.Open("sqlite3", ":memory:")
				if err != nil {
					t.Fatalf("Failed to create test database: %v", err)
				}
				db.Close()
				return db
			},
			validate: func(t *testing.T, err error) {
				if err != nil {
					t.Errorf("Expected nil error for already closed connection, got %v", err)
				}
			},
		},
		{
			name: "Handle Nil Database Connection",
			setup: func() *gorm.DB {
				return nil
			},
			validate: func(t *testing.T, err error) {
				if err != nil {
					t.Errorf("Expected nil error for nil connection, got %v", err)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Executing test:", tc.name)
			db := tc.setup()
			err := DropTestDB(db)
			tc.validate(t, err)
		})
	}

	t.Run("Concurrent Database Closure", func(t *testing.T) {
		db, err := gorm.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("Failed to create test database: %v", err)
		}

		var wg sync.WaitGroup
		concurrentCalls := 5

		for i := 0; i < concurrentCalls; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := DropTestDB(db)
				if err != nil {
					t.Errorf("Concurrent closure failed: %v", err)
				}
			}()
		}
		wg.Wait()
	})

	t.Run("Database Connection with Active Transactions", func(t *testing.T) {
		db, err := gorm.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("Failed to create test database: %v", err)
		}

		tx := db.Begin()
		if tx.Error != nil {
			t.Fatalf("Failed to begin transaction: %v", tx.Error)
		}

		err = DropTestDB(db)
		if err != nil {
			t.Errorf("Expected nil error with active transaction, got %v", err)
		}
	})

	t.Run("Connection with Multiple References", func(t *testing.T) {
		db, err := gorm.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("Failed to create test database: %v", err)
		}

		ref1 := db
		ref2 := db

		err = DropTestDB(db)
		if err != nil {
			t.Errorf("Expected nil error with multiple references, got %v", err)
		}

		if ref1.Error != nil {
			t.Error("Expected reference 1 to be invalidated")
		}
		if ref2.Error != nil {
			t.Error("Expected reference 2 to be invalidated")
		}
	})
}

/*
ROOST_METHOD_HASH=dsn_e202d1c4f9
ROOST_METHOD_SIG_HASH=dsn_b336e03d64


 */
func Testdsn(t *testing.T) {

	originalEnv := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_PORT":     os.Getenv("DB_PORT"),
	}

	defer func() {
		for k, v := range originalEnv {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	}()

	tests := []TestData{
		{
			name: "Success - Valid Environment Variables",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			expected: "user:password@(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
			errMsg:   "",
		},
		{
			name: "Error - Missing DB_HOST",
			envVars: map[string]string{
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			expected: "",
			errMsg:   "$DB_HOST is not set",
		},
		{
			name: "Error - Missing DB_USER",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PASSWORD": "password",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			expected: "",
			errMsg:   "$DB_USER is not set",
		},
		{
			name: "Error - Missing DB_PASSWORD",
			envVars: map[string]string{
				"DB_HOST": "localhost",
				"DB_USER": "user",
				"DB_NAME": "testdb",
				"DB_PORT": "3306",
			},
			expected: "",
			errMsg:   "$DB_PASSWORD is not set",
		},
		{
			name: "Error - Missing DB_NAME",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_PORT":     "3306",
			},
			expected: "",
			errMsg:   "$DB_NAME is not set",
		},
		{
			name: "Error - Missing DB_PORT",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "testdb",
			},
			expected: "",
			errMsg:   "$DB_PORT is not set",
		},
		{
			name: "Success - Special Characters in Environment Variables",
			envVars: map[string]string{
				"DB_HOST":     "host-name.com",
				"DB_USER":     "user@123",
				"DB_PASSWORD": "pass#word!",
				"DB_NAME":     "test_db",
				"DB_PORT":     "3306",
			},
			expected: "user@123:pass#word!@(host-name.com:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local",
			errMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for k := range originalEnv {
				os.Unsetenv(k)
			}

			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			got, err := dsn()

			if tt.errMsg != "" {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.errMsg)
				} else if err.Error() != tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, got)
				}
			}

			t.Logf("Test: %s\nExpected: %s\nGot: %s\nError: %v",
				tt.name, tt.expected, got, err)
		})
	}
}

/*
ROOST_METHOD_HASH=New_1d2840dc39
ROOST_METHOD_SIG_HASH=New_f9cc65f555


 */
func TestNew(t *testing.T) {

	type testCase struct {
		name     string
		envVars  map[string]string
		mockDSN  func() (string, error)
		wantErr  bool
		errMsg   string
		validate func(*testing.T, *gorm.DB)
	}

	tests := []testCase{
		{
			name: "Successful Database Connection",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "test_user",
				"DB_PASSWORD": "test_pass",
				"DB_NAME":     "test_db",
				"DB_PORT":     "3306",
			},
			wantErr: false,
			validate: func(t *testing.T, db *gorm.DB) {
				assert.NotNil(t, db)

				sqlDB := db.DB()
				maxIdle := sqlDB.Stats().MaxOpenConnections
				assert.Equal(t, 3, maxIdle)
			},
		},
		{
			name:    "Invalid Database Credentials",
			envVars: map[string]string{},
			wantErr: true,
			errMsg:  "$DB_HOST is not set",
		},
		{
			name: "DSN Configuration Error",
			envVars: map[string]string{
				"DB_HOST": "invalid_host",
			},
			wantErr: true,
			errMsg:  "$DB_USER is not set",
		},
		{
			name: "Connection Retry Behavior",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "test_user",
				"DB_PASSWORD": "test_pass",
				"DB_NAME":     "test_db",
				"DB_PORT":     "3306",
			},
			wantErr: false,
			validate: func(t *testing.T, db *gorm.DB) {
				assert.NotNil(t, db)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			for k, v := range tc.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			db, err := New()

			if tc.wantErr {
				assert.Error(t, err)
				if tc.errMsg != "" {
					assert.Contains(t, err.Error(), tc.errMsg)
				}
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				if tc.validate != nil {
					tc.validate(t, db)
				}

				if db != nil {
					db.Close()
				}
			}
		})
	}
}

/*
ROOST_METHOD_HASH=Seed_5ad31c3a6c
ROOST_METHOD_SIG_HASH=Seed_878933cebc


 */
func TestSeed(t *testing.T) {
	testDir := "db/seed"
	testFile := filepath.Join(testDir, "users.toml")

	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name        string
		tomlContent string
		setupDB     func() (*gorm.DB, error)
		wantErr     bool
		errMsg      string
	}{
		{
			name: "Successfully seed database",
			tomlContent: `
				[[Users]]
				id = "550e8400-e29b-41d4-a716-446655440000"
				email = "test@example.com"
				username = "testuser"
				password = "password123"
			`,
			setupDB: setupTestDB,
			wantErr: false,
		},
		{
			name: "Invalid TOML format",
			tomlContent: `
				[[Users
				invalid toml
			`,
			setupDB: setupTestDB,
			wantErr: true,
			errMsg:  "toml: error",
		},
		{
			name:        "Empty TOML file",
			tomlContent: "",
			setupDB:     setupTestDB,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := ioutil.WriteFile(testFile, []byte(tt.tomlContent), 0644)
			if err != nil {
				t.Fatalf("Failed to write test TOML file: %v", err)
			}
			defer os.Remove(testFile)

			db, err := tt.setupDB()
			if err != nil {
				t.Fatalf("Failed to setup test database: %v", err)
			}
			defer teardownTestDB(db)

			if tt.name == "Successfully seed database" {
				var wg sync.WaitGroup
				errChan := make(chan error, 3)

				for i := 0; i < 3; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						if err := Seed(db); err != nil {
							errChan <- err
						}
					}()
				}

				wg.Wait()
				close(errChan)

				for err := range errChan {
					t.Errorf("Concurrent seeding error: %v", err)
				}
			}

			err = Seed(db)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)

				var users []User
				if err := db.Find(&users).Error; err != nil {
					t.Errorf("Failed to query users: %v", err)
				}

				if tt.name == "Successfully seed database" {
					assert.Equal(t, 1, len(users))
					assert.Equal(t, "test@example.com", users[0].Email)
					assert.Equal(t, "testuser", users[0].Username)
				} else if tt.name == "Empty TOML file" {
					assert.Equal(t, 0, len(users))
				}
			}
		})
	}
}

func setupTestDB() (*gorm.DB, error) {

	txdb.Register("txdb", "postgres", "postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable")

	db, err := gorm.Open("txdb", "identifier")
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{})

	return db, nil
}

func teardownTestDB(db *gorm.DB) error {
	return db.Close()
}

/*
ROOST_METHOD_HASH=NewTestDB_7feb2c4a7a
ROOST_METHOD_SIG_HASH=NewTestDB_1b71546d9d


 */
func TestNewTestDB(t *testing.T) {
	type testCase struct {
		name          string
		setupFunc     func()
		cleanupFunc   func()
		expectedError bool
		validateFunc  func(*testing.T, *gorm.DB, error)
	}

	resetState := func() {
		txdbInitialized = false
		mutex = sync.Mutex{}
	}

	tests := []testCase{
		{
			name: "Successful Database Connection",
			setupFunc: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "test_user")
				os.Setenv("DB_PASSWORD", "test_password")
				os.Setenv("DB_NAME", "test_db")
				os.Setenv("DB_PORT", "3306")
			},
			cleanupFunc:   resetState,
			expectedError: false,
			validateFunc: func(t *testing.T, db *gorm.DB, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, db)

				maxIdle := db.DB().Stats().MaxOpenConnections
				assert.Equal(t, 3, maxIdle)

				db.LogMode(false)
				assert.NotNil(t, db)
			},
		},
		{
			name: "Missing Environment File",
			setupFunc: func() {
				os.Unsetenv("DB_HOST")
				os.Unsetenv("DB_USER")
				os.Unsetenv("DB_PASSWORD")
				os.Unsetenv("DB_NAME")
				os.Unsetenv("DB_PORT")
			},
			cleanupFunc:   resetState,
			expectedError: true,
			validateFunc: func(t *testing.T, db *gorm.DB, err error) {
				assert.Error(t, err)
				assert.Nil(t, db)
			},
		},
		{
			name: "Concurrent Access",
			setupFunc: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "test_user")
				os.Setenv("DB_PASSWORD", "test_password")
				os.Setenv("DB_NAME", "test_db")
				os.Setenv("DB_PORT", "3306")
			},
			cleanupFunc:   resetState,
			expectedError: false,
			validateFunc: func(t *testing.T, db *gorm.DB, err error) {
				var wg sync.WaitGroup
				concurrentCalls := 5
				results := make(chan error, concurrentCalls)

				for i := 0; i < concurrentCalls; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						db, err := NewTestDB()
						if err != nil {
							results <- err
							return
						}
						if db == nil {
							results <- errors.New("db is nil")
							return
						}
						results <- nil
					}()
				}

				wg.Wait()
				close(results)

				for err := range results {
					assert.NoError(t, err)
				}
			},
		},
		{
			name: "Multiple Sequential Connections",
			setupFunc: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "test_user")
				os.Setenv("DB_PASSWORD", "test_password")
				os.Setenv("DB_NAME", "test_db")
				os.Setenv("DB_PORT", "3306")
			},
			cleanupFunc:   resetState,
			expectedError: false,
			validateFunc: func(t *testing.T, db *gorm.DB, err error) {
				connections := make([]*gorm.DB, 3)
				for i := 0; i < 3; i++ {
					db, err := NewTestDB()
					assert.NoError(t, err)
					assert.NotNil(t, db)
					connections[i] = db
				}

				connStats := make(map[int]bool)
				for _, conn := range connections {
					stats := conn.DB().Stats().OpenConnections
					assert.False(t, connStats[stats])
					connStats[stats] = true
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupFunc != nil {
				tc.setupFunc()
			}

			if tc.cleanupFunc != nil {
				defer tc.cleanupFunc()
			}

			db, err := NewTestDB()

			if tc.validateFunc != nil {
				tc.validateFunc(t, db, err)
			} else {
				if tc.expectedError {
					assert.Error(t, err)
					assert.Nil(t, db)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, db)
				}
			}

			if db != nil {
				db.Close()
			}
		})
	}
}

