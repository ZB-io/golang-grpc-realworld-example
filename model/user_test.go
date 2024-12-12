package model

import (
	"errors"
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

/*
ROOST_METHOD_HASH=ProtoProfile_c70e154ff1
ROOST_METHOD_SIG_HASH=ProtoProfile_def254b98c


*/
func TestUserProtoProfile(t *testing.T) {
	type testCase struct {
		user        User
		following   bool
		expected    *pb.Profile
		description string
	}

	testCases := []testCase{
		{
			user: User{
				Username: "test_user",
				Bio:      "Test bio",
				Image:    "http://test.com/image.png",
			},
			following: false,
			expected: &pb.Profile{
				Username:  "test_user",
				Bio:       "Test bio",
				Image:     "http://test.com/image.png",
				Following: false,
			},
			description: "Typical User Conversion to Profile",
		},
		{
			user: User{
				Username: "test_user",
				Bio:      "Test bio",
				Image:    "http://test.com/image.png",
			},
			following: true,
			expected: &pb.Profile{
				Username:  "test_user",
				Bio:       "Test bio",
				Image:     "http://test.com/image.png",
				Following: true,
			},
			description: "Conversion with Following Set to True",
		},
		{
			user: User{
				Username: "minimal_user",
				Bio:      "",
				Image:    "",
			},
			following: false,
			expected: &pb.Profile{
				Username:  "minimal_user",
				Bio:       "",
				Image:     "",
				Following: false,
			},
			description: "Handle User with Minimal Information",
		},
		{
			user: User{
				Username: "user_with_special_&",
				Bio:      "Special bio #!@%",
				Image:    "http://image.special/?",
			},
			following: false,
			expected: &pb.Profile{
				Username:  "user_with_special_&",
				Bio:       "Special bio #!@%",
				Image:     "http://image.special/?",
				Following: false,
			},
			description: "Handle User with Special Characters in Fields",
		},
		{
			user: User{
				Username: "large_data_user",
				Bio:      string(make([]byte, 1000)),
				Image:    string(make([]byte, 512)),
			},
			following: false,
			expected: &pb.Profile{
				Username:  "large_data_user",
				Bio:       string(make([]byte, 1000)),
				Image:     string(make([]byte, 512)),
				Following: false,
			},
			description: "Large Data Set for User Bio and Image",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			t.Logf("Running Test Case: %s", tc.description)

			result := tc.user.ProtoProfile(tc.following)

			assert.Equal(t, tc.expected, result, "Expected: %v, but got: %v", tc.expected, result)

			t.Logf("Test Case Passed: %s", tc.description)
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
			name: "Normal Case with Valid User and Token",
			user: User{
				Username: "testUser",
				Email:    "test@example.com",
				Bio:      "A test user",
				Image:    "https://example.com/image.png",
			},
			token: "valid_token_123",
			expected: &pb.User{
				Email:    "test@example.com",
				Token:    "valid_token_123",
				Username: "testUser",
				Bio:      "A test user",
				Image:    "https://example.com/image.png",
			},
		},
		{
			name: "User with Default Values",
			user: User{
				Username: "",
				Email:    "",
				Bio:      "",
				Image:    "",
			},
			token: "some_token",
			expected: &pb.User{
				Email:    "",
				Token:    "some_token",
				Username: "",
				Bio:      "",
				Image:    "",
			},
		},
		{
			name: "User with Special Characters",
			user: User{
				Username: "test@#User",
				Email:    "special@@example.com",
				Bio:      "Bio with special chars #$%@!",
				Image:    "https://example.com/special.png",
			},
			token: "!#$%&'*+/=?^_`{|}~",
			expected: &pb.User{
				Email:    "special@@example.com",
				Token:    "!#$%&'*+/=?^_`{|}~",
				Username: "test@#User",
				Bio:      "Bio with special chars #$%@!",
				Image:    "https://example.com/special.png",
			},
		},
		{
			name: "Long Token String",
			user: User{
				Username: "longTokenUser",
				Email:    "long@example.com",
				Bio:      "Long token test",
				Image:    "https://example.com/long.png",
			},
			token: "a really really long token string that could potentially exceed normal lengths",
			expected: &pb.User{
				Email:    "long@example.com",
				Token:    "a really really long token string that could potentially exceed normal lengths",
				Username: "longTokenUser",
				Bio:      "Long token test",
				Image:    "https://example.com/long.png",
			},
		},
		{
			name: "Empty Token",
			user: User{
				Username: "emptyTokenUser",
				Email:    "empty@example.com",
				Bio:      "Empty token test",
				Image:    "https://example.com/empty.png",
			},
			token: "",
			expected: &pb.User{
				Email:    "empty@example.com",
				Token:    "",
				Username: "emptyTokenUser",
				Bio:      "Empty token test",
				Image:    "https://example.com/empty.png",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.ProtoUser(tt.token)
			assert.Equal(t, tt.expected, result, "The result should match the expected pb.User")

			t.Logf("For test case %q: expected %v, got %v", tt.name, tt.expected, result)
		})
	}
}

/*
ROOST_METHOD_HASH=CheckPassword_377b31181b
ROOST_METHOD_SIG_HASH=CheckPassword_e6e0413d83


 */
func TestUserCheckPassword(t *testing.T) {

	var testCases = []struct {
		description    string
		plainText      string
		storedPassword string
		expectation    bool
	}{
		{
			description: "Successful Password Check",
			plainText:   "securePassword123",
			expectation: true,
		},
		{
			description: "Incorrect Password Check",
			plainText:   "wrongPassword",
			expectation: false,
		},
		{
			description: "Empty Password Test",
			plainText:   "",
			expectation: false,
		},
		{
			description: "Long Password Input",
			plainText:   "aVeryLongPasswordStringToTestBufferBoundariesAndPerformance☃☃☃",
			expectation: true,
		},
		{
			description: "Short Password Input",
			plainText:   "a",
			expectation: true,
		},
	}

	successfulHash, _ := bcrypt.GenerateFromPassword([]byte("securePassword123"), bcrypt.DefaultCost)
	longPasswordHash, _ := bcrypt.GenerateFromPassword([]byte("aVeryLongPasswordStringToTestBufferBoundariesAndPerformance☃☃☃"), bcrypt.DefaultCost)
	shortPasswordHash, _ := bcrypt.GenerateFromPassword([]byte("a"), bcrypt.DefaultCost)

	storedPasswords := map[string]string{
		"Successful Password Check": string(successfulHash),
		"Long Password Input":       string(longPasswordHash),
		"Short Password Input":      string(shortPasswordHash),
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {

			storedPassword := storedPasswords[tc.description]
			if tc.description == "Incorrect Password Check" || tc.description == "Empty Password Test" {
				storedPassword = string(successfulHash)
			}

			user := User{
				Password: storedPassword,
			}

			result := user.CheckPassword(tc.plainText)

			t.Logf("Test Case: %s", tc.description)
			t.Logf("Comparing plain text: '%s' with stored hash: '%s'", tc.plainText, storedPassword)

			if result != tc.expectation {
				t.Errorf("Expected CheckPassword to return %v, got %v", tc.expectation, result)
			} else {
				t.Logf("Successfully validated: %v", result)
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
		mockGenerateFn func(password []byte, cost int) ([]byte, error)
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:        "Successfully Hash a Valid Password",
			password:    "validPassword123",
			expectError: false,
		},
		{
			name:           "Handle Empty Password Gracefully",
			password:       "",
			expectError:    true,
			expectedErrMsg: "password should not be empty",
		},
		{
			name:     "Simulate Failure of the Bcrypt Function",
			password: "anyPassword",
			mockGenerateFn: func(password []byte, cost int) ([]byte, error) {
				return nil, errors.New("mock bcrypt failure")
			},
			expectError:    true,
			expectedErrMsg: "mock bcrypt failure",
		},
		{
			name:        "Hashing of Maximum Length Password",
			password:    strings.Repeat("a", 72),
			expectError: false,
		},
		{
			name:        "Non-Alphanumeric Password Content Hashing",
			password:    "valid~Password@123!",
			expectError: false,
		},
		{
			name:        "Retaining Consistency on Identical Passwords",
			password:    "samePassword",
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user := User{
				Password: test.password,
			}

			// originalFunc := bcrypt.GenerateFromPassword
			// if test.mockGenerateFn != nil {
			// 	bcrypt.GenerateFromPassword = test.mockGenerateFn
			// }

			// defer func() { bcrypt.GenerateFromPassword = originalFunc }()

			err := user.HashPassword()

			if test.expectError {
				if err == nil || (test.expectedErrMsg != "" && !strings.Contains(err.Error(), test.expectedErrMsg)) {
					t.Errorf("expected error '%v', got '%v'", test.expectedErrMsg, err)
				} else {
					t.Logf("Correctly received error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else {
					if user.Password == test.password {
						t.Errorf("password was not hashed, remains the same: %s", user.Password)
					} else {
						t.Logf("Password hashed successfully. Original: %s, Hashed: %s", test.password, user.Password)
					}
				}
			}
		})
	}
}

/*
ROOST_METHOD_HASH=Validate_532ff0c623
ROOST_METHOD_SIG_HASH=Validate_663e136f97


 */
func TestUserValidate(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		wantErrs map[string]string
	}{
		{
			name: "Valid user information should pass validation",
			user: User{
				Username: "validUser123",
				Email:    "valid.email@example.com",
				Password: "strongpassword",
			},
			wantErrs: nil,
		},
		{
			name: "Missing username should fail validation",
			user: User{
				Username: "",
				Email:    "user@example.com",
				Password: "password",
			},
			wantErrs: map[string]string{"Username": "cannot be blank"},
		},
		{
			name: "Invalid username format should fail validation",
			user: User{
				Username: "invalid_user!@#",
				Email:    "user@example.com",
				Password: "password",
			},
			wantErrs: map[string]string{"Username": "must be in a valid format"},
		},
		{
			name: "Missing email should fail validation",
			user: User{
				Username: "user123",
				Email:    "",
				Password: "password",
			},
			wantErrs: map[string]string{"Email": "cannot be blank"},
		},
		{
			name: "Invalid email format should fail validation",
			user: User{
				Username: "user123",
				Email:    "user.com",
				Password: "password",
			},
			wantErrs: map[string]string{"Email": "must be a valid email address"},
		},
		{
			name: "Missing password should fail validation",
			user: User{
				Username: "user123",
				Email:    "user@example.com",
				Password: "",
			},
			wantErrs: map[string]string{"Password": "cannot be blank"},
		},
		{
			name: "All fields missing should fail validation",
			user: User{
				Username: "",
				Email:    "",
				Password: "",
			},
			wantErrs: map[string]string{
				"Username": "cannot be blank",
				"Email":    "cannot be blank",
				"Password": "cannot be blank",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if len(tt.wantErrs) == 0 && err != nil {
				t.Errorf("expected no errors, got %v", err)
			} else if len(tt.wantErrs) > 0 {
				if err == nil {
					t.Errorf("expected errors but got none")
				} else {
					errs, ok := err.(validation.Errors)
					if !ok {
						t.Errorf("expected validation errors, got %v", err)
					} else {
						for field, wantErr := range tt.wantErrs {
							if errMsg, exists := errs[field]; !exists || errMsg.Error() != wantErr {
								t.Errorf("expected error %s for field %s, got %v", wantErr, field, errMsg)
							}
						}
					}
				}
			}
			t.Logf("Test %s executed", tt.name)
		})
	}
}

// func (u User) Validate() error {
// 	return validation.ValidateStruct(&u,
// 		validation.Field(
// 			&u.Username,
// 			validation.Required,
// 			validation.Match(regexp.MustCompile("^[a-zA-Z0-9]+$")),
// 		),
// 		validation.Field(
// 			&u.Email,
// 			validation.Required,
// 			is.Email,
// 		),
// 		validation.Field(
// 			&u.Password,
// 			validation.Required,
// 		),
// 	)
// }

