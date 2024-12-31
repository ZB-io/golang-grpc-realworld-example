package model

import (
	"testing"
	"regexp"

	pb "github.com/raahii/golang-grpc-realworld-example/proto"
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
				Username: "johndoe",
				Bio:      "Software developer",
				Image:    "https://example.com/johndoe.jpg",
			},
			following: true,
			want: &pb.Profile{
				Username:  "johndoe",
				Bio:       "Software developer",
				Image:     "https://example.com/johndoe.jpg",
				Following: true,
			},
		},
		{
			name: "Successful Profile Conversion with Following False",
			user: User{
				Username: "janedoe",
				Bio:      "UX Designer",
				Image:    "https://example.com/janedoe.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "janedoe",
				Bio:       "UX Designer",
				Image:     "https://example.com/janedoe.jpg",
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
				Bio:      "こんにちは世界",
				Image:    "https://example.com/世界.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "用户名",
				Bio:       "こんにちは世界",
				Image:     "https://example.com/世界.jpg",
				Following: false,
			},
		},
		{
			name: "Profile Conversion with Maximum Length Strings",
			user: User{
				Username: "a_very_long_username_that_reaches_maximum_length",
				Bio:      "This is a very long bio that is designed to test the maximum length of the bio field in the User struct. It should be long enough to potentially cause issues if there are any length restrictions.",
				Image:    "https://example.com/a_very_long_image_url_that_reaches_maximum_length_for_testing_purposes.jpg",
			},
			following: true,
			want: &pb.Profile{
				Username:  "a_very_long_username_that_reaches_maximum_length",
				Bio:       "This is a very long bio that is designed to test the maximum length of the bio field in the User struct. It should be long enough to potentially cause issues if there are any length restrictions.",
				Image:     "https://example.com/a_very_long_image_url_that_reaches_maximum_length_for_testing_purposes.jpg",
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
				Image:    "http://example.com/image.jpg",
			},
			token: "valid_token",
			expected: &pb.User{
				Email:    "test@example.com",
				Token:    "valid_token",
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
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
				Username: "verylongusernamewithalotofcharacters",
				Bio:      "This is a very long bio with a lot of characters to test the maximum length handling of the ProtoUser method",
				Image:    "https://very-long-domain-name.com/very-long-image-path/very-long-image-name-with-many-characters.jpg",
			},
			token: "very_long_token_string_to_test_maximum_length_handling",
			expected: &pb.User{
				Email:    "verylongemail@verylongdomain.com",
				Token:    "very_long_token_string_to_test_maximum_length_handling",
				Username: "verylongusernamewithalotofcharacters",
				Bio:      "This is a very long bio with a lot of characters to test the maximum length handling of the ProtoUser method",
				Image:    "https://very-long-domain-name.com/very-long-image-path/very-long-image-name-with-many-characters.jpg",
			},
		},
		{
			name: "Conversion with special characters",
			user: User{
				Email:    "special.chars+test@example.com",
				Username: "user_name-123",
				Bio:      "Bio with special chars: !@#$%^&*()_+{}[]|\\:;\"'<>,.?/~`",
				Image:    "http://example.com/image_with_パラメーター.jpg",
			},
			token: "token_with_special_chars!@#$%^&*()",
			expected: &pb.User{
				Email:    "special.chars+test@example.com",
				Token:    "token_with_special_chars!@#$%^&*()",
				Username: "user_name-123",
				Bio:      "Bio with special chars: !@#$%^&*()_+{}[]|\\:;\"'<>,.?/~`",
				Image:    "http://example.com/image_with_パラメーター.jpg",
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
			name:           "Unicode Characters in Password",
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
func TestUserHashPassword(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		expectedError  string
		validateHash   bool
		validateUnique bool
	}{
		{
			name:           "Successfully Hash a Valid Password",
			password:       "validPassword123",
			validateHash:   true,
			validateUnique: true,
		},
		{
			name:          "Attempt to Hash an Empty Password",
			password:      "",
			expectedError: "password should not be empty",
		},
		{
			name:         "Verify Hashed Password is Different from Original",
			password:     "originalPassword",
			validateHash: true,
		},
		{
			name:           "Consistency of Hashing",
			password:       "consistentPassword",
			validateUnique: true,
		},
		{
			name:         "Password Length Limit",
			password:     "a123456789012345678901234567890123456789012345678901234567890123",
			validateHash: true,
		},
		{
			name:         "Unicode Password Handling",
			password:     "パスワード123",
			validateHash: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Password: tt.password}
			err := user.HashPassword()

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("Expected error '%s', got '%v'", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.validateHash {
				if user.Password == tt.password {
					t.Errorf("Hashed password is the same as original")
				}

				match, _ := regexp.MatchString(`^\$2[ayb]\$.{56}$`, user.Password)
				if !match {
					t.Errorf("Hashed password does not match bcrypt format")
				}

				err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tt.password))
				if err != nil {
					t.Errorf("Failed to verify hashed password: %v", err)
				}
			}

			if tt.validateUnique {
				user2 := &User{Password: tt.password}
				err := user2.HashPassword()
				if err != nil {
					t.Errorf("Unexpected error in second hash: %v", err)
				}

				if user.Password == user2.Password {
					t.Errorf("Hashed passwords are not unique")
				}
			}
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
			name:    "All Fields Empty",
			user:    User{},
			wantErr: true,
			errMsg:  "username: cannot be blank; email: cannot be blank; password: cannot be blank.",
		},
		{
			name: "Valid Data with Optional Fields",
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
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("User.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
