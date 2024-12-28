package db

import (
	"os"
	"testing"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func Testdsn(t *testing.T) {
	tests := []struct {
		name          string
		envVars       map[string]string
		expectedDSN   string
		expectedError string
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
			expectedDSN:   "user:password@(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
			expectedError: "",
		},
		{
			name: "Missing DB_HOST Environment Variable",
			envVars: map[string]string{
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			expectedError: "$DB_HOST is not set",
		},
		{
			name: "Missing DB_USER Environment Variable",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PASSWORD": "password",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			expectedError: "$DB_USER is not set",
		},
		{
			name: "Missing DB_PASSWORD Environment Variable",
			envVars: map[string]string{
				"DB_HOST": "localhost",
				"DB_USER": "user",
				"DB_NAME": "testdb",
				"DB_PORT": "3306",
			},
			expectedError: "$DB_PASSWORD is not set",
		},
		{
			name: "Missing DB_NAME Environment Variable",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_PORT":     "3306",
			},
			expectedError: "$DB_NAME is not set",
		},
		{
			name: "Missing DB_PORT Environment Variable",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "testdb",
			},
			expectedError: "$DB_PORT is not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for key, value := range tt.envVars {
				err := os.Setenv(key, value)
				if err != nil {
					t.Fatalf("Unable to set environment variable %s: %v", key, err)
				}
			}

			dsn, err := dsn()

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				} else {
					t.Logf("Correctly received expected error: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error received: %v", err)
			}

			if dsn != tt.expectedDSN {
				t.Errorf("Expected DSN %v, got %v", tt.expectedDSN, dsn)
			} else {
				t.Logf("Successfully received expected DSN: %v", dsn)
			}

			for key := range tt.envVars {
				_ = os.Unsetenv(key)
			}
		})
	}
}
