package handler

import (
	"context"
	"errors"
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

func MockAuthContext(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, "user_id", userID)
}

func MockGetUserID(ctx context.Context) (uint, error) {
	userID, ok := ctx.Value("user_id").(uint)
	if !ok {
		return 0, errors.New("user not authenticated")
	}
	return userID, nil
}

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
				us.On("IsFollowing", &model.User{Username: "currentuser"}, &model.User{Username: "requestuser"}).Return(true, nil)
			},
			ctx: MockAuthContext(context.Background(), 1),
			req: &pb.ShowProfileRequest{Username: "requestuser"},
			expectedResp: &pb.ProfileResponse{
				Profile: &pb.Profile{
					Username:  "requestuser",
					Following: true,
				},
			},
			expectedErrCode: codes.OK,
		},
		{
			name:            "Attempt to retrieve profile with unauthenticated request",
			setupMocks:      func(us *MockUserStore) {},
			ctx:             context.Background(),
			req:             &pb.ShowProfileRequest{Username: "anyuser"},
			expectedResp:    nil,
			expectedErrCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := &MockUserStore{}
			tt.setupMocks(mockUserStore)

			h := &Handler{
				logger: zerolog.Nop(),
				us:     mockUserStore,
			}

			resp, err := h.ShowProfile(tt.ctx, tt.req)

			if tt.expectedErrCode != codes.OK {
				if err == nil {
					t.Errorf("Expected error with code %v, got nil", tt.expectedErrCode)
					return
				}
				st, ok := status.FromError(err)
				if !ok {
					t.Errorf("Expected gRPC error, got %v", err)
					return
				}
				if st.Code() != tt.expectedErrCode {
					t.Errorf("Expected error code %v, got %v", tt.expectedErrCode, st.Code())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if resp == nil {
					t.Error("Expected non-nil response, got nil")
					return
				}
				if resp.Profile.Username != tt.expectedResp.Profile.Username {
					t.Errorf("Expected username %s, got %s", tt.expectedResp.Profile.Username, resp.Profile.Username)
				}
				if resp.Profile.Following != tt.expectedResp.Profile.Following {
					t.Errorf("Expected following status %v, got %v", tt.expectedResp.Profile.Following, resp.Profile.Following)
				}
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
		setupMocks     func(*MockUserStore)
		setupAuth      func() context.Context
		input          *pb.FollowRequest
		expectedOutput *pb.ProfileResponse
		expectedError  error
	}{
		{
			name: "Successfully Follow a User",
			setupMocks: func(mus *MockUserStore) {
				mus.On("GetByID", uint(1)).Return(&model.User{ID: 1, Username: "currentUser"}, nil)
				mus.On("GetByUsername", "requestUser").Return(&model.User{ID: 2, Username: "requestUser"}, nil)
				mus.On("Follow", mock.AnythingOfType("*model.User"), mock.AnythingOfType("*model.User")).Return(nil)
			},
			setupAuth: func() context.Context {
				return MockAuthContext(context.Background(), 1)
			},
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
			mockUserStore := &MockUserStore{}
			tt.setupMocks(mockUserStore)
			h := &Handler{
				logger: zerolog.Nop(),
				us:     mockUserStore,
			}
			ctx := tt.setupAuth()

			got, err := h.FollowUser(ctx, tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, got)
			}

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
		setupMock      func(*MockUserStore)
		setupContext   func() context.Context
		input          *pb.UnfollowRequest
		expectedOutput *pb.ProfileResponse
		expectedError  error
	}{
		{
			name: "Successful Unfollow",
			setupMock: func(m *MockUserStore) {
				currentUser := &model.User{ID: 1, Username: "currentUser"}
				requestUser := &model.User{ID: 2, Username: "requestUser"}
				m.On("GetByID", uint(1)).Return(currentUser, nil)
				m.On("GetByUsername", "requestUser").Return(requestUser, nil)
				m.On("IsFollowing", currentUser, requestUser).Return(true, nil)
				m.On("Unfollow", currentUser, requestUser).Return(nil)
			},
			setupContext: func() context.Context {
				return MockAuthContext(context.Background(), 1)
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
				logger: zerolog.Nop(),
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
