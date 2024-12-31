package auth

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/stretchr/testify/assert"
)

const mockAuthorizationHeader = "authorization"

/*
ROOST_METHOD_HASH=GenerateTokenWithTime_d0df64aa69
ROOST_METHOD_SIG_HASH=GenerateTokenWithTime_72dd09cde6
*/
func TestGenerateTokenWithTime(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	tests := []struct {
		name    string
		id      uint
		time    time.Time
		wantErr bool
	}{
		{
			name:    "Successful Token Generation",
			id:      1,
			time:    time.Now(),
			wantErr: false,
		},
		{
			name:    "Token Generation with Zero User ID",
			id:      0,
			time:    time.Now(),
			wantErr: false,
		},
		{
			name:    "Token Generation with Future Time",
			id:      2,
			time:    time.Now().Add(24 * time.Hour),
			wantErr: false,
		},
		{
			name:    "Token Generation with Past Time",
			id:      3,
			time:    time.Now().Add(-24 * time.Hour),
			wantErr: false,
		},
		{
			name:    "Token Generation with Maximum Uint Value",
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
			if !tt.wantErr {
				if got == "" {
					t.Errorf("GenerateTokenWithTime() returned empty token")
				}

				token, err := jwt.Parse(got, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})

				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}

				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					t.Errorf("Failed to get claims from token")
				}

				if uint(claims["user_id"].(float64)) != tt.id {
					t.Errorf("Token has incorrect user_id. got %v, want %v", uint(claims["user_id"].(float64)), tt.id)
				}

				expTime := time.Unix(int64(claims["exp"].(float64)), 0)
				expectedExpTime := tt.time.Add(time.Hour * 72)
				if expTime.Sub(expectedExpTime) > time.Second {
					t.Errorf("Token has incorrect expiration time. got %v, want %v", expTime, expectedExpTime)
				}
			}
		})
	}
}

func TestGenerateTokenWithTimeConcurrent(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	const numGoroutines = 100
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id uint) {
			_, err := GenerateTokenWithTime(id, time.Now())
			results <- err
		}(uint(i))
	}

	for i := 0; i < numGoroutines; i++ {
		if err := <-results; err != nil {
			t.Errorf("Concurrent GenerateTokenWithTime() returned error: %v", err)
		}
	}
}

func TestGenerateTokenWithTimeInvalidSecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	_, err := GenerateTokenWithTime(1, time.Now())
	if err == nil {
		t.Errorf("GenerateTokenWithTime() did not return error with invalid JWT secret")
	}
}

/*
ROOST_METHOD_HASH=generateToken_2cc40e0108
ROOST_METHOD_SIG_HASH=generateToken_9de4114fe8
*/
func TestgenerateToken(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	os.Setenv("JWT_SECRET", "test_secret")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	now := time.Now()

	tests := []struct {
		name        string
		userID      uint
		currentTime time.Time
		wantErr     bool
		setup       func()
		validate    func(*testing.T, string)
	}{
		{
			name:        "Successful Token Generation",
			userID:      1,
			currentTime: now,
			wantErr:     false,
			validate: func(t *testing.T, token string) {
				if token == "" {
					t.Error("Expected non-empty token, got empty string")
				}
			},
		},
		{
			name:        "Token Expiration Time",
			userID:      2,
			currentTime: now,
			wantErr:     false,
			validate: func(t *testing.T, token string) {
				parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
					return
				}
				if claims, ok := parsedToken.Claims.(*claims); ok {
					expectedExp := now.Add(time.Hour * 72).Unix()
					if claims.ExpiresAt != expectedExp {
						t.Errorf("Expected expiration time %v, got %v", expectedExp, claims.ExpiresAt)
					}
				} else {
					t.Error("Failed to get claims from token")
				}
			},
		},
		{
			name:        "User ID in Token Claims",
			userID:      3,
			currentTime: now,
			wantErr:     false,
			validate: func(t *testing.T, token string) {
				parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
					return
				}
				if claims, ok := parsedToken.Claims.(*claims); ok {
					if claims.UserID != 3 {
						t.Errorf("Expected user ID 3, got %v", claims.UserID)
					}
				} else {
					t.Error("Failed to get claims from token")
				}
			},
		},
		{
			name:        "Invalid JWT Secret",
			userID:      4,
			currentTime: now,
			wantErr:     true,
			setup: func() {
				jwtSecret = []byte{}
			},
		},
		{
			name:        "Zero User ID",
			userID:      0,
			currentTime: now,
			wantErr:     false,
			validate: func(t *testing.T, token string) {
				parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
					return
				}
				if claims, ok := parsedToken.Claims.(*claims); ok {
					if claims.UserID != 0 {
						t.Errorf("Expected user ID 0, got %v", claims.UserID)
					}
				} else {
					t.Error("Failed to get claims from token")
				}
			},
		},
		{
			name:        "Maximum Uint User ID",
			userID:      ^uint(0),
			currentTime: now,
			wantErr:     false,
			validate: func(t *testing.T, token string) {
				parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
					return
				}
				if claims, ok := parsedToken.Claims.(*claims); ok {
					if claims.UserID != ^uint(0) {
						t.Errorf("Expected user ID %v, got %v", ^uint(0), claims.UserID)
					}
				} else {
					t.Error("Failed to get claims from token")
				}
			},
		},
		{
			name:        "Token Signing Method Verification",
			userID:      5,
			currentTime: now,
			wantErr:     false,
			validate: func(t *testing.T, token string) {
				parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, errors.New("unexpected signing method")
					}
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
					return
				}
				if parsedToken.Method != jwt.SigningMethodHS256 {
					t.Errorf("Expected signing method HS256, got %v", parsedToken.Method)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			token, err := generateToken(tt.userID, tt.currentTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, token)
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
				token, err := generateToken(uint(index), now)
				tokens[index] = token
				errors[index] = err
			}(i)
		}

		wg.Wait()

		for i := 0; i < numGoroutines; i++ {
			if errors[i] != nil {
				t.Errorf("Error generating token for goroutine %d: %v", i, errors[i])
			}
			if tokens[i] == "" {
				t.Errorf("Empty token generated for goroutine %d", i)
			}
		}
	})
}

/*
ROOST_METHOD_HASH=GetUserID_f2dd680cb2
ROOST_METHOD_SIG_HASH=GetUserID_e739312e3d
*/
func TestGetUserID(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	os.Setenv("JWT_SECRET", "test_secret")
	jwtSecret = []byte("test_secret")

	tests := []struct {
		name           string
		setupContext   func() context.Context
		expectedUserID uint
		expectedError  string
	}{
		{
			name: "Valid Token with Correct User ID",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				})
				tokenString, _ := token.SignedString(jwtSecret)
				return context.WithValue(context.Background(), mockAuthorizationHeader, "Token "+tokenString)
			},
			expectedUserID: 123,
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()
			userID, err := GetUserID(ctx)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedUserID, userID)
		})
	}
}

func init() {
	grpc_auth.AuthFromMD = mockAuthFromMD
}

func mockAuthFromMD(ctx context.Context, prefix string) (string, error) {
	auth, ok := ctx.Value(mockAuthorizationHeader).(string)
	if !ok {
		return "", errors.New("token not found")
	}
	return auth[len(prefix)+1:], nil
}
