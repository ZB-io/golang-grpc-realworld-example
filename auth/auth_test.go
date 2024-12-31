package auth

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
ROOST_METHOD_HASH=GenerateTokenWithTime_d0df64aa69
ROOST_METHOD_SIG_HASH=GenerateTokenWithTime_72dd09cde6
*/
func TestGenerateTokenWithTime(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	tests := []struct {
		name    string
		id      uint
		t       time.Time
		wantErr bool
	}{
		{
			name:    "Successful Token Generation",
			id:      1,
			t:       time.Now(),
			wantErr: false,
		},
		{
			name:    "Token Generation with Zero User ID",
			id:      0,
			t:       time.Now(),
			wantErr: false,
		},
		{
			name:    "Token Generation with Future Time",
			id:      1,
			t:       time.Now().Add(24 * time.Hour),
			wantErr: false,
		},
		{
			name:    "Token Generation with Past Time",
			id:      1,
			t:       time.Now().Add(-24 * time.Hour),
			wantErr: false,
		},
		{
			name:    "Token Generation with Maximum Uint Value",
			id:      ^uint(0),
			t:       time.Now(),
			wantErr: false,
		},
		{
			name:    "Token Generation with Zero Time",
			id:      1,
			t:       time.Time{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateTokenWithTime(tt.id, tt.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateTokenWithTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("GenerateTokenWithTime() returned empty token")
			}

			if !tt.wantErr {
				token, err := jwt.Parse(got, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
					return
				}

				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					t.Errorf("Failed to get claims from token")
					return
				}

				if uint(claims["user_id"].(float64)) != tt.id {
					t.Errorf("Token user_id claim = %v, want %v", claims["user_id"], tt.id)
				}

				if claims["exp"] == nil || claims["exp"].(float64) == 0 {
					t.Errorf("Token exp claim is not set")
				}
			}
		})
	}

	t.Run("Consistency of Generated Tokens", func(t *testing.T) {
		id := uint(1)
		tokenTime := time.Now()
		token1, err1 := GenerateTokenWithTime(id, tokenTime)
		token2, err2 := GenerateTokenWithTime(id, tokenTime)

		if err1 != nil || err2 != nil {
			t.Errorf("Unexpected errors: %v, %v", err1, err2)
			return
		}

		if token1 != token2 {
			t.Errorf("Tokens are not consistent: %v != %v", token1, token2)
		}
	})
}

/*
ROOST_METHOD_HASH=generateToken_2cc40e0108
ROOST_METHOD_SIG_HASH=generateToken_9de4114fe8
*/
func TestGenerateToken(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	tests := []struct {
		name        string
		userID      uint
		currentTime time.Time
		setupEnv    func()
		wantErr     bool
		validate    func(*testing.T, string, error)
	}{
		{
			name:        "Successful Token Generation",
			userID:      1234,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test_secret")
			},
			wantErr: false,
			validate: func(t *testing.T, token string, err error) {
				if token == "" {
					t.Error("Expected non-empty token, got empty string")
				}
			},
		},
		{
			name:        "Token Expiration Time",
			userID:      5678,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test_secret")
			},
			wantErr: false,
			validate: func(t *testing.T, token string, err error) {
				claims := &claims{}
				_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte("test_secret"), nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				expectedExpiration := time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC).Unix()
				if claims.ExpiresAt != expectedExpiration {
					t.Errorf("Expected expiration %v, got %v", expectedExpiration, claims.ExpiresAt)
				}
			},
		},
		{
			name:        "User ID in Token Claims",
			userID:      9876,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test_secret")
			},
			wantErr: false,
			validate: func(t *testing.T, token string, err error) {
				claims := &claims{}
				_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte("test_secret"), nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				if claims.UserID != 9876 {
					t.Errorf("Expected user ID 9876, got %v", claims.UserID)
				}
			},
		},
		{
			name:        "Error Handling with Empty JWT Secret",
			userID:      1111,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "")
			},
			wantErr: true,
			validate: func(t *testing.T, token string, err error) {
				if token != "" {
					t.Error("Expected empty token, got non-empty string")
				}
				if err == nil {
					t.Error("Expected an error, got nil")
				}
			},
		},
		{
			name:        "Token Signing Method",
			userID:      2222,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test_secret")
			},
			wantErr: false,
			validate: func(t *testing.T, token string, err error) {
				parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return []byte("test_secret"), nil
				})
				if parsedToken.Header["alg"] != "HS256" {
					t.Errorf("Expected signing method HS256, got %v", parsedToken.Header["alg"])
				}
			},
		},
		{
			name:        "Large User ID Handling",
			userID:      ^uint(0),
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test_secret")
			},
			wantErr: false,
			validate: func(t *testing.T, token string, err error) {
				claims := &claims{}
				_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte("test_secret"), nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				if claims.UserID != ^uint(0) {
					t.Errorf("Expected user ID %v, got %v", ^uint(0), claims.UserID)
				}
			},
		},
		{
			name:        "Token Generation with Zero User ID",
			userID:      0,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test_secret")
			},
			wantErr: false,
			validate: func(t *testing.T, token string, err error) {
				claims := &claims{}
				_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte("test_secret"), nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				if claims.UserID != 0 {
					t.Errorf("Expected user ID 0, got %v", claims.UserID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			token, err := generateToken(tt.userID, tt.currentTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.validate(t, token, err)
		})
	}
}

/*
ROOST_METHOD_HASH=GetUserID_f2dd680cb2
ROOST_METHOD_SIG_HASH=GetUserID_e739312e3d
*/
func TestGetUserID(t *testing.T) {
	jwtSecret = []byte("test_secret")

	tests := []struct {
		name        string
		setupCtx    func() context.Context
		expectedID  uint
		expectedErr string
	}{
		{
			name: "Valid Token with Correct User ID",
			setupCtx: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
					UserID: 123,
				})
				tokenString, _ := token.SignedString(jwtSecret)
				ctx := context.Background()
				return mockAppendToOutgoingContext(ctx, "Token", tokenString)
			},
			expectedID:  123,
			expectedErr: "",
		},
		{
			name: "Expired Token",
			setupCtx: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(-time.Hour).Unix(),
					},
					UserID: 456,
				})
				tokenString, _ := token.SignedString(jwtSecret)
				ctx := context.Background()
				return mockAppendToOutgoingContext(ctx, "Token", tokenString)
			},
			expectedID:  0,
			expectedErr: "token expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()

			originalAuthFromMD := grpc_auth.AuthFromMD
			grpc_auth.AuthFromMD = mockAuthFromMD
			defer func() { grpc_auth.AuthFromMD = originalAuthFromMD }()

			gotID, err := GetUserID(ctx)

			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expectedID, gotID)
		})
	}
}

func mockAppendToOutgoingContext(ctx context.Context, key, val string) context.Context {
	return context.WithValue(ctx, key, val)
}

func mockAuthFromMD(ctx context.Context, key string) (string, error) {
	val, ok := ctx.Value(key).(string)
	if !ok {
		return "", errors.New("token not found")
	}
	return val, nil
}
