package handler

import (
	"context"
	"fmt"
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/rs/zerolog"
)









type mockUserService struct {
	mock.Mock
}

func (l *mockLogger) Error() *zerolog.Event {

	return &zerolog.Event{}
}
func (m *mockUserService) Follow(currentUser *User, requestUser *User) error {
	args := m.Called(currentUser, requestUser)
	return args.Error(0)
}
func (m *mockUserService) GetByID(userID int) (*User, error) {
	args := m.Called(userID)
	return args.Get(0).(*User), args.Error(1)
}
func (m *mockUserService) GetByUsername(username string) (*User, error) {
	args := m.Called(username)
	return args.Get(0).(*User), args.Error(1)
}
func (l *mockLogger) Info() *zerolog.Event {

	return &zerolog.Event{}
}
func (u *User) ProtoProfile(following bool) *pb.Profile {
	return &pb.Profile{Username: u.Username, Following: following}
}
func TestHandlerFollowUser(t *testing.T) {

	mockUS := new(mockUserService)

	h := &Handler{
		logger: &mockLogger{},
		us:     mockUS,
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
				mockUS.On("GetByID", 100).Return(nil, status.Error(codes.NotFound, "user not found"))
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
				currentUser := &User{ID: 1, Username: "currentuser"}
				mockUS.On("GetByID", 1).Return(currentUser, nil)
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
				currentUser := &User{ID: 1, Username: "currentuser"}
				mockUS.On("GetByID", 1).Return(currentUser, nil)
				mockUS.On("GetByUsername", "nonexistentuser").Return(nil, status.Error(codes.NotFound, "user not found"))
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
				currentUser := &User{ID: 1, Username: "currentuser"}
				targetUser := &User{ID: 2, Username: "targetuser"}
				mockUS.On("GetByID", 1).Return(currentUser, nil)
				mockUS.On("GetByUsername", "targetuser").Return(targetUser, nil)
				mockUS.On("Follow", currentUser, targetUser).Return(nil)
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
				currentUser := &User{ID: 1, Username: "currentuser"}
				targetUser := &User{ID: 2, Username: "targetuser"}
				mockUS.On("GetByID", 1).Return(currentUser, nil)
				mockUS.On("GetByUsername", "targetuser").Return(targetUser, nil)
				mockUS.On("Follow", currentUser, targetUser).Return(status.Error(codes.Aborted, "failed to follow user"))
			},
			expectedErrorCode: codes.Aborted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			resp, err := h.FollowUser(tt.contextSetup(), tt.requestSetup())
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


