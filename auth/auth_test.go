package auth

import (
	"os"
	"testing"
	"time"
	"math"
	"github.com/dgrijalva/jwt-go"
	"context"
	"errors"
	"fmt"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/metadata"
)

var jwtSecret = []byte("testsecret")
/*
ROOST_METHOD_HASH=GenerateTokenWithTime_d0df64aa69
ROOST_METHOD_SIG_HASH=GenerateTokenWithTime_72dd09cde6


 */
func TestGenerateTokenWithTime(t *testing.T) {
	type testCase struct {
		description string
		id          uint
		time        time.Time
		shouldError bool
	}

	originalJWTSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalJWTSecret)

	testCases := []testCase{
		{
			description: "Valid Token Generation with Correct Inputs",
			id:          123,
			time:        time.Now(),
			shouldError: false,
		},
		{
			description: "Token Generation with a Historical Time Value",
			id:          123,
			time:        time.Now().AddDate(-1, 0, 0),
			shouldError: false,
		},
		{
			description: "Token Generation with Maximum Valid ID",
			id:          math.MaxUint32,
			time:        time.Now(),
			shouldError: false,
		},
		{
			description: "Empty JWT Secret Environment Variable",
			id:          123,
			time:        time.Now(),
			shouldError: true,
		},
		{
			description: "Invalid ID (Zero Value)",
			id:          0,
			time:        time.Now(),
			shouldError: true,
		},
		{
			description: "Handling Future Date for Token Generation",
			id:          123,
			time:        time.Now().AddDate(1, 0, 0),
			shouldError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {

			if tc.description == "Empty JWT Secret Environment Variable" {
				os.Setenv("JWT_SECRET", "")
			} else {
				os.Setenv("JWT_SECRET", "test_secret")
			}

			token, err := GenerateTokenWithTime(tc.id, tc.time)

			if tc.shouldError {
				if err == nil {
					t.Errorf("expected an error but did not get one")
					t.Logf("Scenario: %s - Expected failure due to environment or input conditions", tc.description)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error, but got: %v", err)
					t.Logf("Scenario: %s - Expected success in token generation, got error instead.", tc.description)
				} else if token == "" {
					t.Errorf("expected a valid token, but got an empty string")
					t.Logf("Scenario: %s - Expected non-empty token on valid input", tc.description)
				}
			}

			if err == nil && token != "" {
				t.Logf("Scenario: %s - Successfully generated token: %s", tc.description, token)
			}
		})
	}
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

	tests := []struct {
		name           string
		ctx            context.Context
		expectedUserID uint
		expectedErr    string
	}{
		{
			name:           "Successfully Retrieve User ID from Valid Token",
			ctx:            createTestContext("1", time.Now().Add(1*time.Hour).Unix()),
			expectedUserID: 1,
			expectedErr:    "",
		},
		{
			name:           "Handle Missing Authorization Metadata",
			ctx:            context.Background(),
			expectedUserID: 0,
			expectedErr:    "Request unauthenticated with Token",
		},
		{
			name:           "Handle Invalid Token Error",
			ctx:            createTestContextMalformed(),
			expectedUserID: 0,
			expectedErr:    "invalid token: it's not even a token",
		},
		{
			name:           "Handle Expired Token",
			ctx:            createTestContext("1", time.Now().Add(-1*time.Hour).Unix()),
			expectedUserID: 0,
			expectedErr:    "token expired",
		},
		{
			name:           "Handle Token with Invalid Claims",
			ctx:            createTestContextInvalidClaims(),
			expectedUserID: 0,
			expectedErr:    "invalid token: cannot map token to claims",
		},
		{
			name:           "Token with Invalid Signing Method",
			ctx:            createTestContextInvalidSigningMethod(),
			expectedUserID: 0,
			expectedErr:    "invalid signing method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := GetUserID(tt.ctx)

			if userID != tt.expectedUserID {
				t.Errorf("expected userID: %d, got: %d", tt.expectedUserID, userID)
			}

			if err != nil && err.Error() != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err.Error())
			}

			if err == nil && tt.expectedErr != "" {
				t.Errorf("expected error: %v, got: nil", tt.expectedErr)
			}

			t.Logf("Test Scenario: %s passed\n", tt.name)
		})
	}
}

func createTestContext(userID string, exp int64) context.Context {
	claims := &claims{
		UserID:    1,
		ExpiresAt: exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)

	md := metadata.Pairs("authorization", "Token "+tokenString)
	return metadata.NewIncomingContext(context.Background(), md)
}

func createTestContextInvalidClaims() context.Context {
	standardClaims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, standardClaims)
	tokenString, _ := token.SignedString(jwtSecret)

	md := metadata.Pairs("authorization", "Token "+tokenString)
	return metadata.NewIncomingContext(context.Background(), md)
}

func createTestContextInvalidSigningMethod() context.Context {
	claims := &claims{
		UserID: 1,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)

	md := metadata.Pairs("authorization", "Token "+tokenString)
	return metadata.NewIncomingContext(context.Background(), md)
}

func createTestContextMalformed() context.Context {
	md := metadata.Pairs("authorization", "Token malformed.token.data")
	return metadata.NewIncomingContext(context.Background(), md)
}

