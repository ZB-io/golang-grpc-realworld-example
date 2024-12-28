package model

import (
	"testing"
	"golang.org/x/crypto/bcrypt"
)

func TestCheckPassword(t *testing.T) {
	type testCase struct {
		description string
		plainText   string
		hashedPass  string
		expected    bool
	}

	const bcryptCost = bcrypt.DefaultCost

	correctPasswordPlain := "CorrectPassword123!"
	correctPasswordHash, _ := bcrypt.GenerateFromPassword([]byte(correctPasswordPlain), bcryptCost)

	differentCostHash, _ := bcrypt.GenerateFromPassword([]byte(correctPasswordPlain), bcryptCost+1)

	malformedHash := "$2a$00$invalidandmalformedhashexample1234567890"

	testCases := []testCase{
		{
			description: "Correct Password",
			plainText:   correctPasswordPlain,
			hashedPass:  string(correctPasswordHash),
			expected:    true,
		},
		{
			description: "Incorrect Password",
			plainText:   "WrongPassword!",
			hashedPass:  string(correctPasswordHash),
			expected:    false,
		},
		{
			description: "Empty Plaintext Password",
			plainText:   "",
			hashedPass:  string(correctPasswordHash),
			expected:    false,
		},
		{
			description: "Empty Hashed Password",
			plainText:   correctPasswordPlain,
			hashedPass:  "",
			expected:    false,
		},
		{
			description: "Very Long Plaintext Password",
			plainText:   "aVeryLongPasswordThatExceedsNormalLengthByFarAndIsDefinitelyUnreasonableToUse",
			hashedPass:  string(correctPasswordHash),
			expected:    false,
		},
		{
			description: "Malformed Hashed Password",
			plainText:   correctPasswordPlain,
			hashedPass:  malformedHash,
			expected:    false,
		},
		{
			description: "Correct Hash with Different Hashing Rounds",
			plainText:   correctPasswordPlain,
			hashedPass:  string(differentCostHash),
			expected:    true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			user := &User{Password: tt.hashedPass}
			result := user.CheckPassword(tt.plainText)

			t.Log("Description:", tt.description)
			if result != tt.expected {
				t.Errorf("CheckPassword() = %v; want %v", result, tt.expected)
			} else {
				t.Log("Test passed for:", tt.description)
			}

			t.Logf("Test details: Plaintext=%q, Hashed=%q", tt.plainText, tt.hashedPass)
		})
	}
}
