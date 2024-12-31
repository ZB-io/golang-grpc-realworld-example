package handler

import (
	"context"
	"testing"

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

func (m *MockUserStore) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) IsFollowing(follower, followee *model.User) (bool, error) {
	args := m.Called(follower, followee)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserStore) Follow(follower, followee *model.User) error {
	args := m.Called(follower, followee)
	return args.Error(0)
}

func (m *MockUserStore) Unfollow(follower, followee *model.User) error {
	args := m.Called(follower, followee)
	return args.Error(0)
}

/*
ROOST_METHOD_HASH=ShowProfile_3cf6e3a9fd
ROOST_METHOD_SIG_HASH=ShowProfile_4679c3d9a4
*/
func TestShowProfile(t *testing.T) {
	tests := []struct {
		name            string
		setupMocks      func(*MockUserStore)
		ctx             context.Context
		req             *pb.ShowProfileRequest
		expectedResp    *pb.ProfileResponse
		expectedErrCode codes.Code
	}{
		{
			name: "Successfully retrieve profile for an existing user",
			setupMocks: func(us *MockUserStore) {
				us.On("GetByID", uint(1)).Return(&model.User{Username: "currentuser"}, nil)
				us.On("GetByUsername", "requestuser").Return(&model.User{Username: "requestuser"}, nil)
				us.On("IsFollowing", mock.Anything, mock.Anything).Return(true, nil)
			},
			ctx: context.WithValue(context.Background(), "user_id", uint(1)),
			req: &pb.ShowProfileRequest{Username: "requestuser"},
			expectedResp: &pb.ProfileResponse{
				Profile: &pb.Profile{
					Username:  "requestuser",
					Following: true,
				},
			},
			expectedErrCode: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := new(MockUserStore)
			tt.setupMocks(mockUserStore)

			h := &Handler{
				logger: &zerolog.Logger{},
				us:     mockUserStore,
			}

			resp, err := h.ShowProfile(tt.ctx, tt.req)

			if err != nil {
				if e, ok := status.FromError(err); ok {
					assert.Equal(t, tt.expectedErrCode, e.Code())
				} else {
					t.Errorf("expected gRPC error, got %v", err)
				}
			} else if tt.expectedErrCode != codes.OK {
				t.Errorf("expected error code %v, got nil error", tt.expectedErrCode)
			}

			assert.Equal(t, tt.expectedResp, resp)
			mockUserStore.AssertExpectations(t)
		})
	}
}

/*
ROOST_METHOD_HASH=FollowUser_36d65b5263
ROOST_METHOD_SIG_HASH=FollowUser_bf8ceb04bb
*/
func TestFollowUser(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*MockUserStore)
		ctx            context.Context
		input          *pb.FollowRequest
		expectedOutput *pb.ProfileResponse
		expectedError  error
	}{
		{
			name: "Successfully Follow a User",
			setupMocks: func(mus *MockUserStore) {
				mus.On("GetByID", uint(1)).Return(&model.User{Username: "currentUser"}, nil)
				mus.On("GetByUsername", "requestUser").Return(&model.User{Username: "requestUser"}, nil)
				mus.On("Follow", mock.Anything, mock.Anything).Return(nil)
			},
			ctx:   context.WithValue(context.Background(), "user_id", uint(1)),
			input: &pb.FollowRequest{Username: "requestUser"},
			expectedOutput: &pb.ProfileResponse{
				Profile: &pb.Profile{
					Username:  "requestUser",
					Following: true,
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := new(MockUserStore)
			tt.setupMocks(mockUserStore)

			h := &Handler{
				logger: &zerolog.Logger{},
				us:     mockUserStore,
			}

			output, err := h.FollowUser(tt.ctx, tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedOutput, output)
			mockUserStore.AssertExpectations(t)
		})
	}
}

/*
ROOST_METHOD_HASH=UnfollowUser_843a2807ea
ROOST_METHOD_SIG_HASH=UnfollowUser_a64840f937
*/
func TestUnfollowUser(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*MockUserStore)
		ctx            context.Context
		input          *pb.UnfollowRequest
		expectedOutput *pb.ProfileResponse
		expectedError  error
	}{
		{
			name: "Successful Unfollow",
			setupMocks: func(mus *MockUserStore) {
				mus.On("GetByID", uint(1)).Return(&model.User{Username: "currentUser"}, nil)
				mus.On("GetByUsername", "requestUser").Return(&model.User{Username: "requestUser"}, nil)
				mus.On("IsFollowing", mock.Anything, mock.Anything).Return(true, nil)
				mus.On("Unfollow", mock.Anything, mock.Anything).Return(nil)
			},
			ctx:   context.WithValue(context.Background(), "user_id", uint(1)),
			input: &pb.UnfollowRequest{Username: "requestUser"},
			expectedOutput: &pb.ProfileResponse{
				Profile: &pb.Profile{
					Username:  "requestUser",
					Following: false,
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := new(MockUserStore)
			tt.setupMocks(mockUserStore)

			h := &Handler{
				logger: &zerolog.Logger{},
				us:     mockUserStore,
			}

			output, err := h.UnfollowUser(tt.ctx, tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedOutput, output)
			mockUserStore.AssertExpectations(t)
		})
	}
}
