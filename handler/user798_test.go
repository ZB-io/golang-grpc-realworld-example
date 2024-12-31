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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockUserStore struct {
	mock.Mock
}

type MockAuth struct {
	mock.Mock
}

/*
ROOST_METHOD_HASH=CurrentUser_e3fa631d55
ROOST_METHOD_SIG_HASH=CurrentUser_29413339e9
*/
func TestCurrentUser(t *testing.T) {
	tests := []struct {
		name            string
		setupMocks      func(*MockUserStore, *MockAuth)
		expectedUser    *pb.User
		expectedErrCode codes.Code
	}{
		{
			name: "Successful retrieval of current user",
			setupMocks: func(us *MockUserStore, ma *MockAuth) {
				ma.On("GetUserID", mock.Anything).Return(uint(123), nil)
				us.On("GetByID", uint(123)).Return(&model.User{ID: 123, Username: "testuser"}, nil)
				ma.On("GenerateToken", uint(123)).Return("valid_token", nil)
			},
			expectedUser:    &pb.User{Username: "testuser", Token: "valid_token"},
			expectedErrCode: codes.OK,
		},
		{
			name: "Unauthenticated user attempt",
			setupMocks: func(us *MockUserStore, ma *MockAuth) {
				ma.On("GetUserID", mock.Anything).Return(uint(0), errors.New("unauthenticated"))
			},
			expectedUser:    nil,
			expectedErrCode: codes.Unauthenticated,
		},
		{
			name: "Valid token but non-existent user",
			setupMocks: func(us *MockUserStore, ma *MockAuth) {
				ma.On("GetUserID", mock.Anything).Return(uint(123), nil)
				us.On("GetByID", uint(123)).Return(nil, errors.New("user not found"))
			},
			expectedUser:    nil,
			expectedErrCode: codes.NotFound,
		},
		{
			name: "Database error during user retrieval",
			setupMocks: func(us *MockUserStore, ma *MockAuth) {
				ma.On("GetUserID", mock.Anything).Return(uint(123), nil)
				us.On("GetByID", uint(123)).Return(nil, errors.New("database error"))
			},
			expectedUser:    nil,
			expectedErrCode: codes.NotFound,
		},
		{
			name: "Token generation failure",
			setupMocks: func(us *MockUserStore, ma *MockAuth) {
				ma.On("GetUserID", mock.Anything).Return(uint(123), nil)
				us.On("GetByID", uint(123)).Return(&model.User{ID: 123, Username: "testuser"}, nil)
				ma.On("GenerateToken", uint(123)).Return("", errors.New("token generation failed"))
			},
			expectedUser:    nil,
			expectedErrCode: codes.Aborted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := new(MockUserStore)
			mockAuth := new(MockAuth)
			tt.setupMocks(mockUserStore, mockAuth)

			h := &Handler{
				logger: zerolog.New(zerolog.NewTestWriter(t)),
				us:     mockUserStore,
			}

			ctx := context.Background()
			req := &pb.Empty{}

			ctx = context.WithValue(ctx, auth.AuthKey, mockAuth)

			resp, err := h.CurrentUser(ctx, req)

			if tt.expectedErrCode != codes.OK {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok, "Expected gRPC error")
				assert.Equal(t, tt.expectedErrCode, st.Code(), "Expected error code %v, got %v", tt.expectedErrCode, st.Code())
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedUser != nil {
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.User)
				assert.Equal(t, tt.expectedUser.Username, resp.User.Username)
				assert.Equal(t, tt.expectedUser.Token, resp.User.Token)
			} else {
				assert.Nil(t, resp)
			}

			mockUserStore.AssertExpectations(t)
			mockAuth.AssertExpectations(t)
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
		setupMock      func(*MockUserStore)
		input          *pb.LoginUserRequest
		expectedOutput *pb.UserResponse
		expectedError  error
	}{
		{
			name: "Successful User Login",
			setupMock: func(us *MockUserStore) {
				us.On("GetByEmail", "test@example.com").Return(&model.User{
					Email:    "test@example.com",
					Password: "$2a$10$abcdefghijklmnopqrstuvwxyz123456",
				}, nil)
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
					Token: "mockedToken",
				},
			},
			expectedError: nil,
		},
		{
			name: "Login Attempt with Invalid Email",
			setupMock: func(us *MockUserStore) {
				us.On("GetByEmail", "nonexistent@example.com").Return(nil, errors.New("user not found"))
			},
			input: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "nonexistent@example.com",
					Password: "password",
				},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid email or password"),
		},
		{
			name: "Login Attempt with Incorrect Password",
			setupMock: func(us *MockUserStore) {
				us.On("GetByEmail", "test@example.com").Return(&model.User{
					Email:    "test@example.com",
					Password: "$2a$10$abcdefghijklmnopqrstuvwxyz123456",
				}, nil)
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
			setupMock: func(us *MockUserStore) {
				us.On("GetByEmail", "test@example.com").Return(&model.User{
					Email:    "test@example.com",
					Password: "$2a$10$abcdefghijklmnopqrstuvwxyz123456",
				}, nil)
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
		{
			name: "Login with Empty Credentials",
			setupMock: func(us *MockUserStore) {},
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
			mockUserStore := new(MockUserStore)
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
				return "mockedToken", nil
			}

			result, err := h.LoginUser(context.Background(), tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedOutput != nil {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedOutput.User.Email, result.User.Email)
			} else {
				assert.Nil(t, result)
			}

			mockUserStore.AssertExpectations(t)
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
		setupMock      func(*MockUserStore)
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
			setupMock: func(us *MockUserStore) {
				us.On("Create", mock.AnythingOfType("*model.User")).Return(nil)
			},
			expectedOutput: &pb.UserResponse{
				User: &pb.User{
					Username: "testuser",
					Email:    "test@example.com",
					Token:    "mocked_token",
				},
			},
			expectedError: nil,
		},
		{
			name: "Attempt to Create User with Invalid Input",
			input: &pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Username: "",
					Email:    "invalid_email",
					Password: "short",
				},
			},
			setupMock: func(us *MockUserStore) {},
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
			setupMock: func(us *MockUserStore) {
				us.On("Create", mock.AnythingOfType("*model.User")).Return(errors.New("database error"))
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.Canceled, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := new(MockUserStore)
			tt.setupMock(mockUserStore)

			logger := zerolog.New(zerolog.NewTestWriter(t))
			h := &Handler{
				logger: &logger,
				us:     mockUserStore,
			}

			origGenerateToken := auth.GenerateToken
			defer func() { auth.GenerateToken = origGenerateToken }()
			auth.GenerateToken = func(userID uint) (string, error) {
				return "mocked_token", nil
			}

			output, err := h.CreateUser(context.Background(), tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, output)
			}

			mockUserStore.AssertExpectations(t)
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
		setupMock      func(*MockUserStore, *MockAuth)
		setupContext   func() context.Context
		input          *pb.UpdateUserRequest
		expectedOutput *pb.UserResponse
		expectedError  error
	}{
		{
			name: "Successful User Update",
			setupMock: func(us *MockUserStore, ma *MockAuth) {
				ma.On("GetUserID", mock.Anything).Return(uint(1), nil)
				us.On("GetByID", uint(1)).Return(&model.User{Username: "oldname", Email: "old@email.com"}, nil)
				us.On("Update", mock.AnythingOfType("*model.User")).Return(nil)
			},
			setupContext: func() context.Context {
				return context.Background()
			},
			input: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{
					Username: "newname",
					Email:    "new@email.com",
				},
			},
			expectedOutput: &pb.UserResponse{
				User: &pb.User{
					Username: "newname",
					Email:    "new@email.com",
					Token:    "mock_token",
				},
			},
			expectedError: nil,
		},
		{
			name: "Update with Invalid Email Format",
			setupMock: func(us *MockUserStore, ma *MockAuth) {
				ma.On("GetUserID", mock.Anything).Return(uint(1), nil)
				us.On("GetByID", uint(1)).Return(&model.User{Username: "oldname", Email: "old@email.com"}, nil)
			},
			setupContext: func() context.Context {
				return context.Background()
			},
			input: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{
					Email: "invalid-email",
				},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.InvalidArgument, "validation error"),
		},
		{
			name: "Unauthenticated User Update Attempt",
			setupMock: func(us *MockUserStore, ma *MockAuth) {
				ma.On("GetUserID", mock.Anything).Return(uint(0), errors.New("unauthenticated"))
			},
			setupContext: func() context.Context {
				return context.Background()
			},
			input: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{
					Username: "newname",
				},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "User Not Found During Update",
			setupMock: func(us *MockUserStore, ma *MockAuth) {
				ma.On("GetUserID", mock.Anything).Return(uint(1), nil)
				us.On("GetByID", uint(1)).Return(nil, errors.New("user not found"))
			},
			setupContext: func() context.Context {
				return context.Background()
			},
			input: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{
					Username: "newname",
				},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.NotFound, "not user found"),
		},
		{
			name: "Internal Server Error During Update",
			setupMock: func(us *MockUserStore, ma *MockAuth) {
				ma.On("GetUserID", mock.Anything).Return(uint(1), nil)
				us.On("GetByID", uint(1)).Return(&model.User{Username: "oldname", Email: "old@email.com"}, nil)
				us.On("Update", mock.AnythingOfType("*model.User")).Return(errors.New("database error"))
			},
			setupContext: func() context.Context {
				return context.Background()
			},
			input: &pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{
					Username: "newname",
				},
			},
			expectedOutput: nil,
			expectedError:  status.Error(codes.InvalidArgument, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := new(MockUserStore)
			mockAuth := new(MockAuth)
			tt.setupMock(mockUserStore, mockAuth)

			h := &Handler{
				logger: zerolog.New(zerolog.NewTestWriter(t)),
				us:     mockUserStore,
			}

			ctx := tt.setupContext()
			ctx = context.WithValue(ctx, auth.AuthKey, mockAuth)

			origGenerateToken := auth.GenerateToken
			defer func() { auth.GenerateToken = origGenerateToken }()
			auth.GenerateToken = func(userID uint) (string, error) {
				return "mock_token", nil
			}

			output, err := h.UpdateUser(ctx, tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, tt.expectedOutput.User.Username, output.User.Username)
				assert.Equal(t, tt.expectedOutput.User.Email, output.User.Email)
				assert.Equal(t, tt.expectedOutput.User.Token, output.User.Token)
			}

			mockUserStore.AssertExpectations(t)
			mockAuth.AssertExpectations(t)
		})
	}
}

func (m *MockUserStore) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserStore) GetByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuth) GetUserID(ctx context.Context) (uint, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockAuth) GenerateToken(userID uint) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}
