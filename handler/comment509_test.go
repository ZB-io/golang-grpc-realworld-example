package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
)








/*
ROOST_METHOD_HASH=DeleteComment_452af2f984
ROOST_METHOD_SIG_HASH=DeleteComment_27615e7d69

FUNCTION_DEF=func (h *Handler) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.Empty, error) 

 */
func TestHandlerDeleteComment(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zerolog.New(zerolog.NewTestWriter(t))

	tests := []struct {
		name          string
		setupMocks    func(*store.UserStore, *store.ArticleStore)
		ctx           context.Context
		request       *pb.DeleteCommentRequest
		expectedError error
	}{
		{
			name: "Successful comment deletion",
			setupMocks: func(us *store.UserStore, as *store.ArticleStore) {

				us.EXPECT().GetByID(uint(1)).Return(&model.User{ID: 1}, nil)

				as.EXPECT().GetCommentByID(uint(1)).Return(&model.Comment{
					ID:        1,
					UserID:    1,
					ArticleID: 1,
				}, nil)
				as.EXPECT().DeleteComment(gomock.Any()).Return(nil)
			},
			ctx: context.WithValue(context.Background(), "user_id", uint(1)),
			request: &pb.DeleteCommentRequest{
				Slug: "1",
				Id:   "1",
			},
			expectedError: nil,
		},
		{
			name: "Unauthenticated user",
			setupMocks: func(us *store.UserStore, as *store.ArticleStore) {

			},
			ctx: context.Background(),
			request: &pb.DeleteCommentRequest{
				Slug: "1",
				Id:   "1",
			},
			expectedError: status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Invalid comment ID format",
			setupMocks: func(us *store.UserStore, as *store.ArticleStore) {
				us.EXPECT().GetByID(uint(1)).Return(&model.User{ID: 1}, nil)
			},
			ctx: context.WithValue(context.Background(), "user_id", uint(1)),
			request: &pb.DeleteCommentRequest{
				Slug: "1",
				Id:   "invalid",
			},
			expectedError: status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name: "Comment not found",
			setupMocks: func(us *store.UserStore, as *store.ArticleStore) {
				us.EXPECT().GetByID(uint(1)).Return(&model.User{ID: 1}, nil)
				as.EXPECT().GetCommentByID(uint(1)).Return(nil, errors.New("comment not found"))
			},
			ctx: context.WithValue(context.Background(), "user_id", uint(1)),
			request: &pb.DeleteCommentRequest{
				Slug: "1",
				Id:   "1",
			},
			expectedError: status.Error(codes.InvalidArgument, "failed to get comment"),
		},
		{
			name: "Comment not in article",
			setupMocks: func(us *store.UserStore, as *store.ArticleStore) {
				us.EXPECT().GetByID(uint(1)).Return(&model.User{ID: 1}, nil)
				as.EXPECT().GetCommentByID(uint(1)).Return(&model.Comment{
					ID:        1,
					UserID:    1,
					ArticleID: 2,
				}, nil)
			},
			ctx: context.WithValue(context.Background(), "user_id", uint(1)),
			request: &pb.DeleteCommentRequest{
				Slug: "1",
				Id:   "1",
			},
			expectedError: status.Error(codes.InvalidArgument, "the comment is not in the article"),
		},
		{
			name: "Unauthorized user",
			setupMocks: func(us *store.UserStore, as *store.ArticleStore) {
				us.EXPECT().GetByID(uint(1)).Return(&model.User{ID: 1}, nil)
				as.EXPECT().GetCommentByID(uint(1)).Return(&model.Comment{
					ID:        1,
					UserID:    2,
					ArticleID: 1,
				}, nil)
			},
			ctx: context.WithValue(context.Background(), "user_id", uint(1)),
			request: &pb.DeleteCommentRequest{
				Slug: "1",
				Id:   "1",
			},
			expectedError: status.Error(codes.InvalidArgument, "forbidden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockUserStore := store.NewMockUserStore(ctrl)
			mockArticleStore := store.NewMockArticleStore(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockUserStore, mockArticleStore)
			}

			h := &Handler{
				logger: &logger,
				us:     mockUserStore,
				as:     mockArticleStore,
			}

			_, err := h.DeleteComment(tt.ctx, tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

