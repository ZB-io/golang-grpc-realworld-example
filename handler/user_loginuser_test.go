package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

type User struct {
	gorm.Model
	Username         string    `gorm:"unique_index;not null"`
	Email            string    `gorm:"unique_index;not null"`
	Password         string    `gorm:"not null"`
	Bio              string    `gorm:"not null"`
	Image            string    `gorm:"not null"`
	Follows          []User    `gorm:"many2many:follows;jointable_foreignkey:from_user_id;association_jointable_foreignkey:to_user_id"`
	FavoriteArticles []Article `gorm:"many2many:favorite_articles;"`
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
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
func TestHandlerLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := store.NewMockUserStore(ctrl)
	mockLogger := zerolog.New(os.Stdout).With().Logger()
	handler := &Handler{
		logger: &mockLogger,
		us:     mockUserStore,
	}

	tests := []struct {
		name          string
		setupMocks    func()
		req           *pb.LoginUserRequest
		expectedResp  *pb.UserResponse
		expectedError error
	}{
		{
			name: "Scenario 1: Successful Login",
			setupMocks: func() {
				mockUser := &model.User{
					Email:    "test@example.com",
					Password: "securepasswordhash",
				}
				mockUserStore.EXPECT().GetByEmail("test@example.com").Return(mockUser, nil)

				auth.GenerateToken = func(uint) (string, error) {
					return "validtoken", nil
				}
			},
			req: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "test@example.com",
					Password: "securepassword",
				},
			},
			expectedResp: &pb.UserResponse{
				User: &pb.User{
					Email: "test@example.com",
					Token: "validtoken",
				},
			},
			expectedError: nil,
		},
		{
			name: "Scenario 2: Invalid Email",
			setupMocks: func() {
				mockUserStore.EXPECT().GetByEmail("invalid@example.com").Return(nil, errors.New("user not found"))
			},
			req: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "invalid@example.com",
					Password: "irrelevant",
				},
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.InvalidArgument, "invalid email or password"),
		},
		{
			name: "Scenario 3: Incorrect Password",
			setupMocks: func() {
				mockUser := &model.User{
					Email:    "test@example.com",
					Password: "securepasswordhash",
				}
				mockUserStore.EXPECT().GetByEmail("test@example.com").Return(mockUser, nil)

				mockUser.CheckPassword = func(string) bool {
					return false
				}
			},
			req: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "test@example.com",
					Password: "wrongpassword",
				},
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.InvalidArgument, "invalid email or password"),
		},
		{
			name: "Scenario 4: Token Generation Failure",
			setupMocks: func() {
				mockUser := &model.User{
					Email:    "test@example.com",
					Password: "securepasswordhash",
				}
				mockUserStore.EXPECT().GetByEmail("test@example.com").Return(mockUser, nil)
				auth.GenerateToken = func(uint) (string, error) {
					return "", errors.New("token generation failed")
				}
			},
			req: &pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "test@example.com",
					Password: "securepassword",
				},
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.Aborted, "internal server error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			resp, err := handler.LoginUser(context.Background(), tc.req)

			if tc.expectedResp != nil && resp == nil {
				t.Errorf("Expected non-nil response but got nil")
			}

			if tc.expectedResp == nil && resp != nil {
				t.Errorf("Expected nil response but got non-nil")
			}

			if err != nil && tc.expectedError.Error() != err.Error() {
				t.Errorf("Expected error %v but got %v", tc.expectedError, err)
			}

			if err == nil && tc.expectedError != nil {
				t.Errorf("Expected error %v but got no error", tc.expectedError)
			}
		})
	}
}
