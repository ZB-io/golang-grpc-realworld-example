package handler

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/raahii/golang-grpc-realworld-example/auth/mocks"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestShowProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := mocks.NewMockUserStore(ctrl)
	handler := &Handler{
		logger: Logger(),
		us:     mockUserStore,
	}

	type args struct {
		context  context.Context
		username string
	}

	tests := []struct {
		name    string
		setup   func(args) // Sets up the mock expectations
		args    args
		want    *pb.ProfileResponse // Expected Profile Response
		wantErr error               // Expected error
	}{
		{
			name: "Scenario 1: Valid username and authenticated user",
			setup: func(a args) {
				userID := uint(1)
				currentUser := &User{ID: userID}
				requestedUser := &User{Username: a.username}

				mockUserStore.EXPECT().GetByID(userID).Return(currentUser, nil)
				mockUserStore.EXPECT().GetByUsername(a.username).Return(requestedUser, nil)
				mockUserStore.EXPECT().IsFollowing(currentUser, requestedUser).Return(false, nil)
			},
			args:    args{context.WithValue(context.Background(), "userID", uint(1)), "validuser"},
			want:    &pb.ProfileResponse{Profile: &pb.Profile{Username: "validuser", Following: false}},
			wantErr: nil,
		},
		{
			name: "Scenario 2: Unauthenticated user",
			setup: func(a args) {
				// No setup needed since auth.GetUserID will be mocked
			},
			args:    args{context.Background(), "validuser"},
			want:    nil,
			wantErr: status.Errorf(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Scenario 3: Non-existent username",
			setup: func(a args) {
				userID := uint(1)
				currentUser := &User{ID: userID}
				mockUserStore.EXPECT().GetByID(userID).Return(currentUser, nil)
				mockUserStore.EXPECT().GetByUsername(a.username).Return(nil, sqlmock.ErrNoRows)
			},
			args:    args{context.WithValue(context.Background(), "userID", uint(1)), "nonexistentuser"},
			want:    nil,
			wantErr: status.Error(codes.NotFound, "user was not found"),
		},
		{
			name: "Scenario 4: User not following another user",
			setup: func(a args) {
				userID := uint(1)
				currentUser := &User{ID: userID}
				requestedUser := &User{Username: a.username}

				mockUserStore.EXPECT().GetByID(userID).Return(currentUser, nil)
				mockUserStore.EXPECT().GetByUsername(a.username).Return(requestedUser, nil)
				mockUserStore.EXPECT().IsFollowing(currentUser, requestedUser).Return(false, nil)
			},
			args:    args{context.WithValue(context.Background(), "userID", uint(1)), "someuser"},
			want:    &pb.ProfileResponse{Profile: &pb.Profile{Username: "someuser", Following: false}},
			wantErr: nil,
		},
		{
			name: "Scenario 5: Internal server error on follow status check",
			setup: func(a args) {
				userID := uint(1)
				currentUser := &User{ID: userID}
				requestedUser := &User{Username: a.username}

				mockUserStore.EXPECT().GetByID(userID).Return(currentUser, nil)
				mockUserStore.EXPECT().GetByUsername(a.username).Return(requestedUser, nil)
				mockUserStore.EXPECT().IsFollowing(currentUser, requestedUser).Return(false, errors.New("internal error"))
			},
			args:    args{context.WithValue(context.Background(), "userID", uint(1)), "someuser"},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.args)

			got, err := handler.ShowProfile(tt.args.context, &pb.ShowProfileRequest{Username: tt.args.username})

			if tt.wantErr != nil {
				if status.Code(err) != status.Code(tt.wantErr) {
					t.Errorf("ShowProfile() error = %v, wantErr %v", err, tt.wantErr)
				} else {
					t.Logf("Success: expected error %v received", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("ShowProfile() unexpected error = %v", err)
					return
				}
				if got == nil || got.Profile.Username != tt.want.Profile.Username || got.Profile.Following != tt.want.Profile.Following {
					t.Errorf("ShowProfile() = %v, want %v", got, tt.want)
				} else {
					t.Log("Success: expected response received")
				}
			}
		})
	}
}
