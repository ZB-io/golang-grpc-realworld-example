package auth

import (
	"math"
	"os"
	"testing"
	"time"
)








/*
ROOST_METHOD_HASH=GenerateTokenWithTime_d0df64aa69
ROOST_METHOD_SIG_HASH=GenerateTokenWithTime_72dd09cde6

FUNCTION_DEF=func GenerateTokenWithTime(id uint, t time.Time) (string, error) 

 */
func TestGenerateTokenWithTime(t *testing.T) {

	type testCase struct {
		name        string
		userID      uint
		inputTime   time.Time
		setupEnv    func()
		cleanupEnv  func()
		expectError bool
	}

	setJWTSecret := func(secret string) {
		os.Setenv("JWT_SECRET", secret)

		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}

	cleanJWTSecret := func() {
		os.Unsetenv("JWT_SECRET")
		jwtSecret = []byte{}
	}

	tests := []testCase{
		{
			name:   "Successful token generation with valid ID and current time",
			userID: 1,
			setupEnv: func() {
				setJWTSecret("test-secret")
			},
			cleanupEnv: func() {
				cleanJWTSecret()
			},
			inputTime:   time.Now(),
			expectError: false,
		},
		{
			name:   "Generate token with zero user ID",
			userID: 0,
			setupEnv: func() {
				setJWTSecret("test-secret")
			},
			cleanupEnv: func() {
				cleanJWTSecret()
			},
			inputTime:   time.Now(),
			expectError: false,
		},
		{
			name:   "Generate token with future time",
			userID: 1,
			setupEnv: func() {
				setJWTSecret("test-secret")
			},
			cleanupEnv: func() {
				cleanJWTSecret()
			},
			inputTime:   time.Now().Add(24 * time.Hour),
			expectError: false,
		},
		{
			name:   "Generate token with past time",
			userID: 1,
			setupEnv: func() {
				setJWTSecret("test-secret")
			},
			cleanupEnv: func() {
				cleanJWTSecret()
			},
			inputTime:   time.Now().Add(-24 * time.Hour),
			expectError: false,
		},
		{
			name:   "Generate token with missing JWT secret",
			userID: 1,
			setupEnv: func() {
				cleanJWTSecret()
			},
			cleanupEnv:  func() {},
			inputTime:   time.Now(),
			expectError: true,
		},
		{
			name:   "Generate token with maximum uint value",
			userID: math.MaxUint32,
			setupEnv: func() {
				setJWTSecret("test-secret")
			},
			cleanupEnv: func() {
				cleanJWTSecret()
			},
			inputTime:   time.Now(),
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			tc.setupEnv()
			defer tc.cleanupEnv()

			token, err := GenerateTokenWithTime(tc.userID, tc.inputTime)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if token != "" {
					t.Errorf("Expected empty token but got: %v", token)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if token == "" {
					t.Error("Expected non-empty token but got empty string")
				}

			}

			t.Logf("Test case '%s' completed", tc.name)
		})
	}

	t.Run("Generate multiple tokens sequentially", func(t *testing.T) {
		setJWTSecret("test-secret")
		defer cleanJWTSecret()

		userID := uint(1)
		tokens := make(map[string]bool)

		for i := 0; i < 5; i++ {
			timeStamp := time.Now().Add(time.Duration(i) * time.Hour)
			token, err := GenerateTokenWithTime(userID, timeStamp)

			if err != nil {
				t.Errorf("Failed to generate token %d: %v", i, err)
			}

			if token == "" {
				t.Errorf("Empty token generated for iteration %d", i)
			}

			if tokens[token] {
				t.Errorf("Duplicate token generated: %s", token)
			}

			tokens[token] = true
		}
	})
}

