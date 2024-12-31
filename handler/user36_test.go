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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

/*
ROOST_METHOD_HASH=CurrentUser_e3fa631d55
ROOST_METHOD_SIG_HASH=CurrentUser_29413339e9
*/
func TestCurrentUser(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*store.MockUserStore, *auth.MockAuth)
		expectedUser   *pb.User
		expectedError  error
		expectedLogMsg string
	}{
		{
			name: "Successful retrieval of current user",
			setupMocks: func(us *store.MockUserStore, ma *auth.MockAuth) {
				us.GetByIDFn = func(id uint) (*model.User, error) {
					return &model.User{Username: "testuser"}, nil
				}
				ma.GetUserIDFn = func(ctx context.Context) (uint, error) {
					return 1, nil
				}
				ma.GenerateTokenFn = func(userID uint) (string, error) {
					return "valid_token", nil
				}
			},
			expectedUser: &pb.User{
				Username: "testuser",
				Token:    "valid_token",
			},
			expectedError:  nil,
			expectedLogMsg: "get current user",
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
			expectedLogMsg: "unauthenticated",
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
			expectedLogMsg: "user not found",
		},
		{
			name: "Token generation failure",
			setupMocks: func(us *store.MockUserStore, ma *auth.MockAuth) {
				us.GetByIDFn = func(id uint) (*model.User, error) {
					return &model.User{Username: "testuser"}, nil
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
			expectedLogMsg: "Failed to create token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := &store.MockUserStore{}
			mockAuth := &auth.MockAuth{}

			tt.setupMocks(mockUserStore, mockAuth)

			mockLogger := zerolog.New(zerolog.NewTestWriter(t))

			h := &Handler{
				logger: &mockLogger,
				us:     mockUserStore,
			}

			auth.GetUserID = mockAuth.GetUserIDFn
			auth.GenerateToken = mockAuth.GenerateTokenFn

			resp, err := h.CurrentUser(context.Background(), &pb.Empty{})

			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			} else if tt.expectedError != nil && err == nil {
				t.Errorf("Expected error %v, got nil", tt.expectedError)
			} else if tt.expectedError != nil && err != nil {
				if tt.expectedError.Error() != err.Error() {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
			}

			if tt.expectedUser == nil && resp != nil {
				t.Errorf("Expected nil response, got %v", resp)
			} else if tt.expectedUser != nil && resp == nil {
				t.Errorf("Expected non-nil response, got nil")
			} else if tt.expectedUser != nil && resp != nil {
				if tt.expectedUser.Username != resp.User.Username {
					t.Errorf("Expected username %s, got %s", tt.expectedUser.Username, resp.User.Username)
				}
				if tt.expectedUser.Token != resp.User.Token {
					t.Errorf("Expected token %s, got %s", tt.expectedUser.Token, resp.User.Token)
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
				us.GetByEmailFunc = func(email string) (*model.User, error) {
					return &model.User{
						Email:    "user@example.com",
						Password: "correctpassword",
					}, nil
				}
			},
			input: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "user@example.com",
					Password: "correctpassword",
				},
			},
			expectedOutput: &pb.UserResponse{
				User: &pb.User{
					Email: "user@example.com",
					Token: "mocked-token",
				},
			},
			expectedError: nil,
		},
		{
			name: "Login Attempt with Invalid Email",
			setupMock: func(us *store.MockUserStore) {
				us.GetByEmailFunc = func(email string) (*model.User, error) {
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
				us.GetByEmailFunc = func(email string) (*model.User, error) {
					return &model.User{
						Email:    "user@example.com",
						Password: "correctpassword",
					}, nil
				}
			},
			input: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "user@example.com",
					Password: "wrongpassword",
				},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid email or password"),
		},
		{
			name: "Token Generation Failure",
			setupMock: func(us *store.MockUserStore) {
				us.GetByEmailFunc = func(email string) (*model.User, error) {
					return &model.User{
						Email:    "user@example.com",
						Password: "correctpassword",
					}, nil
				}
			},
			input: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "user@example.com",
					Password: "correctpassword",
				},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.Aborted, "internal server error"),
		},
		{
			name:      "Login with Empty Credentials",
			setupMock: func(us *store.MockUserStore) {},
			input: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "",
					Password: "",
				},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid email or password"),
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

			originalGenerateToken := auth.GenerateToken
			auth.GenerateToken = func(userID uint) (string, error) {
				if tt.name == "Token Generation Failure" {
					return "", errors.New("token generation failed")
				}
				return "mocked-token", nil
			}
			defer func() { auth.GenerateToken = originalGenerateToken }()

			got, err := h.LoginUser(context.Background(), tt.input)

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
					t.Error("expected non-nil output, got nil")
				} else if got.User.Email != tt.expectedOutput.User.Email {
					t.Errorf("expected email %s, got %s", tt.expectedOutput.User.Email, got.User.Email)
				}
			} else if got != nil {
				t.Errorf("expected nil output, got %v", got)
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
		mockUserStore  func() *store.MockUserStore
		mockHashPass   func(*model.User) error
		mockGenToken   func(uint) (string, error)
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
			mockUserStore: func() *store.MockUserStore {
				mock := &store.MockUserStore{}
				mock.CreateFunc = func(user *model.User) error {
					user.ID = 1
					return nil
				}
				return mock
			},
			mockHashPass: func(user *model.User) error {
				return nil
			},
			mockGenToken: func(id uint) (string, error) {
				return "testtoken", nil
			},
			expectedOutput: &pb.UserResponse{
				User: &pb.User{
					Username: "testuser",
					Email:    "test@example.com",
					Token:    "testtoken",
				},
			},
			expectedError: nil,
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
			mockUserStore: func() *store.MockUserStore {
				mock := &store.MockUserStore{}
				mock.CreateFunc = func(user *model.User) error {
					return errors.New("ERROR: duplicate key value violates unique constraint")
				}
				return mock
			},
			mockHashPass: func(user *model.User) error {
				return nil
			},
			mockGenToken: func(id uint) (string, error) {
				return "", nil
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.Canceled, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUS := tt.mockUserStore()

			h := &Handler{
				logger: zerolog.New(zerolog.NewConsoleWriter()),
				us:     mockUS,
			}

			originalHashPassword := model.User.HashPassword
			model.User.HashPassword = tt.mockHashPass
			defer func() { model.User.HashPassword = originalHashPassword }()

			originalGenerateToken := auth.GenerateToken
			auth.GenerateToken = tt.mockGenToken
			defer func() { auth.GenerateToken = originalGenerateToken }()

			result, err := h.CreateUser(context.Background(), tt.input)

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
				} else {
					if result.User.Username != tt.expectedOutput.User.Username {
						t.Errorf("Expected username %s, but got %s", tt.expectedOutput.User.Username, result.User.Username)
					}
					if result.User.Email != tt.expectedOutput.User.Email {
						t.Errorf("Expected email %s, but got %s", tt.expectedOutput.User.Email, result.User.Email)
					}
					if result.User.Token != tt.expectedOutput.User.Token {
						t.Errorf("Expected token %s, but got %s", tt.expectedOutput.User.Token, result.User.Token)
					}
				}
			} else if result != nil {
				t.Error("Expected nil result, but got non-nil")
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
		setupMocks     func(*MockUserStore)
		input          *pb.UpdateUserRequest
		expectedOutput *pb.UserResponse
		expectedError  error
	}{
		{
			name: "Successful User Update with All Fields",
			setupMocks: func(mockUS *MockUserStore) {
				mockUS.On("GetByID", uint(1)).Return(&model.User{Username: "olduser", Email: "old@example.com"}, nil)
				mockUS.On("Update", mock.AnythingOfType("*model.User")).Return(nil)
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
					Token:    "mocked_token",
				},
			},
			expectedError: nil,
		},
		{
			name: "Update User with Empty Request",
			setupMocks: func(mockUS *MockUserStore) {
				mockUS.On("GetByID", uint(1)).Return(&model.User{Username: "user", Email: "user@example.com"}, nil)
			},
			input: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{},
			},
			expectedOutput: &pb.UserResponse{
				User: &pb.User{
					Username: "user",
					Email:    "user@example.com",
					Token:    "mocked_token",
				},
			},
			expectedError: nil,
		},
		{
			name: "Update User with Invalid Authentication",
			setupMocks: func(mockUS *MockUserStore) {
				// No setup needed for this case
			},
			input: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{
					Username: "newuser",
				},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Update Non-existent User",
			setupMocks: func(mockUS *MockUserStore) {
				mockUS.On("GetByID", uint(1)).Return(&model.User{}, errors.New("user not found"))
			},
			input: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{
					Username: "newuser",
				},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.NotFound, "not user found"),
		},
		{
			name: "Update User with Invalid Data",
			setupMocks: func(mockUS *MockUserStore) {
				mockUS.On("GetByID", uint(1)).Return(&model.User{}, nil)
			},
			input: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{
					Email: "invalid-email",
				},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.InvalidArgument, "validation error: invalid email format"),
		},
		{
			name: "Handling Database Update Failure",
			setupMocks: func(mockUS *MockUserStore) {
				mockUS.On("GetByID", uint(1)).Return(&model.User{Username: "user", Email: "user@example.com"}, nil)
				mockUS.On("Update", mock.AnythingOfType("*model.User")).Return(errors.New("database error"))
			},
			input: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{
					Username: "newuser",
				},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.InvalidArgument, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUS := new(MockUserStore)
			tt.setupMocks(mockUS)

			logger := zerolog.New(zerolog.NewTestWriter(t))
			h := &Handler{
				logger: &logger,
				us:     mockUS,
			}

			oldGetUserID := auth.GetUserID
			auth.GetUserID = func(ctx context.Context) (uint, error) {
				if tt.name == "Update User with Invalid Authentication" {
					return 0, errors.New("invalid auth")
				}
				return 1, nil
			}
			defer func() { auth.GetUserID = oldGetUserID }()

			oldGenerateToken := auth.GenerateToken
			auth.GenerateToken = func(userID uint) (string, error) {
				return "mocked_token", nil
			}
			defer func() { auth.GenerateToken = oldGenerateToken }()

			ctx := context.Background()
			response, err := h.UpdateUser(ctx, tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, response)
			}

			mockUS.AssertExpectations(t)
		})
	}
}
