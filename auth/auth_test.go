package auth

import (
	"os"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
ROOST_METHOD_HASH=GenerateTokenWithTime_d0df64aa69
ROOST_METHOD_SIG_HASH=GenerateTokenWithTime_72dd09cde6

FUNCTION_DEF=func GenerateTokenWithTime(id uint, t time.Time) (string, error) 

 */
func TestGenerateMultipleTokens(t *testing.T) {

	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	id := uint(1)
	time1 := time.Now()
	time2 := time1.Add(time.Second)

	token1, err1 := GenerateTokenWithTime(id, time1)
	token2, err2 := GenerateTokenWithTime(id, time2)

	if err1 != nil || err2 != nil {
		t.Errorf("Failed to generate tokens: %v, %v", err1, err2)
		return
	}

	if token1 == token2 {
		t.Errorf("Generated tokens are identical")
	}
}

func TestGenerateTokenWithEmptySecret(t *testing.T) {

	os.Unsetenv("JWT_SECRET")

	_, err := GenerateTokenWithTime(1, time.Now())
	if err == nil {
		t.Errorf("Expected error with empty JWT secret, got nil")
	}
}

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
			got, err := GenerateTokenWithTime(tt.id, tt.time)
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
			}
		})
	}
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
				return context.WithValue(context.Background(), "auth", "Token "+tokenString)
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
				return context.WithValue(context.Background(), "auth", "Token "+tokenString)
			},
			expectedUserID: 0,
			expectedError:  "token expired",
		},
		{
			name: "Malformed Token",
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), "auth", "Token invalid_token")
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
			name: "Token with Invalid Claims Type",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"user_id": 123,
				})
				tokenString, _ := token.SignedString(jwtSecret)
				return context.WithValue(context.Background(), "auth", "Token "+tokenString)
			},
			expectedUserID: 0,
			expectedError:  "invalid token: cannot map token to claims",
		},
		{
			name: "Token with Future NotBefore Claim",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Add(time.Hour).Unix(),
						ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
					},
				})
				tokenString, _ := token.SignedString(jwtSecret)
				return context.WithValue(context.Background(), "auth", "Token "+tokenString)
			},
			expectedUserID: 0,
			expectedError:  "invalid token: couldn't handle this token",
		},
		{
			name: "Token with Incorrect Signature",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					UserID: 123,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				})
				tokenString, _ := token.SignedString([]byte("wrong_secret"))
				return context.WithValue(context.Background(), "auth", "Token "+tokenString)
			},
			expectedUserID: 0,
			expectedError:  "invalid token: couldn't handle this token",
		},
		{
			name: "Valid Token with Maximum Uint Value for UserID",
			setupContext: func() context.Context {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
					UserID: ^uint(0),
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				})
				tokenString, _ := token.SignedString(jwtSecret)
				return context.WithValue(context.Background(), "auth", "Token "+tokenString)
			},
			expectedUserID: ^uint(0),
			expectedError:  "",
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
				} else if !errors.Is(err, status.Error(codes.Unauthenticated, "Request unauthenticated with Token")) && err.Error() != tt.expectedError {
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

	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	os.Setenv("JWT_SECRET", "test-secret")

	tests := []struct {
		name        string
		userID      uint
		currentTime time.Time
		wantErr     bool
	}{
		{
			name:        "Successfully Generate Token for Valid User ID",
			userID:      1234,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:     false,
		},
		{
			name:        "Token Expiration Time Set Correctly",
			userID:      5678,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:     false,
		},
		{
			name:        "User ID Encoded Correctly in Token",
			userID:      5678,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:     false,
		},
		{
			name:        "Generate Tokens for Multiple Users",
			userID:      1111,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:     false,
		},
		{
			name:        "Token Generation with Maximum uint Value",
			userID:      ^uint(0),
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:     false,
		},
		{
			name:        "Consistency of Generated Tokens",
			userID:      9999,
			currentTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := generateToken(tt.userID, tt.currentTime)

			if (err != nil) != tt.wantErr {
				t.Errorf("generateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if token == "" {
					t.Errorf("generateToken() returned an empty token")
				}

				claims := &claims{}
				parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})

				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}

				if !parsedToken.Valid {
					t.Errorf("Token is not valid")
				}

				if claims.UserID != tt.userID {
					t.Errorf("UserID in token %d does not match input %d", claims.UserID, tt.userID)
				}

				expectedExpiration := tt.currentTime.Add(time.Hour * 72).Unix()
				if claims.ExpiresAt != expectedExpiration {
					t.Errorf("Expiration time %v does not match expected %v", claims.ExpiresAt, expectedExpiration)
				}
			}
		})
	}

	t.Run("Error Handling with Empty JWT Secret", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "")
		_, err := generateToken(1234, time.Now())
		if err == nil {
			t.Errorf("Expected error with empty JWT secret, got nil")
		}
	})

	t.Run("Token Generation Performance", func(t *testing.T) {
		start := time.Now()
		for i := uint(0); i < 10000; i++ {
			_, err := generateToken(i, time.Now())
			if err != nil {
				t.Errorf("Error generating token: %v", err)
			}
		}
		duration := time.Since(start)
		t.Logf("Time taken to generate 10,000 tokens: %v", duration)

	})
}

