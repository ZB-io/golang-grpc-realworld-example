package model

import (
	"testing"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"strings"
	"sync"
	"unicode/utf8"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jinzhu/gorm"
)

/*
ROOST_METHOD_HASH=CheckPassword_377b31181b
ROOST_METHOD_SIG_HASH=CheckPassword_e6e0413d83


 */
func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name           string
		storedHash     string
		inputPassword  string
		expectedResult bool
		description    string
	}{
		{
			name:           "Valid Password Match",
			storedHash:     generateHash(t, "correctPassword123"),
			inputPassword:  "correctPassword123",
			expectedResult: true,
			description:    "Testing correct password validation",
		},
		{
			name:           "Invalid Password Mismatch",
			storedHash:     generateHash(t, "correctPassword123"),
			inputPassword:  "wrongPassword123",
			expectedResult: false,
			description:    "Testing incorrect password rejection",
		},
		{
			name:           "Empty Password Check",
			storedHash:     generateHash(t, "somePassword123"),
			inputPassword:  "",
			expectedResult: false,
			description:    "Testing empty password handling",
		},
		{
			name:           "Empty Stored Hash",
			storedHash:     "",
			inputPassword:  "anyPassword123",
			expectedResult: false,
			description:    "Testing empty stored hash handling",
		},
		{
			name:           "Invalid Hash Format",
			storedHash:     "invalidhashformat",
			inputPassword:  "anyPassword123",
			expectedResult: false,
			description:    "Testing invalid hash format handling",
		},
		{
			name:           "Unicode Password Comparison",
			storedHash:     generateHash(t, "パスワード123"),
			inputPassword:  "パスワード123",
			expectedResult: true,
			description:    "Testing Unicode password validation",
		},
		{
			name:           "Maximum Length Password",
			storedHash:     generateHash(t, string(make([]byte, 72))),
			inputPassword:  string(make([]byte, 72)),
			expectedResult: true,
			description:    "Testing maximum length password handling",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.description)

			user := &User{
				Password: tt.storedHash,
			}

			result := user.CheckPassword(tt.inputPassword)

			if result != tt.expectedResult {
				t.Errorf("CheckPassword() = %v, want %v", result, tt.expectedResult)
			}

			if result {
				t.Log("Password verification succeeded as expected")
			} else {
				t.Log("Password verification failed as expected")
			}
		})
	}
}

func generateHash(t *testing.T, password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate hash: %v", err)
	}
	return string(hash)
}

/*
ROOST_METHOD_HASH=HashPassword_ea0347143c
ROOST_METHOD_SIG_HASH=HashPassword_fc69fabec5


 */
func TestHashPassword(t *testing.T) {

	tests := []struct {
		name        string
		password    string
		wantErr     bool
		errMessage  string
		description string
	}{
		{
			name:        "Valid Password",
			password:    "validPassword123",
			wantErr:     false,
			description: "Testing successful hashing of a valid password",
		},
		{
			name:        "Empty Password",
			password:    "",
			wantErr:     true,
			errMessage:  "password should not be empty",
			description: "Testing handling of empty password",
		},
		{
			name:        "Special Characters Password",
			password:    "!@#$%^&*()_+{}[]|\\:;\"'<>,.?/~`",
			wantErr:     false,
			description: "Testing password with special characters",
		},
		{
			name:        "Unicode Password",
			password:    "パスワード123アБВ",
			wantErr:     false,
			description: "Testing password with Unicode characters",
		},
		{
			name:        "Very Long Password",
			password:    strings.Repeat("a", 1024),
			wantErr:     false,
			description: "Testing handling of very long password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			user := &User{
				Password: tt.password,
			}
			originalPassword := user.Password

			err := user.HashPassword()

			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err.Error() != tt.errMessage {
					t.Errorf("HashPassword() error message = %v, want %v", err.Error(), tt.errMessage)
				}
				return
			}

			if user.Password == originalPassword {
				t.Error("HashPassword() did not change the password")
			}

			if !strings.HasPrefix(user.Password, "$2a$") {
				t.Error("HashPassword() did not generate a valid bcrypt hash")
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(originalPassword))
			if err != nil {
				t.Error("HashPassword() generated hash cannot be verified with original password")
			}

			if strings.Contains(tt.name, "Unicode") {
				if !utf8.ValidString(user.Password) {
					t.Error("HashPassword() produced invalid UTF-8 for Unicode password")
				}
			}
		})
	}

	t.Run("Concurrent Password Hashing", func(t *testing.T) {
		const numGoroutines = 10
		var wg sync.WaitGroup
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				user := &User{
					Password: fmt.Sprintf("password%d", i),
				}
				if err := user.HashPassword(); err != nil {
					errors <- err
				}
			}(i)
		}

		wg.Wait()
		close(errors)

		for err := range errors {
			if err != nil {
				t.Errorf("Concurrent HashPassword() failed: %v", err)
			}
		}
	})

	t.Run("Multiple Hash Calls", func(t *testing.T) {
		user := &User{
			Password: "testPassword",
		}

		err := user.HashPassword()
		if err != nil {
			t.Fatalf("First HashPassword() failed: %v", err)
		}
		firstHash := user.Password

		err = user.HashPassword()
		if err != nil {
			t.Fatalf("Second HashPassword() failed: %v", err)
		}

		if firstHash == user.Password {
			t.Error("Multiple HashPassword() calls produced identical hashes")
		}
	})
}

/*
ROOST_METHOD_HASH=Validate_532ff0c623
ROOST_METHOD_SIG_HASH=Validate_663e136f97


 */
func TestUserValidate(t *testing.T) {

	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid User Data",
			user: User{
				Username: "testuser123",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "Missing Username",
			user: User{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "username: cannot be blank",
		},
		{
			name: "Invalid Username Format",
			user: User{
				Username: "test@user#123",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "username: must be in a valid format",
		},
		{
			name: "Missing Email",
			user: User{
				Username: "testuser123",
				Email:    "",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "email: cannot be blank",
		},
		{
			name: "Invalid Email Format",
			user: User{
				Username: "testuser123",
				Email:    "notanemail",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "email: must be a valid email address",
		},
		{
			name: "Missing Password",
			user: User{
				Username: "testuser123",
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: true,
			errMsg:  "password: cannot be blank",
		},
		{
			name: "Multiple Validation Errors",
			user: User{
				Username: "",
				Email:    "notanemail",
				Password: "",
			},
			wantErr: true,
			errMsg:  "multiple validation errors",
		},
		{
			name: "Whitespace Username",
			user: User{
				Username: "   ",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "username: must be in a valid format",
		},
		{
			name: "Very Long Username",
			user: User{
				Username: strings.Repeat("a", 100),
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing scenario: %s", tt.name)

			err := tt.user.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if tt.errMsg != "multiple validation errors" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("User.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
				t.Logf("Validation failed as expected with error: %v", err)
			} else if !tt.wantErr {
				t.Log("Validation passed as expected")
			}
		})
	}
}

