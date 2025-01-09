package model

import (
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"strings"
)








/*
ROOST_METHOD_HASH=ProtoProfile_c70e154ff1
ROOST_METHOD_SIG_HASH=ProtoProfile_def254b98c

FUNCTION_DEF=func (u *User) ProtoProfile(following bool) *pb.Profile 

 */
func TestUserProtoProfile(t *testing.T) {
	tests := []struct {
		name      string
		user      User
		following bool
		want      *pb.Profile
	}{
		{
			name: "Basic Profile Conversion with Following True",
			user: User{
				Username: "johndoe",
				Bio:      "I'm a software developer",
				Image:    "https://example.com/avatar.jpg",
			},
			following: true,
			want: &pb.Profile{
				Username:  "johndoe",
				Bio:       "I'm a software developer",
				Image:     "https://example.com/avatar.jpg",
				Following: true,
			},
		},
		{
			name: "Basic Profile Conversion with Following False",
			user: User{
				Username: "janedoe",
				Bio:      "UX designer",
				Image:    "https://example.com/jane.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "janedoe",
				Bio:       "UX designer",
				Image:     "https://example.com/jane.jpg",
				Following: false,
			},
		},
		{
			name: "Profile Conversion with Empty Bio and Image",
			user: User{
				Username: "emptyuser",
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
		},
		{
			name: "Profile Conversion with Special Characters in Fields",
			user: User{
				Username: "user123_!@#",
				Bio:      "I ❤️ coding! 你好，世界！",
				Image:    "https://example.com/image_😊.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "user123_!@#",
				Bio:       "I ❤️ coding! 你好，世界！",
				Image:     "https://example.com/image_😊.jpg",
				Following: false,
			},
		},
		{
			name: "Profile Conversion with Maximum Length Fields",
			user: User{
				Username: "maxlengthusername1234567890",
				Bio:      "This is a very long bio that reaches the maximum allowed length for testing purposes. It should be preserved entirely in the resulting Profile struct without any truncation.",
				Image:    "https://example.com/very/long/image/url/that/reaches/maximum/allowed/length/for/testing/purposes/image.jpg",
			},
			following: true,
			want: &pb.Profile{
				Username:  "maxlengthusername1234567890",
				Bio:       "This is a very long bio that reaches the maximum allowed length for testing purposes. It should be preserved entirely in the resulting Profile struct without any truncation.",
				Image:     "https://example.com/very/long/image/url/that/reaches/maximum/allowed/length/for/testing/purposes/image.jpg",
				Following: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.ProtoProfile(tt.following)
			assert.Equal(t, tt.want, got)
		})
	}
}


/*
ROOST_METHOD_HASH=ProtoUser_440c1b101c
ROOST_METHOD_SIG_HASH=ProtoUser_fb8c4736ee

FUNCTION_DEF=func (u *User) ProtoUser(token string) *pb.User 

 */
func TestUserProtoUser(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		token    string
		expected *pb.User
	}{
		{
			name: "Valid User with Token",
			user: User{
				Username: "testuser",
				Email:    "test@example.com",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
			},
			token: "validtoken123",
			expected: &pb.User{
				Email:    "test@example.com",
				Token:    "validtoken123",
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
			},
		},
		{
			name: "User with Empty Token",
			user: User{
				Username: "emptytoken",
				Email:    "empty@example.com",
				Bio:      "Empty token bio",
				Image:    "http://example.com/empty.jpg",
			},
			token: "",
			expected: &pb.User{
				Email:    "empty@example.com",
				Token:    "",
				Username: "emptytoken",
				Bio:      "Empty token bio",
				Image:    "http://example.com/empty.jpg",
			},
		},
		{
			name: "User with Empty Fields",
			user: User{
				Username: "emptyfields",
				Email:    "empty@fields.com",
				Bio:      "",
				Image:    "",
			},
			token: "emptyfieldstoken",
			expected: &pb.User{
				Email:    "empty@fields.com",
				Token:    "emptyfieldstoken",
				Username: "emptyfields",
				Bio:      "",
				Image:    "",
			},
		},
		{
			name: "User with Maximum Length Fields",
			user: User{
				Username: "maxlengthuser",
				Email:    "maxlength@example.com",
				Bio:      string(make([]rune, 1000)),
				Image:    "http://example.com/" + string(make([]rune, 500)) + ".jpg",
			},
			token: string(make([]rune, 100)),
			expected: &pb.User{
				Email:    "maxlength@example.com",
				Token:    string(make([]rune, 100)),
				Username: "maxlengthuser",
				Bio:      string(make([]rune, 1000)),
				Image:    "http://example.com/" + string(make([]rune, 500)) + ".jpg",
			},
		},
		{
			name: "User with Unicode Characters",
			user: User{
				Username: "unicode_user_😊",
				Email:    "unicode@example.com",
				Bio:      "This is a bio with unicode: 你好世界",
				Image:    "http://example.com/unicode_🌍.jpg",
			},
			token: "unicode_token_🔑",
			expected: &pb.User{
				Email:    "unicode@example.com",
				Token:    "unicode_token_🔑",
				Username: "unicode_user_😊",
				Bio:      "This is a bio with unicode: 你好世界",
				Image:    "http://example.com/unicode_🌍.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.ProtoUser(tt.token)
			assert.Equal(t, tt.expected, result)
		})
	}
}


/*
ROOST_METHOD_HASH=CheckPassword_377b31181b
ROOST_METHOD_SIG_HASH=CheckPassword_e6e0413d83

FUNCTION_DEF=func (u *User) CheckPassword(plain string) bool 

 */
func TestUserCheckPassword(t *testing.T) {

	hashPassword := func(password string) string {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return string(hashedPassword)
	}

	tests := []struct {
		name           string
		storedPassword string
		inputPassword  string
		expected       bool
	}{
		{
			name:           "Correct Password",
			storedPassword: hashPassword("correctPassword"),
			inputPassword:  "correctPassword",
			expected:       true,
		},
		{
			name:           "Incorrect Password",
			storedPassword: hashPassword("correctPassword"),
			inputPassword:  "wrongPassword",
			expected:       false,
		},
		{
			name:           "Empty Password Input",
			storedPassword: hashPassword("somePassword"),
			inputPassword:  "",
			expected:       false,
		},
		{
			name:           "Empty Stored Password",
			storedPassword: "",
			inputPassword:  "anyPassword",
			expected:       false,
		},
		{
			name:           "Very Long Password Input",
			storedPassword: hashPassword("normalPassword"),
			inputPassword:  string(make([]byte, 1000)),
			expected:       false,
		},
		{
			name:           "Unicode Password",
			storedPassword: hashPassword("パスワード123"),
			inputPassword:  "パスワード123",
			expected:       true,
		},
		{
			name:           "Case Sensitivity Check",
			storedPassword: hashPassword("CaseSensitive"),
			inputPassword:  "casesensitive",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "testuser",
				Email:    "test@example.com",
				Password: tt.storedPassword,
				Bio:      "Test bio",
				Image:    "test.jpg",
			}

			result := user.CheckPassword(tt.inputPassword)
			if result != tt.expected {
				t.Errorf("CheckPassword() = %v, want %v", result, tt.expected)
			}
		})
	}

	t.Run("Timing Attack Resistance", func(t *testing.T) {
		user := &User{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Username: "testuser",
			Email:    "test@example.com",
			Password: hashPassword("timingTestPassword"),
			Bio:      "Test bio",
			Image:    "test.jpg",
		}

		correctStart := time.Now()
		user.CheckPassword("timingTestPassword")
		correctDuration := time.Since(correctStart)

		incorrectStart := time.Now()
		user.CheckPassword("wrongPassword")
		incorrectDuration := time.Since(incorrectStart)

		const acceptableTimingDifference = 10 * time.Millisecond
		if timingDiff := correctDuration - incorrectDuration; timingDiff > acceptableTimingDifference || timingDiff < -acceptableTimingDifference {
			t.Errorf("Potential timing attack vulnerability: correct password time: %v, incorrect password time: %v", correctDuration, incorrectDuration)
		}
	})
}


/*
ROOST_METHOD_HASH=HashPassword_ea0347143c
ROOST_METHOD_SIG_HASH=HashPassword_fc69fabec5

FUNCTION_DEF=func (u *User) HashPassword() error 

 */
func TestUserHashPassword(t *testing.T) {
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
					t.Error("Password was not hashed")
				}
				if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("validPassword123")); err != nil {
					t.Error("Hashed password does not match original")
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
					t.Error("Password was not hashed")
				}
			},
		},
		{
			name:          "Consistency of Hashing",
			password:      "consistentPassword",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				u2 := &User{Password: "consistentPassword"}
				if err := u2.HashPassword(); err != nil {
					t.Errorf("Error hashing second password: %v", err)
				}
				if u.Password == u2.Password {
					t.Error("Hashed passwords should be different due to salt")
				}
			},
		},
		{
			name:          "Password Length After Hashing",
			password:      "checkLengthPassword",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if len(u.Password) != 60 {
					t.Errorf("Expected hashed password length of 60, got %d", len(u.Password))
				}
			},
		},
		{
			name:          "Hashing a Very Long Password",
			password:      strings.Repeat("a", 1000),
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if len(u.Password) != 60 {
					t.Errorf("Expected hashed password length of 60, got %d", len(u.Password))
				}
			},
		},
		{
			name:          "Hashing a Password with Special Characters",
			password:      "P@ssw0rd!@#$%^&*()",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if u.Password == "P@ssw0rd!@#$%^&*()" {
					t.Error("Password was not hashed")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Password: tt.password}
			err := u.HashPassword()

			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("HashPassword() error = %v, expectedError %v", err, tt.expectedError)
			}

			tt.validateResult(t, u, err)
		})
	}
}


/*
ROOST_METHOD_HASH=Validate_532ff0c623
ROOST_METHOD_SIG_HASH=Validate_663e136f97

FUNCTION_DEF=func (u User) Validate() error 

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
				Username: "validuser",
				Email:    "valid@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "Empty Username",
			user: User{
				Username: "",
				Email:    "valid@example.com",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "Username: cannot be blank.",
		},
		{
			name: "Invalid Username Format",
			user: User{
				Username: "invalid@user",
				Email:    "valid@example.com",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "Username: must be in a valid format.",
		},
		{
			name: "Empty Email",
			user: User{
				Username: "validuser",
				Email:    "",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "Email: cannot be blank.",
		},
		{
			name: "Invalid Email Format",
			user: User{
				Username: "validuser",
				Email:    "notanemail",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "Email: must be a valid email address.",
		},
		{
			name: "Empty Password",
			user: User{
				Username: "validuser",
				Email:    "valid@example.com",
				Password: "",
			},
			wantErr: true,
			errMsg:  "Password: cannot be blank.",
		},
		{
			name: "Multiple Validation Errors",
			user: User{
				Username: "invalid@user",
				Email:    "notanemail",
				Password: "",
			},
			wantErr: true,
			errMsg:  "Username: must be in a valid format; Email: must be a valid email address; Password: cannot be blank.",
		},
		{
			name: "Valid Data with Optional Fields",
			user: User{
				Username: "validuser",
				Email:    "valid@example.com",
				Password: "password123",
				Bio:      "This is a bio",
				Image:    "http://example.com/image.jpg",
			},
			wantErr: false,
		},
		{
			name: "Maximum Length Username",
			user: User{
				Username: "verylongusernamebutvalidalphanumeric123456789",
				Email:    "valid@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("User.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

