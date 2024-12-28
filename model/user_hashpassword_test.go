package model

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var bcryptGenerateFromPassword = func(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}


func TestUserHashPassword(t *testing.T) {
	type test struct {
		description      string
		password         string
		expectError      bool
		expectedErrorMsg string
		passwordChange   bool
		mockBcryptErr    error
	}

	tests := []test{
		{
			description:    "Hashing a Non-Empty Password Successfully",
			password:       "password123",
			expectError:    false,
			passwordChange: true,
		},
		{
			description:      "Handling an Empty Password",
			password:         "",
			expectError:      true,
			expectedErrorMsg: "password should not be empty",
			passwordChange:   false,
		},
		{
			description:    "Confirm Password Remains Unchanged on Error",
			password:       "password123",
			expectError:    true,
			mockBcryptErr:  errors.New("bcrypt error"),
			passwordChange: false,
		},
		{
			description:      "Generating Error on Bcrypt Error Propagation",
			password:         "password123",
			expectError:      true,
			expectedErrorMsg: "bcrypt error",
			passwordChange:   false,
			mockBcryptErr:    errors.New("bcrypt error"),
		},
		{
			description:    "Hashing Edge Case with Very Long Password",
			password:       string(make([]byte, 1000)),
			expectError:    false,
			passwordChange: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			user := &User{
				Password: tc.password,
			}

			if tc.mockBcryptErr != nil {
				originalGenerateFromPassword := bcryptGenerateFromPassword
				defer func() { bcryptGenerateFromPassword = originalGenerateFromPassword }()
				bcryptGenerateFromPassword = func(password []byte, cost int) ([]byte, error) {
					return nil, tc.mockBcryptErr
				}
			}

			err := user.HashPassword()

			if tc.expectError {
				assert.Error(t, err)
				if tc.expectedErrorMsg != "" {
					assert.EqualError(t, err, tc.expectedErrorMsg)
				}
				if !tc.passwordChange {
					assert.Equal(t, tc.password, user.Password)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, user.Password)
				assert.NotEqual(t, tc.password, user.Password)
				if tc.passwordChange {
					t.Log("Password changed successfully during hashing.")
				}
			}
		})
	}
}

