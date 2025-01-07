package model

import (
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"unicode/utf8"
	"github.com/go-ozzo/ozzo-validation"
)








/*
ROOST_METHOD_HASH=ProtoProfile_c70e154ff1
ROOST_METHOD_SIG_HASH=ProtoProfile_def254b98c

FUNCTION_DEF=func (u *User) ProtoProfile(following bool) *pb.Profile 

 */
func TestUserProtoProfile(t *testing.T) {

	tests := []struct {
		name      string
		user      *User
		following bool
		want      *pb.Profile
		wantErr   bool
	}{
		{
			name: "Scenario 1: Basic Profile Conversion with Following Status True",
			user: &User{
				Username: "testuser",
				Email:    "test@example.com",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
			},
			following: true,
			want: &pb.Profile{
				Username:  "testuser",
				Bio:       "Test bio",
				Image:     "http://example.com/image.jpg",
				Following: true,
			},
			wantErr: false,
		},
		{
			name: "Scenario 2: Basic Profile Conversion with Following Status False",
			user: &User{
				Username: "testuser2",
				Email:    "test2@example.com",
				Bio:      "Test bio 2",
				Image:    "http://example.com/image2.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "testuser2",
				Bio:       "Test bio 2",
				Image:     "http://example.com/image2.jpg",
				Following: false,
			},
			wantErr: false,
		},
		{
			name: "Scenario 3: Profile Conversion with Empty Bio and Image",
			user: &User{
				Username: "emptyuser",
				Email:    "empty@example.com",
				Bio:      "",
				Image:    "",
			},
			following: true,
			want: &pb.Profile{
				Username:  "emptyuser",
				Bio:       "",
				Image:     "",
				Following: true,
			},
			wantErr: false,
		},
		{
			name: "Scenario 4: Profile Conversion with Special Characters",
			user: &User{
				Username: "special@user#123",
				Email:    "special@example.com",
				Bio:      "Bio with @#$%^&*",
				Image:    "http://example.com/image@123.jpg",
			},
			following: true,
			want: &pb.Profile{
				Username:  "special@user#123",
				Bio:       "Bio with @#$%^&*",
				Image:     "http://example.com/image@123.jpg",
				Following: true,
			},
			wantErr: false,
		},
		{
			name: "Scenario 5: Profile Conversion with Unicode Characters",
			user: &User{
				Username: "用户名",
				Email:    "unicode@example.com",
				Bio:      "バイオ 😊",
				Image:    "http://example.com/イメージ.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "用户名",
				Bio:       "バイオ 😊",
				Image:     "http://example.com/イメージ.jpg",
				Following: false,
			},
			wantErr: false,
		},
		{
			name:      "Scenario 6: Profile Conversion with Nil User",
			user:      nil,
			following: true,
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Starting test:", tt.name)

			if tt.user == nil && !tt.wantErr {
				t.Fatal("Test case error: nil user without expected error")
			}

			if tt.user == nil {
				defer func() {
					if r := recover(); r != nil {
						t.Log("Successfully caught nil user panic")
					}
				}()
			}

			got := tt.user.ProtoProfile(tt.following)

			if tt.wantErr {
				if got != nil {
					t.Errorf("ProtoProfile() = %v, want nil for error case", got)
				}
				return
			}

			assert.Equal(t, tt.want.Username, got.Username, "Username mismatch")
			assert.Equal(t, tt.want.Bio, got.Bio, "Bio mismatch")
			assert.Equal(t, tt.want.Image, got.Image, "Image mismatch")
			assert.Equal(t, tt.want.Following, got.Following, "Following status mismatch")

			t.Log("Successfully completed test:", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=ProtoUser_440c1b101c
ROOST_METHOD_SIG_HASH=ProtoUser_fb8c4736ee

FUNCTION_DEF=func (u *User) ProtoUser(token string) *pb.User 

 */
func TestUserProtoUser(t *testing.T) {

	type testCase struct {
		name      string
		user      *User
		token     string
		expected  *pb.User
		wantPanic bool
	}

	tests := []testCase{
		{
			name: "Scenario 1: Basic User Data Conversion with Token",
			user: &User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			token: "valid-token-123",
			expected: &pb.User{
				Email:    "test@example.com",
				Token:    "valid-token-123",
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
		},
		{
			name: "Scenario 2: Empty Fields Handling",
			user: &User{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "minimaluser",
				Email:    "minimal@example.com",
				Password: "password123",
				Bio:      "",
				Image:    "",
			},
			token: "token-123",
			expected: &pb.User{
				Email:    "minimal@example.com",
				Token:    "token-123",
				Username: "minimaluser",
				Bio:      "",
				Image:    "",
			},
		},
		{
			name: "Scenario 3: Unicode Character Handling",
			user: &User{
				Model: gorm.Model{
					ID:        3,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "用户名",
				Email:    "unicode@example.com",
				Password: "password123",
				Bio:      "バイオ 😊",
				Image:    "https://example.com/イメージ.jpg",
			},
			token: "token-unicode-123",
			expected: &pb.User{
				Email:    "unicode@example.com",
				Token:    "token-unicode-123",
				Username: "用户名",
				Bio:      "バイオ 😊",
				Image:    "https://example.com/イメージ.jpg",
			},
		},
		{
			name: "Scenario 4: Maximum Length Fields",
			user: &User{
				Model: gorm.Model{
					ID:        4,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: string(make([]byte, 255)),
				Email:    "long@example.com",
				Password: "password123",
				Bio:      string(make([]byte, 1000)),
				Image:    string(make([]byte, 1000)),
			},
			token: string(make([]byte, 500)),
			expected: &pb.User{
				Email:    "long@example.com",
				Token:    string(make([]byte, 500)),
				Username: string(make([]byte, 255)),
				Bio:      string(make([]byte, 1000)),
				Image:    string(make([]byte, 1000)),
			},
		},
		{
			name: "Scenario 5: Empty Token Handling",
			user: &User{
				Model: gorm.Model{
					ID:        5,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "emptytoken",
				Email:    "empty@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			token: "",
			expected: &pb.User{
				Email:    "empty@example.com",
				Token:    "",
				Username: "emptytoken",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
		},
		{
			name:      "Scenario 6: Nil User Pointer Handling",
			user:      nil,
			token:     "token-123",
			wantPanic: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Testing:", tc.name)

			if tc.wantPanic {
				assert.Panics(t, func() {
					tc.user.ProtoUser(tc.token)
				}, "Expected panic for nil user pointer")
				return
			}

			result := tc.user.ProtoUser(tc.token)

			assert.Equal(t, tc.expected.Email, result.Email, "Email mismatch")
			assert.Equal(t, tc.expected.Token, result.Token, "Token mismatch")
			assert.Equal(t, tc.expected.Username, result.Username, "Username mismatch")
			assert.Equal(t, tc.expected.Bio, result.Bio, "Bio mismatch")
			assert.Equal(t, tc.expected.Image, result.Image, "Image mismatch")

			t.Log("Test passed successfully")
		})
	}
}


/*
ROOST_METHOD_HASH=CheckPassword_377b31181b
ROOST_METHOD_SIG_HASH=CheckPassword_e6e0413d83

FUNCTION_DEF=func (u *User) CheckPassword(plain string) bool 

 */
func TestUserCheckPassword(t *testing.T) {

	type testCase struct {
		name          string
		storedHash    string
		inputPassword string
		expectedMatch bool
		setupFunction func(*User)
		description   string
	}

	tests := []testCase{
		{
			name: "Valid Password Match",
			setupFunction: func(u *User) {
				plainPass := "correctPassword123"
				hash, err := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)
				if err != nil {
					t.Fatalf("Failed to generate hash: %v", err)
				}
				u.Password = string(hash)
			},
			inputPassword: "correctPassword123",
			expectedMatch: true,
			description:   "Testing correct password validation",
		},
		{
			name: "Invalid Password Mismatch",
			setupFunction: func(u *User) {
				plainPass := "correctPassword123"
				hash, err := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)
				if err != nil {
					t.Fatalf("Failed to generate hash: %v", err)
				}
				u.Password = string(hash)
			},
			inputPassword: "wrongPassword123",
			expectedMatch: false,
			description:   "Testing incorrect password rejection",
		},
		{
			name: "Empty Password Check",
			setupFunction: func(u *User) {
				plainPass := "somePassword123"
				hash, err := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)
				if err != nil {
					t.Fatalf("Failed to generate hash: %v", err)
				}
				u.Password = string(hash)
			},
			inputPassword: "",
			expectedMatch: false,
			description:   "Testing empty password handling",
		},
		{
			name: "Empty Stored Hash",
			setupFunction: func(u *User) {
				u.Password = ""
			},
			inputPassword: "anyPassword123",
			expectedMatch: false,
			description:   "Testing empty stored hash handling",
		},
		{
			name: "Invalid Hash Format",
			setupFunction: func(u *User) {
				u.Password = "invalidhashformat"
			},
			inputPassword: "anyPassword123",
			expectedMatch: false,
			description:   "Testing invalid hash format handling",
		},
		{
			name: "Unicode Password Verification",
			setupFunction: func(u *User) {
				plainPass := "パスワード123"
				hash, err := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)
				if err != nil {
					t.Fatalf("Failed to generate hash: %v", err)
				}
				u.Password = string(hash)
			},
			inputPassword: "パスワード123",
			expectedMatch: true,
			description:   "Testing Unicode password handling",
		},
		{
			name: "Maximum Length Password",
			setupFunction: func(u *User) {

				maxLengthPass := string(make([]byte, 72))
				hash, err := bcrypt.GenerateFromPassword([]byte(maxLengthPass), bcrypt.DefaultCost)
				if err != nil {
					t.Fatalf("Failed to generate hash: %v", err)
				}
				u.Password = string(hash)
			},
			inputPassword: string(make([]byte, 72)),
			expectedMatch: true,
			description:   "Testing maximum length password handling",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			user := &User{}

			if tc.setupFunction != nil {
				tc.setupFunction(user)
			}

			result := user.CheckPassword(tc.inputPassword)

			if result != tc.expectedMatch {
				t.Errorf("Test case '%s' failed: expected match=%v, got match=%v",
					tc.name, tc.expectedMatch, result)
			}

			t.Logf("Test case '%s' passed: %s", tc.name, tc.description)
		})
	}
}


/*
ROOST_METHOD_HASH=HashPassword_ea0347143c
ROOST_METHOD_SIG_HASH=HashPassword_fc69fabec5

FUNCTION_DEF=func (u *User) HashPassword() error 

 */
func TestUserHashPassword(t *testing.T) {

	tests := []struct {
		name        string
		password    string
		wantErr     bool
		errMessage  string
		description string
	}{
		{
			name:        "Valid Password",
			password:    "SecurePass123!",
			wantErr:     false,
			description: "Testing successful password hashing with a valid password",
		},
		{
			name:        "Empty Password",
			password:    "",
			wantErr:     true,
			errMessage:  "password should not be empty",
			description: "Testing error handling for empty password",
		},
		{
			name:        "Special Characters Password",
			password:    "!@#$%^&*()_+{}[]|\\:;\"'<>,.?/~`",
			wantErr:     false,
			description: "Testing password hashing with special characters",
		},
		{
			name:        "Unicode Password",
			password:    "パスワード123アБВ",
			wantErr:     false,
			description: "Testing password hashing with Unicode characters",
		},
		{
			name:        "Very Long Password",
			password:    strings.Repeat("a", 1024),
			wantErr:     false,
			description: "Testing password hashing with a very long input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Logf("Running test scenario: %s", tt.description)

			user := &User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: tt.password,
			}

			originalPassword := user.Password

			err := user.HashPassword()

			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errMessage {
				t.Errorf("HashPassword() error message = %v, want %v", err.Error(), tt.errMessage)
				return
			}

			if !tt.wantErr {

				if user.Password == originalPassword {
					t.Error("HashPassword() did not modify the password")
					return
				}

				if !strings.HasPrefix(user.Password, "$2a$") {
					t.Error("HashPassword() generated hash does not appear to be a bcrypt hash")
					return
				}

				if len(user.Password) < 60 {
					t.Error("HashPassword() generated hash is too short for a bcrypt hash")
					return
				}

				if strings.Contains(tt.name, "Unicode") {
					if !utf8.ValidString(user.Password) {
						t.Error("HashPassword() generated invalid UTF-8 encoding")
						return
					}
				}
			}

			t.Logf("Test case '%s' completed successfully", tt.name)
		})
	}

	t.Run("Multiple Hash Operations", func(t *testing.T) {
		user := &User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "InitialPassword123",
		}

		previousHashes := make([]string, 3)
		for i := 0; i < 3; i++ {
			err := user.HashPassword()
			if err != nil {
				t.Errorf("HashPassword() failed on iteration %d: %v", i+1, err)
				return
			}
			previousHashes[i] = user.Password

			if i > 0 && previousHashes[i] == previousHashes[i-1] {
				t.Error("HashPassword() generated identical hashes on subsequent calls")
				return
			}
		}

		t.Log("Multiple hash operations completed successfully")
	})
}


/*
ROOST_METHOD_HASH=Validate_532ff0c623
ROOST_METHOD_SIG_HASH=Validate_663e136f97

FUNCTION_DEF=func (u User) Validate() error 

 */
func TestUserValidate(t *testing.T) {

	type testCase struct {
		name     string
		user     User
		wantErr  bool
		errField string
	}

	tests := []testCase{
		{
			name: "Valid User",
			user: User{
				Username: "validuser123",
				Email:    "valid@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "Missing Username",
			user: User{
				Email:    "valid@example.com",
				Password: "password123",
			},
			wantErr:  true,
			errField: "username",
		},
		{
			name: "Invalid Username Format",
			user: User{
				Username: "user@name",
				Email:    "valid@example.com",
				Password: "password123",
			},
			wantErr:  true,
			errField: "username",
		},
		{
			name: "Missing Email",
			user: User{
				Username: "validuser123",
				Password: "password123",
			},
			wantErr:  true,
			errField: "email",
		},
		{
			name: "Invalid Email Format",
			user: User{
				Username: "validuser123",
				Email:    "notanemail",
				Password: "password123",
			},
			wantErr:  true,
			errField: "email",
		},
		{
			name: "Missing Password",
			user: User{
				Username: "validuser123",
				Email:    "valid@example.com",
			},
			wantErr:  true,
			errField: "password",
		},
		{
			name: "Multiple Validation Failures",
			user: User{
				Email: "notanemail",
			},
			wantErr: true,
		},
		{
			name: "Boundary Case - Long Username",
			user: User{
				Username: "verylongusername123456789012345678901234567890",
				Email:    "valid@example.com",
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
				if tt.errField != "" {
					if validationErrors, ok := err.(validation.Errors); ok {
						if _, exists := validationErrors[tt.errField]; !exists {
							t.Errorf("Expected validation error for field %s, but got errors: %v",
								tt.errField, validationErrors)
						}
					} else {
						t.Errorf("Expected validation.Errors type, got %T", err)
					}
				}
				t.Logf("Got expected validation error: %v", err)
			}

			if !tt.wantErr && err == nil {
				t.Log("Validation passed as expected")
			}
		})
	}
}

