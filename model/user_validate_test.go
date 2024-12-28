package model

import (
	"errors"
	"regexp"
	"testing"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
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
func TestUserValidate(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		expected error
	}{
		{
			name:     "Valid User Data",
			user:     User{Username: "JohnDoe", Email: "johndoe@example.com", Password: "password123"},
			expected: nil,
		},
		{
			name:     "Invalid Username (Empty)",
			user:     User{Username: "", Email: "johndoe@example.com", Password: "password123"},
			expected: errors.New("Username: cannot be blank."),
		},
		{
			name:     "Invalid Username (Non-Alphanumeric)",
			user:     User{Username: "John*Doe", Email: "johndoe@example.com", Password: "password123"},
			expected: errors.New("Username: must be in a valid format."),
		},
		{
			name:     "Invalid Email Format",
			user:     User{Username: "JohnDoe", Email: "johndoe", Password: "password123"},
			expected: errors.New("Email: must be a valid email address."),
		},
		{
			name:     "Missing Email",
			user:     User{Username: "JohnDoe", Email: "", Password: "password123"},
			expected: errors.New("Email: cannot be blank."),
		},
		{
			name:     "Missing Password",
			user:     User{Username: "JohnDoe", Email: "johndoe@example.com", Password: ""},
			expected: errors.New("Password: cannot be blank."),
		},
		{
			name:     "Edge Case (Minimum Valid Input)",
			user:     User{Username: "J", Email: "a@b.c", Password: "p"},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if err == nil && tt.expected == nil {
				t.Logf("%s: Passed. Expected no error, got no error.", tt.name)
			} else if err != nil && tt.expected != nil && err.Error() == tt.expected.Error() {
				t.Logf("%s: Passed. Expected error '%s', got error '%s'.", tt.name, tt.expected, err)
			} else {
				t.Errorf("%s: Failed. Expected error '%v', got '%v'.", tt.name, tt.expected, err)
			}
		})
	}
}

