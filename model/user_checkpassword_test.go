package model

import (
	"testing"
	"golang.org/x/crypto/bcrypt"
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
func TestCheckPassword(t *testing.T) {
	type testCase struct {
		description string
		user        User
		plainPass   string
		expected    bool
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("correctPassword123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate hashed password: %v", err)
	}

	longPassword := "aVeryLongPasswordContainingLotsAndLotsOfCharactersMoreThanYouWouldUsuallyExpectLikeWayMore1234567890"

	tests := []testCase{
		{
			description: "Correct Password Check",
			user:        User{Password: string(hashedPassword)},
			plainPass:   "correctPassword123",
			expected:    true,
		},
		{
			description: "Incorrect Password Check",
			user:        User{Password: string(hashedPassword)},
			plainPass:   "wrongPassword123",
			expected:    false,
		},
		{
			description: "Empty Password Check",
			user:        User{Password: string(hashedPassword)},
			plainPass:   "",
			expected:    false,
		},
		{
			description: "Password Hash Not Set Scenario",
			user:        User{Password: ""},
			plainPass:   "anyPassword",
			expected:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			t.Logf("Running scenario: %s", tc.description)
			result := tc.user.CheckPassword(tc.plainPass)
			if result != tc.expected {
				t.Errorf("Failed: %s, expected %v, got %v", tc.description, tc.expected, result)
			} else {
				t.Logf("Passed: %s", tc.description)
			}
		})
	}
}
