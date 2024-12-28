package db

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func Testdsn(t *testing.T) {
	tests := []struct {
		name        string
		envValues   map[string]string
		expectedDSN string
		expectedErr error
	}{
		{
			name: "All environment variables set correctly",
			envValues: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "dbname",
				"DB_PORT":     "3306",
			},
			expectedDSN: "user:password@(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
			expectedErr: nil,
		},
		{
			name: "Missing DB_HOST environment variable",
			envValues: map[string]string{
				"DB_HOST":     "",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "dbname",
				"DB_PORT":     "3306",
			},
			expectedDSN: "",
			expectedErr: errors.New("$DB_HOST is not set"),
		},
		{
			name: "Missing DB_USER environment variable",
			envValues: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "",
				"DB_PASSWORD": "password",
				"DB_NAME":     "dbname",
				"DB_PORT":     "3306",
			},
			expectedDSN: "",
			expectedErr: errors.New("$DB_USER is not set"),
		},
		{
			name: "Missing DB_PASSWORD environment variable",
			envValues: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user",
				"DB_PASSWORD": "",
				"DB_NAME":     "dbname",
				"DB_PORT":     "3306",
			},
			expectedDSN: "",
			expectedErr: errors.New("$DB_PASSWORD is not set"),
		},
		{
			name: "Missing DB_NAME environment variable",
			envValues: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "",
				"DB_PORT":     "3306",
			},
			expectedDSN: "",
			expectedErr: errors.New("$DB_NAME is not set"),
		},
		{
			name: "Missing DB_PORT environment variable",
			envValues: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "dbname",
				"DB_PORT":     "",
			},
			expectedDSN: "",
			expectedErr: errors.New("$DB_PORT is not set"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			for key, value := range test.envValues {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			dsn, err := dsn()

			if dsn != test.expectedDSN {
				t.Errorf("expected DSN: %s, got: %s", test.expectedDSN, dsn)
			}
			if (err != nil) != (test.expectedErr != nil) || (err != nil && err.Error() != test.expectedErr.Error()) {
				t.Errorf("expected error: %v, got: %v", test.expectedErr, err)
			}

			t.Logf("Test scenario '%s': Passed", test.name)
		})
	}
}






