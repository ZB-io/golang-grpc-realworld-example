package model

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
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
				Username: "a" + string(make([]byte, 254)),
				Bio:      string(make([]byte, 1000)),
				Image:    "http://example.com/" + string(make([]byte, 980)),
			},
			following: true,
			want: &pb.Profile{
				Username:  "a" + string(make([]byte, 254)),
				Bio:       string(make([]byte, 1000)),
				Image:     "http://example.com/" + string(make([]byte, 980)),
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

func TestUserProtoProfileConsistency(t *testing.T) {
	user := User{
		Username: "consistentuser",
		Bio:      "Consistent bio",
		Image:    "http://example.com/consistent-image.jpg",
	}

	profile1 := user.ProtoProfile(true)
	profile2 := user.ProtoProfile(false)
	profile3 := user.ProtoProfile(true)

	assert.Equal(t, profile1.Username, profile2.Username)
	assert.Equal(t, profile1.Bio, profile2.Bio)
	assert.Equal(t, profile1.Image, profile2.Image)
	assert.True(t, profile1.Following)
	assert.False(t, profile2.Following)
	assert.Equal(t, profile1, profile3)
}

/*
ROOST_METHOD_HASH=ProtoUser_440c1b101c
ROOST_METHOD_SIG_HASH=ProtoUser_fb8c4736ee
*/
func TestUserProtoUser(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		token    string
		expected *pb.User
	}{
		{
			name: "Successful conversion with valid token",
			user: User{
				Email:    "test@example.com",
				Username: "testuser",
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
			name: "Conversion with empty fields",
			user: User{
				Email:    "",
				Username: "",
				Bio:      "",
				Image:    "",
			},
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
			name: "Conversion with maximum length values",
			user: User{
				Email:    "verylongemail@verylongdomain.com",
				Username: "verylongusername1234567890",
				Bio:      "This is a very long bio that tests the maximum length of the bio field in the User struct",
				Image:    "https://example.com/very/long/image/url/that/tests/maximum/length.jpg",
			},
			token: "very_long_token_string_to_test_maximum_length_of_token_field",
			expected: &pb.User{
				Email:    "verylongemail@verylongdomain.com",
				Token:    "very_long_token_string_to_test_maximum_length_of_token_field",
				Username: "verylongusername1234567890",
				Bio:      "This is a very long bio that tests the maximum length of the bio field in the User struct",
				Image:    "https://example.com/very/long/image/url/that/tests/maximum/length.jpg",
			},
		},
		{
			name: "Conversion with special characters",
			user: User{
				Email:    "special.chars+test@例子.com",
				Username: "user_name-123",
				Bio:      "Bio with émojis 😊 and Unicode characters ñ, é, ü",
				Image:    "https://example.com/image_with_$pecial_chars.jpg",
			},
			token: "token_with_$pecial_chars!@#",
			expected: &pb.User{
				Email:    "special.chars+test@例子.com",
				Token:    "token_with_$pecial_chars!@#",
				Username: "user_name-123",
				Bio:      "Bio with émojis 😊 and Unicode characters ñ, é, ü",
				Image:    "https://example.com/image_with_$pecial_chars.jpg",
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

func TestUserProtoUserPerformance(t *testing.T) {
	const numUsers = 10000

	users := make([]User, numUsers)
	for i := 0; i < numUsers; i++ {
		users[i] = User{
			Email:    "test@example.com",
			Username: "testuser",
			Bio:      "Test bio",
			Image:    "https://example.com/image.jpg",
		}
	}

	token := "performance_test_token"

	for i := 0; i < numUsers; i++ {
		result := users[i].ProtoUser(token)
		assert.NotNil(t, result)
	}
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
			storedPassword: hashPassword("パスワード123"),
			inputPassword:  "パスワード123",
			expected:       true,
		},
		{
			name:           "Case Sensitivity Check",
			storedPassword: hashPassword("Password123"),
			inputPassword:  "password123",
			expected:       false,
		},
		{
			name:           "Whitespace Handling",
			storedPassword: hashPassword("password123"),
			inputPassword:  " password123 ",
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
				if !strings.HasPrefix(u.Password, "$2a$") {
					t.Error("Hashed password does not have bcrypt prefix")
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
				if u.Password != "" {
					t.Error("Password should remain empty")
				}
			},
		},
		{
			name:          "Verify Password Field is Updated After Hashing",
			password:      "testPassword456",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("testPassword456")) != nil {
					t.Error("Hashed password does not match original password")
				}
			},
		},
		{
			name:          "Test with Maximum Length Password",
			password:      strings.Repeat("a", 72),
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if u.Password == strings.Repeat("a", 72) {
					t.Error("Password was not hashed")
				}
				if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(strings.Repeat("a", 72))) != nil {
					t.Error("Hashed password does not match original maximum length password")
				}
			},
		},
		{
			name:          "Verify Idempotency of Hashing",
			password:      "idempotentPassword",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				firstHash := u.Password
				err = u.HashPassword()
				if err != nil {
					t.Errorf("Second hash attempt failed: %v", err)
				}
				if u.Password == firstHash {
					t.Error("Password hash should change on second hashing attempt")
				}
			},
		},
		{
			name:          "Test with Special Characters in Password",
			password:      "P@ssw0rd!@#$%^&*()",
			expectedError: nil,
			validateResult: func(t *testing.T, u *User, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if u.Password == "P@ssw0rd!@#$%^&*()" {
					t.Error("Password was not hashed")
				}
				if !strings.HasPrefix(u.Password, "$2a$") {
					t.Error("Hashed password does not have bcrypt prefix")
				}
				if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("P@ssw0rd!@#$%^&*()")) != nil {
					t.Error("Hashed password does not match original password with special characters")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Password: tt.password}
			err := u.HashPassword()
			if !errors.Is(err, tt.expectedError) {
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
			errMsg:  "Username: cannot be blank",
		},
		{
			name: "Invalid Username Format",
			user: User{
				Username: "invalid@user",
				Email:    "valid@example.com",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "Username: must be in a valid format",
		},
		{
			name: "Missing Email",
			user: User{
				Username: "validuser",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "Email: cannot be blank",
		},
		{
			name: "Invalid Email Format",
			user: User{
				Username: "validuser",
				Email:    "invalidemail",
				Password: "password123",
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
			name: "Valid User with Optional Fields",
			user: User{
				Username: "validuser",
				Email:    "valid@example.com",
				Password: "password123",
				Bio:      "User bio",
				Image:    "user-image.jpg",
			},
			wantErr: false,
		},
		{
			name: "Multiple Validation Errors",
			user: User{
				Username: "invalid@user",
				Email:    "invalidemail",
			},
			wantErr: true,
			errMsg:  "Email: must be a valid email address; Password: cannot be blank; Username: must be in a valid format",
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
