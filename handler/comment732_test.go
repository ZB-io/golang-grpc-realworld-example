package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type MockArticleStore struct {
	mock.Mock
}

type MockUserStore struct {
	mock.Mock
}

func (m *MockArticleStore) GetByID(id uint) (*model.Article, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Article), args.Error(1)
}

func (m *MockArticleStore) GetComments(article *model.Article) ([]model.Comment, error) {
	args := m.Called(article)
	return args.Get(0).([]model.Comment), args.Error(1)
}

func (m *MockArticleStore) CreateComment(comment *model.Comment) error {
	args := m.Called(comment)
	return args.Error(0)
}

func (m *MockArticleStore) DeleteComment(comment *model.Comment) error {
	args := m.Called(comment)
	return args.Error(0)
}

func (m *MockArticleStore) GetCommentByID(id uint) (*model.Comment, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Comment), args.Error(1)
}

func (m *MockUserStore) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) IsFollowing(follower, followed *model.User) (bool, error) {
	args := m.Called(follower, followed)
	return args.Bool(0), args.Error(1)
}

func MockAuthContext(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, "user_id", userID)
}

/*
ROOST_METHOD_HASH=DeleteComment_452af2f984
ROOST_METHOD_SIG_HASH=DeleteComment_27615e7d69
*/
func TestDeleteComment(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*MockUserStore, *MockArticleStore)
		req            *pb.DeleteCommentRequest
		expectedError  error
		expectedStatus codes.Code
	}{
		{
			name: "Successfully Delete a Comment",
			setupMocks: func(us *MockUserStore, as *MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetCommentByID", uint(1)).Return(&model.Comment{UserID: 1, ArticleID: 1}, nil)
				as.On("DeleteComment", mock.AnythingOfType("*model.Comment")).Return(nil)
			},
			req: &pb.DeleteCommentRequest{
				Slug: "1",
				Id:   "1",
			},
			expectedError:  nil,
			expectedStatus: codes.OK,
		},
		{
			name: "Attempt to Delete a Comment with Unauthenticated User",
			setupMocks: func(us *MockUserStore, as *MockArticleStore) {
			},
			req: &pb.DeleteCommentRequest{
				Slug: "1",
				Id:   "1",
			},
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
			expectedStatus: codes.Unauthenticated,
		},
		{
			name: "Attempt to Delete a Non-existent Comment",
			setupMocks: func(us *MockUserStore, as *MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{}, nil)
				as.On("GetCommentByID", uint(999)).Return(nil, errors.New("comment not found"))
			},
			req: &pb.DeleteCommentRequest{
				Slug: "1",
				Id:   "999",
			},
			expectedError:  status.Error(codes.InvalidArgument, "failed to get comment"),
			expectedStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUS := &MockUserStore{}
			mockAS := &MockArticleStore{}
			tt.setupMocks(mockUS, mockAS)

			h := &Handler{
				logger: zerolog.New(zerolog.NewConsoleWriter()),
				us:     mockUS,
				as:     mockAS,
			}

			ctx := context.Background()
			if tt.name != "Attempt to Delete a Comment with Unauthenticated User" {
				ctx = MockAuthContext(ctx, 1)
			}

			_, err := h.DeleteComment(ctx, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				statusErr, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, statusErr.Code())
			} else {
				assert.NoError(t, err)
			}

			mockUS.AssertExpectations(t)
			mockAS.AssertExpectations(t)
		})
	}
}

/*
ROOST_METHOD_HASH=GetComments_265127fb6a
ROOST_METHOD_SIG_HASH=GetComments_20efd5abae
*/
func TestGetComments(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*MockArticleStore, *MockUserStore)
		req            *pb.GetCommentsRequest
		ctx            context.Context
		expectedResult *pb.CommentsResponse
		expectedError  error
	}{
		{
			name: "Successfully retrieve comments for a valid article",
			setupMocks: func(mas *MockArticleStore, mus *MockUserStore) {
				mas.On("GetByID", uint(1)).Return(&model.Article{Model: gorm.Model{ID: 1}}, nil)
				mas.On("GetComments", &model.Article{Model: gorm.Model{ID: 1}}).Return([]model.Comment{
					{Model: gorm.Model{ID: 1}, Body: "Comment 1", Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
					{Model: gorm.Model{ID: 2}, Body: "Comment 2", Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
				}, nil)
				mus.On("GetByID", uint(1)).Return(&model.User{Model: gorm.Model{ID: 1}}, nil)
				mus.On("IsFollowing", &model.User{Model: gorm.Model{ID: 1}}, &model.User{Model: gorm.Model{ID: 1}, Username: "user1"}).Return(false, nil)
				mus.On("IsFollowing", &model.User{Model: gorm.Model{ID: 1}}, &model.User{Model: gorm.Model{ID: 2}, Username: "user2"}).Return(true, nil)
			},
			req: &pb.GetCommentsRequest{Slug: "1"},
			ctx: MockAuthContext(context.Background(), uint(1)),
			expectedResult: &pb.CommentsResponse{
				Comments: []*pb.Comment{
					{Id: "1", Body: "Comment 1", Author: &pb.Profile{Username: "user1", Following: false}},
					{Id: "2", Body: "Comment 2", Author: &pb.Profile{Username: "user2", Following: true}},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArticleStore := &MockArticleStore{}
			mockUserStore := &MockUserStore{}
			tt.setupMocks(mockArticleStore, mockUserStore)

			h := &Handler{
				as: mockArticleStore,
				us: mockUserStore,
			}

			result, err := h.GetComments(tt.ctx, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockArticleStore.AssertExpectations(t)
			mockUserStore.AssertExpectations(t)
		})
	}
}

/*
ROOST_METHOD_HASH=CreateComment_c4ccd62dc5
ROOST_METHOD_SIG_HASH=CreateComment_19a3ee5a3b
*/
func TestCreateComment(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*MockUserStore, *MockArticleStore)
		req            *pb.CreateCommentRequest
		expectedResult *pb.CommentResponse
		expectedError  error
	}{
		{
			name: "Successfully Create a Comment",
			setupMocks: func(us *MockUserStore, as *MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{Model: gorm.Model{ID: 1}, Username: "testuser"}, nil)
				as.On("GetByID", uint(1)).Return(&model.Article{Model: gorm.Model{ID: 1}}, nil)
				as.On("CreateComment", mock.AnythingOfType("*model.Comment")).Return(nil)
			},
			req: &pb.CreateCommentRequest{
				Slug:    "1",
				Comment: &pb.CreateCommentRequest_Comment{Body: "Test comment"},
			},
			expectedResult: &pb.CommentResponse{
				Comment: &pb.Comment{
					Body: "Test comment",
					Author: &pb.Profile{
						Username: "testuser",
					},
				},
			},
			expectedError: nil,
		},
		{
			name:       "Attempt to Create Comment with Unauthenticated User",
			setupMocks: func(us *MockUserStore, as *MockArticleStore) {},
			req: &pb.CreateCommentRequest{
				Slug:    "1",
				Comment: &pb.CreateCommentRequest_Comment{Body: "Test comment"},
			},
			expectedResult: nil,
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Attempt to Create Comment for Non-existent Article",
			setupMocks: func(us *MockUserStore, as *MockArticleStore) {
				us.On("GetByID", uint(1)).Return(&model.User{Model: gorm.Model{ID: 1}, Username: "testuser"}, nil)
				as.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)
			},
			req: &pb.CreateCommentRequest{
				Slug:    "999",
				Comment: &pb.CreateCommentRequest_Comment{Body: "Test comment"},
			},
			expectedResult: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid article id"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := &MockUserStore{}
			mockArticleStore := &MockArticleStore{}

			tt.setupMocks(mockUserStore, mockArticleStore)

			h := &Handler{
				logger: zerolog.New(zerolog.NewConsoleWriter()),
				us:     mockUserStore,
				as:     mockArticleStore,
			}

			ctx := context.Background()
			if tt.name != "Attempt to Create Comment with Unauthenticated User" {
				ctx = MockAuthContext(ctx, 1)
			}

			result, err := h.CreateComment(ctx, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Comment.Body, result.Comment.Body)
				assert.Equal(t, tt.expectedResult.Comment.Author.Username, result.Comment.Author.Username)
			}

			mockUserStore.AssertExpectations(t)
			mockArticleStore.AssertExpectations(t)
		})
	}
}
