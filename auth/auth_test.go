package auth

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
)

/*
ROOST_METHOD_HASH=GenerateTokenWithTime_d0df64aa69
ROOST_METHOD_SIG_HASH=GenerateTokenWithTime_72dd09cde6
*/
func TestGenerateTokenWithTime(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	now := time.Now()
	future := now.Add(24 * time.Hour)
	past := now.Add(-24 * time.Hour)
	maxUint := ^uint(0)

	tests := []struct {
		name    string
		id      uint
		t       time.Time
		wantErr bool
	}{
		{
			name:    "Successful Token Generation",
			id:      1,
			t:       now,
			wantErr: false,
		},
		{
			name:    "Token Generation with Zero User ID",
			id:      0,
			t:       now,
			wantErr: true,
		},
		{
			name:    "Token Generation with Future Time",
			id:      1,
			t:       future,
			wantErr: false,
		},
		{
			name:    "Token Generation with Past Time",
			id:      1,
			t:       past,
			wantErr: true,
		},
		{
			name:    "Token Generation with Maximum Uint Value",
			id:      maxUint,
			t:       now,
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
			if !tt.wantErr {
				if got == "" {
					t.Errorf("GenerateTokenWithTime() returned empty token, expected non-empty token")
				}

				token, err := jwt.ParseWithClaims(got, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})

				if err != nil {
					t.Errorf("Failed to parse generated token: %v", err)
					return
				}

				if claims, ok := token.Claims.(*claims); ok && token.Valid {
					if claims.UserID != tt.id {
						t.Errorf("Token UserID = %v, want %v", claims.UserID, tt.id)
					}

					expectedExp := tt.t.Add(72 * time.Hour).Unix()
					if claims.ExpiresAt != expectedExp {
						t.Errorf("Token ExpiresAt = %v, want %v", claims.ExpiresAt, expectedExp)
					}
				} else {
					t.Errorf("Token claims are invalid")
				}
			}
		})
	}
}

func TestGenerateTokenWithTimeMissingSecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", "test_secret")

	_, err := GenerateTokenWithTime(1, time.Now())
	if err == nil {
		t.Errorf("GenerateTokenWithTime() expected error when JWT_SECRET is not set, got nil")
	}
}

func TestGenerateTokenWithTimePerformance(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	now := time.Now()
	numTokens := 1000
	start := time.Now()

	for i := 1; i <= numTokens; i++ {
		_, err := GenerateTokenWithTime(uint(i), now)
		if err != nil {
			t.Errorf("GenerateTokenWithTime() failed on iteration %d: %v", i, err)
			return
		}
	}

	duration := time.Since(start)
	t.Logf("Generated %d tokens in %v", numTokens, duration)

	if duration > 5*time.Second {
		t.Errorf("Token generation took too long: %v for %d tokens", duration, numTokens)
	}
}

func TestVerifyGeneratedToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	userID := uint(123)
	now := time.Now()

	token, err := GenerateTokenWithTime(userID, now)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if claims, ok := parsedToken.Claims.(*claims); ok && parsedToken.Valid {
		if claims.UserID != userID {
			t.Errorf("Token UserID = %v, want %v", claims.UserID, userID)
		}
		expectedExp := now.Add(72 * time.Hour).Unix()
		if claims.ExpiresAt != expectedExp {
			t.Errorf("Token ExpiresAt = %v, want %v", claims.ExpiresAt, expectedExp)
		}
	} else {
		t.Errorf("Token claims are invalid")
	}
}

/*
ROOST_METHOD_HASH=generateToken_2cc40e0108
ROOST_METHOD_SIG_HASH=generateToken_9de4114fe8
*/
func TestGenerateToken(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	os.Setenv("JWT_SECRET", "test_secret")
	jwtSecret = []byte("test_secret")

	tests := []struct {
		name        string
		userID      uint
		currentTime time.Time
		wantErr     bool
		setup       func()
		validate    func(*testing.T, string, error)
	}{
		{
			name:        "Successful Token Generation",
			userID:      1,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:     false,
			validate: func(t *testing.T, token string, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if token == "" {
					t.Error("Expected non-empty token")
				}
			},
		},
		{
			name:        "Token Expiration Time",
			userID:      2,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:     false,
			validate: func(t *testing.T, token string, err error) {
				parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				if claims, ok := parsedToken.Claims.(*claims); ok {
					expectedExpiration := time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC).Unix()
					if claims.ExpiresAt != expectedExpiration {
						t.Errorf("Expected expiration %v, got %v", expectedExpiration, claims.ExpiresAt)
					}
				} else {
					t.Error("Failed to get claims from token")
				}
			},
		},
		{
			name:        "Token Content Verification",
			userID:      3,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:     false,
			validate: func(t *testing.T, token string, err error) {
				parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
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
			name:        "Error Handling with Empty JWT Secret",
			userID:      4,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:     true,
			setup: func() {
				jwtSecret = []byte{}
			},
			validate: func(t *testing.T, token string, err error) {
				if err == nil {
					t.Error("Expected an error, got nil")
				}
				if token != "" {
					t.Errorf("Expected empty token, got %v", token)
				}
			},
		},
		{
			name:        "Token Generation with Zero User ID",
			userID:      0,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:     false,
			validate: func(t *testing.T, token string, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if token == "" {
					t.Error("Expected non-empty token")
				}
				parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
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
			name:        "Token Generation at Time Boundaries",
			userID:      5,
			currentTime: time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.UTC),
			wantErr:     false,
			validate: func(t *testing.T, token string, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				if claims, ok := parsedToken.Claims.(*claims); ok {
					expectedExpiration := time.Date(2024, 1, 3, 23, 59, 59, 999999999, time.UTC).Unix()
					if claims.ExpiresAt != expectedExpiration {
						t.Errorf("Expected expiration %v, got %v", expectedExpiration, claims.ExpiresAt)
					}
				} else {
					t.Error("Failed to get claims from token")
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
			if tt.validate != nil {
				tt.validate(t, token, err)
			}
		})
	}

	t.Run("Concurrent Token Generation", func(t *testing.T) {
		const numGoroutines = 100
		results := make(chan bool, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id uint) {
				token, err := generateToken(id, time.Now())
				results <- (err == nil && token != "")
			}(uint(i))
		}

		for i := 0; i < numGoroutines; i++ {
			if !<-results {
				t.Errorf("Concurrent token generation failed for some goroutine")
				break
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
		name          string
		setupContext  func() context.Context
		expectedID    uint
		expectedError string
	}{
		{
			name: "Valid Token with Correct User ID",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
					UserID: 123,
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    123,
			expectedError: "",
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
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "token expired",
		},
		{
			name: "Malformed Token",
			setupContext: func() context.Context {
				md := metadata.Pairs("authorization", "Token malformed.token.here")
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: it's not even a token",
		},
		{
			name: "Missing Token in Context",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedID:    0,
			expectedError: "Request unauthenticated with Token",
		},
		{
			name: "Token with Invalid Signature",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
					UserID: 123,
				})
				tokenString, _ := token.SignedString([]byte("wrong_secret"))
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: couldn't handle this token",
		},
		{
			name: "Token with Invalid Claims Structure",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"user_id": "not_a_number",
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: cannot map token to claims",
		},
		{
			name: "Valid Token but JWT_SECRET Environment Variable Not Set",
			setupContext: func() context.Context {
				os.Unsetenv("JWT_SECRET")
				jwtSecret = []byte{}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
					UserID: 123,
				})
				tokenString, _ := token.SignedString([]byte("test_secret"))
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: couldn't handle this token",
		},
		{
			name: "Token with Future Not Before Claim",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
						NotBefore: time.Now().Add(time.Hour).Unix(),
					},
					UserID: 123,
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedID:    0,
			expectedError: "invalid token: couldn't handle this token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()
			gotID, err := GetUserID(ctx)

			if gotID != tt.expectedID {
				t.Errorf("GetUserID() gotID = %v, want %v", gotID, tt.expectedID)
			}

			if err == nil && tt.expectedError != "" {
				t.Errorf("GetUserID() error = nil, wantErr %v", tt.expectedError)
			} else if err != nil {
				if tt.expectedError == "" {
					t.Errorf("GetUserID() unexpected error: %v", err)
				} else if !errors.Is(err, errors.New(tt.expectedError)) {
					t.Errorf("GetUserID() error = %v, wantErr %v", err, tt.expectedError)
				}
			}
		})
	}
}
