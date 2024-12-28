package handler

import (
	"context"
	"errors"
	"testing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"fmt"
	"strconv"
	"github.com/raahii/golang-grpc-realworld-example/auth"
)



type MockUserService struct {
	mock.Mock
}




func (m *MockUserService) GetByID(userID uint) (*model.User, error) {
	args := m.Called(userID)
	user, _ := args.Get(0).(*model.User)
	return user, args.Error(1)
}
func (m *MockArticleService) GetFeedArticles(userIDs []uint, limit, offset int) ([]model.Article, error) {
	args := m.Called(userIDs, limit, offset)
	return args.Get(0).([]model.Article), args.Error(1)
}
func (m *MockUserService) GetFollowingUserIDs(user *model.User) ([]uint, error) {
	args := m.Called(user)
	return args.Get(0).([]uint), args.Error(1)
}
func (m *MockArticleService) IsFavorited(article *model.Article, user *model.User) (bool, error) {
	args := m.Called(article, user)
	return args.Bool(0), args.Error(1)
}
func (m *MockUserService) IsFollowing(user *model.User, author *model.User) (bool, error) {
	args := m.Called(user, author)
	return args.Bool(0), args.Error(1)
}
func TestGetFeedArticles(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.GetFeedArticlesRequest
	}

	mockUserService := new(MockUserService)
	mockArticleService := new(MockArticleService)
	handler := &Handler{
		us: mockUserService,
		as: mockArticleService,
	}

	currentUser := &model.User{ID: 1}

	tests := []struct {
		name          string
		args          args
		setupMocks    func()
		expectedError codes.Code
	}{
		{
			name: "Successfully Retrieve Feed Articles for an Authenticated User",
			args: args{
				ctx: context.WithValue(context.Background(), "userID", currentUser.ID),
				req: &pb.GetFeedArticlesRequest{Limit: 10, Offset: 0},
			},
			setupMocks: func() {
				mockUserService.On("GetByID", currentUser.ID).Return(currentUser, nil)
				mockUserService.On("GetFollowingUserIDs", currentUser).Return([]uint{2, 3}, nil)
				mockArticleService.On("GetFeedArticles", []uint{2, 3}, 10, 0).Return([]model.Article{
					{Title: "Test Article", Author: *currentUser},
				}, nil)
				mockArticleService.On("IsFavorited", mock.Anything, currentUser).Return(true, nil)
				mockUserService.On("IsFollowing", currentUser, &currentUser).Return(true, nil)
			},
			expectedError: codes.OK,
		},
		{
			name: "Unauthenticated User Attempts to Retrieve Feed Articles",
			args: args{
				ctx: context.Background(),
				req: &pb.GetFeedArticlesRequest{Limit: 10, Offset: 0},
			},
			setupMocks:    func() {},
			expectedError: codes.Unauthenticated,
		},
		{
			name: "Handle Limit Default Behavior When Not Provided",
			args: args{
				ctx: context.WithValue(context.Background(), "userID", currentUser.ID),
				req: &pb.GetFeedArticlesRequest{Offset: 0},
			},
			setupMocks: func() {
				mockUserService.On("GetByID", currentUser.ID).Return(currentUser, nil)
				mockUserService.On("GetFollowingUserIDs", currentUser).Return([]uint{2, 3}, nil)
				mockArticleService.On("GetFeedArticles", []uint{2, 3}, 20, 0).Return([]model.Article{
					{Title: "Test Article", Author: *currentUser},
				}, nil)
				mockArticleService.On("IsFavorited", mock.Anything, currentUser).Return(true, nil)
				mockUserService.On("IsFollowing", currentUser, &currentUser).Return(true, nil)
			},
			expectedError: codes.OK,
		},
		{
			name: "Retrieve Articles with Favorited Status and Following Status",
			args: args{
				ctx: context.WithValue(context.Background(), "userID", currentUser.ID),
				req: &pb.GetFeedArticlesRequest{Limit: 10, Offset: 0},
			},
			setupMocks: func() {
				author := model.User{ID: 2}
				mockUserService.On("GetByID", currentUser.ID).Return(currentUser, nil)
				mockUserService.On("GetFollowingUserIDs", currentUser).Return([]uint{author.ID}, nil)
				mockArticleService.On("GetFeedArticles", []uint{author.ID}, 10, 0).Return([]model.Article{
					{Title: "Test Article", Author: author},
				}, nil)
				mockArticleService.On("IsFavorited", mock.Anything, currentUser).Return(true, nil)
				mockUserService.On("IsFollowing", currentUser, &author).Return(true, nil)
			},
			expectedError: codes.OK,
		},
		{
			name: "Handle Internal Server Error When Retrieving User Information",
			args: args{
				ctx: context.WithValue(context.Background(), "userID", currentUser.ID),
				req: &pb.GetFeedArticlesRequest{Limit: 10, Offset: 0},
			},
			setupMocks: func() {
				mockUserService.On("GetByID", currentUser.ID).Return(nil, errors.New("user not found"))
			},
			expectedError: codes.NotFound,
		},
		{
			name: "Handling No Follow Network for Current User",
			args: args{
				ctx: context.WithValue(context.Background(), "userID", currentUser.ID),
				req: &pb.GetFeedArticlesRequest{Limit: 10, Offset: 0},
			},
			setupMocks: func() {
				mockUserService.On("GetByID", currentUser.ID).Return(currentUser, nil)
				mockUserService.On("GetFollowingUserIDs", currentUser).Return([]uint{}, nil)
				mockArticleService.On("GetFeedArticles", []uint{}, 10, 0).Return([]model.Article{}, nil)
			},
			expectedError: codes.OK,
		},
		{
			name: "Offset Handling for Large Data Set",
			args: args{
				ctx: context.WithValue(context.Background(), "userID", currentUser.ID),
				req: &pb.GetFeedArticlesRequest{Limit: 10, Offset: 50},
			},
			setupMocks: func() {
				mockUserService.On("GetByID", currentUser.ID).Return(currentUser, nil)
				mockUserService.On("GetFollowingUserIDs", currentUser).Return([]uint{2, 3}, nil)
				mockArticleService.On("GetFeedArticles", []uint{2, 3}, 10, 50).Return([]model.Article{
					{Title: "Test Article", Author: *currentUser},
				}, nil)
				mockArticleService.On("IsFavorited", mock.Anything, currentUser).Return(true, nil)
				mockUserService.On("IsFollowing", currentUser, &currentUser).Return(true, nil)
			},
			expectedError: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			resp, err := handler.GetFeedArticles(tt.args.ctx, tt.args.req)
			if tt.expectedError == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				t.Logf("Expected articles count: %d", resp.ArticlesCount)
			} else {
				assert.Error(t, err)
				assert.Nil(t, resp)
				st, _ := status.FromError(err)
				assert.Equal(t, tt.expectedError, st.Code())
				t.Logf("Expected error code: %v", tt.expectedError)
			}
		})
	}
}



func (m *MockArticleService) GetFeedArticles(userIDs []uint, limit, offset int) ([]model.Article, error) {
	args := m.Called(userIDs, limit, offset)
	return args.Get(0).([]model.Article), args.Error(1)
}
