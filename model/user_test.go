package github

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"golang.org/x/crypto/bcrypt"
	"time"
	"github.com/stretchr/testify/assert"
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
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			following: true,
			want: &pb.Profile{
				Username:  "testuser",
				Bio:       "Test bio",
				Image:     "https://example.com/image.jpg",
				Following: true,
			},
		},
		{
			name: "Basic Profile Conversion with Following False",
			user: User{
				Model:    gorm.Model{ID: 2},
				Username: "anotheruser",
				Email:    "another@example.com",
				Password: "password123",
				Bio:      "Another test bio",
				Image:    "https://example.com/another-image.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "anotheruser",
				Bio:       "Another test bio",
				Image:     "https://example.com/another-image.jpg",
				Following: false,
			},
		},
		{
			name: "Profile Conversion with Empty Fields",
			user: User{
				Model:    gorm.Model{ID: 3},
				Username: "",
				Email:    "empty@example.com",
				Password: "emptypassword",
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
			name: "Profile Conversion with Special Characters",
			user: User{
				Model:    gorm.Model{ID: 4},
				Username: "special_user_😊",
				Email:    "special@example.com",
				Password: "specialpassword",
				Bio:      "Bio with 特殊字符 and 🚀",
				Image:    "https://example.com/image_特殊.jpg",
			},
			following: false,
			want: &pb.Profile{
				Username:  "special_user_😊",
				Bio:       "Bio with 特殊字符 and 🚀",
				Image:     "https://example.com/image_特殊.jpg",
				Following: false,
			},
		},
		{
			name: "Profile Conversion with Maximum Length Fields",
			user: User{
				Model:    gorm.Model{ID: 5},
				Username: "maxlengthusername1234567890",
				Email:    "maxlength@example.com",
				Password: "maxlengthpassword",
				Bio:      "This is a very long bio that reaches the maximum allowed length for testing purposes. It should be preserved entirely in the resulting profile without any truncation.",
				Image:    "https://example.com/very-long-image-url-that-reaches-maximum-length-for-testing-purposes.jpg",
			},
			following: true,
			want: &pb.Profile{
				Username:  "maxlengthusername1234567890",
				Bio:       "This is a very long bio that reaches the maximum allowed length for testing purposes. It should be preserved entirely in the resulting profile without any truncation.",
				Image:     "https://example.com/very-long-image-url-that-reaches-maximum-length-for-testing-purposes.jpg",
				Following: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.ProtoProfile(tt.following)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.ProtoProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserProtoProfileNil(t *testing.T) {
	var u *User = nil
	got := u.ProtoProfile(true)
	if got != nil {
		t.Errorf("Expected nil profile for nil User, got %v", got)
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
			name: "Successful Conversion",
			user: &User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			token: "sample_token",
			expected: &pb.User{
				Email:    "test@example.com",
				Token:    "sample_token",
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
		},
		{
			name: "Empty Fields",
			user: &User{
				Model:    gorm.Model{ID: 2},
				Username: "",
				Email:    "",
				Password: "",
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
			name: "Long String Values",
			user: &User{
				Model:    gorm.Model{ID: 3},
				Username: string(make([]byte, 1000)),
				Email:    string(make([]byte, 1000)) + "@example.com",
				Password: "password",
				Bio:      string(make([]byte, 1000)),
				Image:    "https://example.com/" + string(make([]byte, 1000)),
			},
			token: string(make([]byte, 1000)),
			expected: &pb.User{
				Email:    string(make([]byte, 1000)) + "@example.com",
				Token:    string(make([]byte, 1000)),
				Username: string(make([]byte, 1000)),
				Bio:      string(make([]byte, 1000)),
				Image:    "https://example.com/" + string(make([]byte, 1000)),
			},
		},
		{
			name: "Special Characters",
			user: &User{
				Model:    gorm.Model{ID: 4},
				Username: "user🚀",
				Email:    "special@éxample.com",
				Password: "password",
				Bio:      "Bio with ñ and ü",
				Image:    "https://example.com/image_☺.jpg",
			},
			token: "token_✨",
			expected: &pb.User{
				Email:    "special@éxample.com",
				Token:    "token_✨",
				Username: "user🚀",
				Bio:      "Bio with ñ and ü",
				Image:    "https://example.com/image_☺.jpg",
			},
		},
		{
			name:     "Nil User Pointer",
			user:     nil,
			token:    "token",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.user == nil {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic for nil User pointer, but no panic occurred")
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
ROOST_METHOD_HASH=HashPassword_ea0347143c
ROOST_METHOD_SIG_HASH=HashPassword_fc69fabec5

FUNCTION_DEF=func (u *User) HashPassword() error 

 */
func TestUserHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "Successfully Hash a Valid Password",
			password: "validPassword123",
			wantErr:  false,
		},
		{
			name:     "Attempt to Hash an Empty Password",
			password: "",
			wantErr:  true,
			errMsg:   "password should not be empty",
		},
		{
			name:     "Hash a Very Long Password",
			password: string(make([]byte, 1000)),
			wantErr:  false,
		},
		{
			name:     "Hash a Password with Special Characters",
			password: "P@ssw0rd!@#$%^&*()",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Password: tt.password,
			}

			err := u.HashPassword()

			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err.Error() != tt.errMsg {
					t.Errorf("HashPassword() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {

				if u.Password == tt.password {
					t.Errorf("HashPassword() did not change the password")
				}

				err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(tt.password))
				if err != nil {
					t.Errorf("HashPassword() produced invalid hash: %v", err)
				}
			}
		})
	}

	t.Run("Verify Consistent Hashing with Same Input", func(t *testing.T) {
		password := "samePassword123"
		u1 := &User{Password: password}
		u2 := &User{Password: password}

		err1 := u1.HashPassword()
		err2 := u2.HashPassword()

		if err1 != nil || err2 != nil {
			t.Errorf("HashPassword() unexpected error: %v, %v", err1, err2)
			return
		}

		if u1.Password == u2.Password {
			t.Errorf("HashPassword() produced same hash for same input")
		}
	})

	t.Run("Attempt to Re-hash an Already Hashed Password", func(t *testing.T) {
		u := &User{Password: "originalPassword"}

		err := u.HashPassword()
		if err != nil {
			t.Errorf("First HashPassword() unexpected error: %v", err)
			return
		}

		firstHash := u.Password

		err = u.HashPassword()
		if err != nil {
			t.Errorf("Second HashPassword() unexpected error: %v", err)
			return
		}

		if u.Password == firstHash {
			t.Errorf("Re-hashing did not change the password hash")
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
		plainPassword  string
		expected       bool
	}{
		{
			name:           "Correct Password Verification",
			hashedPassword: "$2a$10$1234567890123456789012uQOHhMGOXyGzsV7QR2z8k3Uf/Ij.qK6",
			plainPassword:  "correctpassword",
			expected:       true,
		},
		{
			name:           "Incorrect Password Rejection",
			hashedPassword: "$2a$10$1234567890123456789012uQOHhMGOXyGzsV7QR2z8k3Uf/Ij.qK6",
			plainPassword:  "wrongpassword",
			expected:       false,
		},
		{
			name:           "Empty Password Handling",
			hashedPassword: "$2a$10$1234567890123456789012uQOHhMGOXyGzsV7QR2z8k3Uf/Ij.qK6",
			plainPassword:  "",
			expected:       false,
		},
		{
			name:           "Hashed Password Comparison",
			hashedPassword: "$2a$10$1234567890123456789012uQOHhMGOXyGzsV7QR2z8k3Uf/Ij.qK6",
			plainPassword:  "correctpassword",
			expected:       true,
		},
		{
			name:           "Case Sensitivity Check",
			hashedPassword: "$2a$10$1234567890123456789012uQOHhMGOXyGzsV7QR2z8k3Uf/Ij.qK6",
			plainPassword:  "CorrectPassword",
			expected:       false,
		},
		{
			name:           "Long Password Handling",
			hashedPassword: "$2a$10$1234567890123456789012uQOHhMGOXyGzsV7QR2z8k3Uf/Ij.qK6",
			plainPassword:  "thisisaverylongpasswordwithmorethanonehundredcharacterstoensurethatthelongpasswordhandlingworksasexpected",
			expected:       false,
		},
		{
			name:           "Unicode Character Handling",
			hashedPassword: "$2a$10$1234567890123456789012uQOHhMGOXyGzsV7QR2z8k3Uf/Ij.qK6",
			plainPassword:  "пароль123",
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

func generateHashedPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
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
				Model:    gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Username: "validuser",
				Email:    "valid@email.com",
				Password: "validpassword",
			},
			wantErr: false,
		},
		{
			name: "Missing Username",
			user: User{
				Model:    gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Email:    "valid@email.com",
				Password: "validpassword",
			},
			wantErr: true,
			errMsg:  "Username: cannot be blank.",
		},
		{
			name: "Invalid Username Format",
			user: User{
				Model:    gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Username: "invalid@user",
				Email:    "valid@email.com",
				Password: "validpassword",
			},
			wantErr: true,
			errMsg:  "Username: must be in a valid format.",
		},
		{
			name: "Missing Email",
			user: User{
				Model:    gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Username: "validuser",
				Password: "validpassword",
			},
			wantErr: true,
			errMsg:  "Email: cannot be blank.",
		},
		{
			name: "Invalid Email Format",
			user: User{
				Model:    gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Username: "validuser",
				Email:    "notanemail",
				Password: "validpassword",
			},
			wantErr: true,
			errMsg:  "Email: must be a valid email address.",
		},
		{
			name: "Missing Password",
			user: User{
				Model:    gorm.Model{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Username: "validuser",
				Email:    "valid@email.com",
			},
			wantErr: true,
			errMsg:  "Password: cannot be blank.",
		},
		{
			name: "All Fields Invalid",
			user: User{
				Model: gorm.Model{ID: 7, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			},
			wantErr: true,
			errMsg:  "Email: cannot be blank; Password: cannot be blank; Username: cannot be blank.",
		},
		{
			name: "Valid Username Edge Case",
			user: User{
				Model:    gorm.Model{ID: 8, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Username: "User123",
				Email:    "valid@email.com",
				Password: "validpassword",
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

