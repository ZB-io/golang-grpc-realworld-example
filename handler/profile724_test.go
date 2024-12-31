package handler

import (
	"context"
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

func (m *MockUserStore) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
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
		contextUserID   uint
		requestUsername string
		wantProfile     *pb.Profile
		wantErr         error
	}{
		{
			name: "Successfully retrieve a user's profile",
			setupMocks: func(us *MockUserStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				us.On("GetByUsername", "test-user").Return(&model.User{Username: "test-user"}, nil)
				us.On("IsFollowing", mock.Anything, mock.Anything).Return(true, nil)
			},
			contextUserID:   1,
			requestUsername: "test-user",
			wantProfile: &pb.Profile{
				Username:  "test-user",
				Following: true,
			},
			wantErr: nil,
		},
		{
			name:            "Attempt to retrieve profile with unauthenticated user",
			setupMocks:      func(us *MockUserStore) {},
			contextUserID:   0,
			requestUsername: "test-user",
			wantProfile:     nil,
			wantErr:         status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := new(MockUserStore)
			tt.setupMocks(mockUserStore)

			h := &Handler{
				logger: &zerolog.Logger{},
				us:     mockUserStore,
			}

			ctx := context.Background()
			if tt.contextUserID != 0 {
				ctx = auth.NewContext(ctx, tt.contextUserID)
			}

			req := &pb.ShowProfileRequest{
				Username: tt.requestUsername,
			}

			got, err := h.ShowProfile(ctx, req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantProfile, got.GetProfile())
			}

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
		setupMock      func(*MockUserStore)
		setupAuth      func() context.Context
		input          *pb.FollowRequest
		expectedOutput *pb.ProfileResponse
		expectedError  error
	}{
		{
			name: "Successfully Follow a User",
			setupMock: func(m *MockUserStore) {
				currentUser := &model.User{Username: "current"}
				requestUser := &model.User{Username: "request"}
				m.On("GetByID", uint(1)).Return(currentUser, nil)
				m.On("GetByUsername", "request").Return(requestUser, nil)
				m.On("Follow", currentUser, requestUser).Return(nil)
			},
			setupAuth: func() context.Context {
				return auth.NewContext(context.Background(), 1)
			},
			input: &pb.FollowRequest{Username: "request"},
			expectedOutput: &pb.ProfileResponse{
				Profile: &pb.Profile{
					Username:  "request",
					Following: true,
				},
			},
			expectedError: nil,
		},
		{
			name: "Attempt to Follow Oneself",
			setupMock: func(m *MockUserStore) {
				currentUser := &model.User{Username: "current"}
				m.On("GetByID", uint(1)).Return(currentUser, nil)
			},
			setupAuth: func() context.Context {
				return auth.NewContext(context.Background(), 1)
			},
			input:          &pb.FollowRequest{Username: "current"},
			expectedOutput: nil,
			expectedError:  status.Error(codes.InvalidArgument, "cannot follow yourself"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUS := new(MockUserStore)
			tt.setupMock(mockUS)

			h := &Handler{
				logger: &zerolog.Logger{},
				us:     mockUS,
			}

			ctx := tt.setupAuth()
			output, err := h.FollowUser(ctx, tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, output)
			}

			mockUS.AssertExpectations(t)
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
		setupMock      func(*MockUserStore)
		setupContext   func() context.Context
		input          *pb.UnfollowRequest
		expectedOutput *pb.ProfileResponse
		expectedError  error
	}{
		{
			name: "Successful Unfollow",
			setupMock: func(m *MockUserStore) {
				currentUser := &model.User{Username: "currentUser"}
				requestUser := &model.User{Username: "requestUser"}
				m.On("GetByID", uint(1)).Return(currentUser, nil)
				m.On("GetByUsername", "requestUser").Return(requestUser, nil)
				m.On("IsFollowing", currentUser, requestUser).Return(true, nil)
				m.On("Unfollow", currentUser, requestUser).Return(nil)
			},
			setupContext: func() context.Context {
				return auth.NewContext(context.Background(), 1)
			},
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
			tt.setupMock(mockUserStore)

			h := &Handler{
				logger: &zerolog.Logger{},
				us:     mockUserStore,
			}

			ctx := tt.setupContext()
			output, err := h.UnfollowUser(ctx, tt.input)

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
