package model

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)

/*
ROOST_METHOD_HASH=ProtoProfile_c70e154ff1
ROOST_METHOD_SIG_HASH=ProtoProfile_def254b98c
*/
func TestUserProtoProfile(t *testing.T) {
	tests := []struct {
		name      string
		user      User
		following bool
		want      *pb.Profile
	}{
		{
			name: "Successful Profile Conversion with Following True",
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
			name: "Successful Profile Conversion with Following False",
			user: User{
				Username: "anotheruser",
				Bio:      "Another bio",
				Image:    "http://example.com/another-image.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "anotheruser",
				Bio:       "Another bio",
				Image:     "http://example.com/another-image.jpg",
				Following: false,
			},
		},
		{
			name: "Profile Conversion with Empty User Fields",
			user: User{
				Username: "",
				Bio:      "",
				Image:    "",
			},
			following: true,
			want: &pb.Profile{
				Username:  "",
				Bio:       "",
				Image:     "",
				Following: true,
			},
		},
		{
			name: "Profile Conversion with Unicode Characters",
			user: User{
				Username: "用户名",
				Bio:      "这是一个测试简介",
				Image:    "http://example.com/图片.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "用户名",
				Bio:       "这是一个测试简介",
				Image:     "http://example.com/图片.jpg",
				Following: false,
			},
		},
		{
			name: "Profile Conversion with Maximum Length Strings",
			user: User{
				Username: strings.Repeat("a", 1000),
				Bio:      strings.Repeat("b", 1000),
				Image:    strings.Repeat("c", 1000),
			},
			following: true,
			want: &pb.Profile{
				Username:  strings.Repeat("a", 1000),
				Bio:       strings.Repeat("b", 1000),
				Image:     strings.Repeat("c", 1000),
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
*/
func TestProtoUser(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		token    string
		expected *pb.User
	}{
		{
			name:  "Conversion with zero values",
			user:  &User{},
			token: "",
			expected: &pb.User{
				Email:    "",
				Token:    "",
				Username: "",
				Bio:      "",
				Image:    "",
			},
		},
		{
			name: "Conversion with whitespace-only fields",
			user: &User{
				Email:    "   ",
				Username: "\t\n",
				Bio:      " ",
				Image:    "  ",
			},
			token: " ",
			expected: &pb.User{
				Email:    "   ",
				Token:    " ",
				Username: "\t\n",
				Bio:      " ",
				Image:    "  ",
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

func TestProtoUserConsistency(t *testing.T) {
	user := &User{
		Email:    "consistency@example.com",
		Username: "consistentuser",
		Bio:      "Consistent bio",
		Image:    "https://example.com/consistent.jpg",
	}
	token := "consistent_token"

	result1 := user.ProtoUser(token)
	result2 := user.ProtoUser(token)
	result3 := user.ProtoUser(token)

	assert.Equal(t, result1, result2, "Multiple calls to ProtoUser should return consistent results")
	assert.Equal(t, result2, result3, "Multiple calls to ProtoUser should return consistent results")
}

func TestProtoUserFieldMapping(t *testing.T) {
	user := &User{
		Email:    "mapping@example.com",
		Username: "mappinguser",
		Bio:      "Mapping bio",
		Image:    "https://example.com/mapping.jpg",
	}
	token := "mapping_token"

	result := user.ProtoUser(token)

	assert.Equal(t, user.Email, result.Email, "Email field should be correctly mapped")
	assert.Equal(t, token, result.Token, "Token should be correctly set")
	assert.Equal(t, user.Username, result.Username, "Username field should be correctly mapped")
	assert.Equal(t, user.Bio, result.Bio, "Bio field should be correctly mapped")
	assert.Equal(t, user.Image, result.Image, "Image field should be correctly mapped")
}

func TestProtoUserNilPointer(t *testing.T) {
	var nilUser *User
	result := nilUser.ProtoUser("token")
	assert.Nil(t, result, "ProtoUser should return nil for a nil User pointer")
}

/*
ROOST_METHOD_HASH=CheckPassword_377b31181b
ROOST_METHOD_SIG_HASH=CheckPassword_e6e0413d83
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
			name:           "Correct Password Verification",
			storedPassword: hashPassword("correctPassword"),
			inputPassword:  "correctPassword",
			expected:       true,
		},
		{
			name:           "Incorrect Password Rejection",
			storedPassword: hashPassword("correctPassword"),
			inputPassword:  "wrongPassword",
			expected:       false,
		},
		{
			name:           "Empty Password Handling",
			storedPassword: hashPassword("somePassword"),
			inputPassword:  "",
			expected:       false,
		},
		{
			name:           "Null Byte in Password",
			storedPassword: hashPassword("password"),
			inputPassword:  "pass\x00word",
			expected:       false,
		},
		{
			name:           "Very Long Password Input",
			storedPassword: hashPassword("normalPassword"),
			inputPassword:  string(make([]byte, 1024*1024)),
			expected:       false,
		},
		{
			name:           "Unicode Password Verification",
			storedPassword: hashPassword("パスワード"),
			inputPassword:  "パスワード",
			expected:       true,
		},
		{
			name:           "Case Sensitivity Check",
			storedPassword: hashPassword("Password"),
			inputPassword:  "password",
			expected:       false,
		},
		{
			name:           "Whitespace Handling",
			storedPassword: hashPassword("password"),
			inputPassword:  " password ",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Password: tt.storedPassword}
			result := user.CheckPassword(tt.inputPassword)
			if result != tt.expected {
				t.Errorf("CheckPassword() = %v, want %v", result, tt.expected)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=HashPassword_ea0347143c
ROOST_METHOD_SIG_HASH=HashPassword_fc69fabec5
*/
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
					t.Error("Password was not hashed")
				}
				if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("validPassword123")); err != nil {
					t.Error("Hashed password does not match original password")
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
			password:      strings.Repeat("a", 1000),
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if u.Password == strings.Repeat("a", 1000) {
					t.Error("Password was not hashed")
				}
				if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(strings.Repeat("a", 1000))); err != nil {
					t.Error("Hashed password does not match original password")
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
					t.Error("Password was not hashed")
				}
				if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("P@ssw0rd!@#$%^&*()_+")); err != nil {
					t.Error("Hashed password does not match original password")
				}
			},
		},
		{
			name:          "Attempt to Re-hash an Already Hashed Password",
			password:      "initialPassword",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error on first hash, got %v", err)
				}
				firstHash := u.Password
				err = u.HashPassword()
				if err != nil {
					t.Errorf("Expected no error on second hash, got %v", err)
				}
				if u.Password == firstHash {
					t.Error("Password hash did not change after second hashing")
				}
				if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("initialPassword")); err != nil {
					t.Error("Re-hashed password does not match original password")
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
			if tt.expectedError != nil && err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("HashPassword() error = %v, expectedError %v", err, tt.expectedError)
			}
			tt.validateResult(t, u, err)
		})
	}
}

/*
ROOST_METHOD_HASH=Validate_532ff0c623
ROOST_METHOD_SIG_HASH=Validate_663e136f97
*/
func TestValidate(t *testing.T) {
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
			errMsg:  "Username: cannot be blank",
		},
		{
			name: "Invalid Username Format",
			user: User{
				Username: "invalid@user",
				Email:    "valid@example.com",
				Password: "validpassword",
			},
			wantErr: true,
			errMsg:  "Username: must be in a valid format",
		},
		{
			name: "Missing Email",
			user: User{
				Username: "validuser",
				Password: "validpassword",
			},
			wantErr: true,
			errMsg:  "Email: cannot be blank",
		},
		{
			name: "Invalid Email Format",
			user: User{
				Username: "validuser",
				Email:    "invalidemail",
				Password: "validpassword",
			},
			wantErr: true,
			errMsg:  "Email: must be a valid email address",
		},
		{
			name: "Missing Password",
			user: User{
				Username: "validuser",
				Email:    "valid@example.com",
			},
			wantErr: true,
			errMsg:  "Password: cannot be blank",
		},
		{
			name: "Multiple Validation Errors",
			user: User{
				Email: "invalidemail",
			},
			wantErr: true,
			errMsg:  "Username: cannot be blank; Email: must be a valid email address; Password: cannot be blank",
		},
		{
			name: "Valid Data with Optional Fields",
			user: User{
				Username: "validuser",
				Email:    "valid@example.com",
				Password: "validpassword",
				Bio:      "",
				Image:    "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
