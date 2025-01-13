package model

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestCheckPassword(t *testing.T) {
	// Helper function to create a hashed password
	hashPassword := func(password string) string {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return string(hashedPassword)
	}

	tests := []struct {
		name           string
		hashedPassword string
		inputPassword  string
		expected       bool
	}{
		{
			name:           "Correct Password Match",
			hashedPassword: hashPassword("correctPassword"),
			inputPassword:  "correctPassword",
			expected:       true,
		},
		{
			name:           "Incorrect Password Mismatch",
			hashedPassword: hashPassword("correctPassword"),
			inputPassword:  "wrongPassword",
			expected:       false,
		},
		{
			name:           "Empty Password Input",
			hashedPassword: hashPassword("somePassword"),
			inputPassword:  "",
			expected:       false,
		},
		{
			name:           "Null Byte in Password",
			hashedPassword: hashPassword("normalPassword"),
			inputPassword:  "normal\x00Password",
			expected:       false,
		},
		{
			name:           "Very Long Password Input",
			hashedPassword: hashPassword("normalPassword"),
			inputPassword:  string(make([]byte, 1024*1024)), // 1MB of data
			expected:       false,
		},
		{
			name:           "Unicode Characters in Password",
			hashedPassword: hashPassword("パスワード123"),
			inputPassword:  "パスワード123",
			expected:       true,
		},
		{
			name:           "Case Sensitivity Check",
			hashedPassword: hashPassword("Password123"),
			inputPassword:  "password123",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Password: tt.hashedPassword}
			result := u.CheckPassword(tt.inputPassword)
			if result != tt.expected {
				t.Errorf("CheckPassword() = %v, want %v", result, tt.expected)
			}
		})
	}
}
