package auth

import (
	"os"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
)



func TestGenerateTokenWithTime(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		t             time.Time
		expectedError bool
		validateFunc  func(string) bool
	}{
		{
			name:          "Valid Token Generation with Correct Inputs",
			id:            1,
			t:             time.Now(),
			expectedError: false,
			validateFunc:  func(token string) bool { return token != "" },
		},
		{
			name:          "Token Generation with a Historical Time Value",
			id:            2,
			t:             time.Now().AddDate(-1, 0, 0),
			expectedError: false,
			validateFunc:  func(token string) bool { return token != "" },
		},
		{
			name:          "Token Generation with Maximum Valid ID",
			id:            ^uint(0),
			t:             time.Now(),
			expectedError: false,
			validateFunc:  func(token string) bool { return token != "" },
		},
		{
			name:          "Empty JWT Secret Environment Variable",
			id:            3,
			t:             time.Now(),
			expectedError: true,
			validateFunc:  nil,
		},
		{
			name:          "Invalid ID (Zero Value)",
			id:            0,
			t:             time.Now(),
			expectedError: true,
			validateFunc:  func(token string) bool { return token == "" },
		},
		{
			name:          "Handling Future Date for Token Generation",
			id:            4,
			t:             time.Now().AddDate(1, 0, 0),
			expectedError: false,
			validateFunc:  func(token string) bool { return token != "" },
		},
	}

	originalJwtSecret := os.Getenv("JWT_SECRET")

	for _, tt := range tests {
		if tt.name == "Empty JWT Secret Environment Variable" {
			os.Setenv("JWT_SECRET", "")
		} else {
			os.Setenv("JWT_SECRET", "defaultSecret")
		}

		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running test case: %s", tt.name)

			token, err := GenerateTokenWithTime(tt.id, tt.t)

			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error: %v, got: %v", tt.expectedError, err)
			}

			if !tt.expectedError && tt.validateFunc != nil && !tt.validateFunc(token) {
				t.Errorf("Token validation failed. Generated token: %s", token)
			}

			if tt.expectedError && err == nil {
				t.Errorf("Expected an error but got none")
			} else if !tt.expectedError {
				t.Logf("Successfully generated token: %s", token)
			}
		})
	}

	os.Setenv("JWT_SECRET", originalJwtSecret)
}




