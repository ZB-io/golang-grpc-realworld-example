package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

/*
ROOST_METHOD_HASH=CurrentUser_e3fa631d55
ROOST_METHOD_SIG_HASH=CurrentUser_29413339e9
*/
func TestCurrentUser(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*store.MockUserStore, *auth.MockAuth)
		expectedUser   *model.User
		expectedError  error
		expectedStatus codes.Code
	}{
		{
			name: "Successful retrieval of current user",
			setupMocks: func(us *store.MockUserStore, ma *auth.MockAuth) {
				us.GetByIDFn = func(id uint) (*model.User, error) {
					return &model.User{ID: 1, Username: "testuser"}, nil
				}
				ma.GetUserIDFn = func(ctx context.Context) (uint, error) {
					return 1, nil
				}
				ma.GenerateTokenFn = func(userID uint) (string, error) {
					return "valid_token", nil
				}
			},
			expectedUser:   &model.User{ID: 1, Username: "testuser"},
			expectedError:  nil,
			expectedStatus: codes.OK,
		},
		{
			name: "Unauthenticated user attempt",
			setupMocks: func(us *store.MockUserStore, ma *auth.MockAuth) {
				ma.GetUserIDFn = func(ctx context.Context) (uint, error) {
					return 0, errors.New("unauthenticated")
				}
			},
			expectedUser:   nil,
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
			expectedStatus: codes.Unauthenticated,
		},
		{
			name: "Valid token but non-existent user",
			setupMocks: func(us *store.MockUserStore, ma *auth.MockAuth) {
				us.GetByIDFn = func(id uint) (*model.User, error) {
					return nil, errors.New("user not found")
				}
				ma.GetUserIDFn = func(ctx context.Context) (uint, error) {
					return 1, nil
				}
			},
			expectedUser:   nil,
			expectedError:  status.Error(codes.NotFound, "user not found"),
			expectedStatus: codes.NotFound,
		},
		{
			name: "Token generation failure",
			setupMocks: func(us *store.MockUserStore, ma *auth.MockAuth) {
				us.GetByIDFn = func(id uint) (*model.User, error) {
					return &model.User{ID: 1, Username: "testuser"}, nil
				}
				ma.GetUserIDFn = func(ctx context.Context) (uint, error) {
					return 1, nil
				}
				ma.GenerateTokenFn = func(userID uint) (string, error) {
					return "", errors.New("token generation failed")
				}
			},
			expectedUser:   nil,
			expectedError:  status.Error(codes.Aborted, "internal server error"),
			expectedStatus: codes.Aborted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := &store.MockUserStore{}
			mockAuth := &auth.MockAuth{}
			mockLogger := zerolog.New(nil)
			tt.setupMocks(mockUserStore, mockAuth)

			h := &Handler{
				logger: &mockLogger,
				us:     mockUserStore,
			}

			auth.GetUserID = mockAuth.GetUserIDFn
			auth.GenerateToken = mockAuth.GenerateTokenFn

			resp, err := h.CurrentUser(context.Background(), &pb.Empty{})

			if tt.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.expectedError)
				}
				if status.Code(err) != tt.expectedStatus {
					t.Errorf("expected status code %v, got %v", tt.expectedStatus, status.Code(err))
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if resp == nil || resp.User == nil {
					t.Fatal("expected non-nil response and user")
				}
				if resp.User.Username != tt.expectedUser.Username {
					t.Errorf("expected user %v, got %v", tt.expectedUser, resp.User)
				}
			}
		})
	}
}

/*
ROOST_METHOD_HASH=LoginUser_079a321a92
ROOST_METHOD_SIG_HASH=LoginUser_e7df23a6bd
*/
func TestLoginUser(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*store.MockUserStore)
		input          *pb.LoginUserRequest
		expectedOutput *pb.UserResponse
		expectedError  error
	}{
		{
			name: "Successful User Login",
			setupMock: func(us *store.MockUserStore) {
				us.GetByEmailFn = func(email string) (*model.User, error) {
					return &model.User{
						Email:    "test@example.com",
						Password: "$2a$10$abcdefghijklmnopqrstuvwxyz123456",
					}, nil
				}
			},
			input: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "test@example.com",
					Password: "correctpassword",
				},
			},
			expectedOutput: &pb.UserResponse{
				User: &pb.User{
					Email: "test@example.com",
					Token: "mocked_token",
				},
			},
			expectedError: nil,
		},
		{
			name: "Login Attempt with Non-existent Email",
			setupMock: func(us *store.MockUserStore) {
				us.GetByEmailFn = func(email string) (*model.User, error) {
					return nil, errors.New("user not found")
				}
			},
			input: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "nonexistent@example.com",
					Password: "anypassword",
				},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid email or password"),
		},
		{
			name: "Login Attempt with Incorrect Password",
			setupMock: func(us *store.MockUserStore) {
				us.GetByEmailFn = func(email string) (*model.User, error) {
					return &model.User{
						Email:    "test@example.com",
						Password: "$2a$10$abcdefghijklmnopqrstuvwxyz123456",
					}, nil
				}
			},
			input: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "test@example.com",
					Password: "wrongpassword",
				},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid email or password"),
		},
		{
			name: "Token Generation Failure",
			setupMock: func(us *store.MockUserStore) {
				us.GetByEmailFn = func(email string) (*model.User, error) {
					return &model.User{
						Email:    "test@example.com",
						Password: "$2a$10$abcdefghijklmnopqrstuvwxyz123456",
					}, nil
				}
			},
			input: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "test@example.com",
					Password: "correctpassword",
				},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.Aborted, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := &store.MockUserStore{}
			tt.setupMock(mockUserStore)

			h := &Handler{
				logger: zerolog.New(zerolog.NewTestWriter(t)),
				us:     mockUserStore,
			}

			origGenerateToken := auth.GenerateToken
			defer func() { auth.GenerateToken = origGenerateToken }()
			auth.GenerateToken = func(userID uint) (string, error) {
				if tt.name == "Token Generation Failure" {
					return "", errors.New("token generation failed")
				}
				return "mocked_token", nil
			}

			result, err := h.LoginUser(context.Background(), tt.input)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, but got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error %v, but got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.expectedOutput != nil {
				if result == nil {
					t.Error("Expected non-nil result, but got nil")
				} else if result.User.Email != tt.expectedOutput.User.Email {
					t.Errorf("Expected email %s, but got %s", tt.expectedOutput.User.Email, result.User.Email)
				}
			} else if result != nil {
				t.Errorf("Expected nil result, but got %v", result)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=CreateUser_f2f8a1c84a
ROOST_METHOD_SIG_HASH=CreateUser_a3af3934da
*/
func TestCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		input          *pb.CreateUserRequest
		mockUserStore  func() *store.UserStore
		mockAuthFunc   func(uint) (string, error)
		expectedOutput *pb.UserResponse
		expectedError  error
	}{
		{
			name: "Successfully Create a New User",
			input: &pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			mockUserStore: func() *store.UserStore {
				return &store.UserStore{
					CreateFunc: func(u *model.User) error {
						u.ID = 1
						return nil
					},
				}
			},
			mockAuthFunc: func(uint) (string, error) {
				return "valid_token", nil
			},
			expectedOutput: &pb.UserResponse{
				User: &pb.User{
					Username: "testuser",
					Email:    "test@example.com",
					Token:    "valid_token",
				},
			},
			expectedError: nil,
		},
		{
			name: "Attempt to Create User with Invalid Data",
			input: &pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Username: "",
					Email:    "invalid_email",
					Password: "short",
				},
			},
			mockUserStore: func() *store.UserStore {
				return &store.UserStore{}
			},
			mockAuthFunc: func(uint) (string, error) {
				return "", nil
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.InvalidArgument, "validation error"),
		},
		{
			name: "Handle Database Error During User Creation",
			input: &pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			mockUserStore: func() *store.UserStore {
				return &store.UserStore{
					CreateFunc: func(u *model.User) error {
						return errors.New("database error")
					},
				}
			},
			mockAuthFunc: func(uint) (string, error) {
				return "", nil
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.Canceled, "internal server error"),
		},
		{
			name: "Handle Token Generation Failure",
			input: &pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			mockUserStore: func() *store.UserStore {
				return &store.UserStore{
					CreateFunc: func(u *model.User) error {
						u.ID = 1
						return nil
					},
				}
			},
			mockAuthFunc: func(uint) (string, error) {
				return "", errors.New("token generation error")
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.Aborted, "internal server error"),
		},
		{
			name: "Attempt to Create User with Existing Username or Email",
			input: &pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Username: "existinguser",
					Email:    "existing@example.com",
					Password: "password123",
				},
			},
			mockUserStore: func() *store.UserStore {
				return &store.UserStore{
					CreateFunc: func(u *model.User) error {
						return gorm.ErrDuplicatedKey
					},
				}
			},
			mockAuthFunc: func(uint) (string, error) {
				return "", nil
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.Canceled, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zerolog.Nop()
			h := &Handler{
				logger: &logger,
				us:     tt.mockUserStore(),
			}

			origGenerateToken := auth.GenerateToken
			auth.GenerateToken = tt.mockAuthFunc
			defer func() { auth.GenerateToken = origGenerateToken }()

			got, err := h.CreateUser(context.Background(), tt.input)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.expectedOutput != nil {
				if got == nil {
					t.Error("expected non-nil response, got nil")
				} else {
					if got.User.Username != tt.expectedOutput.User.Username {
						t.Errorf("expected username %s, got %s", tt.expectedOutput.User.Username, got.User.Username)
					}
					if got.User.Email != tt.expectedOutput.User.Email {
						t.Errorf("expected email %s, got %s", tt.expectedOutput.User.Email, got.User.Email)
					}
					if got.User.Token != tt.expectedOutput.User.Token {
						t.Errorf("expected token %s, got %s", tt.expectedOutput.User.Token, got.User.Token)
					}
				}
			} else if got != nil {
				t.Errorf("expected nil response, got %v", got)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=UpdateUser_6fa4ecf979
ROOST_METHOD_SIG_HASH=UpdateUser_883937d25b
*/
func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*store.MockUserStore)
		input          *pb.UpdateUserRequest
		expectedOutput *pb.UserResponse
		expectedError  error
	}{
		{
			name: "Successful User Update with All Fields",
			setupMocks: func(us *store.MockUserStore) {
				us.GetByIDFn = func(id uint) (*model.User, error) {
					return &model.User{Username: "olduser", Email: "old@example.com"}, nil
				}
				us.UpdateFn = func(u *model.User) error {
					return nil
				}
			},
			input: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{
					Username: "newuser",
					Email:    "new@example.com",
					Password: "newpassword",
					Image:    "newimage.jpg",
					Bio:      "New bio",
				},
			},
			expectedOutput: &pb.UserResponse{
				User: &pb.User{
					Username: "newuser",
					Email:    "new@example.com",
					Image:    "newimage.jpg",
					Bio:      "New bio",
					Token:    "mocked-token",
				},
			},
			expectedError: nil,
		},
		{
			name: "Update User with Invalid Authentication",
			setupMocks: func(us *store.MockUserStore) {
			},
			input: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Update User with Non-existent User ID",
			setupMocks: func(us *store.MockUserStore) {
				us.GetByIDFn = func(id uint) (*model.User, error) {
					return nil, errors.New("user not found")
				}
			},
			input: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.NotFound, "not user found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := &store.MockUserStore{}
			tt.setupMocks(mockUserStore)

			h := &Handler{
				logger: zerolog.Nop(),
				us:     mockUserStore,
			}

			origGetUserID := auth.GetUserID
			auth.GetUserID = func(ctx context.Context) (uint, error) {
				if tt.expectedError != nil && status.Code(tt.expectedError) == codes.Unauthenticated {
					return 0, errors.New("unauthenticated")
				}
				return 1, nil
			}
			defer func() { auth.GetUserID = origGetUserID }()

			origGenerateToken := auth.GenerateToken
			auth.GenerateToken = func(userID uint) (string, error) {
				return "mocked-token", nil
			}
			defer func() { auth.GenerateToken = origGenerateToken }()

			result, err := h.UpdateUser(context.Background(), tt.input)

			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			} else if tt.expectedError != nil && err == nil {
				t.Errorf("Expected error %v, got nil", tt.expectedError)
			} else if tt.expectedError != nil && err != nil {
				if status.Code(tt.expectedError) != status.Code(err) {
					t.Errorf("Expected error code %v, got %v", status.Code(tt.expectedError), status.Code(err))
				}
			}

			if tt.expectedOutput != nil {
				if result == nil {
					t.Error("Expected non-nil result, got nil")
				} else {
					if result.User.Username != tt.expectedOutput.User.Username {
						t.Errorf("Expected username %s, got %s", tt.expectedOutput.User.Username, result.User.Username)
					}
					if result.User.Email != tt.expectedOutput.User.Email {
						t.Errorf("Expected email %s, got %s", tt.expectedOutput.User.Email, result.User.Email)
					}
					if result.User.Token != tt.expectedOutput.User.Token {
						t.Errorf("Expected token %s, got %s", tt.expectedOutput.User.Token, result.User.Token)
					}
				}
			} else if result != nil {
				t.Error("Expected nil result, got non-nil")
			}
		})
	}
}
