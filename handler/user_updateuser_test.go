package handler

import (
	"context"
	"errors"
	"log"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/stretchr/testify/assert"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUpdateUser(t *testing.T) {
	logger := zerolog.New(log.Writer()).With().Logger()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, %v", err)
	}
	defer db.Close()

	userStore := &store.UserStore{DB: db} // Correct field name to "DB" instead of "db"
	handler := &Handler{
		logger: &logger,
		us:     userStore,
	}

	type testcase struct {
		name          string
		setupContext  func() context.Context
		setupRequest  func() *proto.UpdateUserRequest
		setupMocks    func()
		expectedError error
		expectedCode  codes.Code
	}

	tests := []testcase{
		{
			name: "Successful User Update",
			setupContext: func() context.Context {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 1, nil // Mocking GetUserID to return a valid user ID
				}
				return context.TODO()
			},
			setupRequest: func() *proto.UpdateUserRequest {
				return &proto.UpdateUserRequest{
					User: &proto.UpdateUserRequest_User{
						Username: "newusername",
						Email:    "new@example.com",
						Password: "newpassword",
						Image:    "newimage.png",
						Bio:      "New Bio",
					},
				}
			},
			setupMocks: func() {
				mock.ExpectQuery("SELECT").WillReturnRows(
					sqlmock.NewRows([]string{"id", "username", "email", "password", "image", "bio"}).
						AddRow(1, "oldusername", "old@example.com", "oldpassword", "oldimage.png", "Old Bio"))

				mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))

				auth.GenerateToken = func(id uint) (string, error) {
					return "mock-token", nil // Mock token generation
				}
			},
			expectedError: nil,
			expectedCode:  codes.OK,
		},
		{
			name: "Unauthenticated User Request",
			setupContext: func() context.Context {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 0, errors.New("unauthenticated") // Mocking GetUserID to fail
				}
				return context.TODO()
			},
			setupRequest: func() *proto.UpdateUserRequest {
				return &proto.UpdateUserRequest{}
			},
			setupMocks: func() {},
			expectedError: status.Error(codes.Unauthenticated, "unauthenticated"),
			expectedCode:  codes.Unauthenticated,
		},
		{
			name: "User Not Found in Database",
			setupContext: func() context.Context {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 1, nil
				}
				return context.TODO()
			},
			setupRequest: func() *proto.UpdateUserRequest {
				return &proto.UpdateUserRequest{}
			},
			setupMocks: func() {
				mock.ExpectQuery("SELECT").WillReturnError(errors.New("user not found"))
			},
			expectedError: status.Error(codes.NotFound, "not user found"),
			expectedCode:  codes.NotFound,
		},
		{
			name: "Validation Failure on Updated Data",
			setupContext: func() context.Context {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 1, nil
				}
				return context.TODO()
			},
			setupRequest: func() *proto.UpdateUserRequest {
				return &proto.UpdateUserRequest{
					User: &proto.UpdateUserRequest_User{
						Email: "invalid email",
					},
				}
			},
			setupMocks: func() {
				mock.ExpectQuery("SELECT").WillReturnRows(
					sqlmock.NewRows([]string{"id", "username", "email", "password", "image", "bio"}).
						AddRow(1, "username", "old@example.com", "password", "image.png", "Bio"))
			},
			expectedError: status.Error(codes.InvalidArgument, "validation error: email must be a valid email address"),
			expectedCode:  codes.InvalidArgument,
		},
		{
			name: "Password Hashing Error",
			setupContext: func() context.Context {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 1, nil
				}
				return context.TODO()
			},
			setupRequest: func() *proto.UpdateUserRequest {
				return &proto.UpdateUserRequest{
					User: &proto.UpdateUserRequest_User{
						Password: "password",
					},
				}
			},
			setupMocks: func() {
				mock.ExpectQuery("SELECT").WillReturnRows(
					sqlmock.NewRows([]string{"id", "username", "email", "password", "image", "bio"}).
						AddRow(1, "username", "old@example.com", "password", "image.png", "Bio"))

				// TODO: Implement an approach for direct bcrypt mocking if needed
			},
			expectedError: status.Error(codes.Aborted, "internal server error"),
			expectedCode:  codes.Aborted,
		},
		{
			name: "Error Generating Token",
			setupContext: func() context.Context {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 1, nil
				}
				return context.TODO()
			},
			setupRequest: func() *proto.UpdateUserRequest {
				return &proto.UpdateUserRequest{}
			},
			setupMocks: func() {
				mock.ExpectQuery("SELECT").WillReturnRows(
					sqlmock.NewRows([]string{"id", "username", "email", "password", "image", "bio"}).
						AddRow(1, "username", "old@example.com", "oldpassword", "oldimage.png", "Old Bio"))

				mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))

				auth.GenerateToken = func(id uint) (string, error) {
					return "", errors.New("failed to generate token")
				}
			},
			expectedError: status.Error(codes.Aborted, "internal server error"),
			expectedCode:  codes.Aborted,
		},
		{
			name: "Internal Server Error During Update",
			setupContext: func() context.Context {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 1, nil
				}
				return context.TODO()
			},
			setupRequest: func() *proto.UpdateUserRequest {
				return &proto.UpdateUserRequest{}
			},
			setupMocks: func() {
				mock.ExpectQuery("SELECT").WillReturnRows(
					sqlmock.NewRows([]string{"id", "username", "email", "password", "image", "bio"}).
						AddRow(1, "username", "old@example.com", "oldpassword", "oldimage.png", "Old Bio"))

				mock.ExpectExec("UPDATE").WillReturnError(errors.New("failed to update user"))
			},
			expectedError: status.Error(codes.InvalidArgument, "internal server error"),
			expectedCode:  codes.InvalidArgument,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := test.setupContext()
			req := test.setupRequest()
			test.setupMocks()

			_, err := handler.UpdateUser(ctx, req)
			if err != nil && test.expectedError != nil {
				assert.Equal(t, status.Code(err), test.expectedCode)
				assert.Equal(t, strings.TrimSpace(err.Error()), strings.TrimSpace(test.expectedError.Error()))
			} else {
				assert.Equal(t, test.expectedError, err)
			}
		})
	}
}
