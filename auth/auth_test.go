package auth

import (
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"os"
	"context"
	"errors"
	"strings"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/metadata"
)



var grpc_auth_AuthFromMD = AuthFromMD




/*
ROOST_METHOD_HASH=GenerateToken_b7f5ef3740
ROOST_METHOD_SIG_HASH=GenerateToken_d10a3e47a3

FUNCTION_DEF=func GenerateToken(id uint) (string, error) 

 */
func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		userID  uint
		wantErr bool
	}{
		{
			name:    "Successfully Generate Token for Valid User ID",
			userID:  12345,
			wantErr: false,
		},
		{
			name:    "Handle Zero User ID",
			userID:  0,
			wantErr: false,
		},
		{
			name:    "Handle Very Large User ID",
			userID:  ^uint(0),
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
			if !tt.wantErr && token == "" {
				t.Errorf("GenerateToken() returned empty token for userID %v", tt.userID)
			}

			if !tt.wantErr {
				parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
					return
				}
				if claims, ok := parsedToken.Claims.(*claims); ok {
					if claims.UserID != tt.userID {
						t.Errorf("Token claims UserID = %v, want %v", claims.UserID, tt.userID)
					}
					if claims.ExpiresAt == 0 {
						t.Error("Token does not have an expiration time")
					}
				} else {
					t.Error("Failed to extract claims from token")
				}
			}
		})
	}
}

func TestGenerateUniqueTokens(t *testing.T) {
	token1, err1 := GenerateToken(1)
	token2, err2 := GenerateToken(2)

	if err1 != nil || err2 != nil {
		t.Errorf("Error generating tokens: %v, %v", err1, err2)
		return
	}

	if token1 == token2 {
		t.Error("Tokens for different user IDs are not unique")
	}
}

func TestTokenConsistency(t *testing.T) {
	userID := uint(12345)
	token1, err1 := GenerateToken(userID)
	token2, err2 := GenerateToken(userID)

	if err1 != nil || err2 != nil {
		t.Errorf("Error generating tokens: %v, %v", err1, err2)
		return
	}

	if token1 == token2 {
		t.Error("Tokens generated in quick succession for the same user ID are identical")
	}
}

func TestTokenPerformance(t *testing.T) {
	const iterations = 1000
	userID := uint(12345)
	start := time.Now()

	for i := 0; i < iterations; i++ {
		_, err := GenerateToken(userID)
		if err != nil {
			t.Errorf("Error generating token on iteration %d: %v", i, err)
			return
		}
	}

	duration := time.Since(start)
	averageTime := duration / time.Duration(iterations)

	if averageTime > 1*time.Millisecond {
		t.Errorf("Token generation is slower than expected. Average time: %v", averageTime)
	}
}


/*
ROOST_METHOD_HASH=GenerateTokenWithTime_d0df64aa69
ROOST_METHOD_SIG_HASH=GenerateTokenWithTime_72dd09cde6

FUNCTION_DEF=func GenerateTokenWithTime(id uint, t time.Time) (string, error) 

 */
func TestGenerateTokenWithTime(t *testing.T) {

	oldJWTSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", oldJWTSecret)
	os.Setenv("JWT_SECRET", "test_secret")

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
			id:      1,
			time:    time.Now().Add(24 * time.Hour),
			wantErr: false,
		},
		{
			name:    "Token Generation with Past Time",
			id:      1,
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
			if !tt.wantErr && got == "" {
				t.Errorf("GenerateTokenWithTime() returned empty token")
			}

			token, err := jwt.ParseWithClaims(got, &claims{}, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})

			if err != nil {
				t.Errorf("Failed to parse token: %v", err)
				return
			}

			if claims, ok := token.Claims.(*claims); ok && token.Valid {
				if claims.UserID != tt.id {
					t.Errorf("Token UserID = %v, want %v", claims.UserID, tt.id)
				}
				expectedExpiry := tt.time.Add(time.Hour * 72).Unix()
				if claims.ExpiresAt != expectedExpiry {
					t.Errorf("Token ExpiresAt = %v, want %v", claims.ExpiresAt, expectedExpiry)
				}
			} else {
				t.Errorf("Token claims are not valid")
			}
		})
	}
}

func TestGenerateTokenWithTimeInvalidSecret(t *testing.T) {

	oldJWTSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", oldJWTSecret)
	os.Setenv("JWT_SECRET", "")

	_, err := GenerateTokenWithTime(1, time.Now())
	if err == nil {
		t.Errorf("GenerateTokenWithTime() error = nil, wantErr true")
	}
}


/*
ROOST_METHOD_HASH=GetUserID_f2dd680cb2
ROOST_METHOD_SIG_HASH=GetUserID_e739312e3d

FUNCTION_DEF=func GetUserID(ctx context.Context) (uint, error) 

 */
func AuthFromMD(ctx context.Context, expectedScheme string) (string, error) {
	md := metautils.ExtractIncoming(ctx)
	val := md.Get("authorization")
	if val == "" {
		return "", errors.New("Request unauthenticated with " + expectedScheme)
	}
	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", errors.New("Bad authorization string")
	}
	if !strings.EqualFold(splits[0], expectedScheme) {
		return "", errors.New("Request unauthenticated with " + expectedScheme)
	}
	return splits[1], nil
}

func MockToContext(ctx context.Context, token string) context.Context {
	md := metadata.Pairs("authorization", token)
	return metadata.NewIncomingContext(ctx, md)
}

func TestGetUserID(t *testing.T) {
	originalJwtSecret := jwtSecret
	defer func() { jwtSecret = originalJwtSecret }()

	tests := []struct {
		name          string
		setupContext  func() context.Context
		setupEnv      func()
		expectedID    uint
		expectedError string
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
				ctx := context.Background()
				return MockToContext(ctx, "Token "+tokenString)
			},
			expectedID: 123,
		},
		{
			name: "Expired Token",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(-time.Hour).Unix(),
					},
				})
				tokenString, _ := token.SignedString(jwtSecret)
				ctx := context.Background()
				return MockToContext(ctx, "Token "+tokenString)
			},
			expectedError: "token expired",
		},
		{
			name: "Malformed Token",
			setupContext: func() context.Context {
				ctx := context.Background()
				return MockToContext(ctx, "Token malformed.token.here")
			},
			expectedError: "invalid token: it's not even a token",
		},
		{
			name: "Missing Token in Context",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedError: "Request unauthenticated with Token",
		},
		{
			name: "Token with Invalid Signature",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				})
				tokenString, _ := token.SignedString([]byte("wrong_secret"))
				ctx := context.Background()
				return MockToContext(ctx, "Token "+tokenString)
			},
			expectedError: "invalid token: couldn't handle this token",
		},
		{
			name: "Token with Non-Numeric User ID",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"user_id":   "not_a_number",
					"ExpiresAt": time.Now().Add(time.Hour).Unix(),
				})
				tokenString, _ := token.SignedString(jwtSecret)
				ctx := context.Background()
				return MockToContext(ctx, "Token "+tokenString)
			},
			expectedError: "invalid token: cannot map token to claims",
		},
		{
			name: "Valid Token with Future NotBefore Claim",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
						NotBefore: time.Now().Add(time.Hour).Unix(),
					},
				})
				tokenString, _ := token.SignedString(jwtSecret)
				ctx := context.Background()
				return MockToContext(ctx, "Token "+tokenString)
			},
			expectedError: "invalid token: couldn't handle this token",
		},
		{
			name: "Environmental Variable Dependency (JWT_SECRET)",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				})
				tokenString, _ := token.SignedString(jwtSecret)
				ctx := context.Background()
				return MockToContext(ctx, "Token "+tokenString)
			},
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "")
				jwtSecret = []byte("")
			},
			expectedError: "invalid token: couldn't handle this token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnv != nil {
				tt.setupEnv()
			}

			ctx := tt.setupContext()
			gotID, err := GetUserID(ctx)

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("GetUserID() error = %v, expectedError %v", err, tt.expectedError)
				}
			} else if err != nil {
				t.Errorf("GetUserID() unexpected error: %v", err)
			}

			if gotID != tt.expectedID {
				t.Errorf("GetUserID() gotID = %v, want %v", gotID, tt.expectedID)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=generateToken_2cc40e0108
ROOST_METHOD_SIG_HASH=generateToken_9de4114fe8

FUNCTION_DEF=func generateToken(id uint, now time.Time) (string, error) 

 */
func BenchmarkGenerateToken(b *testing.B) {
	id := uint(123)
	now := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := generateToken(id, now)
		if err != nil {
			b.Fatalf("generateToken() error = %v", err)
		}
	}
}

func TestGenerateToken(t *testing.T) {
	originalJWTSecret := jwtSecret
	defer func() {
		jwtSecret = originalJWTSecret
	}()

	tests := []struct {
		name    string
		id      uint
		now     time.Time
		want    string
		wantErr bool
		setup   func()
	}{
		{
			name:    "Successfully Generate Token for Valid User ID",
			id:      123,
			now:     time.Now(),
			wantErr: false,
			setup:   func() {},
		},
		{
			name:    "Token Expiration Time Set Correctly",
			id:      456,
			now:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
			setup:   func() {},
		},
		{
			name:    "Token Contains Correct User ID",
			id:      789,
			now:     time.Now(),
			wantErr: false,
			setup:   func() {},
		},
		{
			name:    "Generate Token with Zero User ID",
			id:      0,
			now:     time.Now(),
			wantErr: false,
			setup:   func() {},
		},
		{
			name:    "Token Generation with Maximum uint Value",
			id:      ^uint(0),
			now:     time.Now(),
			wantErr: false,
			setup:   func() {},
		},
		{
			name:    "Token Generation with Empty JWT Secret",
			id:      123,
			now:     time.Now(),
			wantErr: true,
			setup: func() {
				jwtSecret = []byte{}
			},
		},
		{
			name:    "Consistency of Generated Tokens",
			id:      123,
			now:     time.Now(),
			wantErr: false,
			setup:   func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := generateToken(tt.id, tt.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == "" {
					t.Errorf("generateToken() returned empty token")
				}

				token, err := jwt.ParseWithClaims(got, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})

				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
					return
				}

				if claims, ok := token.Claims.(*claims); ok && token.Valid {
					if claims.UserID != tt.id {
						t.Errorf("Token UserID = %v, want %v", claims.UserID, tt.id)
					}

					if tt.name == "Token Expiration Time Set Correctly" {
						expectedExpiration := tt.now.Add(time.Hour * 72).Unix()
						if claims.ExpiresAt != expectedExpiration {
							t.Errorf("Token ExpiresAt = %v, want %v", claims.ExpiresAt, expectedExpiration)
						}
					}
				} else {
					t.Errorf("Token claims are not valid")
				}
			}

			if tt.name == "Consistency of Generated Tokens" {
				secondToken, _ := generateToken(tt.id, tt.now.Add(time.Second))
				if got == secondToken {
					t.Errorf("Generated tokens are identical, expected different tokens")
				}
			}
		})
	}
}

