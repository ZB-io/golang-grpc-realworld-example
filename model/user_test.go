package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	reflect "reflect"

	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"golang.org/x/crypto/bcrypt"
)

var bcryptGenerateFromPassword = func(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

/*
ROOST_METHOD_HASH=ProtoProfile_c70e154ff1
ROOST_METHOD_SIG_HASH=ProtoProfile_def254b98c
*/
func TestProtoProfile(t *testing.T) {
	type testCase struct {
		description    string
		user           User
		following      bool
		expectedResult *pb.Profile
	}

	testCases := []testCase{
		{
			description: "Successful Profile Generation",
			user: User{
				Username: "john_doe",
				Bio:      "Software Developer",
				Image:    "http://example.com/image.jpg",
			},
			following: true,
			expectedResult: &pb.Profile{
				Username:  "john_doe",
				Bio:       "Software Developer",
				Image:     "http://example.com/image.jpg",
				Following: true,
			},
		},
		{
			description: "Empty User Profile",
			user: User{
				Username: "",
				Bio:      "",
				Image:    "",
			},
			following: false,
			expectedResult: &pb.Profile{
				Username:  "",
				Bio:       "",
				Image:     "",
				Following: false,
			},
		},
		{
			description: "Following Status Toggle",
			user: User{
				Username: "user",
				Bio:      "A regular user",
				Image:    "http://example.com/user.jpg",
			},
			following: true,
			expectedResult: &pb.Profile{
				Username:  "user",
				Bio:       "A regular user",
				Image:     "http://example.com/user.jpg",
				Following: true,
			},
		},
		{
			description: "Null Image Field",
			user: User{
				Username: "no_image_user",
				Bio:      "No image URL provided",
				Image:    "",
			},
			following: false,
			expectedResult: &pb.Profile{
				Username:  "no_image_user",
				Bio:       "No image URL provided",
				Image:     "",
				Following: false,
			},
		},
	}

	for i, tc := range testCases {
		t.Logf("Running test case %d: %s", i+1, tc.description)

		result := tc.user.ProtoProfile(tc.following)

		if result.Username != tc.expectedResult.Username {
			t.Errorf("Test case %d failed, expected Username: %v, got: %v", i+1, tc.expectedResult.Username, result.Username)
		}
		if result.Bio != tc.expectedResult.Bio {
			t.Errorf("Test case %d failed, expected Bio: %v, got: %v", i+1, tc.expectedResult.Bio, result.Bio)
		}
		if result.Image != tc.expectedResult.Image {
			t.Errorf("Test case %d failed, expected Image: %v, got: %v", i+1, tc.expectedResult.Image, result.Image)
		}
		if result.Following != tc.expectedResult.Following {
			t.Errorf("Test case %d failed, expected Following: %v, got: %v", i+1, tc.expectedResult.Following, result.Following)
		}

		t.Logf("Test case %d succeeded", i+1)
	}
}

/*
ROOST_METHOD_HASH=ProtoUser_440c1b101c
ROOST_METHOD_SIG_HASH=ProtoUser_fb8c4736ee
*/
func TestProtoUser(t *testing.T) {
	tests := []struct {
		name          string
		user          *User
		token         string
		expectedProto *pb.User
		shouldPanic   bool
	}{
		{
			name: "Convert User Struct to ProtoUser with Valid Data",
			user: &User{
				Email:    "test@example.com",
				Username: "testuser",
				Bio:      "Hello, world!",
				Image:    "http://example.com/image.jpg",
			},
			token: "validToken123",
			expectedProto: &pb.User{
				Email:    "test@example.com",
				Token:    "validToken123",
				Username: "testuser",
				Bio:      "Hello, world!",
				Image:    "http://example.com/image.jpg",
			},
			shouldPanic: false,
		},
		{
			name: "Handle Empty Token while Converting User Struct",
			user: &User{
				Email:    "test@example.com",
				Username: "testuser",
				Bio:      "Hello, world!",
				Image:    "http://example.com/image.jpg",
			},
			token: "",
			expectedProto: &pb.User{
				Email:    "test@example.com",
				Token:    "",
				Username: "testuser",
				Bio:      "Hello, world!",
				Image:    "http://example.com/image.jpg",
			},
			shouldPanic: false,
		},
		{
			name: "Convert User Struct with Empty User Fields",
			user: &User{
				Email:    "",
				Username: "",
				Bio:      "",
				Image:    "",
			},
			token: "validToken123",
			expectedProto: &pb.User{
				Email:    "",
				Token:    "validToken123",
				Username: "",
				Bio:      "",
				Image:    "",
			},
			shouldPanic: false,
		},
		{
			name: "Convert User Struct with Special Characters in Fields",
			user: &User{
				Email:    "special!@#example.com",
				Username: "specialUser!@#",
				Bio:      "Bio with special chars!@#",
				Image:    "http://example.com/image@.jpg",
			},
			token: "specialToken@123",
			expectedProto: &pb.User{
				Email:    "special!@#example.com",
				Token:    "specialToken@123",
				Username: "specialUser!@#",
				Bio:      "Bio with special chars!@#",
				Image:    "http://example.com/image@.jpg",
			},
			shouldPanic: false,
		},
		{
			name: "Convert User Struct with Long String Fields",
			user: &User{
				Email:    "verylongemailexample@example.com",
				Username: "averylongusername",
				Bio:      "A very long bio that contains many words and characters to simulate a lengthy entry in the user profile description.",
				Image:    "http://example.com/averylongimageurlthatcontainssufficientlymanycharacters.jpg",
			},
			token: "standardToken123",
			expectedProto: &pb.User{
				Email:    "verylongemailexample@example.com",
				Token:    "standardToken123",
				Username: "averylongusername",
				Bio:      "A very long bio that contains many words and characters to simulate a lengthy entry in the user profile description.",
				Image:    "http://example.com/averylongimageurlthatcontainssufficientlymanycharacters.jpg",
			},
			shouldPanic: false,
		},
		{
			name: "Validate Field Mapping Consistency Across Multiple Calls",
			user: &User{
				Email:    "consistent@example.com",
				Username: "consistentUser",
				Bio:      "Consistent Bio",
				Image:    "http://example.com/consistent.jpg",
			},
			token: "consistentToken",
			expectedProto: &pb.User{
				Email:    "consistent@example.com",
				Token:    "consistentToken",
				Username: "consistentUser",
				Bio:      "Consistent Bio",
				Image:    "http://example.com/consistent.jpg",
			},
			shouldPanic: false,
		},
		{
			name:          "Pass Nil User Struct to ProtoUser",
			user:          nil,
			token:         "validToken123",
			expectedProto: nil,
			shouldPanic:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("The code did not panic when it should have")
					}
				}()
			}

			result := tt.user.ProtoUser(tt.token)
			if !tt.shouldPanic {
				assert.True(t, reflect.DeepEqual(result, tt.expectedProto), "Expected %v, but got %v", tt.expectedProto, result)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=CheckPassword_377b31181b
ROOST_METHOD_SIG_HASH=CheckPassword_e6e0413d83
*/
func TestCheckPassword(t *testing.T) {
	type testCase struct {
		description string
		plainText   string
		hashedPass  string
		expected    bool
	}

	const bcryptCost = bcrypt.DefaultCost

	correctPasswordPlain := "CorrectPassword123!"
	correctPasswordHash, _ := bcrypt.GenerateFromPassword([]byte(correctPasswordPlain), bcryptCost)

	differentCostHash, _ := bcrypt.GenerateFromPassword([]byte(correctPasswordPlain), bcryptCost+1)

	malformedHash := "$2a$00$invalidandmalformedhashexample1234567890"

	testCases := []testCase{
		{
			description: "Correct Password",
			plainText:   correctPasswordPlain,
			hashedPass:  string(correctPasswordHash),
			expected:    true,
		},
		{
			description: "Incorrect Password",
			plainText:   "WrongPassword!",
			hashedPass:  string(correctPasswordHash),
			expected:    false,
		},
		{
			description: "Empty Plaintext Password",
			plainText:   "",
			hashedPass:  string(correctPasswordHash),
			expected:    false,
		},
		{
			description: "Empty Hashed Password",
			plainText:   correctPasswordPlain,
			hashedPass:  "",
			expected:    false,
		},
		{
			description: "Very Long Plaintext Password",
			plainText:   "aVeryLongPasswordThatExceedsNormalLengthByFarAndIsDefinitelyUnreasonableToUse",
			hashedPass:  string(correctPasswordHash),
			expected:    false,
		},
		{
			description: "Malformed Hashed Password",
			plainText:   correctPasswordPlain,
			hashedPass:  malformedHash,
			expected:    false,
		},
		{
			description: "Correct Hash with Different Hashing Rounds",
			plainText:   correctPasswordPlain,
			hashedPass:  string(differentCostHash),
			expected:    true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			user := &User{Password: tt.hashedPass}
			result := user.CheckPassword(tt.plainText)

			t.Log("Description:", tt.description)
			if result != tt.expected {
				t.Errorf("CheckPassword() = %v; want %v", result, tt.expected)
			} else {
				t.Log("Test passed for:", tt.description)
			}

			t.Logf("Test details: Plaintext=%q, Hashed=%q", tt.plainText, tt.hashedPass)
		})
	}
}

/*
ROOST_METHOD_HASH=HashPassword_ea0347143c
ROOST_METHOD_SIG_HASH=HashPassword_fc69fabec5
*/
func TestUserHashPassword(t *testing.T) {
	type test struct {
		description      string
		password         string
		expectError      bool
		expectedErrorMsg string
		passwordChange   bool
		mockBcryptErr    error
	}

	tests := []test{
		{
			description:    "Hashing a Non-Empty Password Successfully",
			password:       "password123",
			expectError:    false,
			passwordChange: true,
		},
		{
			description:      "Handling an Empty Password",
			password:         "",
			expectError:      true,
			expectedErrorMsg: "password should not be empty",
			passwordChange:   false,
		},
		{
			description:    "Confirm Password Remains Unchanged on Error",
			password:       "password123",
			expectError:    true,
			mockBcryptErr:  errors.New("bcrypt error"),
			passwordChange: false,
		},
		{
			description:      "Generating Error on Bcrypt Error Propagation",
			password:         "password123",
			expectError:      true,
			expectedErrorMsg: "bcrypt error",
			passwordChange:   false,
			mockBcryptErr:    errors.New("bcrypt error"),
		},
		{
			description:    "Hashing Edge Case with Very Long Password",
			password:       string(make([]byte, 1000)),
			expectError:    false,
			passwordChange: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			user := &User{
				Password: tc.password,
			}

			if tc.mockBcryptErr != nil {
				originalGenerateFromPassword := bcryptGenerateFromPassword
				defer func() { bcryptGenerateFromPassword = originalGenerateFromPassword }()
				bcryptGenerateFromPassword = func(password []byte, cost int) ([]byte, error) {
					return nil, tc.mockBcryptErr
				}
			}

			err := user.HashPassword()

			if tc.expectError {
				assert.Error(t, err)
				if tc.expectedErrorMsg != "" {
					assert.EqualError(t, err, tc.expectedErrorMsg)
				}
				if !tc.passwordChange {
					assert.Equal(t, tc.password, user.Password)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, user.Password)
				assert.NotEqual(t, tc.password, user.Password)
				if tc.passwordChange {
					t.Log("Password changed successfully during hashing.")
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
	type testCase struct {
		desc     string
		user     User
		expected error
	}

	testCases := []testCase{
		{
			desc: "Scenario 1: Validate a User with Correct Information",
			user: User{
				Username: "ValidUser123",
				Email:    "valid.email@example.com",
				Password: "securepassword",
			},
			expected: nil,
		},
		{
			desc: "Scenario 2: Username Contains Invalid Characters",
			user: User{
				Username: "InvalidUser!@#$",
				Email:    "valid.email@example.com",
				Password: "securepassword",
			},
			expected: validation.NewError("validation_match_error", "InvalidUser!@#$ does not match pattern [a-zA-Z0-9]+"),
		},
		{
			desc: "Scenario 3: Missing Email Field",
			user: User{
				Username: "ValidUser123",
				Email:    "",
				Password: "securepassword",
			},
			expected: validation.NewError("validation_required", "Email is required"),
		},
		{
			desc: "Scenario 4: Invalid Email Format",
			user: User{
				Username: "ValidUser123",
				Email:    "invalid-email",
				Password: "securepassword",
			},
			expected: validation.NewError("validation_invalid_format", "Email is not a valid email"),
		},
		{
			desc: "Scenario 5: Missing Password",
			user: User{
				Username: "ValidUser123",
				Email:    "valid.email@example.com",
				Password: "",
			},
			expected: validation.NewError("validation_required", "Password is required"),
		},
		{
			desc: "Scenario 6: All Validation Fields Missing",
			user: User{
				Username: "",
				Email:    "",
				Password: "",
			},
			expected: validation.Errors{
				"Username": validation.NewError("validation_required", "Username is required"),
				"Email":    validation.NewError("validation_required", "Email is required"),
				"Password": validation.NewError("validation_required", "Password is required"),
			},
		},
		{
			desc: "Scenario 7: Edge Case - Minimum Length Username",
			user: User{
				Username: "A",
				Email:    "valid.email@example.com",
				Password: "securepassword",
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Logf("Running: %s", tc.desc)
			err := tc.user.Validate()

			if tc.expected == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			} else if tc.expected != nil {
				assertError(t, err, tc.expected.Error())
			}
		})
	}
}

func assertError(t *testing.T, err error, expected string) {
	if err == nil || err.Error() != expected {
		t.Fatalf("expected error: %v, got: %v", expected, err)
	}
}
