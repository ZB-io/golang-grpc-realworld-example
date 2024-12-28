package model

import (
	"errors"
	"regexp"
	"testing"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)



func TestValidate(t *testing.T) {
	type testCase struct {
		desc     string
		user     User
		expected error
	}

	testCases := []testCase{
		{
			desc: "Scenario 1: Validate a User with Correct Information",
			user: User{
				Username: "ValidUser123",
				Email:    "valid.email@example.com",
				Password: "securepassword",
			},
			expected: nil,
		},
		{
			desc: "Scenario 2: Username Contains Invalid Characters",
			user: User{
				Username: "InvalidUser!@#$",
				Email:    "valid.email@example.com",
				Password: "securepassword",
			},
			expected: validation.NewError("validation_match_error", "InvalidUser!@#$ does not match pattern [a-zA-Z0-9]+"),
		},
		{
			desc: "Scenario 3: Missing Email Field",
			user: User{
				Username: "ValidUser123",
				Email:    "",
				Password: "securepassword",
			},
			expected: validation.NewError("validation_required", "Email is required"),
		},
		{
			desc: "Scenario 4: Invalid Email Format",
			user: User{
				Username: "ValidUser123",
				Email:    "invalid-email",
				Password: "securepassword",
			},
			expected: validation.NewError("validation_invalid_format", "Email is not a valid email"),
		},
		{
			desc: "Scenario 5: Missing Password",
			user: User{
				Username: "ValidUser123",
				Email:    "valid.email@example.com",
				Password: "",
			},
			expected: validation.NewError("validation_required", "Password is required"),
		},
		{
			desc: "Scenario 6: All Validation Fields Missing",
			user: User{
				Username: "",
				Email:    "",
				Password: "",
			},
			expected: validation.Errors{
				"Username": validation.NewError("validation_required", "Username is required"),
				"Email":    validation.NewError("validation_required", "Email is required"),
				"Password": validation.NewError("validation_required", "Password is required"),
			},
		},
		{
			desc: "Scenario 7: Edge Case - Minimum Length Username",
			user: User{
				Username: "A",
				Email:    "valid.email@example.com",
				Password: "securepassword",
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Logf("Running: %s", tc.desc)
			err := tc.user.Validate()

			if tc.expected == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			} else if tc.expected != nil {
				assertError(t, err, tc.expected.Error())
			}
		})
	}
}
func assertError(t *testing.T, err error, expected string) {
	if err == nil || err.Error() != expected {
		t.Fatalf("expected error: %v, got: %v", expected, err)
	}
}
