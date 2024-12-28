package auth

import (
	"os"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
)




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






