package auth

import (
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"os"
	"context"
	"google.golang.org/grpc/metadata"
	"encoding/base64"
	"encoding/json"
	"math"
	"strings"
)









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
func TestGetUserID(t *testing.T) {

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
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedUserID: 123,
			expectedError:  "",
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
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedUserID: 0,
			expectedError:  "token expired",
		},
		{
			name: "Malformed Token",
			setupContext: func() context.Context {
				md := metadata.Pairs("authorization", "Token malformed.token.string")
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedUserID: 0,
			expectedError:  "invalid token: it's not even a token",
		},
		{
			name: "Missing Token in Context",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedUserID: 0,
			expectedError:  "Request unauthenticated with Token",
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
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedUserID: 0,
			expectedError:  "invalid token: couldn't handle this token",
		},
		{
			name: "Token with Future Not Before Claim",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Add(time.Hour).Unix(),
						ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
					},
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedUserID: 0,
			expectedError:  "invalid token: couldn't handle this token",
		},
		{
			name: "Token with Invalid Claims Type",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"user_id": 123,
					"exp":     time.Now().Add(time.Hour).Unix(),
				})
				tokenString, _ := token.SignedString(jwtSecret)
				md := metadata.Pairs("authorization", "Token "+tokenString)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedUserID: 0,
			expectedError:  "invalid token: cannot map token to claims",
		},
	}

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
					t.Errorf("Expected error %q, got %q", tt.expectedError, err.Error())
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
		name         string
		userID       uint
		now          time.Time
		setupFunc    func()
		wantErr      bool
		validateFunc func(t *testing.T, token string)
	}{
		{
			name:   "Successful Token Generation",
			userID: 123,
			now:    time.Now(),
			validateFunc: func(t *testing.T, token string) {
				if token == "" {
					t.Error("Expected non-empty token, got empty string")
				}
			},
		},
		{
			name:   "Token Expiration Time",
			userID: 456,
			now:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			validateFunc: func(t *testing.T, token string) {
				claims := validateAndExtractClaims(t, token)
				expectedExp := time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC).Unix()
				if claims.ExpiresAt != expectedExp {
					t.Errorf("Expected expiration time %v, got %v", expectedExp, claims.ExpiresAt)
				}
			},
		},
		{
			name:   "User ID in Claims",
			userID: 789,
			now:    time.Now(),
			validateFunc: func(t *testing.T, token string) {
				claims := validateAndExtractClaims(t, token)
				if claims.UserID != 789 {
					t.Errorf("Expected user ID 789, got %d", claims.UserID)
				}
			},
		},
		{
			name:   "Invalid JWT Secret",
			userID: 101,
			now:    time.Now(),
			setupFunc: func() {
				jwtSecret = []byte{}
			},
			wantErr: true,
		},
		{
			name:   "Large User ID",
			userID: math.MaxUint32,
			now:    time.Now(),
			validateFunc: func(t *testing.T, token string) {
				claims := validateAndExtractClaims(t, token)
				if claims.UserID != math.MaxUint32 {
					t.Errorf("Expected user ID %d, got %d", math.MaxUint32, claims.UserID)
				}
			},
		},
		{
			name:   "Token Signing Method",
			userID: 202,
			now:    time.Now(),
			validateFunc: func(t *testing.T, token string) {
				parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return nil, nil
				})
				if _, ok := parsedToken.Method.(*jwt.SigningMethodHMAC); !ok {
					t.Errorf("Expected signing method HMAC, got %v", parsedToken.Method)
				}
			},
		},
		{
			name:   "Consistency of Generated Tokens",
			userID: 303,
			now:    time.Now(),
			validateFunc: func(t *testing.T, token1 string) {
				token2, err := generateToken(303, time.Now())
				if err != nil {
					t.Fatalf("Failed to generate second token: %v", err)
				}
				if token1 == token2 {
					t.Error("Expected different tokens, got identical tokens")
				}
			},
		},
		{
			name:   "Error Handling for Time Issues",
			userID: 404,
			now:    time.Now().Add(365 * 24 * time.Hour),
			validateFunc: func(t *testing.T, token string) {
				claims := validateAndExtractClaims(t, token)
				if claims.ExpiresAt <= time.Now().Unix() {
					t.Error("Expected future expiration time")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			token, err := generateToken(tt.userID, tt.now)

			if (err != nil) != tt.wantErr {
				t.Errorf("generateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validateFunc != nil {
				tt.validateFunc(t, token)
			}
		})
	}
}

func validateAndExtractClaims(t *testing.T, tokenString string) *claims {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		t.Fatalf("Invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("Failed to decode token payload: %v", err)
	}

	var claims claims
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		t.Fatalf("Failed to unmarshal claims: %v", err)
	}

	return &claims
}

