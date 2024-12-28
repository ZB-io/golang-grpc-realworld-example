package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnfavoriteArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock Dependencies
	authMock := auth.NewMockService(ctrl)
	userMock := model.NewMockUserService(ctrl)
	articleMock := model.NewMockArticleService(ctrl)

	// Assuming Handler and its methods utilize these dependencies
	h := &Handler{
		us: userMock,
		as: articleMock,
	}

	tests := []struct {
		name           string
		prepareMocks   func()
		input          *pb.UnfavoriteArticleRequest
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Unauthenticated User Attempts to Unfavorite Article",
			prepareMocks: func() {
				authMock.EXPECT().GetUserID(gomock.Any()).Return(0, errors.New("unauthenticated"))
			},
			input:          &pb.UnfavoriteArticleRequest{Slug: "1"},
			expectedErr:    status.Error(codes.Unauthenticated, "unauthenticated"),
			expectedStatus: codes.Unauthenticated,
		},
		{
			name: "User Not Found in the System",
			prepareMocks: func() {
				authMock.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				userMock.EXPECT().GetByID(1).Return(nil, errors.New("user not found"))
			},
			input:          &pb.UnfavoriteArticleRequest{Slug: "1"},
			expectedErr:    status.Error(codes.NotFound, "not user found"),
			expectedStatus: codes.NotFound,
		},
		{
			name: "Invalid Slug Provided",
			prepareMocks: func() {
				authMock.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				userMock.EXPECT().GetByID(1).Return(&model.User{}, nil)
			},
			input:          &pb.UnfavoriteArticleRequest{Slug: "invalid"},
			expectedErr:    status.Error(codes.InvalidArgument, "invalid article id"),
			expectedStatus: codes.InvalidArgument,
		},
		{
			name: "Article Not Found by Article Service",
			prepareMocks: func() {
				authMock.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				userMock.EXPECT().GetByID(1).Return(&model.User{}, nil)
				articleMock.EXPECT().GetByID(uint(1)).Return(nil, errors.New("article not found"))
			},
			input:          &pb.UnfavoriteArticleRequest{Slug: "1"},
			expectedErr:    status.Error(codes.InvalidArgument, "invalid article id"),
			expectedStatus: codes.InvalidArgument,
		},
		{
			name: "Favorite Removal Failure",
			prepareMocks: func() {
				authMock.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				userMock.EXPECT().GetByID(1).Return(&model.User{}, nil)
				articleMock.EXPECT().GetByID(uint(1)).Return(&model.Article{}, nil)
				articleMock.EXPECT().DeleteFavorite(gomock.Any(), gomock.Any()).Return(errors.New("deletion failed"))
			},
			input:          &pb.UnfavoriteArticleRequest{Slug: "1"},
			expectedErr:    status.Error(codes.InvalidArgument, "failed to remove favorite"),
			expectedStatus: codes.InvalidArgument,
		},
		{
			name: "Successful Unfavorite Operation",
			prepareMocks: func() {
				authMock.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				userMock.EXPECT().GetByID(1).Return(&model.User{}, nil)
				articleMock.EXPECT().GetByID(uint(1)).Return(&model.Article{}, nil)
				articleMock.EXPECT().DeleteFavorite(gomock.Any(), gomock.Any()).Return(nil)
				userMock.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(false, nil)
			},
			input:          &pb.UnfavoriteArticleRequest{Slug: "1"},
			expectedErr:    nil,
			expectedStatus: codes.OK,
		},
		{
			name: "Issue When Checking Following Status",
			prepareMocks: func() {
				authMock.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				userMock.EXPECT().GetByID(1).Return(&model.User{}, nil)
				articleMock.EXPECT().GetByID(uint(1)).Return(&model.Article{}, nil)
				articleMock.EXPECT().DeleteFavorite(gomock.Any(), gomock.Any()).Return(nil)
				userMock.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(false, errors.New("following status error"))
			},
			input:          &pb.UnfavoriteArticleRequest{Slug: "1"},
			expectedErr:    status.Error(codes.NotFound, "internal server error"),
			expectedStatus: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks()

			resp, err := h.UnfavoriteArticle(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Nil(t, resp)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
				t.Logf("Expected error: %v, got: %v", tt.expectedErr, err)
			} else {
				assert.NotNil(t, resp)
				assert.Nil(t, err)
				t.Logf("Expected success response, got: %v", resp)
			}
		})
	}
}
