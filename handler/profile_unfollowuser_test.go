package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)

type User struct {
	ID       int
	Username string
}
func TestHandlerUnfollowUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := NewMockAuth(ctrl)
	mockUserService := NewMockUserService(ctrl)
	logger := NewMockLogger(ctrl)

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
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				mockUserService.EXPECT().GetByID(1).Return(&User{ID: 1, Username: "current"}, nil)
				mockUserService.EXPECT().GetByUsername("toUnfollow").Return(&User{ID: 2, Username: "toUnfollow"}, nil)
				mockUserService.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(true, nil)
				mockUserService.EXPECT().Unfollow(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError:   codes.OK,
			expectedProfile: &pb.Profile{Username: "toUnfollow"},
		},
		{
			name: "Unauthenticated User Error",
			prepareMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(0, errors.New("unauthenticated"))
			},
			expectedError: codes.Unauthenticated,
		},
		{
			name: "Current User Not Found",
			prepareMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				mockUserService.EXPECT().GetByID(1).Return(nil, errors.New("user not found"))
			},
			expectedError: codes.NotFound,
		},
		{
			name: "Request User Not Found",
			prepareMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				mockUserService.EXPECT().GetByID(1).Return(&User{ID: 1, Username: "current"}, nil)
				mockUserService.EXPECT().GetByUsername("toUnfollow").Return(nil, errors.New("user not found"))
			},
			expectedError: codes.NotFound,
		},
		{
			name: "Attempt to Unfollow Self",
			prepareMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				mockUserService.EXPECT().GetByID(1).Return(&User{ID: 1, Username: "current"}, nil)
			},
			expectedError: codes.InvalidArgument,
		},
		{
			name: "User Not Following the Target User",
			prepareMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				mockUserService.EXPECT().GetByID(1).Return(&User{ID: 1, Username: "current"}, nil)
				mockUserService.EXPECT().GetByUsername("toUnfollow").Return(&User{ID: 2, Username: "toUnfollow"}, nil)
				mockUserService.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(false, nil)
			},
			expectedError: codes.PermissionDenied,
		},
		{
			name: "Unfollow Operation Failure",
			prepareMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				mockUserService.EXPECT().GetByID(1).Return(&User{ID: 1, Username: "current"}, nil)
				mockUserService.EXPECT().GetByUsername("toUnfollow").Return(&User{ID: 2, Username: "toUnfollow"}, nil)
				mockUserService.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(true, nil)
				mockUserService.EXPECT().Unfollow(gomock.Any(), gomock.Any()).Return(errors.New("unfollow failure"))
			},
			expectedError: codes.Aborted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks()
			h := &Handler{us: mockUserService, logger: logger}

			resp, err := h.UnfollowUser(context.TODO(), &pb.UnfollowRequest{Username: "toUnfollow"})

			if err != nil {
				if status.Code(err) != tt.expectedError {
					t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
				}
				t.Logf("Error occurred: %v", err)
			} else {
				if tt.expectedProfile != nil && resp.Profile.GetUsername() != tt.expectedProfile.GetUsername() {
					t.Errorf("expected profile: %v, got: %v", tt.expectedProfile, resp.Profile.GetUsername())
				}
				t.Logf("Successful response: %+v", resp)
			}
		})
	}
}


