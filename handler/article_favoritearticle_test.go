package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandlerFavoriteArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer mockDB.Close()

	mockUserStore := store.NewMockUserStore(ctrl)
	mockArticleStore := store.NewMockArticleStore(ctrl)
	h := &Handler{
		logger: &zerolog.Logger{},
		us:     mockUserStore,
		as:     mockArticleStore,
	}

	tests := []struct {
		name      string
		prepare   func()
		req       *pb.FavoriteArticleRequest
		wantErr   codes.Code
		wantResponse *pb.ArticleResponse
	}{
		{
			name: "Successful Article Favorite",
			prepare: func() {
				mockUser := &model.User{ID: 1}
				mockArticle := &model.Article{ID: 1, Title: "Title"}

				auth.GetUserID = func(ctx context.Context) (uint, error) { return 1, nil } // mock GetUserID
				mockUserStore.EXPECT().GetByID(uint(1)).Return(mockUser, nil)
				mockArticleStore.EXPECT().GetByID(uint(1)).Return(mockArticle, nil)
				mockArticleStore.EXPECT().AddFavorite(mockArticle, mockUser).Return(nil)
				mockUserStore.EXPECT().IsFollowing(mockUser, &mockArticle.Author).Return(false, nil)
			},
			req: &pb.FavoriteArticleRequest{Slug: "1"},
			wantErr: codes.OK,
			wantResponse: &pb.ArticleResponse{
				Article: &pb.Article{
					Slug:           "1",
					Title:          "Title",
					Favorited:      true,
					FavoritesCount: 1,
					// Populate other fields as needed
				},
			},
		},
		{
			name: "Unauthenticated User",
			prepare: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 0, errors.New("unauthenticated user")
				}
			},
			req:       &pb.FavoriteArticleRequest{Slug: "1"},
			wantErr:   codes.Unauthenticated,
			wantResponse: nil,
		},
		{
			name: "Non-existent User",
			prepare: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) { return 1, nil }
				mockUserStore.EXPECT().GetByID(uint(1)).Return(nil, errors.New("user not found"))
			},
			req:       &pb.FavoriteArticleRequest{Slug: "1"},
			wantErr:   codes.NotFound,
			wantResponse: nil,
		},
		{
			name: "Invalid Slug Format",
			prepare: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) { return 1, nil }
			},
			req:       &pb.FavoriteArticleRequest{Slug: "abc"},
			wantErr:   codes.InvalidArgument,
			wantResponse: nil,
		},
		{
			name: "Article Not Found",
			prepare: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) { return 1, nil }
				mockUserStore.EXPECT().GetByID(uint(1)).Return(&model.User{ID: 1}, nil)
				mockArticleStore.EXPECT().GetByID(uint(1)).Return(nil, errors.New("article not found"))
			},
			req:       &pb.FavoriteArticleRequest{Slug: "1"},
			wantErr:   codes.InvalidArgument,
			wantResponse: nil,
		},
		{
			name: "Failure to Add Favorite",
			prepare: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) { return 1, nil }
				mockUserStore.EXPECT().GetByID(uint(1)).Return(&model.User{ID: 1}, nil)
				mockArticleStore.EXPECT().GetByID(uint(1)).Return(&model.Article{ID: 1}, nil)
				mockArticleStore.EXPECT().AddFavorite(gomock.Any(), gomock.Any()).Return(errors.New("insertion error"))
			},
			req:       &pb.FavoriteArticleRequest{Slug: "1"},
			wantErr:   codes.InvalidArgument,
			wantResponse: nil,
		},
		{
			name: "Failed to Retrieve Following Status",
			prepare: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) { return 1, nil }
				mockUserStore.EXPECT().GetByID(uint(1)).Return(&model.User{ID: 1}, nil)
				mockArticleStore.EXPECT().GetByID(uint(1)).Return(&model.Article{ID: 1}, nil)
				mockArticleStore.EXPECT().AddFavorite(gomock.Any(), gomock.Any()).Return(nil)
				mockUserStore.EXPECT().IsFollowing(gomock.Any(), gomock.Any()).Return(false, errors.New("follow status error"))
			},
			req:       &pb.FavoriteArticleRequest{Slug: "1"},
			wantErr:   codes.NotFound,
			wantResponse: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			resp, err := h.FavoriteArticle(context.Background(), tt.req)

			if code := status.Code(err); code != tt.wantErr {
				t.Errorf("expected error code %v, got %v", tt.wantErr, code)
			}

			if tt.wantErr == codes.OK && !compareArticleResponse(resp, tt.wantResponse) {
				t.Errorf("expected response %+v, got %+v", tt.wantResponse, resp)
			}
		})
	}
}

// Helper function to compare ArticleResponse objects
func compareArticleResponse(a, b *pb.ArticleResponse) bool {
	// Simplistic comparison logic; extend as necessary for full field comparison
	return a.Article.Slug == b.Article.Slug
}
