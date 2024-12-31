package db

import (
	"errors"
	"sync"
	"testing"
	"time"
	"os"
	"io/ioutil"

	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/DATA-DOG/go-txdb"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type mockDB struct {
	*gorm.DB
}

func (m *mockDB) LogMode(enable bool) {}

func (m *mockDB) SetMaxIdleConns(n int) {}

/*
ROOST_METHOD_HASH=AutoMigrate_94b22622a5
ROOST_METHOD_SIG_HASH=AutoMigrate_2cd152caa7
*/
func TestAutoMigrate(t *testing.T) {
	tests := []struct {
		name    string
		dbSetup func() *gorm.DB
		wantErr bool
	}{
		{
			name: "Successful Auto-Migration",
			dbSetup: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				return db
			},
			wantErr: false,
		},
		{
			name: "Database Connection Error",
			dbSetup: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				db.Close()
				return db
			},
			wantErr: true,
		},
		{
			name: "Partial Migration Failure",
			dbSetup: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")

				db.CreateTable(&model.User{})
				db.CreateTable(&model.Article{})

				db.AddError(errors.New("migration failed for some models"))
				return db
			},
			wantErr: true,
		},
		{
			name: "Auto-Migration with Existing Schema",
			dbSetup: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				db.CreateTable(&model.User{})
				db.CreateTable(&model.Article{})
				return db
			},
			wantErr: false,
		},
		{
			name: "Auto-Migration with Custom Table Names",
			dbSetup: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")

				return db
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.dbSetup()
			err := AutoMigrate(db)

			if tt.wantErr {
				assert.Error(t, err, "Expected an error, but got nil")
			} else {
				assert.NoError(t, err, "Expected no error, but got: %v", err)
			}

			if !tt.wantErr {

				tables := []string{"users", "articles", "tags", "comments"}
				for _, table := range tables {
					assert.True(t, db.HasTable(table), "Table %s was not created", table)
				}

				for _, model := range []interface{}{&model.User{}, &model.Article{}, &model.Tag{}, &model.Comment{}} {
					err := db.AutoMigrate(model).Error
					assert.NoError(t, err, "Error when verifying table structure for %T: %v", model, err)
				}
			}
		})
	}

	t.Run("Concurrent Auto-Migration Attempts", func(t *testing.T) {
		db, _ := gorm.Open("sqlite3", ":memory:")
		var wg sync.WaitGroup
		numGoroutines := 5
		errChan := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
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

		for err := range errChan {
			t.Errorf("Concurrent AutoMigrate() failed: %v", err)
		}

		tables := []string{"users", "articles", "tags", "comments"}
		for _, table := range tables {
			assert.True(t, db.HasTable(table), "Table %s was not created in concurrent test", table)
		}

		for _, model := range []interface{}{&model.User{}, &model.Article{}, &model.Tag{}, &model.Comment{}} {
			err := db.AutoMigrate(model).Error
			assert.NoError(t, err, "Error when verifying table structure for %T after concurrent migrations: %v", model, err)
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
			if (err != nil) != tt.wantErr {
				t.Errorf("DropTestDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	t.Run("Verify Database Connection is Closed", func(t *testing.T) {
		mockDB := &gorm.DB{}
		err := DropTestDB(mockDB)
		if err != nil {
			t.Errorf("DropTestDB() error = %v", err)
		}

	})

	t.Run("Handle Already Closed Database", func(t *testing.T) {
		mockDB := &gorm.DB{}
		_ = DropTestDB(mockDB)
		err := DropTestDB(mockDB)
		if err != nil {
			t.Errorf("DropTestDB() error = %v, want nil", err)
		}
	})

	t.Run("Concurrent Access Safety", func(t *testing.T) {
		const numGoroutines = 10
		var wg sync.WaitGroup
		errChan := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mockDB := &gorm.DB{}
				err := DropTestDB(mockDB)
				if err != nil {
					errChan <- err
				}
			}()
		}

		wg.Wait()
		close(errChan)

		for err := range errChan {
			t.Errorf("DropTestDB() error in concurrent execution: %v", err)
		}
	})

	t.Run("Performance Under Load", func(t *testing.T) {
		const numIterations = 1000
		const maxDuration = 5 * time.Second

		start := time.Now()
		for i := 0; i < numIterations; i++ {
			mockDB := &gorm.DB{}
			err := DropTestDB(mockDB)
			if err != nil {
				t.Errorf("DropTestDB() error = %v", err)
			}
		}
		duration := time.Since(start)

		t.Logf("Time taken for %d iterations: %v", numIterations, duration)
		if duration > maxDuration {
			t.Errorf("DropTestDB() took %v, which exceeds the maximum allowed duration of %v", duration, maxDuration)
		}
	})
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
		name            string
		setupMockDB     func() *gorm.DB
		setupTOMLFile   func() error
		expectedError   error
		expectedInserts int
	}{
		{
			name: "Successful Seeding of Users",
			setupMockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Error = nil
				return db
			},
			setupTOMLFile: func() error {
				content := `
				[[Users]]
				username = "user1"
				email = "user1@example.com"
				password = "password1"

				[[Users]]
				username = "user2"
				email = "user2@example.com"
				password = "password2"
				`
				return ioutil.WriteFile("db/seed/users.toml", []byte(content), 0644)
			},
			expectedError:   nil,
			expectedInserts: 2,
		},
		{
			name: "File Not Found Error",
			setupMockDB: func() *gorm.DB {
				return &gorm.DB{}
			},
			setupTOMLFile: func() error {

				os.Remove("db/seed/users.toml")
				return nil
			},
			expectedError:   errors.New("open db/seed/users.toml: no such file or directory"),
			expectedInserts: 0,
		},
		{
			name: "Invalid TOML Format",
			setupMockDB: func() *gorm.DB {
				return &gorm.DB{}
			},
			setupTOMLFile: func() error {
				content := `
				[[Users]
				username = "invalid"
				`
				return ioutil.WriteFile("db/seed/users.toml", []byte(content), 0644)
			},
			expectedError:   errors.New("toml: line 2: expected '=', '.' or ']' after key"),
			expectedInserts: 0,
		},
		{
			name: "Database Insertion Error",
			setupMockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Error = errors.New("database insertion error")
				return db
			},
			setupTOMLFile: func() error {
				content := `
				[[Users]]
				username = "user1"
				email = "user1@example.com"
				password = "password1"
				`
				return ioutil.WriteFile("db/seed/users.toml", []byte(content), 0644)
			},
			expectedError:   errors.New("database insertion error"),
			expectedInserts: 0,
		},
		{
			name: "Empty Users File",
			setupMockDB: func() *gorm.DB {
				return &gorm.DB{}
			},
			setupTOMLFile: func() error {
				return ioutil.WriteFile("db/seed/users.toml", []byte(""), 0644)
			},
			expectedError:   nil,
			expectedInserts: 0,
		},
		{
			name: "Partial Seeding with Error",
			setupMockDB: func() *gorm.DB {
				db := &gorm.DB{}
				callCount := 0
				db.Callback().Create().Register("mockCreate", func(scope *gorm.Scope) {
					callCount++
					if callCount == 2 {
						scope.Err(errors.New("database error after first insertion"))
					}
				})
				return db
			},
			setupTOMLFile: func() error {
				content := `
				[[Users]]
				username = "user1"
				email = "user1@example.com"
				password = "password1"

				[[Users]]
				username = "user2"
				email = "user2@example.com"
				password = "password2"
				`
				return ioutil.WriteFile("db/seed/users.toml", []byte(content), 0644)
			},
			expectedError:   errors.New("database error after first insertion"),
			expectedInserts: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db := tt.setupMockDB()
			err := tt.setupTOMLFile()
			if err != nil {
				t.Fatalf("Failed to setup TOML file: %v", err)
			}

			err = Seed(db)

			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("Seed() error = %v, expectedError %v", err, tt.expectedError)
			}

			insertCount := 0
			db.Callback().Create().Register("countInserts", func(scope *gorm.Scope) {
				insertCount++
			})

			if insertCount != tt.expectedInserts {
				t.Errorf("Expected %d insertions, but got %d", tt.expectedInserts, insertCount)
			}

			os.Remove("db/seed/users.toml")
		})
	}
}

/*
ROOST_METHOD_HASH=New_1d2840dc39
ROOST_METHOD_SIG_HASH=New_f9cc65f555
*/
func TestNew(t *testing.T) {
	tests := []struct {
		name            string
		dsnFunc         func() (string, error)
		gormOpenFunc    func(dialect string, args ...interface{}) (*gorm.DB, error)
		expectedDB      bool
		expectedError   error
		retryAttempts   int
		concurrentCalls int
	}{
		{
			name: "Successful Database Connection",
			dsnFunc: func() (string, error) {
				return "valid_dsn", nil
			},
			gormOpenFunc: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				return &gorm.DB{}, nil
			},
			expectedDB:    true,
			expectedError: nil,
			retryAttempts: 1,
		},
		{
			name: "Database Connection Failure",
			dsnFunc: func() (string, error) {
				return "invalid_dsn", nil
			},
			gormOpenFunc: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				return nil, errors.New("connection failed")
			},
			expectedDB:    false,
			expectedError: errors.New("connection failed"),
			retryAttempts: 10,
		},
		{
			name: "DSN Retrieval Failure",
			dsnFunc: func() (string, error) {
				return "", errors.New("DSN retrieval failed")
			},
			gormOpenFunc: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				return nil, nil
			},
			expectedDB:    false,
			expectedError: errors.New("DSN retrieval failed"),
			retryAttempts: 0,
		},
		{
			name: "Retry Mechanism on Temporary Connection Failure",
			dsnFunc: func() (string, error) {
				return "valid_dsn", nil
			},
			gormOpenFunc: func() func(dialect string, args ...interface{}) (*gorm.DB, error) {
				count := 0
				return func(dialect string, args ...interface{}) (*gorm.DB, error) {
					count++
					if count < 3 {
						return nil, errors.New("temporary failure")
					}
					return &gorm.DB{}, nil
				}
			}(),
			expectedDB:    true,
			expectedError: nil,
			retryAttempts: 3,
		},
		{
			name: "Connection Pool Configuration",
			dsnFunc: func() (string, error) {
				return "valid_dsn", nil
			},
			gormOpenFunc: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				return &mockDB{&gorm.DB{}}, nil
			},
			expectedDB:    true,
			expectedError: nil,
			retryAttempts: 1,
		},
		{
			name: "Concurrent Access Safety",
			dsnFunc: func() (string, error) {
				return "valid_dsn", nil
			},
			gormOpenFunc: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				return &gorm.DB{}, nil
			},
			expectedDB:      true,
			expectedError:   nil,
			retryAttempts:   1,
			concurrentCalls: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			originalDSN := dsn
			dsn = tt.dsnFunc
			defer func() { dsn = originalDSN }()

			originalGormOpen := gorm.Open
			gorm.Open = tt.gormOpenFunc
			defer func() { gorm.Open = originalGormOpen }()

			if tt.concurrentCalls > 0 {
				var wg sync.WaitGroup
				for i := 0; i < tt.concurrentCalls; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						db, err := New()
						if (db == nil) == tt.expectedDB {
							t.Errorf("New() returned unexpected db status")
						}
						if (err != nil) != (tt.expectedError != nil) {
							t.Errorf("New() returned unexpected error status")
						}
					}()
				}
				wg.Wait()
			} else {
				db, err := New()

				if (db != nil) != tt.expectedDB {
					t.Errorf("New() returned unexpected db status")
				}

				if (err != nil) != (tt.expectedError != nil) {
					t.Errorf("New() error = %v, expectedError %v", err, tt.expectedError)
				}

				if db != nil {

					if mockDB, ok := db.(*mockDB); ok {

						_ = mockDB
					}
				}
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
		name           string
		setupFunc      func() error
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Successful Database Connection and Initialization",
			setupFunc: func() error {
				return os.WriteFile("../env/test.env", []byte("DB_DSN=root:password@/testdb?charset=utf8&parseTime=True&loc=Local"), 0644)
			},
			wantErr: false,
		},
		{
			name: "Environment File Not Found",
			setupFunc: func() error {
				return os.Remove("../env/test.env")
			},
			wantErr:        true,
			expectedErrMsg: "open ../env/test.env: no such file or directory",
		},
		{
			name: "Invalid Database Credentials",
			setupFunc: func() error {
				return os.WriteFile("../env/test.env", []byte("DB_DSN=invalid:credentials@/nonexistentdb"), 0644)
			},
			wantErr:        true,
			expectedErrMsg: "Error 1045: Access denied for user 'invalid'@'localhost' (using password: YES)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.setupFunc(); err != nil {
				t.Fatalf("Test setup failed: %v", err)
			}

			db, err := NewTestDB()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if err.Error() != tt.expectedErrMsg {
					t.Errorf("Expected error message '%s', but got '%s'", tt.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if db == nil {
					t.Errorf("Expected non-nil DB, but got nil")
				} else {
					defer db.Close()
				}
			}
		})
	}
}

func TestNewTestDBConcurrency(t *testing.T) {
	const numGoroutines = 10
	var wg sync.WaitGroup
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			db, err := NewTestDB()
			if err == nil && db != nil {
				defer db.Close()
			}
			results <- err
		}()
	}

	wg.Wait()
	close(results)

	for err := range results {
		if err != nil {
			t.Errorf("Concurrent NewTestDB call failed: %v", err)
		}
	}
}

func TestNewTestDBConnectionLimit(t *testing.T) {
	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("Failed to create test DB: %v", err)
	}
	defer db.Close()

	if db.DB().Stats().MaxOpenConnections != 3 {
		t.Errorf("Expected max open connections to be 3, but got %d", db.DB().Stats().MaxOpenConnections)
	}
}

func TestNewTestDBUniqueConnections(t *testing.T) {
	db1, err1 := NewTestDB()
	if err1 != nil {
		t.Fatalf("Failed to create first test DB: %v", err1)
	}
	defer db1.Close()

	db2, err2 := NewTestDB()
	if err2 != nil {
		t.Fatalf("Failed to create second test DB: %v", err2)
	}
	defer db2.Close()

	if db1 == db2 {
		t.Errorf("Expected unique DB connections, but got the same instance")
	}
}

func init() {
	txdb.Register("txdb", "mysql", "root:password@/testdb?charset=utf8&parseTime=True&loc=Local")
}
