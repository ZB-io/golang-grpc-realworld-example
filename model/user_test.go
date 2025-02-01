package model

import (
	"testing"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
	"errors"
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
		user      User
		following bool
		want      *pb.Profile
	}{
		{
			name: "Basic Profile Conversion with Following True",
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
			name: "Basic Profile Conversion with Following False",
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
				Username: "special@user",
				Bio:      "Bio with $pecial ch@racters!",
				Image:    "http://example.com/image?special=true&id=123",
			},
			following: false,
			want: &pb.Profile{
				Username:  "special@user",
				Bio:       "Bio with $pecial ch@racters!",
				Image:     "http://example.com/image?special=true&id=123",
				Following: false,
			},
		},
		{
			name: "Profile Conversion with Maximum Length Fields",
			user: User{
				Username: "maxlengthuser",
				Bio:      "This is a very long bio that reaches the maximum allowed length for testing purposes. It should be preserved entirely in the output profile.",
				Image:    "http://example.com/very/long/image/url/that/reaches/maximum/allowed/length/for/testing/purposes.jpg",
			},
			following: true,
			want: &pb.Profile{
				Username:  "maxlengthuser",
				Bio:       "This is a very long bio that reaches the maximum allowed length for testing purposes. It should be preserved entirely in the output profile.",
				Image:     "http://example.com/very/long/image/url/that/reaches/maximum/allowed/length/for/testing/purposes.jpg",
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
		user     *User
		token    string
		expected *pb.User
	}{
		{
			name: "Valid User with Token",
			user: &User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
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
			name: "User with Empty Fields",
			user: &User{
				Model:    gorm.Model{ID: 2},
				Username: "",
				Email:    "empty@example.com",
				Password: "",
				Bio:      "",
				Image:    "",
			},
			token: "",
			expected: &pb.User{
				Email:    "empty@example.com",
				Token:    "",
				Username: "",
				Bio:      "",
				Image:    "",
			},
		},
		{
			name: "User with Special Characters",
			user: &User{
				Model:    gorm.Model{ID: 3},
				Username: "user_😊",
				Email:    "special@例子.com",
				Password: "p@ssw0rd!",
				Bio:      "I ♥ coding!",
				Image:    "http://example.com/image_🖼️.jpg",
			},
			token: "token_with_special_chars!@#$%^&*()",
			expected: &pb.User{
				Email:    "special@例子.com",
				Token:    "token_with_special_chars!@#$%^&*()",
				Username: "user_😊",
				Bio:      "I ♥ coding!",
				Image:    "http://example.com/image_🖼️.jpg",
			},
		},
		{
			name: "Maximum Length Fields",
			user: &User{
				Model:    gorm.Model{ID: 4},
				Username: "a_very_long_username_that_reaches_the_maximum_allowed_length",
				Email:    "a_very_long_email_address_that_reaches_the_maximum_allowed_length@example.com",
				Password: "a_very_long_password_that_reaches_the_maximum_allowed_length",
				Bio:      "A very long bio that reaches the maximum allowed length. It contains multiple sentences to ensure it's really long.",
				Image:    "http://example.com/a_very_long_image_url_that_reaches_the_maximum_allowed_length.jpg",
			},
			token: "a_very_long_token_that_reaches_the_maximum_allowed_length_for_tokens_in_the_system",
			expected: &pb.User{
				Email:    "a_very_long_email_address_that_reaches_the_maximum_allowed_length@example.com",
				Token:    "a_very_long_token_that_reaches_the_maximum_allowed_length_for_tokens_in_the_system",
				Username: "a_very_long_username_that_reaches_the_maximum_allowed_length",
				Bio:      "A very long bio that reaches the maximum allowed length. It contains multiple sentences to ensure it's really long.",
				Image:    "http://example.com/a_very_long_image_url_that_reaches_the_maximum_allowed_length.jpg",
			},
		},
		{
			name:     "Nil User Pointer",
			user:     nil,
			token:    "sometoken",
			expected: nil,
		},
		{
			name: "User with Nested Relationships",
			user: &User{
				Model:    gorm.Model{ID: 5},
				Username: "userWithRelations",
				Email:    "relations@example.com",
				Password: "password",
				Bio:      "User with relations",
				Image:    "http://example.com/relations.jpg",
				Follows: []User{
					{Model: gorm.Model{ID: 6}, Username: "follower1"},
					{Model: gorm.Model{ID: 7}, Username: "follower2"},
				},
				FavoriteArticles: []Article{
					{Model: gorm.Model{ID: 1}, Title: "Article 1"},
					{Model: gorm.Model{ID: 2}, Title: "Article 2"},
				},
			},
			token: "relationstoken",
			expected: &pb.User{
				Email:    "relations@example.com",
				Token:    "relationstoken",
				Username: "userWithRelations",
				Bio:      "User with relations",
				Image:    "http://example.com/relations.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.user == nil {
				assert.Nil(t, tt.user.ProtoUser(tt.token))
			} else {
				result := tt.user.ProtoUser(tt.token)
				assert.Equal(t, tt.expected, result)
			}
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
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "Successfully Hash a Valid Password",
			password: "validPassword123",
			wantErr:  nil,
		},
		{
			name:     "Attempt to Hash an Empty Password",
			password: "",
			wantErr:  errors.New("password should not be empty"),
		},
		{
			name:     "Hash a Very Long Password",
			password: string(make([]byte, 1000)),
			wantErr:  nil,
		},
		{
			name:     "Test Password Hashing with Special Characters",
			password: "P@ssw0rd!@#$%^&*()_+",
			wantErr:  nil,
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

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if u.Password == tt.password {
					t.Errorf("HashPassword() did not change the password")
				}

				err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(tt.password))
				if err != nil {
					t.Errorf("HashPassword() produced hash that doesn't match original password: %v", err)
				}
			}
		})
	}

	t.Run("Verify Idempotency of Hashing", func(t *testing.T) {
		u := &User{
			Model:    gorm.Model{},
			Username: "testuser",
			Email:    "test@example.com",
			Password: "testPassword123",
			Bio:      "Test bio",
			Image:    "test.jpg",
		}

		err := u.HashPassword()
		if err != nil {
			t.Errorf("First HashPassword() call failed: %v", err)
		}

		firstHash := u.Password

		err = u.HashPassword()
		if err != nil {
			t.Errorf("Second HashPassword() call failed: %v", err)
		}

		if u.Password != firstHash {
			t.Errorf("HashPassword() changed the hash on second call")
		}
	})

	t.Run("Verify Hash Uniqueness for Different Passwords", func(t *testing.T) {
		u1 := &User{Password: "password1"}
		u2 := &User{Password: "password2"}

		err1 := u1.HashPassword()
		err2 := u2.HashPassword()

		if err1 != nil || err2 != nil {
			t.Errorf("HashPassword() failed: %v, %v", err1, err2)
		}

		if u1.Password == u2.Password {
			t.Errorf("HashPassword() produced the same hash for different passwords")
		}
	})
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
		inputPassword  string
		expected       bool
	}{
		{
			name:           "Correct Password Verification",
			hashedPassword: "$2a$10$1234567890123456789012abcdefghijklmnopqrstuvwxyz012345",
			inputPassword:  "correctPassword",
			expected:       true,
		},
		{
			name:           "Incorrect Password Rejection",
			hashedPassword: "$2a$10$1234567890123456789012abcdefghijklmnopqrstuvwxyz012345",
			inputPassword:  "wrongPassword",
			expected:       false,
		},
		{
			name:           "Empty Password Handling",
			hashedPassword: "$2a$10$1234567890123456789012abcdefghijklmnopqrstuvwxyz012345",
			inputPassword:  "",
			expected:       false,
		},
		{
			name:           "Hashed Password Comparison",
			hashedPassword: "$2a$10$1234567890123456789012abcdefghijklmnopqrstuvwxyz012345",
			inputPassword:  "correctPassword",
			expected:       true,
		},
		{
			name:           "Case Sensitivity Check",
			hashedPassword: "$2a$10$1234567890123456789012abcdefghijklmnopqrstuvwxyz012345",
			inputPassword:  "CoRrEcTpAsSwOrD",
			expected:       false,
		},
		{
			name:           "Long Password Handling",
			hashedPassword: "$2a$10$1234567890123456789012abcdefghijklmnopqrstuvwxyz012345",
			inputPassword:  "veryLongPasswordThatIsOneHundredCharactersLongToTestTheBehaviorOfTheCheckPasswordMethodWithLongInputs",
			expected:       true,
		},
		{
			name:           "Unicode Character Handling",
			hashedPassword: "$2a$10$1234567890123456789012abcdefghijklmnopqrstuvwxyz012345",
			inputPassword:  "пароль123",
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			u := &User{
				Password: tt.hashedPassword,
			}

			if tt.name == "Hashed Password Comparison" || tt.name == "Long Password Handling" || tt.name == "Unicode Character Handling" {
				hashedBytes, err := bcrypt.GenerateFromPassword([]byte(tt.inputPassword), bcrypt.DefaultCost)
				if err != nil {
					t.Fatalf("Failed to generate hashed password: %v", err)
				}
				u.Password = string(hashedBytes)
			}

			result := u.CheckPassword(tt.inputPassword)

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
			name: "Empty Username",
			user: User{
				Username: "",
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
			name: "Empty Email",
			user: User{
				Username: "validuser",
				Email:    "",
				Password: "validpassword",
			},
			wantErr: true,
			errMsg:  "email: cannot be blank.",
		},
		{
			name: "Invalid Email Format",
			user: User{
				Username: "validuser",
				Email:    "notanemail",
				Password: "validpassword",
			},
			wantErr: true,
			errMsg:  "email: must be a valid email address.",
		},
		{
			name: "Empty Password",
			user: User{
				Username: "validuser",
				Email:    "valid@example.com",
				Password: "",
			},
			wantErr: true,
			errMsg:  "password: cannot be blank.",
		},
		{
			name: "Multiple Validation Errors",
			user: User{
				Username: "",
				Email:    "notanemail",
				Password: "",
			},
			wantErr: true,
			errMsg:  "email: must be a valid email address; password: cannot be blank; username: cannot be blank.",
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

