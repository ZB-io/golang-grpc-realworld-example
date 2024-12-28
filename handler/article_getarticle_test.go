package handler

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"errors"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



type MockArticleService struct {
	mock.Mock
}



func (m *MockArticleService) GetByID(id uint) (*model.Article, error) {
	args := m.Called(id)
	article, _ := args.Get(0).(*model.Article)
	return article, args.Error(1)
}
func (m *MockArticleService) IsFavorited(article *model.Article, user *model.User) (bool, error) {
	args := m.Called(article, user)
	return args.Bool(0), args.Error(1)
}
func (m *MockUserService) IsFollowing(user *model.User, author *model.User) (bool, error) {
	args := m.Called(user, author)
	return args.Bool(0), args.Error(1)
}
func TestGetArticle(t *testing.T) {
	mockArticleService := new(MockArticleService)
	mockUserService := new(MockUserService)
	h := &Handler{
		as: mockArticleService,
		us: mockUserService,
	}

	tests := []struct {
		name            string
		slug            string
		userID          interface{}
		mockSetup       func()
		expectedError   codes.Code
		expectedArticle *pb.Article
	}{
		{
			name:   "Scenario 1: Valid Article Retrieval Without User Context",
			slug:   "1",
			userID: nil,
			mockSetup: func() {
				mockArticleService.On("GetByID", uint(1)).Return(&model.Article{}, nil).Once()
			},
			expectedError:   codes.OK,
			expectedArticle: &pb.Article{},
		},
		{
			name:   "Scenario 2: Valid Article Retrieval With User Context",
			slug:   "1",
			userID: 1,
			mockSetup: func() {
				mockArticleService.On("GetByID", uint(1)).Return(&model.Article{}, nil).Once()
				mockUserService.On("GetByID", uint(1)).Return(&model.User{}, nil).Once()
				mockArticleService.On("IsFavorited", mock.Anything, mock.Anything).Return(true, nil).Once()
				mockUserService.On("IsFollowing", mock.Anything, mock.Anything).Return(true, nil).Once()
			},
			expectedError:   codes.OK,
			expectedArticle: &pb.Article{},
		},
		{
			name:            "Scenario 3: Invalid Slug Format Leading to Error",
			slug:            "invalidSlug",
			userID:          nil,
			mockSetup:       func() {},
			expectedError:   codes.InvalidArgument,
			expectedArticle: nil,
		},
		{
			name:   "Scenario 4: Nonexistent Article Slug",
			slug:   "9999",
			userID: nil,
			mockSetup: func() {
				mockArticleService.On("GetByID", uint(9999)).Return(nil, errors.New("article not found")).Once()
			},
			expectedError:   codes.InvalidArgument,
			expectedArticle: nil,
		},
		{
			name:   "Scenario 5: Authenticated User Not Found",
			slug:   "1",
			userID: 1,
			mockSetup: func() {
				mockArticleService.On("GetByID", uint(1)).Return(&model.Article{}, nil).Once()
				mockUserService.On("GetByID", uint(1)).Return(nil, errors.New("user not found")).Once()
			},
			expectedError:   codes.NotFound,
			expectedArticle: nil,
		},
		{
			name:   "Scenario 6: Failure in Checking Favorited Status",
			slug:   "1",
			userID: 1,
			mockSetup: func() {
				mockArticleService.On("GetByID", uint(1)).Return(&model.Article{}, nil).Once()
				mockUserService.On("GetByID", uint(1)).Return(&model.User{}, nil).Once()
				mockArticleService.On("IsFavorited", mock.Anything, mock.Anything).Return(false, errors.New("error checking favorited status")).Once()
			},
			expectedError:   codes.Aborted,
			expectedArticle: nil,
		},
		{
			name:   "Scenario 7: Error on Following Status Check",
			slug:   "1",
			userID: 1,
			mockSetup: func() {
				mockArticleService.On("GetByID", uint(1)).Return(&model.Article{}, nil).Once()
				mockUserService.On("GetByID", uint(1)).Return(&model.User{}, nil).Once()
				mockArticleService.On("IsFavorited", mock.Anything, mock.Anything).Return(true, nil).Once()
				mockUserService.On("IsFollowing", mock.Anything, mock.Anything).Return(false, errors.New("error checking following status")).Once()
			},
			expectedError:   codes.NotFound,
			expectedArticle: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			ctx := context.Background()
			if tt.userID != nil {
				ctx = auth.NewContext(ctx, int(tt.userID.(int)))
			}
			req := &pb.GetArticleRequest{Slug: tt.slug}

			resp, err := h.GetArticle(ctx, req)

			if tt.expectedError == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedArticle, resp.Article)
			} else {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				if ok {
					assert.Equal(t, tt.expectedError, st.Code())
				}
				assert.Nil(t, resp)
			}
		})
	}
}


