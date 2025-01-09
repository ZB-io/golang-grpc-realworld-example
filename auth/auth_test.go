package auth

import (
	"math"
	"os"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"encoding/base64"
	"strings"
	"github.com/stretchr/testify/assert"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)








/*
ROOST_METHOD_HASH=GenerateTokenWithTime_d0df64aa69
ROOST_METHOD_SIG_HASH=GenerateTokenWithTime_72dd09cde6

FUNCTION_DEF=func GenerateTokenWithTime(id uint, t time.Time) (string, error) 

 */
func TestGenerateTokenWithTime(t *testing.T) {
	originalJWTSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalJWTSecret)

	tests := []struct {
		name        string
		userID      uint
		inputTime   time.Time
		setupEnv    func()
		expectToken bool
		expectErr   bool
	}{
		{
			name:        "Successful Token Generation",
			userID:      12345,
			inputTime:   time.Now(),
			setupEnv:    func() { os.Setenv("JWT_SECRET", "test_secret") },
			expectToken: true,
			expectErr:   false,
		},
		{
			name:        "Error Handling for Empty JWT Secret",
			userID:      12345,
			inputTime:   time.Now(),
			setupEnv:    func() { os.Setenv("JWT_SECRET", "") },
			expectToken: false,
			expectErr:   true,
		},
		{
			name:        "Token Generation with Minimum Valid User ID",
			userID:      1,
			inputTime:   time.Now(),
			setupEnv:    func() { os.Setenv("JWT_SECRET", "test_secret") },
			expectToken: true,
			expectErr:   false,
		},
		{
			name:        "Token Generation with Maximum User ID",
			userID:      math.MaxUint32,
			inputTime:   time.Now(),
			setupEnv:    func() { os.Setenv("JWT_SECRET", "test_secret") },
			expectToken: true,
			expectErr:   false,
		},
		{
			name:        "Token Generation with Zero User ID",
			userID:      0,
			inputTime:   time.Now(),
			setupEnv:    func() { os.Setenv("JWT_SECRET", "test_secret") },
			expectToken: false,
			expectErr:   true,
		},
		{
			name:        "Token Generation with Past Time",
			userID:      12345,
			inputTime:   time.Now().Add(-24 * time.Hour),
			setupEnv:    func() { os.Setenv("JWT_SECRET", "test_secret") },
			expectToken: true,
			expectErr:   false,
		},
		{
			name:        "Token Generation with Future Time",
			userID:      12345,
			inputTime:   time.Now().Add(24 * time.Hour),
			setupEnv:    func() { os.Setenv("JWT_SECRET", "test_secret") },
			expectToken: true,
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()

			token, err := GenerateTokenWithTime(tt.userID, tt.inputTime)

			if tt.expectErr && err == nil {
				t.Errorf("Expected an error, but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.expectToken && token == "" {
				t.Errorf("Expected a non-empty token, but got an empty string")
			}
			if !tt.expectToken && token != "" {
				t.Errorf("Expected an empty token, but got: %s", token)
			}

			if tt.expectToken {

				claims := &claims{}
				_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				if claims.UserID != tt.userID {
					t.Errorf("Expected UserID %d, but got %d", tt.userID, claims.UserID)
				}
				expectedExp := tt.inputTime.Add(time.Hour * 24 * 7).Unix()
				if claims.ExpiresAt != expectedExp {
					t.Errorf("Expected ExpiresAt %d, but got %d", expectedExp, claims.ExpiresAt)
				}
			}
		})
	}
}

func TestGenerateTokenWithTimeConsistency(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Setenv("JWT_SECRET", "")

	userID := uint(12345)
	inputTime := time.Now()

	token1, err1 := GenerateTokenWithTime(userID, inputTime)
	if err1 != nil {
		t.Fatalf("Unexpected error on first call: %v", err1)
	}

	token2, err2 := GenerateTokenWithTime(userID, inputTime)
	if err2 != nil {
		t.Fatalf("Unexpected error on second call: %v", err2)
	}

	if token1 != token2 {
		t.Errorf("Expected consistent tokens, but got different tokens:\nToken1: %s\nToken2: %s", token1, token2)
	}
}


/*
ROOST_METHOD_HASH=GenerateToken_b7f5ef3740
ROOST_METHOD_SIG_HASH=GenerateToken_d10a3e47a3

FUNCTION_DEF=func GenerateToken(id uint) (string, error) 

 */
func TestGenerateToken(t *testing.T) {

	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	tests := []struct {
		name        string
		userID      uint
		setupEnv    func()
		wantErr     bool
		validateJWT func(*testing.T, string)
	}{
		{
			name:   "Successful Token Generation",
			userID: 1,
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test_secret")
			},
			wantErr: false,
			validateJWT: func(t *testing.T, token string) {
				claims := &claims{}
				_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte("test_secret"), nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				if claims.UserID != 1 {
					t.Errorf("Expected UserID 1, got %d", claims.UserID)
				}
				if claims.ExpiresAt <= time.Now().Unix() {
					t.Error("Token has already expired")
				}
			},
		},
		{
			name:   "Error Handling with Invalid User ID",
			userID: 0,
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test_secret")
			},
			wantErr: true,
		},
		{
			name:   "Token Expiration",
			userID: 1,
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test_secret")
			},
			wantErr: false,
			validateJWT: func(t *testing.T, token string) {
				time.Sleep(2 * time.Second)
				claims := &claims{}
				_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte("test_secret"), nil
				})
				if err == nil {
					t.Error("Expected token to be expired, but it's still valid")
				}
			},
		},
		{
			name:   "Consistency of Generated Tokens",
			userID: 1,
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test_secret")
			},
			wantErr: false,
			validateJWT: func(t *testing.T, token1 string) {
				token2, err := GenerateToken(1)
				if err != nil {
					t.Errorf("Failed to generate second token: %v", err)
				}
				if token1 == token2 {
					t.Error("Expected different tokens, but they are the same")
				}
			},
		},
		{
			name:   "Handling of JWT_SECRET Environment Variable",
			userID: 1,
			setupEnv: func() {
				os.Unsetenv("JWT_SECRET")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			token, err := GenerateToken(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.validateJWT != nil {
				tt.validateJWT(t, token)
			}
		})
	}
}

func TestGenerateTokenPerformance(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	numTokens := 1000
	start := time.Now()
	for i := 1; i <= numTokens; i++ {
		_, err := GenerateToken(uint(i))
		if err != nil {
			t.Errorf("Failed to generate token %d: %v", i, err)
		}
	}
	duration := time.Since(start)
	averageTime := duration / time.Duration(numTokens)
	t.Logf("Average time per token generation: %v", averageTime)

	if averageTime > 1*time.Millisecond {
		t.Errorf("Token generation too slow. Average time: %v", averageTime)
	}
}


/*
ROOST_METHOD_HASH=generateToken_2cc40e0108
ROOST_METHOD_SIG_HASH=generateToken_9de4114fe8

FUNCTION_DEF=func generateToken(id uint, now time.Time) (string, error) 

 */
func TestGenerateToken(t *testing.T) {

	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)
	os.Setenv("JWT_SECRET", "test_secret")

	tests := []struct {
		name        string
		userID      uint
		currentTime time.Time
		wantErr     bool
		validate    func(*testing.T, string, error)
	}{
		{
			name:        "Valid Token Generation",
			userID:      1,
			currentTime: time.Now(),
			wantErr:     false,
			validate: func(t *testing.T, token string, err error) {
				assert.NotEmpty(t, token)
				assert.NoError(t, err)
			},
		},
		{
			name:        "Verify Token Expiration Time",
			userID:      2,
			currentTime: time.Now(),
			wantErr:     false,
			validate: func(t *testing.T, token string, err error) {
				assert.NoError(t, err)
				claims := &claims{}
				_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				assert.NoError(t, err)
				expectedExpTime := time.Now().Add(time.Hour * 72).Unix()
				assert.InDelta(t, expectedExpTime, claims.ExpiresAt, 1)
			},
		},
		{
			name:        "Verify Token Claims",
			userID:      3,
			currentTime: time.Now(),
			wantErr:     false,
			validate: func(t *testing.T, token string, err error) {
				assert.NoError(t, err)
				claims := &claims{}
				_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				assert.NoError(t, err)
				assert.Equal(t, uint(3), claims.UserID)
			},
		},
		{
			name:        "Handle Zero User ID",
			userID:      0,
			currentTime: time.Now(),
			wantErr:     false,
			validate: func(t *testing.T, token string, err error) {
				assert.NotEmpty(t, token)
				assert.NoError(t, err)
			},
		},
		{
			name:        "Test with Different Time Zones",
			userID:      4,
			currentTime: time.Now().In(time.FixedZone("GMT+8", 8*60*60)),
			wantErr:     false,
			validate: func(t *testing.T, token string, err error) {
				assert.NoError(t, err)
				claims := &claims{}
				_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				assert.NoError(t, err)
				expectedExpTime := time.Now().Add(time.Hour * 72).Unix()
				assert.InDelta(t, expectedExpTime, claims.ExpiresAt, 1)
			},
		},
		{
			name:        "Verify Token Signing Method",
			userID:      5,
			currentTime: time.Now(),
			wantErr:     false,
			validate: func(t *testing.T, token string, err error) {
				assert.NoError(t, err)
				parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				assert.NoError(t, err)
				assert.Equal(t, jwt.SigningMethodHS256, parsedToken.Method)
			},
		},
		{
			name:        "Test with Maximum uint Value",
			userID:      ^uint(0),
			currentTime: time.Now(),
			wantErr:     false,
			validate: func(t *testing.T, token string, err error) {
				assert.NotEmpty(t, token)
				assert.NoError(t, err)
			},
		},
		{
			name:        "Verify Token Format",
			userID:      6,
			currentTime: time.Now(),
			wantErr:     false,
			validate: func(t *testing.T, token string, err error) {
				assert.NoError(t, err)
				parts := strings.Split(token, ".")
				assert.Equal(t, 3, len(parts))
				_, err = base64.RawURLEncoding.DecodeString(parts[0])
				assert.NoError(t, err)
				_, err = base64.RawURLEncoding.DecodeString(parts[1])
				assert.NoError(t, err)
			},
		},
		{
			name:        "Test with Empty JWT Secret",
			userID:      7,
			currentTime: time.Now(),
			wantErr:     true,
			validate: func(t *testing.T, token string, err error) {
				os.Setenv("JWT_SECRET", "")
				assert.Error(t, err)
				assert.Empty(t, token)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := generateToken(tt.userID, tt.currentTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.validate(t, gotToken, err)
		})
	}
}

func TestGenerateTokenPerformance(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	userID := uint(1)
	now := time.Now()
	iterations := 1000

	startTime := time.Now()
	for i := 0; i < iterations; i++ {
		_, err := generateToken(userID, now)
		assert.NoError(t, err)
	}
	duration := time.Since(startTime)

	averageTime := duration / time.Duration(iterations)
	t.Logf("Average time to generate token: %v", averageTime)
	assert.Less(t, averageTime, 1*time.Millisecond, "Token generation is too slow")
}


/*
ROOST_METHOD_HASH=GetUserID_f2dd680cb2
ROOST_METHOD_SIG_HASH=GetUserID_e739312e3d

FUNCTION_DEF=func GetUserID(ctx context.Context) (uint, error) 

 */
func TestGetUserID(t *testing.T) {

	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	testSecret := "test_secret"
	os.Setenv("JWT_SECRET", testSecret)
	jwtSecret = []byte(testSecret)

	tests := []struct {
		name           string
		setupContext   func() context.Context
		expectedUserID uint
		expectedError  string
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()
			userID, err := GetUserID(ctx)

			if userID != tt.expectedUserID {
				t.Errorf("Expected user ID %d, got %d", tt.expectedUserID, userID)
			}

			if tt.expectedError == "" && err != nil {
				t.Errorf("Expected no error, got %v", err)
			} else if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error %q, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {

					if statusErr, ok := status.FromError(err); ok {
						if statusErr.Code() != codes.Unauthenticated || statusErr.Message() != tt.expectedError {
							t.Errorf("Expected error %q, got %q", tt.expectedError, statusErr.Message())
						}
					} else {
						t.Errorf("Expected error %q, got %q", tt.expectedError, err.Error())
					}
				}
			}
		})
	}
}

