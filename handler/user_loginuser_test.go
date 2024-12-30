package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestHandlerLoginUser(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer mockDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to initialize gorm DB: %s", err)
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: t})
	userStore := &store.UserStore{DB: gormDB}
	handler := &Handler{logger: &logger, us: userStore}

	testCases := []struct {
		name          string
		req           *proto.LoginUserRequest
		setupMocks    func()
		expectedResp  *proto.UserResponse
		expectedError error
	}{
		{
			name: "Successful Login",
			req: &proto.LoginUserRequest{
				User: &proto.LoginUserRequest_User{
					Email:    "test@example.com",
					Password: "validpassword",
				},
			},
			setupMocks: func() {
				user := &model.User{Email: "test@example.com", Password: "$2a$12$eW5pR69pjoi0.LP1LbcG0e.uJ1QBh5AXQG/tex6M958lom1MJ5G8y", ID: 1}
				userStore.GetByEmailMock = func(email string) (*model.User, error) {
					if email == "test@example.com" {
						return user, nil
					}
					return nil, fmt.Errorf("user not found")
				}
				auth.GenerateTokenMock = func(id uint) (string, error) {
					return "validtoken", nil
				}
			},
			expectedResp: &proto.UserResponse{
				User: &proto.User{
					Email:    "test@example.com",
					Token:    "validtoken",
					Username: "testuser",
					Bio:      "Bio",
					Image:    "Image.jpg",
				},
			},
			expectedError: nil,
		},
		{
			name: "Invalid Email",
			req: &proto.LoginUserRequest{
				User: &proto.LoginUserRequest_User{
					Email:    "invalid@example.com",
					Password: "irrelevant",
				},
			},
			setupMocks: func() {
				userStore.GetByEmailMock = func(email string) (*model.User, error) {
					return nil, fmt.Errorf("user not found")
				}
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.InvalidArgument, "invalid email or password"),
		},
		{
			name: "Incorrect Password",
			req: &proto.LoginUserRequest{
				User: &proto.LoginUserRequest_User{
					Email:    "test@example.com",
					Password: "wrongpassword",
				},
			},
			setupMocks: func() {
				user := &model.User{Email: "test@example.com", Password: "$2a$12$eW5pR69pjoi0.LP1LbcG0e.uJ1QBh5AXQG/tex6M958lom1MJ5G8y"}
				userStore.GetByEmailMock = func(email string) (*model.User, error) {
					if email == "test@example.com" {
						return user, nil
					}
					return nil, fmt.Errorf("user not found")
				}
				auth.GenerateTokenMock = func(id uint) (string, error) {
					return "validtoken", nil
				}
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.InvalidArgument, "invalid email or password"),
		},
		{
			name: "Token Generation Failure",
			req: &proto.LoginUserRequest{
				User: &proto.LoginUserRequest_User{
					Email:    "test@example.com",
					Password: "validpassword",
				},
			},
			setupMocks: func() {
				user := &model.User{Email: "test@example.com", Password: "$2a$12$eW5pR69pjoi0.LP1LbcG0e.uJ1QBh5AXQG/tex6M958lom1MJ5G8y", ID: 1}
				userStore.GetByEmailMock = func(email string) (*model.User, error) {
					if email == "test@example.com" {
						return user, nil
					}
					return nil, fmt.Errorf("user not found")
				}
				auth.GenerateTokenMock = func(id uint) (string, error) {
					return "", fmt.Errorf("token generation error")
				}
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.Aborted, "internal server error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			resp, err := handler.LoginUser(context.Background(), tc.req)

			if resp != nil && tc.expectedResp != nil && resp.User.Token != tc.expectedResp.User.Token {
				t.Fatalf("expected response token %v, got %v", tc.expectedResp.User.Token, resp.User.Token)
			}

			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Fatalf("expected error %v, got %v", tc.expectedError, err)
			}

			if err == nil && tc.expectedError != nil {
				t.Fatalf("expected error %v, got none", tc.expectedError)
			}

			if err != nil && tc.expectedError == nil {
				t.Fatalf("did not expect error, got %v", err)
			}

			if err == nil && tc.expectedResp == nil {
				t.Fatalf("expected nil response, got %v", resp)
			}

			if err == nil && resp.User.Email != tc.expectedResp.User.Email {
				t.Fatalf("expected email %v, got %v", tc.expectedResp.User.Email, resp.User.Email)
			}
		})
	}
}
