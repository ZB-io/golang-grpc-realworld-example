package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockUserStore struct {
	mock.Mock
}

func (m *mockUserStore) GetByID(userID uint) (*model.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserStore) IsFollowing(currentUser *model.User, author *model.User) (bool, error) {
	args := m.Called(currentUser, author)
	return args.Bool(0), args.Error(1)
}

type mockArticleStore struct {
	mock.Mock
}

func (m *mockArticleStore) Create(article *model.Article) error {
	args := m.Called(article)
	return args.Error(0)
}

func TestCreateArticle(t *testing.T) {
	// Define test scenarios
	tests := []struct {
		name               string
		context            context.Context
		request            *pb.CreateAritcleRequest
		mockUserStore      func(mu *mockUserStore)
		mockArticleStore   func(ma *mockArticleStore)
		expectedResponse   *pb.ArticleResponse
		expectedError      error
	}{
		{
			name:    "Successful Article Creation",
			context: createAuthenticatedContext(1),
			request: &pb.CreateAritcleRequest{
				Article: &pb.CreateAritcleRequest_Article{
					Title:       "Test Article",
					Description: "Description for test article",
					Body:        "Body of the test article",
					TagList:     []string{"tag1", "tag2"},
				},
			},
			mockUserStore: func(mu *mockUserStore) {
				mu.On("GetByID", uint(1)).Return(&model.User{ID: 1, Username: "testuser"}, nil)
				mu.On("IsFollowing", mock.Anything, mock.Anything).Return(false, nil)
			},
			mockArticleStore: func(ma *mockArticleStore) {
				ma.On("Create", mock.Anything).Return(nil)
			},
			expectedResponse: &pb.ArticleResponse{
				Article: &pb.Article{
					Title:       "Test Article",
					Description: "Description for test article",
					Body:        "Body of the test article",
					TagList:     []string{"tag1", "tag2"},
				},
			},
			expectedError: nil,
		},
		{
			name:    "Unauthenticated User",
			context: context.Background(), // No authenticated context
			request: &pb.CreateAritcleRequest{
				Article: &pb.CreateAritcleRequest_Article{},
			},
			mockUserStore:    func(mu *mockUserStore) {},
			mockArticleStore: func(ma *mockArticleStore) {},
			expectedResponse: nil,
			expectedError:    status.Errorf(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name:    "User Not Found",
			context: createAuthenticatedContext(2),
			request: &pb.CreateAritcleRequest{
				Article: &pb.CreateAritcleRequest_Article{},
			},
			mockUserStore: func(mu *mockUserStore) {
				mu.On("GetByID", uint(2)).Return(nil, errors.New("user not found"))
			},
			mockArticleStore: func(ma *mockArticleStore) {},
			expectedResponse: nil,
			expectedError:    status.Error(codes.NotFound, "user not found"),
		},
		{
			name:    "Article Validation Error",
			context: createAuthenticatedContext(1),
			request: &pb.CreateAritcleRequest{
				Article: &pb.CreateAritcleRequest_Article{
					Title: "",
				},
			},
			mockUserStore: func(mu *mockUserStore) {
				mu.On("GetByID", uint(1)).Return(&model.User{ID: 1, Username: "testuser"}, nil)
			},
			mockArticleStore: func(ma *mockArticleStore) {},
			expectedResponse: nil,
			expectedError:    status.Error(codes.InvalidArgument, "validation error"),
		},
		{
			name:    "Article Store Failure",
			context: createAuthenticatedContext(1),
			request: &pb.CreateAritcleRequest{
				Article: &pb.CreateAritcleRequest_Article{
					Title:       "Test Article",
					Description: "Description for test article",
					Body:        "Body of the test article",
					TagList:     []string{"tag1", "tag2"},
				},
			},
			mockUserStore: func(mu *mockUserStore) {
				mu.On("GetByID", uint(1)).Return(&model.User{ID: 1, Username: "testuser"}, nil)
			},
			mockArticleStore: func(ma *mockArticleStore) {
				ma.On("Create", mock.Anything).Return(errors.New("store failure"))
			},
			expectedResponse: nil,
			expectedError:    status.Error(codes.Canceled, "Failed to create user."),
		},
		{
			name:    "Check Following Status Error",
			context: createAuthenticatedContext(1),
			request: &pb.CreateAritcleRequest{
				Article: &pb.CreateAritcleRequest_Article{
					Title:       "Test Article",
					Description: "Description for test article",
					Body:        "Body of the test article",
					TagList:     []string{"tag1", "tag2"},
				},
			},
			mockUserStore: func(mu *mockUserStore) {
				mu.On("GetByID", uint(1)).Return(&model.User{ID: 1, Username: "testuser"}, nil)
				mu.On("IsFollowing", mock.Anything, mock.Anything).Return(false, errors.New("following check error"))
			},
			mockArticleStore: func(ma *mockArticleStore) {
				ma.On("Create", mock.Anything).Return(nil)
			},
			expectedResponse: nil,
			expectedError:    status.Error(codes.NotFound, "internal server error"),
		},
		{
			name:    "Tag List Handling",
			context: createAuthenticatedContext(1),
			request: &pb.CreateAritcleRequest{
				Article: &pb.CreateAritcleRequest_Article{
					Title:       "Test Article",
					Description: "Description for test article",
					Body:        "Body of the test article",
					TagList:     []string{}, // Empty tag list
				},
			},
			mockUserStore: func(mu *mockUserStore) {
				mu.On("GetByID", uint(1)).Return(&model.User{ID: 1, Username: "testuser"}, nil)
				mu.On("IsFollowing", mock.Anything, mock.Anything).Return(false, nil)
			},
			mockArticleStore: func(ma *mockArticleStore) {
				ma.On("Create", mock.Anything).Return(nil)
			},
			expectedResponse: &pb.ArticleResponse{
				Article: &pb.Article{
					Title:       "Test Article",
					Description: "Description for test article",
					Body:        "Body of the test article",
					TagList:     []string{},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := &mockUserStore{}
			ma := &mockArticleStore{}

			tt.mockUserStore(mu)
			tt.mockArticleStore(ma)

			h := &Handler{
				us: mu,
				as: ma,
			}

			resp, err := h.CreateArticle(tt.context, tt.request)

			if tt.expectedError != nil {
				assert.Nil(t, resp)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NotNil(t, resp)
				assert.True(t, proto.Equal(tt.expectedResponse, resp))
			}

			mu.AssertExpectations(t)
			ma.AssertExpectations(t)
		})
	}
}

func createAuthenticatedContext(userID uint) context.Context {
	md := metadata.New(map[string]string{"authorization": "Token " + strconv.Itoa(int(userID))})
	return metadata.NewIncomingContext(context.Background(), md)
}

// Note:
// - Unexported methods or private logic that require testing might need additional refactoring for testability.
// - Effective mocking is crucial in scenarios where external dependencies (e.g., databases, authentication services) are involved.
// - The use of `"github.com/stretchr/testify/mock"` and `"github.com/DATA-DOG/go-sqlmock"` may help in mocking database and other service interactions efficiently.
