package auth

import (
	"testing"
	"time"
	"math"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"fmt"
	"os"
	"sync"
	"context"
	grpc_metadata "google.golang.org/grpc/metadata"
)

/*
ROOST_METHOD_HASH=generateToken_2cc40e0108
ROOST_METHOD_SIG_HASH=generateToken_9de4114fe8


 */
func TestgenerateToken(t *testing.T) {

	type testCase struct {
		name        string
		id          uint
		currentTime time.Time
		wantErr     bool
		errMsg      string
	}

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []testCase{
		{
			name:        "Successful Token Generation",
			id:          123,
			currentTime: fixedTime,
			wantErr:     false,
		},
		{
			name:        "Zero Value User ID",
			id:          0,
			currentTime: fixedTime,
			wantErr:     false,
		},
		{
			name:        "Maximum Value User ID",
			id:          math.MaxUint32,
			currentTime: fixedTime,
			wantErr:     false,
		},
		{
			name:        "Zero Time Value",
			id:          123,
			currentTime: time.Time{},
			wantErr:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			token, err := generateToken(tc.id, tc.currentTime)

			if (err != nil) != tc.wantErr {
				t.Errorf("generateToken() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {

				if token == "" {
					t.Error("generateToken() returned empty token")
					return
				}

				parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})

				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
					return
				}

				if claims, ok := parsedToken.Claims.(*claims); ok {

					if claims.UserID != tc.id {
						t.Errorf("Token UserID = %v, want %v", claims.UserID, tc.id)
					}

					expectedExp := tc.currentTime.Add(time.Hour * 72).Unix()
					if claims.ExpiresAt != expectedExp {
						t.Errorf("Token ExpiresAt = %v, want %v", claims.ExpiresAt, expectedExp)
					}

					t.Logf("Successfully validated token for ID: %d", tc.id)
				} else {
					t.Error("Failed to assert token claims type")
				}
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GenerateToken_b7f5ef3740
ROOST_METHOD_SIG_HASH=GenerateToken_d10a3e47a3


 */
func TestGenerateToken(t *testing.T) {

	type testCase struct {
		name          string
		userID        uint
		setupEnv      func()
		cleanupEnv    func()
		expectedError bool
		validate      func(t *testing.T, token string, err error)
	}

	validateToken := func(token string, expectedID uint) error {
		parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil {
			return err
		}

		if claims, ok := parsedToken.Claims.(*claims); ok {
			if claims.UserID != expectedID {
				return fmt.Errorf("expected user ID %d, got %d", expectedID, claims.UserID)
			}
			return nil
		}
		return errors.New("invalid token claims")
	}

	tests := []testCase{
		{
			name:   "Successful Token Generation",
			userID: 1,
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test-secret")
			},
			cleanupEnv: func() {
				os.Unsetenv("JWT_SECRET")
			},
			expectedError: false,
			validate: func(t *testing.T, token string, err error) {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if token == "" {
					t.Error("expected non-empty token")
				}
				if err := validateToken(token, 1); err != nil {
					t.Errorf("token validation failed: %v", err)
				}
			},
		},
		{
			name:   "Zero User ID",
			userID: 0,
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test-secret")
			},
			cleanupEnv: func() {
				os.Unsetenv("JWT_SECRET")
			},
			expectedError: true,
			validate: func(t *testing.T, token string, err error) {
				if err == nil {
					t.Error("expected error for zero user ID")
				}
				if token != "" {
					t.Error("expected empty token for zero user ID")
				}
			},
		},
		{
			name:   "Maximum uint Value",
			userID: math.MaxUint,
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test-secret")
			},
			cleanupEnv: func() {
				os.Unsetenv("JWT_SECRET")
			},
			expectedError: false,
			validate: func(t *testing.T, token string, err error) {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if err := validateToken(token, math.MaxUint); err != nil {
					t.Errorf("token validation failed: %v", err)
				}
			},
		},
		{
			name:   "Missing JWT Secret",
			userID: 1,
			setupEnv: func() {
				os.Unsetenv("JWT_SECRET")
			},
			cleanupEnv:    func() {},
			expectedError: true,
			validate: func(t *testing.T, token string, err error) {
				if err == nil {
					t.Error("expected error for missing JWT secret")
				}
				if token != "" {
					t.Error("expected empty token when JWT secret is missing")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			tc.setupEnv()
			defer tc.cleanupEnv()

			token, err := GenerateToken(tc.userID)

			tc.validate(t, token, err)
		})
	}

	t.Run("Multiple Sequential Tokens", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")
		defer os.Unsetenv("JWT_SECRET")

		tokens := make([]string, 3)
		for i := 0; i < 3; i++ {
			token, err := GenerateToken(1)
			if err != nil {
				t.Errorf("failed to generate token %d: %v", i, err)
			}
			tokens[i] = token
		}

		for i := 0; i < len(tokens); i++ {
			for j := i + 1; j < len(tokens); j++ {
				if tokens[i] == tokens[j] {
					t.Errorf("tokens %d and %d are identical", i, j)
				}
			}
		}
	})

	t.Run("Concurrent Token Generation", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")
		defer os.Unsetenv("JWT_SECRET")

		var wg sync.WaitGroup
		tokenChan := make(chan string, 10)
		errChan := make(chan error, 10)

		for i := uint(1); i <= 10; i++ {
			wg.Add(1)
			go func(id uint) {
				defer wg.Done()
				token, err := GenerateToken(id)
				if err != nil {
					errChan <- err
					return
				}
				tokenChan <- token
			}(i)
		}

		wg.Wait()
		close(tokenChan)
		close(errChan)

		for err := range errChan {
			t.Errorf("concurrent generation error: %v", err)
		}

		tokens := make([]string, 0)
		for token := range tokenChan {
			tokens = append(tokens, token)
		}

		if len(tokens) != 10 {
			t.Errorf("expected 10 tokens, got %d", len(tokens))
		}
	})

	t.Run("Token Expiration", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")
		defer os.Unsetenv("JWT_SECRET")

		token, err := GenerateToken(1)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil {
			t.Fatalf("failed to parse token: %v", err)
		}

		if claims, ok := parsedToken.Claims.(*claims); ok {
			expectedExp := time.Now().Add(time.Hour * 72).Unix()
			if claims.ExpiresAt < time.Now().Unix() {
				t.Error("token is already expired")
			}
			if claims.ExpiresAt > expectedExp+60 {
				t.Error("token expiration time is too far in the future")
			}
		} else {
			t.Error("failed to get token claims")
		}
	})
}

/*
ROOST_METHOD_HASH=GenerateTokenWithTime_d0df64aa69
ROOST_METHOD_SIG_HASH=GenerateTokenWithTime_72dd09cde6


 */
func TestGenerateTokenWithTime(t *testing.T) {
	type testCase struct {
		name        string
		id          uint
		timestamp   time.Time
		expectError bool
		errorMsg    string
	}

	now := time.Now()

	tests := []testCase{
		{
			name:        "Successful Token Generation",
			id:          1,
			timestamp:   now,
			expectError: false,
		},
		{
			name:        "Zero ID",
			id:          0,
			timestamp:   now,
			expectError: true,
			errorMsg:    "invalid user ID",
		},
		{
			name:        "Future Timestamp",
			id:          1,
			timestamp:   now.Add(24 * time.Hour),
			expectError: false,
		},
		{
			name:        "Past Timestamp",
			id:          1,
			timestamp:   now.Add(-24 * time.Hour),
			expectError: true,
			errorMsg:    "invalid timestamp",
		},
		{
			name:        "Maximum uint Value",
			id:          math.MaxUint32,
			timestamp:   now,
			expectError: false,
		},
		{
			name:        "Zero Time",
			id:          1,
			timestamp:   time.Time{},
			expectError: true,
			errorMsg:    "invalid time",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token, err := GenerateTokenWithTime(tc.id, tc.timestamp)

			t.Logf("Test: %s", tc.name)
			t.Logf("Input ID: %d", tc.id)
			t.Logf("Input Timestamp: %v", tc.timestamp)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if err != nil && err.Error() != tc.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tc.errorMsg, err.Error())
				}
				if token != "" {
					t.Errorf("Expected empty token but got: %s", token)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if token == "" {
					t.Error("Expected non-empty token but got empty string")
				}

				parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})

				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				if !parsedToken.Valid {
					t.Error("Token is invalid")
				}
			}
		})
	}

	t.Run("Multiple Sequential Tokens", func(t *testing.T) {
		tokens := make(map[string]bool)
		validID := uint(1)

		for i := 0; i < 5; i++ {
			token, err := GenerateTokenWithTime(validID, now.Add(time.Duration(i)*time.Minute))
			if err != nil {
				t.Errorf("Failed to generate token in sequence %d: %v", i, err)
			}

			if tokens[token] {
				t.Error("Generated duplicate token")
			}
			tokens[token] = true
		}
	})
}

/*
ROOST_METHOD_HASH=GetUserID_f2dd680cb2
ROOST_METHOD_SIG_HASH=GetUserID_e739312e3d


 */
func TestGetUserID(t *testing.T) {
	tests := []struct {
		name          string
		setupContext  func() context.Context
		expectedID    uint
		expectedError string
	}{
		{
			name: "Valid Token with Valid Claims",
			setupContext: func() context.Context {
				claims := &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(jwtSecret)
				md := grpc_metadata.Pairs("Token", tokenString)
				return grpc_metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    123,
			expectedError: "",
		},
		{
			name: "Expired Token",
			setupContext: func() context.Context {
				claims := &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(-time.Hour).Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(jwtSecret)
				md := grpc_metadata.Pairs("Token", tokenString)
				return grpc_metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "token expired",
		},
		{
			name: "Malformed Token",
			setupContext: func() context.Context {
				md := grpc_metadata.Pairs("Token", "malformed.token.string")
				return grpc_metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: it's not even a token",
		},
		{
			name: "Missing Authentication Metadata",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedID:    0,
			expectedError: "metadata not found",
		},
		{
			name: "Invalid Claims Type",
			setupContext: func() context.Context {
				token := jwt.New(jwt.SigningMethodHS256)
				tokenString, _ := token.SignedString(jwtSecret)
				md := grpc_metadata.Pairs("Token", tokenString)
				return grpc_metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: cannot map token to claims",
		},
		{
			name: "Future Token",
			setupContext: func() context.Context {
				claims := &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
						NotBefore: time.Now().Add(time.Hour).Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(jwtSecret)
				md := grpc_metadata.Pairs("Token", tokenString)
				return grpc_metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: couldn't handle this token",
		},
		{
			name: "Invalid Signature",
			setupContext: func() context.Context {
				claims := &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte("wrong-secret"))
				md := grpc_metadata.Pairs("Token", tokenString)
				return grpc_metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: couldn't handle this token",
		},
		{
			name: "Zero UserID in Valid Token",
			setupContext: func() context.Context {
				claims := &claims{
					UserID: 0,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(jwtSecret)
				md := grpc_metadata.Pairs("Token", tokenString)
				return grpc_metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running test case: %s", tt.name)

			ctx := tt.setupContext()
			gotID, err := GetUserID(ctx)

			if gotID != tt.expectedID {
				t.Errorf("GetUserID() got UserID = %v, want %v", gotID, tt.expectedID)
			}

			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("GetUserID() unexpected error = %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("GetUserID() expected error containing %q, got nil", tt.expectedError)
				} else if !contains(err.Error(), tt.expectedError) {
					t.Errorf("GetUserID() expected error containing %q, got %q", tt.expectedError, err.Error())
				}
			}

			t.Logf("Test case completed: %s", tt.name)
		})
	}
}

func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[:len(substr)] == substr
}

