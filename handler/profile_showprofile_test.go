package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/auth"
)








func MockAuthGetUserID(ctx context.Context) (string, error) {

	return "", nil
}
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {

	return &MockLogger{}
}
func NewMockUserStore(ctrl *gomock.Controller) *MockUserStore {

	return &MockUserStore{}
}

func TestHandlerShowProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := NewMockUserStore(ctrl)
	mockLogger := NewMockLogger(ctrl)

	handler := &Handler{
		us:     mockUserStore,
		logger: mockLogger,
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
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(&User{}, nil)
				mockUserStore.EXPECT().GetByUsername(gomock.Any()).Return(&User{}, nil)
				mockUserStore.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(true, nil)
				auth.MockAuthGetUserID(gomock.Any()).Return("valid-user-id", nil)
			},
			expectedError: nil,
		},
		{
			name: "Unauthenticated User",
			mockSetup: func() {
				auth.MockAuthGetUserID(gomock.Any()).Return("", status.Error(codes.Unauthenticated, "unauthenticated"))
			},
			expectedError: status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Current User Not Found",
			mockSetup: func() {
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(nil, sql.ErrNoRows)
				auth.MockAuthGetUserID(gomock.Any()).Return("valid-user-id", nil)
			},
			expectedError: status.Error(codes.NotFound, "user not found"),
		},
		{
			name: "Requested User Not Found",
			mockSetup: func() {
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(&User{}, nil)
				mockUserStore.EXPECT().GetByUsername(gomock.Any()).Return(nil, sql.ErrNoRows)
				auth.MockAuthGetUserID(gomock.Any()).Return("valid-user-id", nil)
			},
			expectedError: status.Error(codes.NotFound, "user was not found"),
		},
		{
			name: "Following Status Retrieval Fails",
			mockSetup: func() {
				mockUserStore.EXPECT().GetByID(gomock.Any()).Return(&User{}, nil)
				mockUserStore.EXPECT().GetByUsername(gomock.Any()).Return(&User{}, nil)
				mockUserStore.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(false, fmt.Errorf("failed"))
				auth.MockAuthGetUserID(gomock.Any()).Return("valid-user-id", nil)
			},
			expectedError: status.Error(codes.Internal, "internal server error"),
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


func (h *Handler) ShowProfile(ctx context.Context, req *pb.ShowProfileRequest) (*pb.ProfileResponse, error) {
	h.logger.Info().Msg("show profile")

	userID, err := auth.GetUserID(ctx)
	if err != nil {
		h.logger.Error().Msg("unauthenticated")
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	currentUser, err := h.us.GetByID(userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("current user not found")
		return nil, status.Error(codes.NotFound, "user not found")
	}

	requestUser, err := h.us.GetByUsername(req.GetUsername())
	if err != nil {
		msg := "user was not found"
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.NotFound, msg)
	}

	following, err := h.us.IsFollowing(currentUser, requestUser)
	if err != nil {
		msg := "failed to get following status"
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.ProfileResponse{Profile: &pb.Profile{}}, nil
}
