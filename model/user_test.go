package model

import (
	"testing"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"errors"
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
			name: "Successfully create a Profile with following set to true",
			user: User{
				Username: "testuser",
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
		},
		{
			name: "Successfully create a Profile with following set to false",
			user: User{
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "testuser",
				Bio:       "Test bio",
				Image:     "http://example.com/image.jpg",
				Following: false,
			},
		},
		{
			name: "Create a Profile with empty Bio and Image fields",
			user: User{
				Username: "testuser",
				Bio:      "",
				Image:    "",
			},
			following: true,
			want: &pb.Profile{
				Username:  "testuser",
				Bio:       "",
				Image:     "",
				Following: true,
			},
		},
		{
			name: "Create a Profile with maximum length strings",
			user: User{
				Username: "testuser_with_very_long_username_that_is_allowed",
				Bio:      "This is a very long bio that contains the maximum allowed characters for testing purposes. It should be handled correctly by the ProtoProfile method without any truncation or modification.",
				Image:    "https://example.com/very/long/image/url/that/contains/maximum/allowed/characters/for/testing/purposes/image.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "testuser_with_very_long_username_that_is_allowed",
				Bio:       "This is a very long bio that contains the maximum allowed characters for testing purposes. It should be handled correctly by the ProtoProfile method without any truncation or modification.",
				Image:     "https://example.com/very/long/image/url/that/contains/maximum/allowed/characters/for/testing/purposes/image.jpg",
				Following: false,
			},
		},
		{
			name: "Handling of Unicode characters in User fields",
			user: User{
				Username: "用户名",
				Bio:      "这是一个包含Unicode字符的简介",
				Image:    "http://example.com/图片.jpg",
			},
			following: true,
			want: &pb.Profile{
				Username:  "用户名",
				Bio:       "这是一个包含Unicode字符的简介",
				Image:     "http://example.com/图片.jpg",
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

	t.Run("Consistency across multiple calls", func(t *testing.T) {
		user := User{
			Username: "testuser",
			Bio:      "Test bio",
			Image:    "http://example.com/image.jpg",
		}
		profile1 := user.ProtoProfile(true)
		profile2 := user.ProtoProfile(true)
		profile3 := user.ProtoProfile(true)

		assert.Equal(t, profile1, profile2)
		assert.Equal(t, profile2, profile3)
	})

	t.Run("Verify immutability of the original User object", func(t *testing.T) {
		user := User{
			Username: "testuser",
			Bio:      "Test bio",
			Image:    "http://example.com/image.jpg",
		}
		originalUser := user

		user.ProtoProfile(true)
		user.ProtoProfile(false)

		assert.Equal(t, originalUser, user)
	})
}


/*
ROOST_METHOD_HASH=ProtoUser_440c1b101c
ROOST_METHOD_SIG_HASH=ProtoUser_fb8c4736ee

FUNCTION_DEF=func (u *User) ProtoUser(token string) *pb.User 

 */
func TestUserProtoUser(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		token    string
		expected *pb.User
	}{
		{
			name: "Successful Conversion",
			user: &User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
			},
			token: "testtoken",
			expected: &pb.User{
				Email:    "test@example.com",
				Token:    "testtoken",
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
			},
		},
		{
			name: "Empty Fields Handling",
			user: &User{
				Model:    gorm.Model{ID: 2},
				Username: "emptyuser",
				Email:    "empty@example.com",
				Password: "password",
				Bio:      "",
				Image:    "",
			},
			token: "emptytoken",
			expected: &pb.User{
				Email:    "empty@example.com",
				Token:    "emptytoken",
				Username: "emptyuser",
				Bio:      "",
				Image:    "",
			},
		},
		{
			name: "Long String Values",
			user: &User{
				Model:    gorm.Model{ID: 3},
				Username: string(make([]byte, 1000)),
				Email:    "long@example.com",
				Password: "password",
				Bio:      string(make([]byte, 1000)),
				Image:    string(make([]byte, 1000)),
			},
			token: string(make([]byte, 1000)),
			expected: &pb.User{
				Email:    "long@example.com",
				Token:    string(make([]byte, 1000)),
				Username: string(make([]byte, 1000)),
				Bio:      string(make([]byte, 1000)),
				Image:    string(make([]byte, 1000)),
			},
		},
		{
			name: "Special Characters in Fields",
			user: &User{
				Model:    gorm.Model{ID: 4},
				Username: "user😊",
				Email:    "special@例子.com",
				Password: "password",
				Bio:      "Bio with ñ and é",
				Image:    "http://example.com/image_ñ.jpg",
			},
			token: "token_with_special_chars_@#$%^&*",
			expected: &pb.User{
				Email:    "special@例子.com",
				Token:    "token_with_special_chars_@#$%^&*",
				Username: "user😊",
				Bio:      "Bio with ñ and é",
				Image:    "http://example.com/image_ñ.jpg",
			},
		},
		{
			name:     "Null User Pointer Handling",
			user:     nil,
			token:    "validtoken",
			expected: nil,
		},
		{
			name: "Token Omission",
			user: &User{
				Model:    gorm.Model{ID: 5},
				Username: "tokenlessuser",
				Email:    "tokenless@example.com",
				Password: "password",
				Bio:      "No token bio",
				Image:    "http://example.com/no_token_image.jpg",
			},
			token: "",
			expected: &pb.User{
				Email:    "tokenless@example.com",
				Token:    "",
				Username: "tokenlessuser",
				Bio:      "No token bio",
				Image:    "http://example.com/no_token_image.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.user == nil {

				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic for nil user, but it didn't happen")
					}
				}()
			}

			got := tt.user.ProtoUser(tt.token)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ProtoUser() = %v, want %v", got, tt.expected)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=CheckPassword_377b31181b
ROOST_METHOD_SIG_HASH=CheckPassword_e6e0413d83

FUNCTION_DEF=func (u *User) CheckPassword(plain string) bool 

 */
func TestUserCheckPassword(t *testing.T) {
	tests := []struct {
		name           string
		hashedPassword string
		plainPassword  string
		expected       bool
	}{
		{
			name:           "Correct Password Verification",
			hashedPassword: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			plainPassword:  "password123",
			expected:       true,
		},
		{
			name:           "Incorrect Password Rejection",
			hashedPassword: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			plainPassword:  "wrongpassword",
			expected:       false,
		},
		{
			name:           "Empty Password Handling",
			hashedPassword: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			plainPassword:  "",
			expected:       false,
		},
		{
			name:           "Hashed Password Mismatch",
			hashedPassword: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			plainPassword:  "password123",
			expected:       true,
		},
		{
			name:           "Long Password Handling",
			hashedPassword: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			plainPassword:  "verylongpasswordthatisover100characterslong1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_+",
			expected:       false,
		},
		{
			name:           "Unicode Password Verification",
			hashedPassword: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			plainPassword:  "パスワード123",
			expected:       false,
		},
		{
			name:           "Null Byte Handling",
			hashedPassword: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			plainPassword:  "password\x00123",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Password: tt.hashedPassword,
			}
			result := user.CheckPassword(tt.plainPassword)
			if result != tt.expected {
				t.Errorf("CheckPassword() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUserCheckPasswordPerformance(t *testing.T) {

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := &User{
		Password: string(hashedPassword),
	}

	for i := 0; i < 1000; i++ {
		if !user.CheckPassword("correctpassword") {
			t.Errorf("CheckPassword() failed for correct password on iteration %d", i)
		}
	}

	for i := 0; i < 1000; i++ {
		if user.CheckPassword("incorrectpassword") {
			t.Errorf("CheckPassword() succeeded for incorrect password on iteration %d", i)
		}
	}
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
				if len(u.Password) != 60 {
					t.Errorf("Expected hashed password length to be 60, got %d", len(u.Password))
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
					t.Error("Hashed password is same as original")
				}
			},
		},
		{
			name:          "Consistency of Hashing",
			password:      "samePassword",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				u2 := &User{Password: "samePassword"}
				_ = u2.HashPassword()
				if u.Password == u2.Password {
					t.Error("Hashed passwords are identical for same input")
				}
			},
		},
		{
			name:          "Password Length After Hashing",
			password:      "testPassword",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if len(u.Password) != 60 {
					t.Errorf("Expected hashed password length to be 60, got %d", len(u.Password))
				}
			},
		},
		{
			name:          "Hashing a Very Long Password",
			password:      string(make([]byte, 1000)),
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if len(u.Password) != 60 {
					t.Errorf("Expected hashed password length to be 60, got %d", len(u.Password))
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
				return
			}

			if tt.validateResult != nil {
				tt.validateResult(t, u, err)
			}
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
			name: "Missing Username",
			user: User{
				Email:    "valid@example.com",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "username: cannot be blank.",
		},
		{
			name: "Invalid Username Format",
			user: User{
				Username: "invalid@user",
				Email:    "valid@example.com",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "username: must be in a valid format.",
		},
		{
			name: "Missing Email",
			user: User{
				Username: "validuser",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "email: cannot be blank.",
		},
		{
			name: "Invalid Email Format",
			user: User{
				Username: "validuser",
				Email:    "notanemail",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "email: must be a valid email address.",
		},
		{
			name: "Missing Password",
			user: User{
				Username: "validuser",
				Email:    "valid@example.com",
			},
			wantErr: true,
			errMsg:  "password: cannot be blank.",
		},
		{
			name: "Multiple Validation Errors",
			user: User{
				Username: "invalid@user",
				Email:    "notanemail",
			},
			wantErr: true,
			errMsg:  "email: must be a valid email address; password: cannot be blank; username: must be in a valid format.",
		},
		{
			name: "Valid Data with Optional Fields Empty",
			user: User{
				Username: "validuser",
				Email:    "valid@example.com",
				Password: "password123",
				Bio:      "",
				Image:    "",
			},
			wantErr: false,
		},
		{
			name: "Valid Data with Optional Fields Filled",
			user: User{
				Username: "validuser",
				Email:    "valid@example.com",
				Password: "password123",
				Bio:      "User bio",
				Image:    "user-image.jpg",
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
			if tt.wantErr {
				if err.Error() != tt.errMsg {
					t.Errorf("User.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

