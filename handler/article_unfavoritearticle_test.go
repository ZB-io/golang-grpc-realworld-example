package handler

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/stretchr/testify/mock"
)

// Mocking the UserStore and ArticleStore using testify/mock
type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) IsFollowing(a *model.User, b *model.User) (bool, error) {
	args := m.Called(a, b)
	return args.Bool(0), args.Error(1)
}

type MockArticleStore struct {
	mock.Mock
}

func (m *MockArticleStore) GetByID(id uint) (*model.Article, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Article), args.Error(1)
}

func (m *MockArticleStore) DeleteFavorite(a *model.Article, u *model.User) error {
	args := m.Called(a, u)
	return args.Error(0)
}

func TestHandlerUnfavoriteArticle(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *proto.UnfavoriteArticleRequest
	}
	tests := []struct {
		name       string
		setupMocks func(us *MockUserStore, as *MockArticleStore, logger *zerolog.Logger)
		args       args
		want       *proto.ArticleResponse
		wantErr    error
	}{
		{
			name: "Scenario 1: Successfully Unfavoriting an Article",
			setupMocks: func(us *MockUserStore, as *MockArticleStore, logger *zerolog.Logger) {
				// Mock user retrieval
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)

				// Mock article retrieval
				as.On("GetByID", uint(2)).Return(&model.Article{
					ID:             2,
					FavoritedUsers: []model.User{{ID: 1}},
				}, nil)

				// Mock delete favorite
				as.On("DeleteFavorite", &model.Article{ID: 2}, &model.User{ID: 1}).Return(nil)
			},
			args: args{
				ctx: context.WithValue(context.Background(), auth.UserKey, uint(1)),
				req: &proto.UnfavoriteArticleRequest{Slug: "2"},
			},
			want: &proto.ArticleResponse{Article: &proto.Article{Favorited: false, FavoritesCount: 0}},
			wantErr: nil,
		},
		{
			name: "Scenario 2: Handling an Unauthenticated User",
			setupMocks: func(us *MockUserStore, as *MockArticleStore, logger *zerolog.Logger) {
				logger.Error().Err(errors.New("unauthenticated")).Msg("unauthenticated")
			},
			args: args{
				ctx: context.Background(),
				req: &proto.UnfavoriteArticleRequest{Slug: "2"},
			},
			want:    nil,
			wantErr: status.Errorf(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Scenario 3: User Not Found in Database",
			setupMocks: func(us *MockUserStore, as *MockArticleStore, logger *zerolog.Logger) {
				us.On("GetByID", uint(1)).Return(nil, sql.ErrNoRows)
				logger.Error().Err(sql.ErrNoRows).Msg("not user found")
			},
			args: args{
				ctx: context.WithValue(context.Background(), auth.UserKey, uint(1)),
				req: &proto.UnfavoriteArticleRequest{Slug: "2"},
			},
			want:    nil,
			wantErr: status.Error(codes.NotFound, "not user found"),
		},
		{
			name: "Scenario 4: Slug Conversion Error",
			setupMocks: func(us *MockUserStore, as *MockArticleStore, logger *zerolog.Logger) {
				logger.Error().Err(errors.New("conversion error")).Msg("conversion error")
			},
			args: args{
				ctx: context.WithValue(context.Background(), auth.UserKey, uint(1)),
				req: &proto.UnfavoriteArticleRequest{Slug: "invalid_slug"},
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name: "Scenario 5: Article Not Found with Valid Slug",
			setupMocks: func(us *MockUserStore, as *MockArticleStore, logger *zerolog.Logger) {
				as.On("GetByID", uint(2)).Return(nil, sql.ErrNoRows)
				logger.Error().Err(sql.ErrNoRows).Msg("requested article not found")
			},
			args: args{
				ctx: context.WithValue(context.Background(), auth.UserKey, uint(1)),
				req: &proto.UnfavoriteArticleRequest{Slug: "2"},
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name: "Scenario 6: Failure to Remove Favorite Status",
			setupMocks: func(us *MockUserStore, as *MockArticleStore, logger *zerolog.Logger) {
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", uint(2)).Return(&model.Article{ID: 2}, nil)
				as.On("DeleteFavorite", &model.Article{ID: 2}, &model.User{ID: 1}).Return(errors.New("delete error"))
				logger.Error().Err(errors.New("delete error")).Msg("failed to remove favorite")
			},
			args: args{
				ctx: context.WithValue(context.Background(), auth.UserKey, uint(1)),
				req: &proto.UnfavoriteArticleRequest{Slug: "2"},
			},
			want: nil,
			wantErr: status.Error(codes.InvalidArgument, "failed to remove favorite"),
		},
		{
			name: "Scenario 7: Failure to Check Following Status",
			setupMocks: func(us *MockUserStore, as *MockArticleStore, logger *zerolog.Logger) {
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", uint(2)).Return(&model.Article{
					ID:    2,
					Author: model.User{ID: 3},
				}, nil)
				us.On("IsFollowing", &model.User{ID: 1}, &model.User{ID: 3}).Return(false, errors.New("isFollowing error"))
				logger.Error().Err(errors.New("isFollowing error")).Msg("failed to get following status")
			},
			args: args{
				ctx: context.WithValue(context.Background(), auth.UserKey, uint(1)),
				req: &proto.UnfavoriteArticleRequest{Slug: "2"},
			},
			want: nil,
			wantErr: status.Error(codes.NotFound, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the Logger, UserStore, and ArticleStore mocks
			us := new(MockUserStore)
			as := new(MockArticleStore)
			logger := zerolog.New(nil)
			h := &Handler{logger: &logger, us: us, as: as}

			tt.setupMocks(us, as, &logger)

			got, err := h.UnfavoriteArticle(tt.args.ctx, tt.args.req)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UnfavoriteArticle() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !equalArticleResponse(got, tt.want) {
				t.Errorf("UnfavoriteArticle() got = %v, want %v", got, tt.want)
			}

			// Log the detailed test results for clarity
			if err != nil {
				t.Log("Test failed with error: ", err)
			} else {
				t.Logf("Test successful, response: %+v", got)
			}
		})
	}
}

// equalArticleResponse provides utility to compare expected and received ArticleResponse
func equalArticleResponse(a, b *proto.ArticleResponse) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Article.Favorited == b.Article.Favorited && a.Article.FavoritesCount == b.Article.FavoritesCount
}
