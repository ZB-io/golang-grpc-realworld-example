package model

import (
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"reflect"
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
			name: "Successful Profile Conversion with Following True",
			user: User{
				Model:    gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
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
			name: "Successful Profile Conversion with Following False",
			user: User{
				Model:    gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
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
			name: "Profile Conversion with Empty User Fields",
			user: User{
				Model:    gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
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
			name: "Profile Conversion with Maximum Length Strings",
			user: User{
				Model:    gorm.Model{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Username: string(make([]byte, 1000)),
				Email:    "long@example.com",
				Password: "longpassword",
				Bio:      string(make([]byte, 1000)),
				Image:    string(make([]byte, 1000)),
			},
			following: false,
			want: &pb.Profile{
				Username:  string(make([]byte, 1000)),
				Bio:       string(make([]byte, 1000)),
				Image:     string(make([]byte, 1000)),
				Following: false,
			},
		},
		{
			name: "Profile Conversion with Special Characters",
			user: User{
				Model:    gorm.Model{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Username: "special_user_😊",
				Email:    "special@example.com",
				Password: "specialpassword",
				Bio:      "Bio with 特殊字符 and 🚀",
				Image:    "https://example.com/image_特殊.jpg",
			},
			following: true,
			want: &pb.Profile{
				Username:  "special_user_😊",
				Bio:       "Bio with 特殊字符 and 🚀",
				Image:     "https://example.com/image_特殊.jpg",
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

func TestUserProtoProfileMultipleCalls(t *testing.T) {
	user := User{
		Model:    gorm.Model{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Username: "multicalluser",
		Email:    "multicall@example.com",
		Password: "multicallpassword",
		Bio:      "Multi-call test bio",
		Image:    "https://example.com/multicall-image.jpg",
	}

	profile1 := user.ProtoProfile(true)
	assert.Equal(t, &pb.Profile{
		Username:  "multicalluser",
		Bio:       "Multi-call test bio",
		Image:     "https://example.com/multicall-image.jpg",
		Following: true,
	}, profile1)

	profile2 := user.ProtoProfile(false)
	assert.Equal(t, &pb.Profile{
		Username:  "multicalluser",
		Bio:       "Multi-call test bio",
		Image:     "https://example.com/multicall-image.jpg",
		Following: false,
	}, profile2)

	profile3 := user.ProtoProfile(true)
	assert.Equal(t, &pb.Profile{
		Username:  "multicalluser",
		Bio:       "Multi-call test bio",
		Image:     "https://example.com/multicall-image.jpg",
		Following: true,
	}, profile3)
}


/*
ROOST_METHOD_HASH=ProtoUser_440c1b101c
ROOST_METHOD_SIG_HASH=ProtoUser_fb8c4736ee

FUNCTION_DEF=func (u *User) ProtoUser(token string) *pb.User 

 */
func TestUserProtoUser(t *testing.T) {
	tests := []struct {
		name  string
		user  *User
		token string
		want  *pb.User
	}{
		{
			name: "Successful Conversion",
			user: &User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
			},
			token: "testtoken",
			want: &pb.User{
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
				Model: gorm.Model{
					ID:        2,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "emptyuser",
				Email:    "empty@example.com",
				Password: "password",
				Bio:      "",
				Image:    "",
			},
			token: "emptytoken",
			want: &pb.User{
				Email:    "empty@example.com",
				Token:    "emptytoken",
				Username: "emptyuser",
				Bio:      "",
				Image:    "",
			},
		},
		{
			name: "Token Handling",
			user: &User{
				Model: gorm.Model{
					ID:        3,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "tokenuser",
				Email:    "token@example.com",
				Password: "password",
				Bio:      "Token bio",
				Image:    "http://example.com/token.jpg",
			},
			token: "uniquetoken123",
			want: &pb.User{
				Email:    "token@example.com",
				Token:    "uniquetoken123",
				Username: "tokenuser",
				Bio:      "Token bio",
				Image:    "http://example.com/token.jpg",
			},
		},
		{
			name:  "Null User Handling",
			user:  nil,
			token: "nulltoken",
			want:  nil,
		},
		{
			name: "Long String Handling",
			user: &User{
				Model: gorm.Model{
					ID:        4,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: string(make([]byte, 1000)),
				Email:    string(make([]byte, 1000)) + "@example.com",
				Password: "password",
				Bio:      string(make([]byte, 1000)),
				Image:    "http://example.com/" + string(make([]byte, 1000)),
			},
			token: "longtoken",
			want: &pb.User{
				Email:    string(make([]byte, 1000)) + "@example.com",
				Token:    "longtoken",
				Username: string(make([]byte, 1000)),
				Bio:      string(make([]byte, 1000)),
				Image:    "http://example.com/" + string(make([]byte, 1000)),
			},
		},
		{
			name: "Special Character Handling",
			user: &User{
				Model: gorm.Model{
					ID:        5,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "special_user_😊",
				Email:    "special@example.com",
				Password: "password",
				Bio:      "Bio with 特殊字符 and emojis 🚀",
				Image:    "http://example.com/special_image_☺️.jpg",
			},
			token: "special_token_🔑",
			want: &pb.User{
				Email:    "special@example.com",
				Token:    "special_token_🔑",
				Username: "special_user_😊",
				Bio:      "Bio with 特殊字符 and emojis 🚀",
				Image:    "http://example.com/special_image_☺️.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.user == nil {

				if got := tt.user.ProtoUser(tt.token); got != tt.want {
					t.Errorf("ProtoUser() = %v, want %v", got, tt.want)
				}
				return
			}

			got := tt.user.ProtoUser(tt.token)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProtoUser() = %v, want %v", got, tt.want)
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

	createHashedPassword := func(password string) string {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return string(hashedPassword)
	}

	tests := []struct {
		name           string
		hashedPassword string
		inputPassword  string
		expected       bool
	}{
		{
			name:           "Correct Password Validation",
			hashedPassword: createHashedPassword("correctPassword"),
			inputPassword:  "correctPassword",
			expected:       true,
		},
		{
			name:           "Incorrect Password Rejection",
			hashedPassword: createHashedPassword("correctPassword"),
			inputPassword:  "wrongPassword",
			expected:       false,
		},
		{
			name:           "Empty Password Handling",
			hashedPassword: createHashedPassword("somePassword"),
			inputPassword:  "",
			expected:       false,
		},
		{
			name:           "Null Byte in Password",
			hashedPassword: createHashedPassword("normalPassword"),
			inputPassword:  "normal\x00Password",
			expected:       false,
		},
		{
			name:           "Very Long Password Handling",
			hashedPassword: createHashedPassword("normalPassword"),
			inputPassword:  string(make([]byte, 1024*1024)),
			expected:       false,
		},
		{
			name:           "Unicode Password Validation",
			hashedPassword: createHashedPassword("パスワード123"),
			inputPassword:  "パスワード123",
			expected:       true,
		},
		{
			name:           "Case Sensitivity Check",
			hashedPassword: createHashedPassword("CaseSensitivePassword"),
			inputPassword:  "casesensitivepassword",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Password: tt.hashedPassword,
			}
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

FUNCTION_DEF=func (u *User) HashPassword() error 

 */
func TestUserHashPassword(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		expectedErrMsg string
	}{
		{
			name:           "Successfully Hash a Valid Password",
			password:       "validPassword123",
			expectedErrMsg: "",
		},
		{
			name:           "Attempt to Hash an Empty Password",
			password:       "",
			expectedErrMsg: "password should not be empty",
		},
		{
			name:           "Hash a Very Long Password",
			password:       string(make([]byte, 1000)),
			expectedErrMsg: "",
		},
		{
			name:           "Hash a Password with Special Characters",
			password:       "P@ssw0rd!@#$%^&*()_+",
			expectedErrMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Model:    gorm.Model{},
				Password: tt.password,
			}

			err := user.HashPassword()

			if tt.expectedErrMsg != "" {
				if err == nil {
					t.Errorf("Expected error with message '%s', but got nil", tt.expectedErrMsg)
				} else if err.Error() != tt.expectedErrMsg {
					t.Errorf("Expected error message '%s', but got '%s'", tt.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}

				if user.Password == tt.password {
					t.Errorf("Password was not hashed")
				}

				err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tt.password))
				if err != nil {
					t.Errorf("Hashed password does not match original: %v", err)
				}
			}
		})
	}

	t.Run("Consistency of Hashing", func(t *testing.T) {
		password := "samePassword123"
		user1 := &User{Password: password}
		user2 := &User{Password: password}

		_ = user1.HashPassword()
		_ = user2.HashPassword()

		if user1.Password == user2.Password {
			t.Errorf("Hashed passwords should be different due to salt")
		}
	})

	t.Run("Password Length After Hashing", func(t *testing.T) {
		passwords := []string{"short", "medium_length_password", string(make([]byte, 100))}
		var hashedLengths []int

		for _, pass := range passwords {
			user := &User{Password: pass}
			_ = user.HashPassword()
			hashedLengths = append(hashedLengths, len(user.Password))
		}

		for i := 1; i < len(hashedLengths); i++ {
			if hashedLengths[i] != hashedLengths[0] {
				t.Errorf("Hashed password lengths are not consistent")
				break
			}
		}
	})
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
				Email:    "valid@email.com",
				Password: "validpassword",
			},
			wantErr: false,
		},
		{
			name: "Missing Username",
			user: User{
				Email:    "valid@email.com",
				Password: "validpassword",
			},
			wantErr: true,
			errMsg:  "Username: cannot be blank.",
		},
		{
			name: "Invalid Username Format",
			user: User{
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
				Username: "validuser",
				Password: "validpassword",
			},
			wantErr: true,
			errMsg:  "Email: cannot be blank.",
		},
		{
			name: "Invalid Email Format",
			user: User{
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
				Username: "validuser",
				Email:    "valid@email.com",
			},
			wantErr: true,
			errMsg:  "Password: cannot be blank.",
		},
		{
			name:    "Multiple Validation Errors",
			user:    User{},
			wantErr: true,
			errMsg:  "Username: cannot be blank; Email: cannot be blank; Password: cannot be blank.",
		},
		{
			name: "Valid Username Edge Case - Very Long",
			user: User{
				Username: "a" + string(make([]byte, 99)),
				Email:    "valid@email.com",
				Password: "validpassword",
			},
			wantErr: false,
		},
		{
			name: "Valid Username Edge Case - Numbers Only",
			user: User{
				Username: "123456",
				Email:    "valid@email.com",
				Password: "validpassword",
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
				if err == nil {
					t.Errorf("User.Validate() expected error, got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("User.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

