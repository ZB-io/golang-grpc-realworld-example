package model

import (
	"errors"
	"testing"
	"golang.org/x/crypto/bcrypt"
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
func TestUserHashPassword(t *testing.T) {
	type test struct {
		name        string
		user        User
		expectError error
		validate    func(t *testing.T, user User, err error)
	}

	tests := []test{
		{
			name: "Successful Password Hashing",
			user: User{
				Password: "secure123",
			},
			expectError: nil,
			validate: func(t *testing.T, user User, err error) {
				assert.Nil(t, err, "expected no error during hashing")
				assert.NotEqual(t, "secure123", user.Password, "hashed password should not match original")
				assert.NotEqual(t, "", user.Password, "hashed password should not be empty")
				t.Log("Password hashing altered original password successfully")
			},
		},
		{
			name: "Hashing an Empty Password",
			user: User{
				Password: "",
			},
			expectError: errors.New("password should not be empty"),
			validate: func(t *testing.T, user User, err error) {
				assert.NotNil(t, err, "expected an error for empty password")
				assert.Equal(t, "password should not be empty", err.Error(), "error message should match expected")
				t.Log("Proper error returned for empty password")
			},
		},
		{
			name: "Bcrypt Error Handling",
			user: User{
				Password: "secure123",
			},
			expectError: errors.New("bcrypt error simulation"),
			validate: func(t *testing.T, user User, err error) {
				assert.NotNil(t, err, "expected an error from bcrypt failure")
				assert.Equal(t, "bcrypt error simulation", err.Error(), "error message should match mocked error")
				t.Log("Handled bcrypt error as expected")
			},
		},
		{
			name: "Consistency of Hashing Outcome",
			user: User{
				Password: "secure123",
			},
			expectError: nil,
			validate: func(t *testing.T, user User, err error) {
				assert.Nil(t, err, "expected no error during hashing")
				err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("secure123"))
				assert.Nil(t, err, "hashed password should successfully match the original password")
				t.Log("Password consistently verifiable against original input")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.user.HashPassword()
			tc.validate(t, tc.user, err)
		})
	}
}
