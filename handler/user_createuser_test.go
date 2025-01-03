package handler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandlerCreateUser(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(mockUserStore *store.MockUserStore, mockAuth *auth.MockAuth)
		input         *proto.CreateUserRequest
		expectedError error
		expectedUser  *proto.UserResponse
	}{
		{
			name: "Successful User Creation with Valid Input",
			setupMocks: func(mockUserStore *store.MockUserStore, mockAuth *auth.MockAuth) {
				mockUserStore.EXPECT().Create(gomock.Any()).Return(nil)
				mockAuth.EXPECT().GenerateToken(gomock.Any()).Return("sampleToken", nil)
			},
			input: &proto.CreateUserRequest{
				User: &proto.CreateUserRequest_User{
					Username: "validUser",
					Email:    "user@example.com",
					Password: "securepassword",
				},
			},
			expectedError: nil,
			expectedUser: &proto.UserResponse{
				User: &proto.User{
					Email:    "user@example.com",
					Token:    "sampleToken",
					Username: "validUser",
				},
			},
		},
		{
			name: "Validation Error on Invalid User Data",
			setupMocks: func(mockUserStore *store.MockUserStore, mockAuth *auth.MockAuth) {
				// No user store or auth calls expected due to validation error
			},
			input: &proto.CreateUserRequest{
				User: &proto.CreateUserRequest_User{
					Username: "",
					Email:    "invalidemail",
					Password: "",
				},
			},
			expectedError: status.Error(codes.InvalidArgument, "validation error"),
			expectedUser:  nil,
		},
		{
			name: "Password Hashing Failure",
			setupMocks: func(mockUserStore *store.MockUserStore, mockAuth *auth.MockAuth) {
				originalHashPassword := model.HashPassword
				model.HashPassword = func(u *model.User) error {
					return errors.New("hashing error")
				}
				defer func() { model.HashPassword = originalHashPassword }()
			},
			input: &proto.CreateUserRequest{
				User: &proto.CreateUserRequest_User{
					Username: "user",
					Email:    "user@domain.com",
					Password: "password",
				},
			},
			expectedError: status.Error(codes.Aborted, "internal server error"),
			expectedUser:  nil,
		},
		{
			name: "User Creation Failure in Store",
			setupMocks: func(mockUserStore *store.MockUserStore, mockAuth *auth.MockAuth) {
				mockUserStore.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("DB Error"))
			},
			input: &proto.CreateUserRequest{
				User: &proto.CreateUserRequest_User{
					Username: "user",
					Email:    "user@example.com",
					Password: "password",
				},
			},
			expectedError: status.Error(codes.Canceled, "internal server error"),
			expectedUser:  nil,
		},
		{
			name: "Token Generation Failure",
			setupMocks: func(mockUserStore *store.MockUserStore, mockAuth *auth.MockAuth) {
				mockUserStore.EXPECT().Create(gomock.Any()).Return(nil)
				mockAuth.EXPECT().GenerateToken(gomock.Any()).Return("", errors.New("token error"))
			},
			input: &proto.CreateUserRequest{
				User: &proto.CreateUserRequest_User{
					Username: "user",
					Email:    "user@example.com",
					Password: "password",
				},
			},
			expectedError: status.Error(codes.Aborted, "internal server error"),
			expectedUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserStore := store.NewMockUserStore(ctrl)
			mockAuth := auth.NewMockAuth(ctrl)
			logger := zerolog.New(os.Stdout)

			handler := &Handler{
				logger: &logger,
				us:     mockUserStore,
			}
			tt.setupMocks(mockUserStore, mockAuth)

			ctx := context.TODO()
			userResponse, err := handler.CreateUser(ctx, tt.input)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.expectedUser, userResponse)
		})
	}
}
