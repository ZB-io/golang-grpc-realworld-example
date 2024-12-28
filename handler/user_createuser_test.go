package handler

import (
	"context"
	"errors"
	"testing"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/rs/zerolog"
)

type ExpectedBegin struct {
	commonExpectation
	delay time.Duration
}

type ExpectedCommit struct {
	commonExpectation
}

type ExpectedExec struct {
	queryBasedExpectation
	result driver.Result
	delay  time.Duration
}

type ExpectedRollback struct {
	commonExpectation
}


type Controller struct {
	// T should only be called within a generated mock. It is not intended to
	// be used in user code and may be changed in future versions. T is the
	// TestReporter passed in when creating the Controller via NewController.
	// If the TestReporter does not implement a TestHelper it will be wrapped
	// with a nopTestHelper.
	T             TestHelper
	mu            sync.Mutex
	expectedCalls *callSet
	finished      bool
}

type UserStore struct {
	db *gorm.DB
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}


type Controller struct {
	// T should only be called within a generated mock. It is not intended to
	// be used in user code and may be changed in future versions. T is the
	// TestReporter passed in when creating the Controller via NewController.
	// If the TestReporter does not implement a TestHelper it will be wrapped
	// with a nopTestHelper.
	T             TestHelper
	mu            sync.Mutex
	expectedCalls *callSet
	finished      bool
}

type Logger struct {
	w       LevelWriter
	level   Level
	sampler Sampler
	context []byte
	hooks   []Hook
}


type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestHandlerCreateUser(t *testing.T) {
	type testCase struct {
		name         string
		request      *pb.CreateUserRequest
		setupMocks   func(mock sqlmock.Sqlmock)
		expectedResp *pb.UserResponse
		expectedErr  error
	}

	tests := []testCase{
		{
			name: "Successful User Creation with Valid Input",
			request: &pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Email:    "valid@example.com",
					Username: "validuser",
					Password: "securepassword",
				},
			},
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WithArgs(sqlmock.AnyArg(), "validuser", "valid@example.com", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedResp: &pb.UserResponse{
				User: &pb.User{
					Email:    "valid@example.com",
					Username: "validuser",
					Token:    "sampletoken",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Validation Error on Invalid User Data",
			request: &pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Email:    "",
					Username: "invaliduser",
					Password: "securepassword",
				},
			},
			setupMocks:   func(mock sqlmock.Sqlmock) {},
			expectedResp: nil,
			expectedErr:  status.Error(codes.InvalidArgument, "validation error"),
		},
		{
			name: "Password Hashing Failure",
			request: &pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Email:    "valid@example.com",
					Username: "validuser",
					Password: "",
				},
			},
			setupMocks: func(mock sqlmock.Sqlmock) {

			},
			expectedResp: nil,
			expectedErr:  status.Error(codes.Aborted, "internal server error"),
		},
		{
			name: "User Creation Failure in Store",
			request: &pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Email:    "valid@example.com",
					Username: "validuser",
					Password: "securepassword",
				},
			},
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WithArgs(sqlmock.AnyArg(), "validuser", "valid@example.com", sqlmock.AnyArg()).
					WillReturnError(fmt.Errorf("duplicate key value violates unique constraint"))
				mock.ExpectRollback()
			},
			expectedResp: nil,
			expectedErr:  status.Error(codes.Canceled, "internal server error"),
		},
		{
			name: "Token Generation Failure",
			request: &pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Email:    "valid@example.com",
					Username: "validuser",
					Password: "securepassword",
				},
			},
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WithArgs(sqlmock.AnyArg(), "validuser", "valid@example.com", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				auth.GenerateToken = func(id uint) (string, error) {
					return "", errors.New("token generation failed")
				}
			},
			expectedResp: nil,
			expectedErr:  status.Error(codes.Aborted, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mock, ctrl := setUp(t)
			defer ctrl.Finish()

			tt.setupMocks(*mock)

			resp, err := handler.CreateUser(context.Background(), tt.request)

			if resp != nil && tt.expectedResp != nil {
				if resp.User.Email != tt.expectedResp.User.Email ||
					resp.User.Username != tt.expectedResp.User.Username {
					t.Errorf("expected response %v, got %v", tt.expectedResp, resp)
				}
			}

			if err != nil && tt.expectedErr != nil {
				st, _ := status.FromError(err)
				expectedSt, _ := status.FromError(tt.expectedErr)
				if st.Code() != expectedSt.Code() || st.Message() != expectedSt.Message() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			}
		})
	}
}
func setUp(t *testing.T) (*Handler, *sqlmock.Sqlmock, *gomock.Controller) {
	mockCtrl := gomock.NewController(t)
	logger := zerolog.New(nil)
	db, mock, _ := sqlmock.New()
	userStore := store.NewUserStore(db)
	handler := &Handler{
		logger: &logger,
		us:     userStore,
	}
	return handler, &mock, mockCtrl
}
