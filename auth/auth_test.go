package auth

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/stretchr/testify/assert"
)

/*
ROOST_METHOD_HASH=GenerateTokenWithTime_d0df64aa69
ROOST_METHOD_SIG_HASH=GenerateTokenWithTime_72dd09cde6
*/
func TestGenerateTokenWithTime(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		t             time.Time
		expectedError bool
		validateFunc  func(string) bool
	}{
		{
			name:          "Valid Token Generation with Correct Inputs",
			id:            1,
			t:             time.Now(),
			expectedError: false,
			validateFunc:  func(token string) bool { return token != "" },
		},
		{
			name:          "Token Generation with a Historical Time Value",
			id:            2,
			t:             time.Now().AddDate(-1, 0, 0),
			expectedError: false,
			validateFunc:  func(token string) bool { return token != "" },
		},
		{
			name:          "Token Generation with Maximum Valid ID",
			id:            ^uint(0),
			t:             time.Now(),
			expectedError: false,
			validateFunc:  func(token string) bool { return token != "" },
		},
		{
			name:          "Empty JWT Secret Environment Variable",
			id:            3,
			t:             time.Now(),
			expectedError: true,
			validateFunc:  nil,
		},
		{
			name:          "Invalid ID (Zero Value)",
			id:            0,
			t:             time.Now(),
			expectedError: true,
			validateFunc:  func(token string) bool { return token == "" },
		},
		{
			name:          "Handling Future Date for Token Generation",
			id:            4,
			t:             time.Now().AddDate(1, 0, 0),
			expectedError: false,
			validateFunc:  func(token string) bool { return token != "" },
		},
	}

	originalJwtSecret := os.Getenv("JWT_SECRET")

	for _, tt := range tests {
		if tt.name == "Empty JWT Secret Environment Variable" {
			os.Setenv("JWT_SECRET", "")
		} else {
			os.Setenv("JWT_SECRET", "defaultSecret")
		}

		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running test case: %s", tt.name)

			token, err := GenerateTokenWithTime(tt.id, tt.t)

			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error: %v, got: %v", tt.expectedError, err)
			}

			if !tt.expectedError && tt.validateFunc != nil && !tt.validateFunc(token) {
				t.Errorf("Token validation failed. Generated token: %s", token)
			}

			if tt.expectedError && err == nil {
				t.Errorf("Expected an error but got none")
			} else if !tt.expectedError {
				t.Logf("Successfully generated token: %s", token)
			}
		})
	}

	os.Setenv("JWT_SECRET", originalJwtSecret)
}

/*
ROOST_METHOD_HASH=generateToken_2cc40e0108
ROOST_METHOD_SIG_HASH=generateToken_9de4114fe8
*/
func TestgenerateToken(t *testing.T) {
	originalJWTSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalJWTSecret)

	tests := []struct {
		name       string
		userID     uint
		now        time.Time
		expectErr  bool
		verifyFunc func(t *testing.T, token string, err error)
	}{
		{
			name:   "Successful Token Generation",
			userID: 1,
			now:    time.Now(),
			verifyFunc: func(t *testing.T, token string, err error) {
				if err != nil {
					t.Fatalf("Expected no error but got: %v", err)
				}
				if token == "" {
					t.Fatal("Expected a non-empty token")
				}
			},
		},
		{
			name:      "Invalid JWT Secret",
			userID:    1,
			now:       time.Now(),
			expectErr: true,
			verifyFunc: func(t *testing.T, token string, err error) {
				if err == nil {
					t.Fatalf("Expected an error due to invalid JWT secret but got none")
				}
			},
		},
		{
			name:   "Expiration Time Calculation",
			userID: 1,
			now:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			verifyFunc: func(t *testing.T, token string, err error) {
				if err != nil {
					t.Fatalf("Failed to generate token: %v", err)
				}
				parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if err != nil {
					t.Fatalf("Failed to parse token: %v", err)
				}
				claims, ok := parsedToken.Claims.(*claims)
				if !ok || !parsedToken.Valid {
					t.Fatal("Token claims are not valid")
				}
				expectedExpireAt := time.Date(2023, 10, 4, 0, 0, 0, 0, time.UTC).Unix()
				if claims.ExpiresAt != expectedExpireAt {
					t.Fatalf("Expected ExpiresAt %v, got %v", expectedExpireAt, claims.ExpiresAt)
				}
			},
		},
		{
			name:   "Boundary Testing with Minimal User ID",
			userID: 0,
			now:    time.Now(),
			verifyFunc: func(t *testing.T, token string, err error) {
				if err != nil {
					t.Fatalf("Expected no error but got: %v", err)
				}
				if token == "" {
					t.Fatal("Expected a non-empty token")
				}
			},
		},
		{
			name:   "Error Handling for Large User ID",
			userID: ^uint(0),
			now:    time.Now(),
			verifyFunc: func(t *testing.T, token string, err error) {
				if err != nil {
					t.Fatalf("Expected no error but got: %v", err)
				}
				if token == "" {
					t.Fatal("Expected a non-empty token")
				}
			},
		},
		{
			name:   "System Time Manipulation",
			userID: 1,
			now:    time.Unix(-2208988800, 0),
			verifyFunc: func(t *testing.T, token string, err error) {
				if err != nil {
					t.Fatalf("Expected no error but got: %v", err)
				}
				if token == "" {
					t.Fatal("Expected a non-empty token")
				}
				parsedToken, parseErr := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if parseErr != nil {
					t.Fatalf("Failed to parse token: %v", parseErr)
				}
				if !parsedToken.Valid {
					t.Fatal("Expected token to be valid, but it was invalid")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "Invalid JWT Secret" {
				os.Setenv("JWT_SECRET", "")
			} else {
				os.Setenv("JWT_SECRET", "myValidSecret")
			}
			token, err := generateToken(test.userID, test.now)
			test.verifyFunc(t, token, err)
		})
	}
}

/*
ROOST_METHOD_HASH=GetUserID_f2dd680cb2
ROOST_METHOD_SIG_HASH=GetUserID_e739312e3d
*/
func TestGetUserID(t *testing.T) {
	type testCase struct {
		desc           string
		token          string
		expectedUserID uint
		expectedError  error
		setupToken     func() string
	}

	validUserID := uint(123)

	testCases := []testCase{
		{
			desc: "Successfully Retrieve User ID from Valid Token",
			setupToken: func() string {
				claims := &claims{
					UserID: validUserID,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(jwtSecret)
				return tokenString
			},
			expectedUserID: validUserID,
			expectedError:  nil,
		},
		{
			desc: "Handle Missing Authorization Metadata",
			setupToken: func() string { return "" },
			expectedUserID: 0,
			expectedError:  grpc_auth.ErrUnauthenticated,
		},
		{
			desc: "Handle Invalid Token Error",
			setupToken: func() string { return "invalid.token.string" },
			expectedUserID: 0,
			expectedError:  errors.New("invalid token: it's not even a token"),
		},
		{
			desc: "Handle Expired Token",
			setupToken: func() string {
				claims := &claims{
					UserID: validUserID,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(jwtSecret)
				return tokenString
			},
			expectedUserID: 0,
			expectedError:  errors.New("token expired"),
		},
		{
			desc: "Handle Token with Invalid Claims",
			setupToken: func() string {
				type invalidClaims struct {
					Name string `json:"name"`
					jwt.StandardClaims
				}
				claims := &invalidClaims{
					Name: "UnauthorizedAccess",
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(jwtSecret)
				return tokenString
			},
			expectedUserID: 0,
			expectedError:  errors.New("invalid token: cannot map token to claims"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tokenString := tc.setupToken()
			if tokenString != "" {
				md := map[string]string{
					"authorization": "Token " + tokenString,
				}
				ctx := context.WithValue(context.Background(), grpc_auth.AuthHeaderPayloadKey{}, md)
				userID, err := GetUserID(ctx)

				assert.Equal(t, tc.expectedUserID, userID)
				if err != nil {
					assert.EqualError(t, err, tc.expectedError.Error())
				} else {
					assert.Nil(t, tc.expectedError)
				}
			} else {
				ctx := context.Background()
				userID, err := GetUserID(ctx)
				assert.Equal(t, tc.expectedUserID, userID)
				if err != nil {
					assert.EqualError(t, err, tc.expectedError.Error())
				} else {
					assert.Nil(t, tc.expectedError)
				}
			}
		})
	}
}
