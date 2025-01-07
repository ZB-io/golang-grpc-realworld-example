package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)





type MockUserStore struct {
	mock.Mock
}


/*
ROOST_METHOD_HASH=ShowProfile_3cf6e3a9fd
ROOST_METHOD_SIG_HASH=ShowProfile_4679c3d9a4

FUNCTION_DEF=func (h *Handler) ShowProfile(ctx context.Context, req *pb.ShowProfileRequest) (*pb.ProfileResponse, error) 

 */
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

func (m *MockUserStore) IsFollowing(a *model.User, b *model.User) (bool, error) {
	args := m.Called(a, b)
	return args.Bool(0), args.Error(1)
}

func TestHandlerShowProfile(t *testing.T) {

	tests := []struct {
		name           string
		setupMocks     func(*MockUserStore)
		setupContext   func() context.Context
		input          *pb.ShowProfileRequest
		expectedOutput *pb.ProfileResponse
		expectedError  error
	}{
		{
			name: "Successful Profile Retrieval",
			setupMocks: func(us *MockUserStore) {
				currentUser := &model.User{ID: 1, Username: "current"}
				requestUser := &model.User{ID: 2, Username: "requested"}
				us.On("GetByID", uint(1)).Return(currentUser, nil)
				us.On("GetByUsername", "requested").Return(requestUser, nil)
				us.On("IsFollowing", currentUser, requestUser).Return(true, nil)
			},
			setupContext: func() context.Context {

				return context.Background()
			},
			input: &pb.ShowProfileRequest{Username: "requested"},
			expectedOutput: &pb.ProfileResponse{
				Profile: &pb.Profile{
					Username:  "requested",
					Following: true,
				},
			},
			expectedError: nil,
		},
		{
			name: "Unauthenticated Request",
			setupMocks: func(us *MockUserStore) {

			},
			setupContext: func() context.Context {

				return context.Background()
			},
			input:          &pb.ShowProfileRequest{Username: "any"},
			expectedOutput: nil,
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Current User Not Found",
			setupMocks: func(us *MockUserStore) {
				us.On("GetByID", uint(1)).Return(nil, errors.New("user not found"))
			},
			setupContext: func() context.Context {

				return context.Background()
			},
			input:          &pb.ShowProfileRequest{Username: "requested"},
			expectedOutput: nil,
			expectedError:  status.Error(codes.NotFound, "user not found"),
		},
		{
			name: "Requested Profile User Not Found",
			setupMocks: func(us *MockUserStore) {
				currentUser := &model.User{ID: 1, Username: "current"}
				us.On("GetByID", uint(1)).Return(currentUser, nil)
				us.On("GetByUsername", "nonexistent").Return(nil, errors.New("not found"))
			},
			setupContext: func() context.Context {

				return context.Background()
			},
			input:          &pb.ShowProfileRequest{Username: "nonexistent"},
			expectedOutput: nil,
			expectedError:  status.Error(codes.NotFound, "user was not found"),
		},
		{
			name: "Following Status Error",
			setupMocks: func(us *MockUserStore) {
				currentUser := &model.User{ID: 1, Username: "current"}
				requestUser := &model.User{ID: 2, Username: "requested"}
				us.On("GetByID", uint(1)).Return(currentUser, nil)
				us.On("GetByUsername", "requested").Return(requestUser, nil)
				us.On("IsFollowing", currentUser, requestUser).Return(false, errors.New("db error"))
			},
			setupContext: func() context.Context {

				return context.Background()
			},
			input:          &pb.ShowProfileRequest{Username: "requested"},
			expectedOutput: nil,
			expectedError:  status.Error(codes.NotFound, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockUS := new(MockUserStore)
			tt.setupMocks(mockUS)

			logger := zerolog.New(zerolog.NewTestWriter(t))
			h := &Handler{
				logger: &logger,
				us:     mockUS,
			}

			ctx := tt.setupContext()
			response, err := h.ShowProfile(ctx, tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, status.Code(tt.expectedError), status.Code(err))
				assert.Contains(t, err.Error(), status.Convert(tt.expectedError).Message())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, response)
			}

			mockUS.AssertExpectations(t)
		})
	}
}

