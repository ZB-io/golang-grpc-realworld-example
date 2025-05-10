package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Mock stores
type UserStoreMock struct {
	mock.Mock
	store.UserStore
}

type ArticleStoreMock struct {
	mock.Mock
	store.ArticleStore
}

// Mock function implementations
func (m *UserStoreMock) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserStoreMock) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserStoreMock) IsFollowing(a *model.User, b *model.User) (bool, error) {
	args := m.Called(a, b)
	return args.Bool(0), args.Error(1)
}

func (m *ArticleStoreMock) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([]model.Article, error) {
	args := m.Called(tagName, username, favoritedBy, limit, offset)
	return args.Get(0).([]model.Article), args.Error(1)
}

func (m *ArticleStoreMock) IsFavorited(a *model.Article, u *model.User) (bool, error) {
	args := m.Called(a, u)
	return args.Bool(0), args.Error(1)
}

func TestGetArticles(t *testing.T) {
	mockUserStore := new(UserStoreMock)
	mockArticleStore := new(ArticleStoreMock)
	logger := zerolog.New(zerolog.ConsoleWriter{Out: t.Log()}) // Properly initialize zerolog.Logger
	handler := &Handler{logger: &logger, us: mockUserStore, as: mockArticleStore}

	// Define test cases
	testCases := []struct {
		name     string
		request  *proto.GetArticlesRequest
		setup    func()
		expected func(t *testing.T, resp *proto.ArticlesResponse, err error)
	}{
		{
			name: "Ensure Limits Default Correctly to 20",
			request: &proto.GetArticlesRequest{
				Limit: 0,
			},
			setup: func() {
				mockArticleStore.On("GetArticles", "", "", (*model.User)(nil), int64(20), int64(0)).
					Return([]model.Article{model.Article{}, model.Article{}, model.Article{}, model.Article{}}, nil)
			},
			expected: func(t *testing.T, resp *proto.ArticlesResponse, err error) {
				assert.NoError(t, err)
				assert.LessOrEqual(t, resp.ArticlesCount, int32(20))
			},
		},
		{
			name: "Retrieving Favorited Articles by User",
			request: &proto.GetArticlesRequest{
				Favorited: "user1",
			},
			setup: func() {
				mockUserStore.On("GetByUsername", "user1").
					Return(&model.User{UserID: 1}, nil) // Correct field
				mockArticleStore.On("GetArticles", "", "", &model.User{UserID: 1}, int64(20), int64(0)).
					Return([]model.Article{model.Article{}}, nil)
			},
			expected: func(t *testing.T, resp *proto.ArticlesResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int32(1), resp.ArticlesCount)
			},
		},
		{
			name: "Handling Non-existent 'Favorited' User",
			request: &proto.GetArticlesRequest{
				Favorited: "nonexistent",
			},
			setup: func() {
				mockUserStore.On("GetByUsername", "nonexistent").
					Return((*model.User)(nil), errors.New("user not found"))
				mockArticleStore.On("GetArticles", "", "", (*model.User)(nil), int64(20), int64(0)).
					Return([]model.Article{}, nil)
			},
			expected: func(t *testing.T, resp *proto.ArticlesResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int32(0), resp.ArticlesCount)
			},
		},
		{
			name: "Handling Database Error During Article Retrieval",
			request: &proto.GetArticlesRequest{
				Limit: 10,
			},
			setup: func() {
				mockArticleStore.On("GetArticles", "", "", (*model.User)(nil), int64(10), int64(0)).
					Return([]model.Article{}, errors.New("db error"))
			},
			expected: func(t *testing.T, resp *proto.ArticlesResponse, err error) {
				assert.Error(t, err)
				assert.Equal(t, codes.Aborted, status.Code(err))
			},
		},
		{
			name: "User Not Found for Favorited Articles Processing",
			request: &proto.GetArticlesRequest{
				Favorited: "someuser",
			},
			setup: func() {
				mockUserStore.On("GetByUsername", "someuser").
					Return(&model.User{UserID: 1}, nil) // Correct field
				mockUserStore.On("GetByID", uint(123)).
					Return((*model.User)(nil), errors.New("not found"))
			},
			expected: func(t *testing.T, resp *proto.ArticlesResponse, err error) {
				assert.Error(t, err)
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Verify Function Behavior with Valid Inputs",
			request: &proto.GetArticlesRequest{
				Limit:     5,
				Favorited: "user1",
			},
			setup: func() {
				mockUserStore.On("GetByUsername", "user1").
					Return(&model.User{UserID: 1}, nil) // Correct field
				mockArticleStore.On("GetArticles", "", "", &model.User{UserID: 1}, int64(5), int64(0)).
					Return([]model.Article{model.Article{}, model.Article{}}, nil)
			},
			expected: func(t *testing.T, resp *proto.ArticlesResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int32(2), resp.ArticlesCount)
			},
		},
		{
			name: "Confirm Article Favoriting Status is Handled Properly",
			request: &proto.GetArticlesRequest{
				Limit:     5,
				Favorited: "user1",
			},
			setup: func() {
				mockUserStore.On("GetByUsername", "user1").
					Return(&model.User{UserID: 1}, nil) // Correct field
				article := model.Article{}
				mockArticleStore.On("GetArticles", "", "", &model.User{UserID: 1}, int64(5), int64(0)).
					Return([]model.Article{article}, nil)
				mockArticleStore.On("IsFavorited", &article, mock.Anything).
					Return(true, nil)
			},
			expected: func(t *testing.T, resp *proto.ArticlesResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int32(1), resp.ArticlesCount)
			},
		},
		{
			name: "Check Following Status is Accurately Retrieved",
			request: &proto.GetArticlesRequest{
				Limit: 5,
			},
			setup: func() {
				article := model.Article{Author: model.User{UserID: 1}} // Correct field
				mockArticleStore.On("GetArticles", "", "", (*model.User)(nil), int64(5), int64(0)).
					Return([]model.Article{article}, nil)
				mockUserStore.On("IsFollowing", mock.Anything, &article.Author).
					Return(true, nil)
			},
			expected: func(t *testing.T, resp *proto.ArticlesResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int32(1), resp.ArticlesCount)
				assert.True(t, resp.Articles[0].Author.Following)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			resp, err := handler.GetArticles(context.Background(), tc.request)
			tc.expected(t, resp, err)
		})
	}
}
