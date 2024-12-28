package handler

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFavoriteArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockArticleService := NewMockArticleService(ctrl)

	h := &Handler{
		us: mockUserService,
		as: mockArticleService,
	}

	userID := 1
	articleID := 1
	existingArticle := &model.Article{ID: uint(articleID)}
	validContext := auth.ContextWithUserID(context.Background(), userID)

	tests := []struct {
		name       string
		ctx        context.Context
		req        *pb.FavoriteArticleRequest
		setupMocks func()
		assertion  func(*pb.ArticleResponse, error)
	}{
		{
			name:       "Unauthorized User Attempting to Favorite an Article",
			ctx:        context.Background(),
			req:        &pb.FavoriteArticleRequest{Slug: "1"},
			setupMocks: func() {},
			assertion: func(resp *pb.ArticleResponse, err error) {
				if status.Code(err) != codes.Unauthenticated {
					t.Fatalf("expected Unauthenticated error, got %v", err)
				}
			},
		},
		{
			name: "User Not Found in System",
			ctx:  validContext,
			req:  &pb.FavoriteArticleRequest{Slug: "1"},
			setupMocks: func() {
				mockUserService.EXPECT().GetByID(userID).Return(nil, fmt.Errorf("user not found"))
			},
			assertion: func(resp *pb.ArticleResponse, err error) {
				if status.Code(err) != codes.NotFound {
					t.Fatalf("expected NotFound error, got %v", err)
				}
			},
		},
		{
			name: "Invalid Article Slug Format",
			ctx:  validContext,
			req:  &pb.FavoriteArticleRequest{Slug: "abcd"},
			setupMocks: func() {
				mockUserService.EXPECT().GetByID(userID).Return(&model.User{ID: uint(userID)}, nil)
			},
			assertion: func(resp *pb.ArticleResponse, err error) {
				if status.Code(err) != codes.InvalidArgument {
					t.Fatalf("expected InvalidArgument error, got %v", err)
				}
			},
		},
		{
			name: "Article Not Found by Given Slug",
			ctx:  validContext,
			req:  &pb.FavoriteArticleRequest{Slug: "12345"},
			setupMocks: func() {
				mockUserService.EXPECT().GetByID(userID).Return(&model.User{ID: uint(userID)}, nil)
				mockArticleService.EXPECT().GetByID(uint(12345)).Return(nil, fmt.Errorf("article not found"))
			},
			assertion: func(resp *pb.ArticleResponse, err error) {
				if status.Code(err) != codes.InvalidArgument {
					t.Fatalf("expected InvalidArgument error, got %v", err)
				}
			},
		},
		{
			name: "Successful Article Favoriting",
			ctx:  validContext,
			req:  &pb.FavoriteArticleRequest{Slug: strconv.Itoa(articleID)},
			setupMocks: func() {
				mockUserService.EXPECT().GetByID(userID).Return(&model.User{ID: uint(userID)}, nil)
				mockArticleService.EXPECT().GetByID(uint(articleID)).Return(existingArticle, nil)
				mockArticleService.EXPECT().AddFavorite(existingArticle, gomock.Any()).Return(nil)
				mockUserService.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(true, nil)
			},
			assertion: func(resp *pb.ArticleResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if resp.Article == nil || !resp.Article.Favorited {
					t.Fatalf("expected article to be favorited, got %v", resp)
				}
			},
		},
		{
			name: "Error While Adding a Favorite",
			ctx:  validContext,
			req:  &pb.FavoriteArticleRequest{Slug: strconv.Itoa(articleID)},
			setupMocks: func() {
				mockUserService.EXPECT().GetByID(userID).Return(&model.User{ID: uint(userID)}, nil)
				mockArticleService.EXPECT().GetByID(uint(articleID)).Return(existingArticle, nil)
				mockArticleService.EXPECT().AddFavorite(existingArticle, gomock.Any()).Return(fmt.Errorf("database error"))
			},
			assertion: func(resp *pb.ArticleResponse, err error) {
				if status.Code(err) != codes.InvalidArgument {
					t.Fatalf("expected InvalidArgument error, got %v", err)
				}
			},
		},
		{
			name: "Error Determining Following Status",
			ctx:  validContext,
			req:  &pb.FavoriteArticleRequest{Slug: strconv.Itoa(articleID)},
			setupMocks: func() {
				mockUserService.EXPECT().GetByID(userID).Return(&model.User{ID: uint(userID)}, nil)
				mockArticleService.EXPECT().GetByID(uint(articleID)).Return(existingArticle, nil)
				mockArticleService.EXPECT().AddFavorite(existingArticle, gomock.Any()).Return(nil)
				mockUserService.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(false, fmt.Errorf("error fetching follow status"))
			},
			assertion: func(resp *pb.ArticleResponse, err error) {
				if status.Code(err) != codes.NotFound {
					t.Fatalf("expected NotFound error, got %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			resp, err := h.FavoriteArticle(tt.ctx, tt.req)
			tt.assertion(resp, err)
			t.Logf("Scenario '%s' finished", tt.name)
		})
	}
}



