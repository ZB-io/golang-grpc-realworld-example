package model

import (
	"errors"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		expectedError  error
		validateResult func(*testing.T, *User, error)
	}{
		{
			name:          "Successfully Hash a Valid Password",
			password:      "validPassword123",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if u.Password == "validPassword123" {
					t.Errorf("Password was not hashed")
				}
				if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("validPassword123")); err != nil {
					t.Errorf("Hashed password does not match original: %v", err)
				}
			},
		},
		{
			name:          "Attempt to Hash an Empty Password",
			password:      "",
			expectedError: errors.New("password should not be empty"),
			validateResult: func(t *testing.T, u *User, err error) {
				if err == nil || err.Error() != "password should not be empty" {
					t.Errorf("Expected error 'password should not be empty', got %v", err)
				}
			},
		},
		{
			name:          "Verify Hashed Password is Different from Original",
			password:      "originalPassword",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if u.Password == "originalPassword" {
					t.Errorf("Hashed password is same as original")
				}
			},
		},
		{
			name:          "Consistent Hashing for the Same Password",
			password:      "samePassword",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				u2 := &User{Password: "samePassword"}
				if err := u2.HashPassword(); err != nil {
					t.Errorf("Error hashing second password: %v", err)
				}
				if u.Password == u2.Password {
					t.Errorf("Hashed passwords are the same for two instances")
				}
			},
		},
		{
			name:          "Hash a Very Long Password",
			password:      strings.Repeat("a", 1000),
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if len(u.Password) == 0 {
					t.Errorf("Password was not hashed")
				}
			},
		},
		{
			name:          "Hash a Password with Special Characters",
			password:      "P@ssw0rd!@#$%^&*()_+",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if u.Password == "P@ssw0rd!@#$%^&*()_+" {
					t.Errorf("Password was not hashed")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Password: tt.password}
			err := u.HashPassword()
			if (err != nil) != (tt.expectedError != nil) {
				t.Errorf("HashPassword() error = %v, expectedError %v", err, tt.expectedError)
			}
			tt.validateResult(t, u, err)
		})
	}
}
