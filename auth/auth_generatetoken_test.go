package auth

import (
	"os"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestgenerateToken(t *testing.T) {

	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	testCases := []struct {
		Name    string
		UserID  uint
		Now     time.Time
		SetEnv  bool
		EnvVal  string
		IsValid bool
	}{
		{
			Name:    "Scenario 1: Successful Token Generation",
			UserID:  123,
			Now:     time.Now(),
			SetEnv:  true,
			EnvVal:  "secret",
			IsValid: true,
		},
		{
			Name:    "Scenario 2: Token Expiration Claims Correctness",
			UserID:  456,
			Now:     time.Now(),
			SetEnv:  true,
			EnvVal:  "secret",
			IsValid: true,
		},
		{
			Name:    "Scenario 3: Invalid Signing Key",
			UserID:  789,
			Now:     time.Now(),
			SetEnv:  false,
			IsValid: false,
		},
		{
			Name:    "Scenario 4: Incorrect UserID Input",
			UserID:  0,
			Now:     time.Now(),
			SetEnv:  true,
			EnvVal:  "secret",
			IsValid: true,
		},
		{
			Name:    "Scenario 5: Future Time Input",
			UserID:  101,
			Now:     time.Now().Add(time.Hour * 24),
			SetEnv:  true,
			EnvVal:  "secret",
			IsValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.SetEnv {
				os.Setenv("JWT_SECRET", tc.EnvVal)
			} else {
				os.Setenv("JWT_SECRET", "")
			}

			tokenStr, err := generateToken(tc.UserID, tc.Now)
			if tc.IsValid {
				assert.NoError(t, err, "Expected error not to occur")
				assert.NotEmpty(t, tokenStr, "Token should not be empty")
				token, parseErr := jwt.ParseWithClaims(tokenStr, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(tc.EnvVal), nil
				})
				assert.NoError(t, parseErr, "Expected token parsing to succeed")
				if claims, ok := token.Claims.(*claims); ok && token.Valid {
					if tc.Name == "Scenario 2: Token Expiration Claims Correctness" {
						expectedExpiresAt := tc.Now.Add(72 * time.Hour).Unix()
						assert.Equal(t, expectedExpiresAt, claims.ExpiresAt, "Expiration time does not match expected value")
					}
				} else {
					t.Logf("Token claims invalid or the token itself is invalid")
				}
			} else {
				assert.Error(t, err, "Expected an error due to invalid parameters")
				assert.Empty(t, tokenStr, "Token should be empty on failure")
			}
		})
	}

}
