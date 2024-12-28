package handler

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetByID(userID int) (*model.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) IsFollowing(currentUser *model.User, author *model.User) (bool, error) {
	args := m.Called(currentUser, author)
	return args.Bool(0), args.Error(1)
}

type MockArticleService struct {
	mock.Mock
}

func (m *MockArticleService) Create(article *model.Article) error {
	args := m.Called(article)
	return args.Error(0)
}

func TestCreateArticle(t *testing.T) {
	tests := []struct {
		name          string
		prepareMock   func(us *MockUserService, as *MockArticleService)
		setupContext  func() context.Context
		request       *pb.CreateAritcleRequest
		expectedError error
	}{
		{
			name: "Scenario 1 - Successfully Create an Article",
			prepareMock: func(us *MockUserService, as *MockArticleService) {
				us.On("GetByID", 1).Return(&model.User{ID: 1}, nil)
				as.On("Create", mock.AnythingOfType("*model.Article")).Return(nil)
			},
			setupContext: func() context.Context {
				ctx := context.Background()
				ctx = auth.SetUserID(ctx, 1)
				return ctx
			},
			request: &pb.CreateAritcleRequest{
				Article: &pb.Article{
					Title:       "Test Title",
					Description: "Test Description",
					Body:        "Test Body",
					TagList:     []string{"tag1", "tag2"},
				},
			},
			expectedError: nil,
		},
		{
			name: "Scenario 2 - Missing User Authentication",
			prepareMock: func(us *MockUserService, as *MockArticleService) {},
			setupContext: func() context.Context {
				return context.Background()
			},
			request: &pb.CreateAritcleRequest{
				Article: &pb.Article{Title: "Test Title"},
			},
			expectedError: status.Errorf(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Scenario 3 - User Not Found",
			prepareMock: func(us *MockUserService, as *MockArticleService) {
				us.On("GetByID", 999).Return(nil, status.Error(codes.NotFound, "user not found"))
			},
			setupContext: func() context.Context {
				ctx := context.Background()
				ctx = auth.SetUserID(ctx, 999)
				return ctx
			},
			request: &pb.CreateAritcleRequest{
				Article: &pb.Article{Title: "Test Title"},
			},
			expectedError: status.Error(codes.NotFound, "user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := new(MockUserService)
			as := new(MockArticleService)
			tt.prepareMock(us, as)

			h := &Handler{
				us: us,
				as: as,
			}

			ctx := tt.setupContext()
			resp, err := h.CreateArticle(ctx, tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, status.Convert(err).Err())
				t.Log("Expected error:", err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				t.Log("Article created successfully:", resp)
			}

			us.AssertExpectations(t)
			as.AssertExpectations(t)
		})
	}
}
