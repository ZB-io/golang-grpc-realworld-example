package auth

import (
	"os"
	"strconv"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
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
func TestGenerateTokenWithTime(t *testing.T) {

	tests := []struct {
		name      string
		id        uint
		tokenTime time.Time
		expectErr bool
		setupEnv  func()
	}{
		{
			name:      "Scenario 1: Generating a Token with Current Time",
			id:        123,
			tokenTime: time.Now(),
			expectErr: false,
			setupEnv: func() {

				os.Setenv("JWT_SECRET", "testsecret")
			},
		},
		{
			name:      "Scenario 2: Token Generation with a Zero User ID",
			id:        0,
			tokenTime: time.Now(),
			expectErr: true,
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "testsecret")
			},
		},
		{
			name:      "Scenario 3: Token Generation with a Past Date",
			id:        456,
			tokenTime: time.Now().AddDate(-1, 0, 0),
			expectErr: false,
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "testsecret")
			},
		},
		{
			name:      "Scenario 4: Token Generation with Future Date",
			id:        789,
			tokenTime: time.Now().AddDate(0, 1, 0),
			expectErr: false,
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "testsecret")
			},
		},
		{
			name:      "Scenario 5: Token Generation Without JWT Secret",
			id:        101,
			tokenTime: time.Now(),
			expectErr: true,
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "")
			},
		},
		{
			name:      "Scenario 6: Token Generation with Maximum User ID Value",
			id:        ^uint(0),
			tokenTime: time.Now(),
			expectErr: false,
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "testsecret")
			},
		},
		{
			name:      "Scenario 7: Error Handling for Invalid Time Formats",
			id:        112,
			tokenTime: time.Time{},
			expectErr: true,
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "testsecret")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			tc.setupEnv()

			token, err := GenerateTokenWithTime(tc.id, tc.tokenTime)

			if tc.expectErr && err == nil {
				t.Errorf("%s: expected an error but got a token %v", tc.name, token)
			} else if !tc.expectErr && err != nil {
				t.Errorf("%s: expected a token but got error %v", tc.name, err)
			} else if !tc.expectErr && token == "" {
				t.Errorf("%s: expected a non-empty token", tc.name)
			} else if err == nil {

				parsedToken, parseErr := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})
				if parseErr != nil || !parsedToken.Valid {
					t.Errorf("%s: error parsing token %v, parseErr: %v", tc.name, token, parseErr)
				} else {
					t.Logf("%s: token generated successfully", tc.name)
				}
			}
		})
	}
}
