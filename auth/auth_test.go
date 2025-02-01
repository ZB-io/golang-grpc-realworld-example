package auth

import (
	"math"
	"os"
	"sync"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
)








/*
ROOST_METHOD_HASH=GenerateToken_b7f5ef3740
ROOST_METHOD_SIG_HASH=GenerateToken_d10a3e47a3

FUNCTION_DEF=func GenerateToken(id uint) (string, error) 

*/
func TestGenerateToken(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	tests := []struct {
		name    string
		userID  uint
		wantErr bool
		setup   func()
		verify  func(*testing.T, string)
	}{
		{
			name:    "Successful Token Generation",
			userID:  1,
			wantErr: false,
			verify: func(t *testing.T, token string) {
				if token == "" {
					t.Error("Expected non-empty token")
				}
			},
		},
		{
			name:    "Token Generation with Zero User ID",
			userID:  0,
			wantErr: true,
		},
		{
			name:    "Token Generation with Maximum uint Value",
			userID:  math.MaxUint32,
			wantErr: false,
			verify: func(t *testing.T, token string) {
				if token == "" {
					t.Error("Expected non-empty token")
				}
			},
		},
		{
			name:    "Verification of Token Content",
			userID:  42,
			wantErr: false,
			verify: func(t *testing.T, token string) {
				claims := &claims{}
				_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				if claims.UserID != 42 {
					t.Errorf("Expected UserID 42, got %d", claims.UserID)
				}
			},
		},
		{
			name:    "Token Expiration",
			userID:  1,
			wantErr: false,
			verify: func(t *testing.T, token string) {
				claims := &claims{}
				_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				if claims.ExpiresAt <= time.Now().Unix() {
					t.Error("Token should not be expired")
				}
			},
		},
		{
			name:    "Error Handling with Invalid JWT Secret",
			userID:  1,
			wantErr: true,
			setup: func() {
				os.Setenv("JWT_SECRET", "")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			got, err := GenerateToken(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.verify != nil {
				tt.verify(t, got)
			}
		})
	}

	t.Run("Concurrent Token Generation", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 100
		tokens := make([]string, numGoroutines)
		errors := make([]error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				token, err := GenerateToken(uint(index))
				tokens[index] = token
				errors[index] = err
			}(i)
		}

		wg.Wait()

		for i, err := range errors {
			if err != nil {
				t.Errorf("Goroutine %d failed: %v", i, err)
			}
			if tokens[i] == "" {
				t.Errorf("Goroutine %d returned empty token", i)
			}
		}

		uniqueTokens := make(map[string]bool)
		for _, token := range tokens {
			if uniqueTokens[token] {
				t.Error("Duplicate token found")
			}
			uniqueTokens[token] = true
		}
	})
}


/*
ROOST_METHOD_HASH=GenerateTokenWithTime_d0df64aa69
ROOST_METHOD_SIG_HASH=GenerateTokenWithTime_72dd09cde6

FUNCTION_DEF=func GenerateTokenWithTime(id uint, t time.Time) (string, error) 

*/
func TestGenerateTokenWithTime(t *testing.T) {

	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	os.Setenv("JWT_SECRET", "test_secret")

	tests := []struct {
		name    string
		id      uint
		time    time.Time
		wantErr bool
	}{
		{
			name:    "Valid User ID and Current Time",
			id:      1,
			time:    time.Now(),
			wantErr: false,
		},
		{
			name:    "Zero User ID",
			id:      0,
			time:    time.Now(),
			wantErr: true,
		},
		{
			name:    "Future Time",
			id:      2,
			time:    time.Now().Add(24 * time.Hour),
			wantErr: false,
		},
		{
			name:    "Past Time",
			id:      3,
			time:    time.Now().Add(-24 * time.Hour),
			wantErr: true,
		},
		{
			name:    "Maximum Uint Value for User ID",
			id:      ^uint(0),
			time:    time.Now(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateTokenWithTime(tt.id, tt.time)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateTokenWithTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("GenerateTokenWithTime() returned empty token")
			}
			if !tt.wantErr {

				token, err := jwt.ParseWithClaims(got, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				if claims, ok := token.Claims.(*claims); ok && token.Valid {
					if claims.UserID != tt.id {
						t.Errorf("Token has incorrect UserID. got = %v, want = %v", claims.UserID, tt.id)
					}
					expectedExpiry := tt.time.Add(time.Hour * 72).Unix()
					if claims.ExpiresAt != expectedExpiry {
						t.Errorf("Token has incorrect ExpiresAt. got = %v, want = %v", claims.ExpiresAt, expectedExpiry)
					}
				} else {
					t.Errorf("Token claims are invalid")
				}
			}
		})
	}

	t.Run("Empty JWT Secret", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "")
		_, err := GenerateTokenWithTime(1, time.Now())
		if err == nil {
			t.Errorf("GenerateTokenWithTime() did not return an error with empty JWT secret")
		}
	})

	t.Run("Generate Multiple Tokens Sequentially", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test_secret")
		tokens := make(map[string]bool)
		for i := uint(1); i <= 10; i++ {
			token, err := GenerateTokenWithTime(i, time.Now())
			if err != nil {
				t.Errorf("Failed to generate token for user %d: %v", i, err)
			}
			if tokens[token] {
				t.Errorf("Generated duplicate token: %s", token)
			}
			tokens[token] = true
		}
	})
}


/*
ROOST_METHOD_HASH=generateToken_2cc40e0108
ROOST_METHOD_SIG_HASH=generateToken_9de4114fe8

FUNCTION_DEF=func generateToken(id uint, now time.Time) (string, error) 

*/
func TestGenerateMultipleTokens(t *testing.T) {
	jwtSecret = []byte("test_secret")
	userID := uint(1234)
	baseTime := time.Now()
	tokenSet := make(map[string]bool)

	for i := 0; i < 100; i++ {
		token, err := generateToken(userID, baseTime.Add(time.Duration(i)*time.Millisecond))
		if err != nil {
			t.Errorf("Failed to generate token: %v", err)
		}
		if tokenSet[token] {
			t.Errorf("Duplicate token generated: %s", token)
		}
		tokenSet[token] = true
	}
}

func TestGenerateToken(t *testing.T) {

	fixedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	originalSecret := jwtSecret
	defer func() { jwtSecret = originalSecret }()

	tests := []struct {
		name        string
		userID      uint
		jwtSecret   []byte
		expectedErr bool
		validate    func(*testing.T, string)
	}{
		{
			name:        "Successfully Generate Token for Valid User ID",
			userID:      1234,
			jwtSecret:   []byte("test_secret"),
			expectedErr: false,
			validate: func(t *testing.T, token string) {
				if token == "" {
					t.Error("Expected non-empty token, got empty string")
				}
			},
		},
		{
			name:        "Generate Token with Zero User ID",
			userID:      0,
			jwtSecret:   []byte("test_secret"),
			expectedErr: false,
			validate: func(t *testing.T, token string) {
				if token == "" {
					t.Error("Expected non-empty token, got empty string")
				}
			},
		},
		{
			name:        "Verify Token Expiration Time",
			userID:      5678,
			jwtSecret:   []byte("test_secret"),
			expectedErr: false,
			validate: func(t *testing.T, token string) {
				claims := &claims{}
				_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				expectedExp := fixedTime.Add(time.Hour * 72).Unix()
				if claims.ExpiresAt != expectedExp {
					t.Errorf("Expected expiration %d, got %d", expectedExp, claims.ExpiresAt)
				}
			},
		},
		{
			name:        "Generate Token with Maximum uint Value",
			userID:      math.MaxUint32,
			jwtSecret:   []byte("test_secret"),
			expectedErr: false,
			validate: func(t *testing.T, token string) {
				if token == "" {
					t.Error("Expected non-empty token, got empty string")
				}
			},
		},
		{
			name:        "Attempt to Generate Token with Empty JWT Secret",
			userID:      1234,
			jwtSecret:   []byte{},
			expectedErr: true,
			validate: func(t *testing.T, token string) {
				if token != "" {
					t.Error("Expected empty token, got non-empty string")
				}
			},
		},
		{
			name:        "Verify Token Signing Method",
			userID:      9012,
			jwtSecret:   []byte("test_secret"),
			expectedErr: false,
			validate: func(t *testing.T, token string) {
				parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if _, ok := parsedToken.Method.(*jwt.SigningMethodHMAC); !ok {
					t.Errorf("Expected signing method HMAC, got %v", parsedToken.Method)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret = tt.jwtSecret

			token, err := generateToken(tt.userID, fixedTime)

			if (err != nil) != tt.expectedErr {
				t.Errorf("generateToken() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if tt.validate != nil {
				tt.validate(t, token)
			}
		})
	}
}

