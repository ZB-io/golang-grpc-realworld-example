package auth

import (
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"math"
	"os"
	"strings"
	"sync"
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

/*
ROOST_METHOD_HASH=generateToken_2cc40e0108
ROOST_METHOD_SIG_HASH=generateToken_9de4114fe8


 */
func TestgenerateToken(t *testing.T) {

	tests := []struct {
		name       string
		id         uint
		timeInput  time.Time
		wantErr    bool
		validateFn func(*testing.T, string, uint, time.Time)
	}{
		{
			name:      "Successful Token Generation",
			id:        1,
			timeInput: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:   false,
			validateFn: func(t *testing.T, token string, id uint, inputTime time.Time) {
				if token == "" {
					t.Error("expected non-empty token")
				}

				parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})

				if err != nil {
					t.Errorf("failed to parse token: %v", err)
				}

				if claims, ok := parsedToken.Claims.(*claims); ok {
					if claims.ID != id {
						t.Errorf("expected ID %d, got %d", id, claims.ID)
					}

					expectedExp := inputTime.Add(time.Hour * 72).Unix()
					if claims.ExpiresAt != expectedExp {
						t.Errorf("expected expiration %d, got %d", expectedExp, claims.ExpiresAt)
					}
				}
			},
		},
		{
			name:      "Zero Value User ID",
			id:        0,
			timeInput: time.Now(),
			wantErr:   false,
			validateFn: func(t *testing.T, token string, id uint, _ time.Time) {
				parsedToken, _ := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})

				if claims, ok := parsedToken.Claims.(*claims); ok && claims.ID != 0 {
					t.Error("expected ID to be 0")
				}
			},
		},
		{
			name:      "Maximum uint Value",
			id:        ^uint(0),
			timeInput: time.Now(),
			wantErr:   false,
			validateFn: func(t *testing.T, token string, id uint, _ time.Time) {
				parsedToken, _ := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})

				if claims, ok := parsedToken.Claims.(*claims); ok && claims.ID != id {
					t.Errorf("expected ID %d, got %d", id, claims.ID)
				}
			},
		},
		{
			name:      "Different Time Zones",
			id:        1,
			timeInput: time.Now().In(time.FixedZone("GMT+8", 8*60*60)),
			wantErr:   false,
			validateFn: func(t *testing.T, token string, _ uint, inputTime time.Time) {
				parsedToken, _ := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})

				if claims, ok := parsedToken.Claims.(*claims); ok {
					expectedExp := inputTime.Add(time.Hour * 72).Unix()
					if claims.ExpiresAt != expectedExp {
						t.Errorf("expected expiration %d, got %d", expectedExp, claims.ExpiresAt)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Executing test case: %s", tt.name)

			token, err := generateToken(tt.id, tt.timeInput)

			if (err != nil) != tt.wantErr {
				t.Errorf("generateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validateFn != nil {
				tt.validateFn(t, token, tt.id, tt.timeInput)
			}

			t.Logf("Test case completed: %s", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=GenerateToken_b7f5ef3740
ROOST_METHOD_SIG_HASH=GenerateToken_d10a3e47a3


 */
func TestGenerateToken(t *testing.T) {

	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	os.Setenv("JWT_SECRET", "test-secret-key")

	tests := []struct {
		name    string
		userID  uint
		wantErr bool
		setup   func()
		cleanup func()
	}{
		{
			name:    "Valid user ID",
			userID:  1,
			wantErr: false,
			setup:   func() {},
			cleanup: func() {},
		},
		{
			name:    "Zero user ID",
			userID:  0,
			wantErr: true,
			setup:   func() {},
			cleanup: func() {},
		},
		{
			name:    "Maximum uint value",
			userID:  math.MaxUint,
			wantErr: false,
			setup:   func() {},
			cleanup: func() {},
		},
		{
			name:    "Missing JWT secret",
			userID:  1,
			wantErr: true,
			setup: func() {
				os.Unsetenv("JWT_SECRET")
			},
			cleanup: func() {
				os.Setenv("JWT_SECRET", "test-secret-key")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()

			token, err := GenerateToken(tt.userID)

			t.Logf("Test case: %s", tt.name)
			t.Logf("Input userID: %d", tt.userID)
			t.Logf("Generated token: %s", token)
			if err != nil {
				t.Logf("Error: %v", err)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if token == "" {
					t.Error("Expected non-empty token, got empty string")
				}

				parts := strings.Split(token, ".")
				if len(parts) != 3 {
					t.Error("Token does not follow JWT format")
				}
			}
		})
	}
}

func TestGenerateTokenConcurrent(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Setenv("JWT_SECRET", "")

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
		t.Errorf("Concurrent token generation error: %v", err)
	}

	tokens := make(map[string]bool)
	for token := range tokenChan {
		if tokens[token] {
			t.Error("Generated duplicate token in concurrent execution")
		}
		tokens[token] = true
	}
}

func TestGenerateTokenMultiple(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Setenv("JWT_SECRET", "")

	tokens := make(map[string]bool)
	userID := uint(1)

	for i := 0; i < 5; i++ {
		token, err := GenerateToken(userID)
		if err != nil {
			t.Errorf("Failed to generate token on iteration %d: %v", i, err)
			continue
		}

		if tokens[token] {
			t.Error("Generated duplicate token")
		}
		tokens[token] = true
	}
}

/*
ROOST_METHOD_HASH=GenerateTokenWithTime_d0df64aa69
ROOST_METHOD_SIG_HASH=GenerateTokenWithTime_72dd09cde6


 */
func TestGenerateTokenWithTime(t *testing.T) {

	tests := []struct {
		name        string
		id          uint
		timestamp   time.Time
		expectError bool
		scenario    string
	}{
		{
			name:        "Successful Token Generation",
			id:          1,
			timestamp:   time.Now(),
			expectError: false,
			scenario:    "Scenario 1: Valid ID and current time",
		},
		{
			name:        "Zero ID",
			id:          0,
			timestamp:   time.Now(),
			expectError: true,
			scenario:    "Scenario 2: Invalid zero ID",
		},
		{
			name:        "Future Timestamp",
			id:          1,
			timestamp:   time.Now().Add(24 * time.Hour),
			expectError: false,
			scenario:    "Scenario 3: Future timestamp",
		},
		{
			name:        "Past Timestamp",
			id:          1,
			timestamp:   time.Now().Add(-24 * time.Hour),
			expectError: false,
			scenario:    "Scenario 4: Past timestamp",
		},
		{
			name:        "Maximum uint Value",
			id:          math.MaxUint32,
			timestamp:   time.Now(),
			expectError: false,
			scenario:    "Scenario 5: Maximum uint value",
		},
		{
			name:        "Zero Time",
			id:          1,
			timestamp:   time.Time{},
			expectError: true,
			scenario:    "Scenario 6: Zero time value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Executing %s", tt.scenario)

			token, err := GenerateTokenWithTime(tt.id, tt.timestamp)

			if tt.expectError {
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

			t.Logf("Token generation result: %v", err == nil)
		})
	}

	t.Run("Concurrent Token Generation", func(t *testing.T) {
		t.Log("Executing Scenario 7: Concurrent token generation")

		const numGoroutines = 10
		var wg sync.WaitGroup
		tokens := make([]string, numGoroutines)
		errors := make([]error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				token, err := GenerateTokenWithTime(uint(index+1), time.Now())
				tokens[index] = token
				errors[index] = err
			}(i)
		}

		wg.Wait()

		tokenMap := make(map[string]bool)
		for i, token := range tokens {
			if errors[i] != nil {
				t.Errorf("Goroutine %d failed with error: %v", i, errors[i])
			}
			if token == "" {
				t.Errorf("Goroutine %d generated empty token", i)
			}
			if tokenMap[token] {
				t.Errorf("Duplicate token detected: %v", token)
			}
			tokenMap[token] = true
		}

		t.Log("Concurrent token generation completed")
	})
}

/*
ROOST_METHOD_HASH=GetUserID_f2dd680cb2
ROOST_METHOD_SIG_HASH=GetUserID_e739312e3d


 */
func TestGetUserID(t *testing.T) {

	jwtSecret = []byte("test-secret-key")

	tests := []struct {
		name          string
		setupContext  func() context.Context
		expectedID    uint
		expectedError string
	}{
		{
			name: "Valid Token with Valid Claims",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
					UserID: 123,
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.New(map[string]string{
					"Token": tokenString,
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    123,
			expectedError: "",
		},
		{
			name: "Missing Authentication Token",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedID:    0,
			expectedError: "metadata not found",
		},
		{
			name: "Expired Token",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(-time.Hour).Unix(),
					},
					UserID: 123,
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.New(map[string]string{
					"Token": tokenString,
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "token expired",
		},
		{
			name: "Malformed Token",
			setupContext: func() context.Context {
				md := metadata.New(map[string]string{
					"Token": "invalid-token-string",
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: it's not even a token",
		},
		{
			name: "Invalid Token Signature",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
					UserID: 123,
				})
				tokenString, _ := token.SignedString([]byte("wrong-secret"))
				md := metadata.New(map[string]string{
					"Token": tokenString,
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: couldn't handle this token",
		},
		{
			name: "Invalid Claims Structure",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"exp": time.Now().Add(time.Hour).Unix(),
					"uid": "invalid-type",
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.New(map[string]string{
					"Token": tokenString,
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: cannot map token to claims",
		},
		{
			name: "Future Token",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
						NotBefore: time.Now().Add(time.Hour).Unix(),
					},
					UserID: 123,
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.New(map[string]string{
					"Token": tokenString,
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "token expired",
		},
		{
			name: "Zero UserID in Valid Token",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
					UserID: 0,
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.New(map[string]string{
					"Token": tokenString,
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Testing scenario:", tt.name)

			ctx := tt.setupContext()
			userID, err := GetUserID(ctx)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Equal(t, uint(0), userID)
				t.Logf("Expected error received: %v", err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, userID)
				t.Logf("Successfully retrieved UserID: %d", userID)
			}
		})
	}
}

