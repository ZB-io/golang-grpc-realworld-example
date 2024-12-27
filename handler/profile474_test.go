package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/rs/zerolog"
)

/*
ROOST_METHOD_HASH=ShowProfile_3cf6e3a9fd
ROOST_METHOD_SIG_HASH=ShowProfile_4679c3d9a4
*/

func TestHandlerShowProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := NewMockUserStore(ctrl)
	mockLogger := zerolog.NopLogger{} // Or instantiate according to your logging setup

	handler := &Handler{
		us:     mockUserStore,
		logger: &mockLogger,
	}

	type testCase struct {
		name          string
		mockSetup     func()
		expectedError error
	}

	testCases := []testCase{
		{
			name: "Successfully Retrieve User Profile",
			mockSetup: func() {
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(&model.User{}, nil)
				mockUserStore.EXPECT().GetByUsername(gomock.Any()).Return(&model.User{}, nil)
				mockUserStore.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(true, nil)
				auth.GetUserID = func(ctx context.Context) (string, error) { return "valid-user-id", nil }
			},
			expectedError: nil,
		},
		{
			name: "Unauthenticated User",
			mockSetup: func() {
				auth.GetUserID = func(ctx context.Context) (string, error) {
					return "", status.Error(codes.Unauthenticated, "unauthenticated")
				}
			},
			expectedError: status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Current User Not Found",
			mockSetup: func() {
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(nil, sql.ErrNoRows)
				auth.GetUserID = func(ctx context.Context) (string, error) { return "valid-user-id", nil }
			},
			expectedError: status.Error(codes.NotFound, "user not found"),
		},
		{
			name: "Requested User Not Found",
			mockSetup: func() {
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(&model.User{}, nil)
				mockUserStore.EXPECT().GetByUsername(gomock.Any()).Return(nil, sql.ErrNoRows)
				auth.GetUserID = func(ctx context.Context) (string, error) { return "valid-user-id", nil }
			},
			expectedError: status.Error(codes.NotFound, "user was not found"),
		},
		{
			name: "Following Status Retrieval Fails",
			mockSetup: func() {
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(&model.User{}, nil)
				mockUserStore.EXPECT().GetByUsername(gomock.Any()).Return(&model.User{}, nil)
				mockUserStore.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(false, fmt.Errorf("failed"))
				auth.GetUserID = func(ctx context.Context) (string, error) { return "valid-user-id", nil }
			},
			expectedError: status.Error(codes.NotFound, "internal server error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			req := &pb.ShowProfileRequest{Username: "testuser"}
			resp, err := handler.ShowProfile(context.Background(), req)

			if tc.expectedError != nil {
				assert.Nil(t, resp)
				assert.Equal(t, tc.expectedError, err)
				t.Log(err.Error())
			} else {
				assert.NotNil(t, resp)
				assert.NoError(t, err)
				t.Log("Profile successfully retrieved")
			}
		})
	}
}

/*
ROOST_METHOD_HASH=FollowUser_36d65b5263
ROOST_METHOD_SIG_HASH=FollowUser_bf8ceb04bb
 */

func TestHandlerFollowUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := NewMockUserStore(ctrl)
	mockLogger := zerolog.NopLogger{} // Or instantiate according to your logging setup

	handler := &Handler{
		us:     mockUserStore,
		logger: &mockLogger,
	}

	tests := []struct {
		name              string
		contextSetup      func() context.Context
		requestSetup      func() *pb.FollowRequest
		setupMocks        func()
		expectedErrorCode codes.Code
	}{
		{
			name: "unauthenticated user attempt",
			contextSetup: func() context.Context {
				return context.TODO()
			},
			requestSetup: func() *pb.FollowRequest {
				return &pb.FollowRequest{Username: "targetuser"}
			},
			setupMocks:        func() {},
			expectedErrorCode: codes.Unauthenticated,
		},
		{
			name: "current user not found",
			contextSetup: func() context.Context {
				ctx := context.TODO()
				ctx = context.WithValue(ctx, auth.UserIDKey, 100)
				return ctx
			},
			requestSetup: func() *pb.FollowRequest {
				return &pb.FollowRequest{Username: "targetuser"}
			},
			setupMocks: func() {
				mockUserStore.EXPECT().GetByID(100).Return(nil, status.Error(codes.NotFound, "user not found"))
			},
			expectedErrorCode: codes.NotFound,
		},
		{
			name: "attempt to follow oneself",
			contextSetup: func() context.Context {
				ctx := context.TODO()
				ctx = context.WithValue(ctx, auth.UserIDKey, 1)
				return ctx
			},
			requestSetup: func() *pb.FollowRequest {
				return &pb.FollowRequest{Username: "currentuser"}
			},
			setupMocks: func() {
				currentUser := &model.User{ID: 1, Username: "currentuser"}
				mockUserStore.EXPECT().GetByID(1).Return(currentUser, nil)
			},
			expectedErrorCode: codes.InvalidArgument,
		},
		{
			name: "target user not found",
			contextSetup: func() context.Context {
				ctx := context.TODO()
				ctx = context.WithValue(ctx, auth.UserIDKey, 1)
				return ctx
			},
			requestSetup: func() *pb.FollowRequest {
				return &pb.FollowRequest{Username: "nonexistentuser"}
			},
			setupMocks: func() {
				currentUser := &model.User{ID: 1, Username: "currentuser"}
				mockUserStore.EXPECT().GetByID(1).Return(currentUser, nil)
				mockUserStore.EXPECT().GetByUsername("nonexistentuser").Return(nil, status.Error(codes.NotFound, "user not found"))
			},
			expectedErrorCode: codes.NotFound,
		},
		{
			name: "successful follow operation",
			contextSetup: func() context.Context {
				ctx := context.TODO()
				ctx = context.WithValue(ctx, auth.UserIDKey, 1)
				return ctx
			},
			requestSetup: func() *pb.FollowRequest {
				return &pb.FollowRequest{Username: "targetuser"}
			},
			setupMocks: func() {
				currentUser := &model.User{ID: 1, Username: "currentuser"}
				targetUser := &model.User{ID: 2, Username: "targetuser"}
				mockUserStore.EXPECT().GetByID(1).Return(currentUser, nil)
				mockUserStore.EXPECT().GetByUsername("targetuser").Return(targetUser, nil)
				mockUserStore.EXPECT().Follow(currentUser, targetUser).Return(nil)
			},
			expectedErrorCode: codes.OK,
		},
		{
			name: "failed follow operation due to service error",
			contextSetup: func() context.Context {
				ctx := context.TODO()
				ctx = context.WithValue(ctx, auth.UserIDKey, 1)
				return ctx
			},
			requestSetup: func() *pb.FollowRequest {
				return &pb.FollowRequest{Username: "targetuser"}
			},
			setupMocks: func() {
				currentUser := &model.User{ID: 1, Username: "currentuser"}
				targetUser := &model.User{ID: 2, Username: "targetuser"}
				mockUserStore.EXPECT().GetByID(1).Return(currentUser, nil)
				mockUserStore.EXPECT().GetByUsername("targetuser").Return(targetUser, nil)
				mockUserStore.EXPECT().Follow(currentUser, targetUser).Return(status.Error(codes.Aborted, "failed to follow user"))
			},
			expectedErrorCode: codes.Aborted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			resp, err := handler.FollowUser(tt.contextSetup(), tt.requestSetup())
			if tt.expectedErrorCode == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, resp)
			} else {
				require.Error(t, err)
				require.Nil(t, resp)
				st, _ := status.FromError(err)
				require.Equal(t, tt.expectedErrorCode, st.Code())
			}
		})
	}
}

/*
ROOST_METHOD_HASH=UnfollowUser_843a2807ea
ROOST_METHOD_SIG_HASH=UnfollowUser_a64840f937
*/

func TestHandlerUnfollowUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserStore(ctrl) // Correct to the actual interface used
	logger := zerolog.NopLogger{} // Or properly configure according to your logging package

	type scenario struct {
		name            string
		prepareMocks    func()
		expectedError   codes.Code
		expectedProfile *pb.Profile
	}

	tests := []scenario{
		{
			name: "Successful Unfollow Operation",
			prepareMocks: func() {
				auth.GetUserID = func(ctx context.Context) (string, error) { return "1", nil }
				mockUserService.EXPECT().GetByID(1).Return(&model.User{ID: 1, Username: "current"}, nil)
				mockUserService.EXPECT().GetByUsername("toUnfollow").Return(&model.User{ID: 2, Username: "toUnfollow"}, nil)
				mockUserService.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(true, nil)
				mockUserService.EXPECT().Unfollow(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError:   codes.OK,
			expectedProfile: &pb.Profile{Username: "toUnfollow"},
		},
		{
			name: "Unauthenticated User Error",
			prepareMocks: func() {
				auth.GetUserID = func(ctx context.Context) (string, error) { return "", errors.New("unauthenticated") }
			},
			expectedError: codes.Unauthenticated,
		},
		{
			name: "Current User Not Found",
			prepareMocks: func() {
				auth.GetUserID = func(ctx context.Context) (string, error) { return "1", nil }
				mockUserService.EXPECT().GetByID(1).Return(nil, errors.New("user not found"))
			},
			expectedError: codes.NotFound,
		},
		{
			name: "Request User Not Found",
			prepareMocks: func() {
				auth.GetUserID = func(ctx context.Context) (string, error) { return "1", nil }
				mockUserService.EXPECT().GetByID(1).Return(&model.User{ID: 1, Username: "current"}, nil)
				mockUserService.EXPECT().GetByUsername("toUnfollow").Return(nil, errors.New("user not found"))
			},
			expectedError: codes.NotFound,
		},
		{
			name: "Attempt to Unfollow Self",
			prepareMocks: func() {
				auth.GetUserID = func(ctx context.Context) (string, error) { return "1", nil }
				mockUserService.EXPECT().GetByID(1).Return(&model.User{ID: 1, Username: "current"}, nil)
			},
			expectedError: codes.InvalidArgument,
		},
		{
			name: "User Not Following the Target User",
			prepareMocks: func() {
				auth.GetUserID = func(ctx context.Context) (string, error) { return "1", nil }
				mockUserService.EXPECT().GetByID(1).Return(&model.User{ID: 1, Username: "current"}, nil)
				mockUserService.EXPECT().GetByUsername("toUnfollow").Return(&model.User{ID: 2, Username: "toUnfollow"}, nil)
				mockUserService.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(false, nil)
			},
			expectedError: codes.PermissionDenied,
		},
		{
			name: "Unfollow Operation Failure",
			prepareMocks: func() {
				auth.GetUserID = func(ctx context.Context) (string, error) { return "1", nil }
				mockUserService.EXPECT().GetByID(1).Return(&model.User{ID: 1, Username: "current"}, nil)
				mockUserService.EXPECT().GetByUsername("toUnfollow").Return(&model.User{ID: 2, Username: "toUnfollow"}, nil)
				mockUserService.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(true, nil)
				mockUserService.EXPECT().Unfollow(gomock.Any(), gomock.Any()).Return(errors.New("unfollow failure"))
			},
			expectedError: codes.Aborted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks()
			handler := &Handler{us: mockUserService, logger: &logger}

			resp, err := handler.UnfollowUser(context.TODO(), &pb.UnfollowRequest{Username: "toUnfollow"})

			if err != nil {
				if status.Code(err) != tt.expectedError {
					t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
				}
				t.Logf("Error occurred: %v", err)
			} else {
				if tt.expectedProfile != nil && resp.Profile.GetUsername() != tt.expectedProfile.GetUsername() {
					t.Errorf("expected profile: %v, got: %v", tt.expectedProfile.GetUsername(), resp.Profile.GetUsername())
				}
				t.Logf("Successful response: %+v", resp)
			}
		})
	}
}
