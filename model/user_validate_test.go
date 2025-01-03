package model

import (
	"errors"
	"regexp"
	"testing"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jinzhu/gorm"
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
		expected string
	}{
		{
			name:     "Valid User Data",
			user:     User{Username: "JohnDoe", Email: "johndoe@example.com", Password: "password123"},
			expected: "",
		},
		{
			name:     "Invalid Username (Empty)",
			user:     User{Username: "", Email: "johndoe@example.com", Password: "password123"},
			expected: "Username: cannot be blank.",
		},
		{
			name:     "Invalid Username (Non-Alphanumeric)",
			user:     User{Username: "John*Doe", Email: "johndoe@example.com", Password: "password123"},
			expected: "Username: must be in a valid format.",
		},
		{
			name:     "Invalid Email Format",
			user:     User{Username: "JohnDoe", Email: "johndoe", Password: "password123"},
			expected: "Email: must be a valid email address.",
		},
		{
			name:     "Missing Email",
			user:     User{Username: "JohnDoe", Email: "", Password: "password123"},
			expected: "Email: cannot be blank.",
		},
		{
			name:     "Missing Password",
			user:     User{Username: "JohnDoe", Email: "johndoe@example.com", Password: ""},
			expected: "Password: cannot be blank.",
		},
		{
			name:     "Edge Case (Minimum Valid Input)",
			user:     User{Username: "J", Email: "a@b.c", Password: "p"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			errStr := ""
			if err != nil {
				errStr = err.Error()
			}
			if errStr == tt.expected {
				t.Logf("%s: Passed.", tt.name)
			} else {
				t.Errorf("%s: Failed. Expected error '%v', got '%v'.", tt.name, tt.expected, errStr)
			}
		})
	}
}
