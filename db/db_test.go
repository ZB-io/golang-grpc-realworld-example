package undefined

import (
	"os"
	"testing"
	"errors"
	"fmt"
	"io/ioutil"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)





type mockDB struct {
	*gorm.DB
	createError  error
	createdUsers []model.User
}


/*
ROOST_METHOD_HASH=dsn_e202d1c4f9
ROOST_METHOD_SIG_HASH=dsn_b336e03d64

FUNCTION_DEF=func dsn() (string, error) 

 */
func TestDsn(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectedDSN string
		expectedErr string
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
			expectedDSN: "user:password@(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
			expectedErr: "",
		},
		{
			name: "Missing DB_HOST Environment Variable",
			envVars: map[string]string{
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			expectedDSN: "",
			expectedErr: "$DB_HOST is not set",
		},
		{
			name: "Missing DB_USER Environment Variable",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PASSWORD": "password",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			expectedDSN: "",
			expectedErr: "$DB_USER is not set",
		},
		{
			name: "Missing DB_PASSWORD Environment Variable",
			envVars: map[string]string{
				"DB_HOST": "localhost",
				"DB_USER": "user",
				"DB_NAME": "testdb",
				"DB_PORT": "3306",
			},
			expectedDSN: "",
			expectedErr: "$DB_PASSWORD is not set",
		},
		{
			name: "Missing DB_NAME Environment Variable",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_PORT":     "3306",
			},
			expectedDSN: "",
			expectedErr: "$DB_NAME is not set",
		},
		{
			name: "Missing DB_PORT Environment Variable",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "testdb",
			},
			expectedDSN: "",
			expectedErr: "$DB_PORT is not set",
		},
		{
			name: "Special Characters in Environment Variables",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user",
				"DB_PASSWORD": "p@ssw0rd!#$%",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			expectedDSN: "user:p@ssw0rd!#$%@(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			os.Clearenv()

			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			gotDSN, err := dsn()

			if tt.expectedErr != "" {
				if err == nil {
					t.Errorf("Expected error: %s, but got nil", tt.expectedErr)
				} else if err.Error() != tt.expectedErr {
					t.Errorf("Expected error: %s, but got: %s", tt.expectedErr, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %s", err.Error())
				}
			}

			if gotDSN != tt.expectedDSN {
				t.Errorf("Expected DSN: %s, but got: %s", tt.expectedDSN, gotDSN)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=Seed_5ad31c3a6c
ROOST_METHOD_SIG_HASH=Seed_878933cebc

FUNCTION_DEF=func Seed(db *gorm.DB) error 

 */
func (m *mockDB) Create(value interface{}) *gorm.DB {
	if m.createError != nil {
		return &gorm.DB{Error: m.createError}
	}
	m.createdUsers = append(m.createdUsers, *value.(*model.User))
	return m.DB
}

func TestSeed(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func() *mockDB
		setupFile     func() error
		expectedError error
		expectedUsers int
	}{
		{
			name: "Successful Seeding of Users",
			setupMock: func() *mockDB {
				return newMockDB()
			},
			setupFile: func() error {
				content := `
[[Users]]
username = "user1"
email = "user1@example.com"
[[Users]]
username = "user2"
email = "user2@example.com"
`
				return ioutil.WriteFile("db/seed/users.toml", []byte(content), 0644)
			},
			expectedError: nil,
			expectedUsers: 2,
		},
		{
			name: "File Not Found Error",
			setupMock: func() *mockDB {
				return newMockDB()
			},
			setupFile: func() error {
				return os.Remove("db/seed/users.toml")
			},
			expectedError: errors.New("open db/seed/users.toml: no such file or directory"),
			expectedUsers: 0,
		},
		{
			name: "Invalid TOML Format",
			setupMock: func() *mockDB {
				return newMockDB()
			},
			setupFile: func() error {
				content := `
[[Users]
invalid = "toml"
`
				return ioutil.WriteFile("db/seed/users.toml", []byte(content), 0644)
			},
			expectedError: errors.New("toml: line 2: expected '.' or ']' after table key"),
			expectedUsers: 0,
		},
		{
			name: "Database Insertion Error",
			setupMock: func() *mockDB {
				mock := newMockDB()
				mock.createError = errors.New("database insertion error")
				return mock
			},
			setupFile: func() error {
				content := `
[[Users]]
username = "user1"
email = "user1@example.com"
`
				return ioutil.WriteFile("db/seed/users.toml", []byte(content), 0644)
			},
			expectedError: errors.New("database insertion error"),
			expectedUsers: 0,
		},
		{
			name: "Empty Users File",
			setupMock: func() *mockDB {
				return newMockDB()
			},
			setupFile: func() error {
				content := ""
				return ioutil.WriteFile("db/seed/users.toml", []byte(content), 0644)
			},
			expectedError: nil,
			expectedUsers: 0,
		},
		{
			name: "Large Number of Users",
			setupMock: func() *mockDB {
				return newMockDB()
			},
			setupFile: func() error {
				content := ""
				for i := 0; i < 10000; i++ {
					content += fmt.Sprintf(`
[[Users]]
username = "user%d"
email = "user%d@example.com"
`, i, i)
				}
				return ioutil.WriteFile("db/seed/users.toml", []byte(content), 0644)
			},
			expectedError: nil,
			expectedUsers: 10000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := tt.setupMock()
			err := tt.setupFile()
			if err != nil {
				t.Fatalf("Failed to setup file: %v", err)
			}

			err = Seed(mockDB.DB)

			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
			}

			if len(mockDB.createdUsers) != tt.expectedUsers {
				t.Errorf("Expected %d users to be created, got %d", tt.expectedUsers, len(mockDB.createdUsers))
			}

			os.Remove("db/seed/users.toml")
		})
	}
}

func newMockDB() *mockDB {
	return &mockDB{
		DB: &gorm.DB{},
	}
}

