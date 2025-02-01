package auth

import (
	"os"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"math"
	"strings"
)








/*
ROOST_METHOD_HASH=GenerateToken_b7f5ef3740
ROOST_METHOD_SIG_HASH=GenerateToken_d10a3e47a3

FUNCTION_DEF=func GenerateToken(id uint) (string, error) 

*/
func TestGenerateToken(t *testing.T) {

	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	os.Setenv("JWT_SECRET", "test_secret")

	tests := []struct {
		name    string
		userID  uint
		wantErr bool
	}{
		{
			name:    "Successfully Generate Token for Valid User ID",
			userID:  1234,
			wantErr: false,
		},
		{
			name:    "Attempt to Generate Token with Zero User ID",
			userID:  0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if token == "" {
					t.Errorf("GenerateToken() returned an empty token")
				}

				claims := &claims{}
				_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})

				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}

				if claims.UserID != tt.userID {
					t.Errorf("Token UserID = %v, want %v", claims.UserID, tt.userID)
				}

				if claims.ExpiresAt == 0 {
					t.Errorf("Token does not have an expiration time")
				}
			}
		})
	}
}

func TestGenerateTokenPerformance(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	start := time.Now()
	for i := 0; i < 1000; i++ {
		_, err := GenerateToken(uint(i))
		if err != nil {
			t.Errorf("Failed to generate token: %v", err)
		}
	}
	duration := time.Since(start)

	if duration > 1*time.Second {
		t.Errorf("Token generation took too long: %v", duration)
	}
}

func TestGenerateTokenUniqueness(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	token1, err1 := GenerateToken(1234)
	token2, err2 := GenerateToken(5678)

	if err1 != nil || err2 != nil {
		t.Errorf("Failed to generate tokens: %v, %v", err1, err2)
	}

	if token1 == token2 {
		t.Errorf("Tokens for different user IDs are not unique")
	}
}

func TestGenerateTokenWithoutSecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	_, err := GenerateToken(1234)
	if err == nil {
		t.Errorf("Expected error when JWT_SECRET is not set, but got nil")
	}
}


/*
ROOST_METHOD_HASH=GenerateTokenWithTime_d0df64aa69
ROOST_METHOD_SIG_HASH=GenerateTokenWithTime_72dd09cde6

FUNCTION_DEF=func GenerateTokenWithTime(id uint, t time.Time) (string, error) 

*/
func TestGenerateTokenWithTime(t *testing.T) {

	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

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
			wantErr: false,
		},
		{
			name:    "Future Time",
			id:      1,
			time:    time.Now().Add(24 * time.Hour),
			wantErr: false,
		},
		{
			name:    "Past Time",
			id:      1,
			time:    time.Now().Add(-24 * time.Hour),
			wantErr: false,
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
			token1, err1 := GenerateTokenWithTime(tt.id, tt.time)
			if (err1 != nil) != tt.wantErr {
				t.Errorf("GenerateTokenWithTime() error = %v, wantErr %v", err1, tt.wantErr)
				return
			}
			if !tt.wantErr && token1 == "" {
				t.Errorf("GenerateTokenWithTime() returned empty token")
				return
			}

			if !tt.wantErr {
				token2, err2 := GenerateTokenWithTime(tt.id, tt.time.Add(time.Second))
				if err2 != nil {
					t.Errorf("GenerateTokenWithTime() second call error = %v", err2)
					return
				}
				if token1 == token2 {
					t.Errorf("GenerateTokenWithTime() generated identical tokens for different times")
				}
			}

			if !tt.wantErr {
				claims := &claims{}
				token, err := jwt.ParseWithClaims(token1, claims, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(os.Getenv("JWT_SECRET")), nil
				})

				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
					return
				}

				if !token.Valid {
					t.Errorf("Token is invalid")
					return
				}

				if claims.UserID != tt.id {
					t.Errorf("Token claims UserID = %v, want %v", claims.UserID, tt.id)
				}

				if claims.IssuedAt != tt.time.Unix() {
					t.Errorf("Token claims IssuedAt = %v, want %v", claims.IssuedAt, tt.time.Unix())
				}
			}
		})
	}
}


/*
ROOST_METHOD_HASH=generateToken_2cc40e0108
ROOST_METHOD_SIG_HASH=generateToken_9de4114fe8

FUNCTION_DEF=func generateToken(id uint, now time.Time) (string, error) 

*/
func TestGenerateToken(t *testing.T) {
	originalJWTSecret := jwtSecret
	defer func() { jwtSecret = originalJWTSecret }()

	tests := []struct {
		name        string
		userID      uint
		currentTime time.Time
		wantErr     bool
		checkFunc   func(*testing.T, string)
	}{
		{
			name:        "Valid User ID",
			userID:      1,
			currentTime: time.Now(),
			wantErr:     false,
			checkFunc: func(t *testing.T, token string) {
				if token == "" {
					t.Error("Expected non-empty token")
				}
			},
		},
		{
			name:        "Verify Token Expiration",
			userID:      2,
			currentTime: time.Now(),
			wantErr:     false,
			checkFunc: func(t *testing.T, token string) {
				claims := &claims{}
				_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				expectedExp := time.Now().Add(72 * time.Hour).Unix()
				if claims.ExpiresAt < expectedExp-5 || claims.ExpiresAt > expectedExp+5 {
					t.Errorf("Expiration time not within expected range")
				}
			},
		},
		{
			name:        "Zero User ID",
			userID:      0,
			currentTime: time.Now(),
			wantErr:     false,
			checkFunc: func(t *testing.T, token string) {
				if token == "" {
					t.Error("Expected non-empty token even with zero user ID")
				}
			},
		},
		{
			name:        "Maximum Uint Value",
			userID:      math.MaxUint,
			currentTime: time.Now(),
			wantErr:     false,
			checkFunc: func(t *testing.T, token string) {
				if token == "" {
					t.Error("Expected non-empty token with maximum uint value")
				}
			},
		},
		{
			name:        "Empty JWT Secret",
			userID:      1,
			currentTime: time.Now(),
			wantErr:     true,
			checkFunc: func(t *testing.T, token string) {
				jwtSecret = []byte{}
			},
		},
		{
			name:        "Far Future Time",
			userID:      1,
			currentTime: time.Unix(1<<63-1, 0),
			wantErr:     false,
			checkFunc: func(t *testing.T, token string) {
				claims := &claims{}
				_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				if claims.ExpiresAt <= time.Now().Unix() {
					t.Error("Expiration time should be in the future")
				}
			},
		},
		{
			name:        "Verify Token Structure and Claims",
			userID:      42,
			currentTime: time.Now(),
			wantErr:     false,
			checkFunc: func(t *testing.T, token string) {
				parts := strings.Split(token, ".")
				if len(parts) != 3 {
					t.Error("Token should have three parts")
				}

				claims := &claims{}
				parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				if claims.UserID != 42 {
					t.Errorf("Expected UserID 42, got %d", claims.UserID)
				}
				if parsedToken.Method.Alg() != "HS256" {
					t.Errorf("Expected HS256 signing method, got %s", parsedToken.Method.Alg())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.checkFunc != nil {
				tt.checkFunc(t, "")
			}

			got, err := generateToken(tt.userID, tt.currentTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, got)
			}
		})
	}
}

