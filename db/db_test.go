package db

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DATA-DOG/go-txdb"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			name: "Successful Auto-Migration",
			dbSetup: func() (*gorm.DB, error) {
				db, _ := gorm.Open("sqlite3", ":memory:")
				return db, nil
			},
			wantErr: false,
		},
		{
			name: "Database Connection Error",
			dbSetup: func() (*gorm.DB, error) {
				return nil, errors.New("database connection error")
			},
			wantErr: true,
		},
		{
			name: "Partial Migration Failure",
			dbSetup: func() (*gorm.DB, error) {
				db, _ := gorm.Open("sqlite3", ":memory:")
				db.Close()
				return db, nil
			},
			wantErr: true,
		},
		{
			name: "Concurrent Auto-Migration Attempts",
			dbSetup: func() (*gorm.DB, error) {
				db, _ := gorm.Open("sqlite3", ":memory:")
				return db, nil
			},
			wantErr: false,
		},
		{
			name: "Auto-Migration with Existing Schema",
			dbSetup: func() (*gorm.DB, error) {
				db, _ := gorm.Open("sqlite3", ":memory:")
				db.CreateTable(&model.User{}, &model.Article{}, &model.Tag{}, &model.Comment{})
				return db, nil
			},
			wantErr: false,
		},
		{
			name: "Auto-Migration with Custom Model Configurations",
			dbSetup: func() (*gorm.DB, error) {
				db, _ := gorm.Open("sqlite3", ":memory:")
				return db, nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := tt.dbSetup()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Unexpected error in test setup: %v", err)
				}
				return
			}
			defer db.Close()

			if tt.name == "Concurrent Auto-Migration Attempts" {
				var wg sync.WaitGroup
				errChan := make(chan error, 5)
				for i := 0; i < 5; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						err := AutoMigrate(db)
						if err != nil {
							errChan <- err
						}
					}()
				}
				wg.Wait()
				close(errChan)

				if tt.wantErr && len(errChan) == 0 {
					t.Errorf("Expected errors in concurrent migration, but got none")
				}
				if !tt.wantErr && len(errChan) > 0 {
					t.Errorf("Unexpected errors in concurrent migration: %v", <-errChan)
				}
			} else {
				err := AutoMigrate(db)
				if (err != nil) != tt.wantErr {
					t.Errorf("AutoMigrate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			if !tt.wantErr {
				tables := []string{"users", "articles", "tags", "comments"}
				for _, table := range tables {
					if !db.HasTable(table) {
						t.Errorf("Table %s was not created", table)
					}
				}

				verifyTableStructure(t, db, &model.User{}, "users")
				verifyTableStructure(t, db, &model.Article{}, "articles")
				verifyTableStructure(t, db, &model.Tag{}, "tags")
				verifyTableStructure(t, db, &model.Comment{}, "comments")
			}
		})
	}
}

func verifyTableStructure(t *testing.T, db *gorm.DB, model interface{}, tableName string) {
	t.Helper()
	columns, err := db.Table(tableName).Columns()
	if err != nil {
		t.Errorf("Error getting columns for table %s: %v", tableName, err)
		return
	}

	expectedFields := db.NewScope(model).Fields()
	for _, field := range expectedFields {
		if !field.IsIgnored {
			found := false
			for _, column := range columns {
				if column.Name() == field.DBName {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected column %s not found in table %s", field.DBName, tableName)
			}
		}
	}
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
		setup   func(*gorm.DB)
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
		{
			name: "Handle Already Closed Database",
			db:   &gorm.DB{},
			setup: func(db *gorm.DB) {
				db.Error = errors.New("sql: database is closed")
			},
			wantErr: false,
		},
		{
			name: "Concurrent Access During Closure",
			db:   &gorm.DB{},
			setup: func(db *gorm.DB) {
				var wg sync.WaitGroup
				wg.Add(5)
				for i := 0; i < 5; i++ {
					go func() {
						defer wg.Done()
						_ = db.Close()
					}()
				}
				wg.Wait()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.db)
			}

			err := DropTestDB(tt.db)

			if (err != nil) != tt.wantErr {
				t.Errorf("DropTestDB() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.name == "Handle Nil Database Pointer" && err != nil {
				t.Errorf("Expected no error for nil database, got %v", err)
			}
		})
	}
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
				"DB_HOST":     "local@host",
				"DB_USER":     "user!name",
				"DB_PASSWORD": "pass#word$",
				"DB_NAME":     "test@db",
				"DB_PORT":     "3306",
			},
			expected: "user!name:pass#word$@(local@host:3306)/test@db?charset=utf8mb4&parseTime=True&loc=Local",
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
		expectedError error
		expectedUsers int
		cleanupFunc   func()
	}{
		{
			name: "Successful Seeding of Users",
			setupFunc: func() (*gorm.DB, error) {
				mockDB, _ := gorm.Open("sqlite3", ":memory:")
				users := struct {
					Users []model.User
				}{
					Users: []model.User{
						{Username: "user1"},
						{Username: "user2"},
						{Username: "user3"},
					},
				}
				file, _ := os.Create("db/seed/users.toml")
				defer file.Close()
				toml.NewEncoder(file).Encode(users)
				return mockDB, nil
			},
			expectedError: nil,
			expectedUsers: 3,
			cleanupFunc: func() {
				os.Remove("db/seed/users.toml")
			},
		},
		{
			name: "Non-existent TOML File",
			setupFunc: func() (*gorm.DB, error) {
				mockDB, _ := gorm.Open("sqlite3", ":memory:")
				os.Remove("db/seed/users.toml")
				return mockDB, nil
			},
			expectedError: errors.New("open db/seed/users.toml: no such file or directory"),
			expectedUsers: 0,
			cleanupFunc:   func() {},
		},
		{
			name: "Malformed TOML File",
			setupFunc: func() (*gorm.DB, error) {
				mockDB, _ := gorm.Open("sqlite3", ":memory:")
				ioutil.WriteFile("db/seed/users.toml", []byte("invalid toml content"), 0644)
				return mockDB, nil
			},
			expectedError: errors.New("toml: line 1: unexpected EOF"),
			expectedUsers: 0,
			cleanupFunc: func() {
				os.Remove("db/seed/users.toml")
			},
		},
		{
			name: "Database Insertion Failure",
			setupFunc: func() (*gorm.DB, error) {
				mockDB, _ := gorm.Open("sqlite3", ":memory:")
				mockDB.AddError(errors.New("database insertion error"))
				users := struct {
					Users []model.User
				}{
					Users: []model.User{
						{Username: "user1"},
						{Username: "user2"},
						{Username: "user3"},
					},
				}
				file, _ := os.Create("db/seed/users.toml")
				defer file.Close()
				toml.NewEncoder(file).Encode(users)
				return mockDB, nil
			},
			expectedError: errors.New("database insertion error"),
			expectedUsers: 0,
			cleanupFunc: func() {
				os.Remove("db/seed/users.toml")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := tt.setupFunc()
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			defer tt.cleanupFunc()

			err = Seed(mockDB)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			var count int
			mockDB.Model(&model.User{}).Count(&count)
			assert.Equal(t, tt.expectedUsers, count)
		})
	}
}

/*
ROOST_METHOD_HASH=New_1d2840dc39
ROOST_METHOD_SIG_HASH=New_f9cc65f555
*/
func TestNew(t *testing.T) {
	originalDsn := dsn
	defer func() { dsn = originalDsn }()

	originalSleep := time.Sleep
	defer func() { time.Sleep = originalSleep }()
	time.Sleep = func(d time.Duration) {}

	tests := []struct {
		name           string
		dsnFunc        func() (string, error)
		gormOpenFunc   func(dialect string, args ...interface{}) (*gorm.DB, error)
		expectedDB     bool
		expectedError  error
		expectedCalls  int
		validateConfig func(*testing.T, *gorm.DB)
	}{
		{
			name: "Successful Database Connection",
			dsnFunc: func() (string, error) {
				return "valid_dsn", nil
			},
			gormOpenFunc: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				db, _, _ := sqlmock.New()
				return &gorm.DB{DB: db}, nil
			},
			expectedDB:    true,
			expectedError: nil,
			expectedCalls: 1,
			validateConfig: func(t *testing.T, db *gorm.DB) {
				sqlDB := db.DB()
				if sqlDB == nil {
					t.Error("Expected SQL DB to be set")
				}
			},
		},
		{
			name: "Database Connection Failure",
			dsnFunc: func() (string, error) {
				return "valid_dsn", nil
			},
			gormOpenFunc: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				return nil, errors.New("connection failed")
			},
			expectedDB:    false,
			expectedError: errors.New("connection failed"),
			expectedCalls: 10,
		},
		{
			name: "DSN Retrieval Failure",
			dsnFunc: func() (string, error) {
				return "", errors.New("DSN error")
			},
			expectedDB:    false,
			expectedError: errors.New("DSN error"),
			expectedCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn = tt.dsnFunc
			originalGormOpen := gorm.Open
			defer func() { gorm.Open = originalGormOpen }()

			calls := 0
			gorm.Open = func(dialect string, args ...interface{}) (*gorm.DB, error) {
				calls++
				return tt.gormOpenFunc(dialect, args...)
			}

			db, err := New()

			if (db != nil) != tt.expectedDB {
				t.Errorf("New() returned unexpected db status, got: %v, want: %v", db != nil, tt.expectedDB)
			}

			if err != nil && tt.expectedError != nil {
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("New() returned unexpected error, got: %v, want: %v", err, tt.expectedError)
				}
			} else if (err != nil) != (tt.expectedError != nil) {
				t.Errorf("New() returned unexpected error status, got: %v, want: %v", err, tt.expectedError)
			}

			if calls != tt.expectedCalls {
				t.Errorf("New() made unexpected number of calls, got: %d, want: %d", calls, tt.expectedCalls)
			}

			if tt.validateConfig != nil && db != nil {
				tt.validateConfig(t, db)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=NewTestDB_7feb2c4a7a
ROOST_METHOD_SIG_HASH=NewTestDB_1b71546d9d
*/
func TestNewTestDB(t *testing.T) {
	if err := os.Rename("../env/test.env", "../env/test.env.bak"); err != nil {
		t.Fatalf("Failed to backup test.env: %v", err)
	}
	defer os.Rename("../env/test.env.bak", "../env/test.env")

	tests := []struct {
		name         string
		setupFunc    func() error
		wantErr      bool
		validateFunc func(*testing.T, *gorm.DB, error)
	}{
		{
			name: "Successful Database Connection and Initialization",
			setupFunc: func() error {
				return os.WriteFile("../env/test.env", []byte("DB_DSN=validdsn"), 0644)
			},
			wantErr: false,
			validateFunc: func(t *testing.T, db *gorm.DB, err error) {
				if db == nil {
					t.Error("Expected non-nil DB, got nil")
				}
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			},
		},
		{
			name: "Environment File Not Found",
			setupFunc: func() error {
				return os.Remove("../env/test.env")
			},
			wantErr: true,
			validateFunc: func(t *testing.T, db *gorm.DB, err error) {
				if db != nil {
					t.Error("Expected nil DB, got non-nil")
				}
				if err == nil {
					t.Error("Expected an error, got nil")
				}
			},
		},
		{
			name: "Invalid Database Credentials",
			setupFunc: func() error {
				return os.WriteFile("../env/test.env", []byte("DB_DSN=invaliddsn"), 0644)
			},
			wantErr: true,
			validateFunc: func(t *testing.T, db *gorm.DB, err error) {
				if db != nil {
					t.Error("Expected nil DB, got non-nil")
				}
				if err == nil {
					t.Error("Expected an error, got nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mutex = sync.Mutex{}
			txdbInitialized = false

			if tt.setupFunc != nil {
				if err := tt.setupFunc(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			db, err := NewTestDB()

			if (err != nil) != tt.wantErr {
				t.Errorf("NewTestDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, db, err)
			}

			if db != nil {
				db.Close()
			}
		})
	}
}

func init() {
	txdb.Register("txdb", "mysql", "mockdsn")
}
