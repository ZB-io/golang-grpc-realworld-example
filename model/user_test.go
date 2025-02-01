package model

import (
	"strings"
	"testing"
	"github.com/jinzhu/gorm"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"errors"
	"regexp"
	"golang.org/x/crypto/bcrypt"
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
	}{
		{
			name: "Basic Profile Conversion with Following True",
			user: &User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
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
			name: "Basic Profile Conversion with Following False",
			user: &User{
				Model:    gorm.Model{ID: 2},
				Username: "anotheruser",
				Email:    "another@example.com",
				Password: "password123",
				Bio:      "Another test bio",
				Image:    "http://example.com/another-image.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "anotheruser",
				Bio:       "Another test bio",
				Image:     "http://example.com/another-image.jpg",
				Following: false,
			},
		},
		{
			name: "Profile Conversion with Empty Bio and Image",
			user: &User{
				Model:    gorm.Model{ID: 3},
				Username: "emptyuser",
				Email:    "empty@example.com",
				Password: "emptypass",
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
			user: &User{
				Model:    gorm.Model{ID: 4},
				Username: "special@user",
				Email:    "special@example.com",
				Password: "pass!@#$%^&*()",
				Bio:      "Bio with ñ and 你好",
				Image:    "http://example.com/image-with-ñ.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "special@user",
				Bio:       "Bio with ñ and 你好",
				Image:     "http://example.com/image-with-ñ.jpg",
				Following: false,
			},
		},
		{
			name: "Profile Conversion with Maximum Length Fields",
			user: &User{
				Model:    gorm.Model{ID: 5},
				Username: strings.Repeat("a", 50),
				Email:    "max@example.com",
				Password: "maxpassword",
				Bio:      strings.Repeat("b", 500),
				Image:    strings.Repeat("i", 200),
			},
			following: true,
			want: &pb.Profile{
				Username:  strings.Repeat("a", 50),
				Bio:       strings.Repeat("b", 500),
				Image:     strings.Repeat("i", 200),
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

	t.Run("Profile Conversion with Nil User", func(t *testing.T) {
		var nilUser *User
		assert.Panics(t, func() { nilUser.ProtoProfile(true) })
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
		user     User
		token    string
		expected *pb.User
	}{
		{
			name: "Valid User with Token",
			user: User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			token: "valid_token",
			expected: &pb.User{
				Email:    "test@example.com",
				Token:    "valid_token",
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
		},
		{
			name: "User with Empty Token",
			user: User{
				Model:    gorm.Model{ID: 2},
				Username: "emptytoken",
				Email:    "empty@example.com",
				Password: "password",
				Bio:      "Empty token bio",
				Image:    "https://example.com/empty.jpg",
			},
			token: "",
			expected: &pb.User{
				Email:    "empty@example.com",
				Token:    "",
				Username: "emptytoken",
				Bio:      "Empty token bio",
				Image:    "https://example.com/empty.jpg",
			},
		},
		{
			name: "User with Empty Fields",
			user: User{
				Model:    gorm.Model{ID: 3},
				Username: "",
				Email:    "empty@fields.com",
				Password: "password",
				Bio:      "",
				Image:    "",
			},
			token: "empty_fields_token",
			expected: &pb.User{
				Email:    "empty@fields.com",
				Token:    "empty_fields_token",
				Username: "",
				Bio:      "",
				Image:    "",
			},
		},
		{
			name: "User with Maximum Length Fields",
			user: User{
				Model:    gorm.Model{ID: 4},
				Username: "maxlengthusername",
				Email:    "maxlength@example.com",
				Password: "password",
				Bio:      "This is a very long bio that reaches the maximum allowed length for testing purposes.",
				Image:    "https://example.com/very/long/image/url/that/reaches/maximum/length.jpg",
			},
			token: "max_length_token",
			expected: &pb.User{
				Email:    "maxlength@example.com",
				Token:    "max_length_token",
				Username: "maxlengthusername",
				Bio:      "This is a very long bio that reaches the maximum allowed length for testing purposes.",
				Image:    "https://example.com/very/long/image/url/that/reaches/maximum/length.jpg",
			},
		},
		{
			name: "User with Special Characters",
			user: User{
				Model:    gorm.Model{ID: 5},
				Username: "special_user_😊",
				Email:    "special@例子.com",
				Password: "password",
				Bio:      "Bio with special chars: àáâãäå",
				Image:    "https://example.com/image_🖼️.jpg",
			},
			token: "special_token_🔑",
			expected: &pb.User{
				Email:    "special@例子.com",
				Token:    "special_token_🔑",
				Username: "special_user_😊",
				Bio:      "Bio with special chars: àáâãäå",
				Image:    "https://example.com/image_🖼️.jpg",
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
				if !regexp.MustCompile(`^\$2a\$`).MatchString(u.Password) {
					t.Error("Password hash does not start with bcrypt identifier")
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
			name:          "Hash a Very Long Password",
			password:      string(make([]byte, 1000)),
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if !regexp.MustCompile(`^\$2a\$`).MatchString(u.Password) {
					t.Error("Password hash does not start with bcrypt identifier")
				}
			},
		},
		{
			name:          "Hash a Password with Special Characters",
			password:      "P@ssw0rd!@#$%^&*()",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if !regexp.MustCompile(`^\$2a\$`).MatchString(u.Password) {
					t.Error("Password hash does not start with bcrypt identifier")
				}
			},
		},
		{
			name:          "Attempt to Re-hash an Already Hashed Password",
			password:      "alreadyHashedPassword",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				firstHash := u.Password
				err = u.HashPassword()
				if err != nil {
					t.Errorf("Expected no error on second hash, got %v", err)
				}
				if firstHash == u.Password {
					t.Error("Password hash did not change on second hash")
				}
			},
		},
		{
			name:          "Hash a Password at the Minimum Allowed Length",
			password:      "min",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if !regexp.MustCompile(`^\$2a\$`).MatchString(u.Password) {
					t.Error("Password hash does not start with bcrypt identifier")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Model:    gorm.Model{},
				Username: "testuser",
				Email:    "test@example.com",
				Password: tt.password,
				Bio:      "Test bio",
				Image:    "test.jpg",
			}

			err := u.HashPassword()

			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("HashPassword() error = %v, expectedError %v", err, tt.expectedError)
			}

			tt.validateResult(t, u, err)
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
			name:           "Correct Password Match",
			storedPassword: hashPassword("correctPassword"),
			inputPassword:  "correctPassword",
			expected:       true,
		},
		{
			name:           "Incorrect Password Mismatch",
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
			name:           "Hashed Password is Empty",
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
			name:           "Unicode Characters in Password",
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
				Password: tt.storedPassword,
			}
			result := user.CheckPassword(tt.inputPassword)
			if result != tt.expected {
				t.Errorf("CheckPassword() = %v, want %v", result, tt.expected)
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
				Password: "validpassword",
			},
			wantErr: false,
		},
		{
			name: "Missing Username",
			user: User{
				Email:    "valid@example.com",
				Password: "validpassword",
			},
			wantErr: true,
			errMsg:  "username: cannot be blank.",
		},
		{
			name: "Invalid Username Format",
			user: User{
				Username: "invalid@user",
				Email:    "valid@example.com",
				Password: "validpassword",
			},
			wantErr: true,
			errMsg:  "username: must be in a valid format.",
		},
		{
			name: "Missing Email",
			user: User{
				Username: "validuser",
				Password: "validpassword",
			},
			wantErr: true,
			errMsg:  "email: cannot be blank.",
		},
		{
			name: "Invalid Email Format",
			user: User{
				Username: "validuser",
				Email:    "invalidemail",
				Password: "validpassword",
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
			name:    "All Fields Invalid",
			user:    User{},
			wantErr: true,
			errMsg:  "email: cannot be blank; password: cannot be blank; username: cannot be blank.",
		},
		{
			name: "Valid Username Edge Case",
			user: User{
				Username: "a",
				Email:    "valid@example.com",
				Password: "validpassword",
			},
			wantErr: false,
		},
		{
			name: "Long Valid Inputs",
			user: User{
				Username: "verylongusernamebutvalidalphanumeric123456789",
				Email:    "very.long.email.address@very.long.domain.name.com",
				Password: "verylongpasswordbutvalid1234567890!@#$%^&*()_+",
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

