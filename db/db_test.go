package undefined

import (
	"os"
	"testing"
)








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
			name: "Correct DSN Format",
			envVars: map[string]string{
				"DB_HOST":     "testhost",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
				"DB_PORT":     "5432",
			},
			expectedDSN: "testuser:testpass@(testhost:5432)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
			expectedErr: "",
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
				if err == nil || err.Error() != tt.expectedErr {
					t.Errorf("dsn() error = %v, expectedErr %v", err, tt.expectedErr)
				}
			} else if err != nil {
				t.Errorf("dsn() unexpected error: %v", err)
			}

			if gotDSN != tt.expectedDSN {
				t.Errorf("dsn() gotDSN = %v, want %v", gotDSN, tt.expectedDSN)
			}
		})
	}
}

