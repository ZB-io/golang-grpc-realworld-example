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
func TestHandlerCurrentUser(t *testing.T) {

	type args struct {
		context context.Context
		request *pb.Empty
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := zerolog.New(os.Stdout)
	mockUserStore := store.NewMockUserStore(ctrl)

	h := &Handler{
		logger: &mockLogger,
		us:     mockUserStore,
	}

	tests := []struct {
		name          string
		args          args
		mockSetup     func()
		expectedResp  *pb.UserResponse
		expectedError error
	}{
		{
			name: "Successfully Retrieve Current User",
			args: args{
				context: context.TODO(),
				request: &pb.Empty{},
			},
			mockSetup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 1, nil
				}
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(&model.User{
					ID:       1,
					Email:    "test@example.com",
					Username: "testuser",
				}, nil)
				auth.GenerateToken = func(id uint) (string, error) {
					return "testtoken", nil
				}
			},
			expectedResp: &pb.UserResponse{
				User: &pb.User{
					Email:    "test@example.com",
					Token:    "testtoken",
					Username: "testuser",
				},
			},
			expectedError: nil,
		},
		{
			name: "Unauthenticated User Request",
			args: args{
				context: context.TODO(),
				request: &pb.Empty{},
			},
			mockSetup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 0, errors.New("unauthenticated")
				}
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "User Not Found",
			args: args{
				context: context.TODO(),
				request: &pb.Empty{},
			},
			mockSetup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 2, nil
				}
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(nil, errors.New("record not found"))
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.NotFound, "user not found"),
		},
		{
			name: "Token Generation Failure",
			args: args{
				context: context.TODO(),
				request: &pb.Empty{},
			},
			mockSetup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 1, nil
				}
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(&model.User{
					ID:       1,
					Email:    "test@example.com",
					Username: "testuser",
				}, nil)
				auth.GenerateToken = func(id uint) (string, error) {
					return "", errors.New("token generation failed")
				}
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.Aborted, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup()

			resp, err := h.CurrentUser(tt.args.context, tt.args.request)

			if resp != nil && tt.expectedResp != nil {
				if resp.User.Email != tt.expectedResp.User.Email ||
					resp.User.Token != tt.expectedResp.User.Token ||
					resp.User.Username != tt.expectedResp.User.Username {
					t.Errorf("unexpected response user info: got=%v, want=%v", resp.User, tt.expectedResp.User)
				}
			}

			if (err != nil) && (tt.expectedError == nil || err.Error() != tt.expectedError.Error()) {
				t.Errorf("unexpected error: got=%v, want=%v", err, tt.expectedError)
			}
		})
	}
}
