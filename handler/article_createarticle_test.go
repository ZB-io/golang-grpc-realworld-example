package handler

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestHandlerCreateArticle tests the CreateArticle function
func TestHandlerCreateArticle(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T, us *store.UserStore, as *store.ArticleStore)
		request     *pb.CreateAritcleRequest
		expectError codes.Code
	}{
		{
			name: "Valid Article Creation",
			setup: func(t *testing.T, us *store.UserStore, as *store.ArticleStore) {
				user := &model.User{Model: gorm.Model{ID: 1}, Username: "test_user"}
				us.On("GetByID", uint(1)).Return(user, nil)
				as.On("Create", mock.AnythingOfType("*model.Article")).Return(nil)
			},
			request: &pb.CreateAritcleRequest{
				Article: &pb.CreateAritcleRequest_Article{
					Title:       "Test Title",
					Description: "Test Description",
					Body:        "Test Body",
					TagList:     []string{"tag1", "tag2"},
				},
			},
			expectError: codes.OK,
		},
		{
			name: "Unauthenticated User",
			setup: func(t *testing.T, us *store.UserStore, as *store.ArticleStore) {
				// No specific setup needed for this test
			},
			request: &pb.CreateAritcleRequest{
				Article: &pb.CreateAritcleRequest_Article{
					Title:       "Test Title",
					Description: "Test Description",
					Body:        "Test Body",
					TagList:     []string{"tag1", "tag2"},
				},
			},
			expectError: codes.Unauthenticated,
		},
		{
			name: "User Not Found in Store",
			setup: func(t *testing.T, us *store.UserStore, as *store.ArticleStore) {
				us.On("GetByID", uint(1)).Return(nil, status.Error(codes.NotFound, "user not found"))
			},
			request: &pb.CreateAritcleRequest{
				Article: &pb.CreateAritcleRequest_Article{
					Title:       "Test Title",
					Description: "Test Description",
					Body:        "Test Body",
					TagList:     []string{"tag1", "tag2"},
				},
			},
			expectError: codes.NotFound,
		},
		{
			name: "Article Validation Failure",
			setup: func(t *testing.T, us *store.UserStore, as *store.ArticleStore) {
				user := &model.User{Model: gorm.Model{ID: 1}, Username: "test_user"}
				us.On("GetByID", uint(1)).Return(user, nil)
			},
			request: &pb.CreateAritcleRequest{
				Article: &pb.CreateAritcleRequest_Article{
					Title:       "", // Invalid as title is required
					Description: "Test Description",
					Body:        "Test Body",
					TagList:     []string{"tag1", "tag2"},
				},
			},
			expectError: codes.InvalidArgument,
		},
		{
			name: "Article Store Creation Failure",
			setup: func(t *testing.T, us *store.UserStore, as *store.ArticleStore) {
				user := &model.User{Model: gorm.Model{ID: 1}, Username: "test_user"}
				us.On("GetByID", uint(1)).Return(user, nil)
				as.On("Create", mock.AnythingOfType("*model.Article")).Return(status.Error(codes.Canceled, "storage error"))
			},
			request: &pb.CreateAritcleRequest{
				Article: &pb.CreateAritcleRequest_Article{
					Title:       "Test Title",
					Description: "Test Description",
					Body:        "Test Body",
					TagList:     []string{"tag1", "tag2"},
				},
			},
			expectError: codes.Canceled,
		},
		// TODO: Complete the implementation for "Get Following Status Error"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			logger := zerolog.New(os.Stdout)
			us := new(mockUserStore)
			as := new(mockArticleStore)

			h := &Handler{
				logger: &logger,
				us:     us,
				as:     as,
			}

			tt.setup(t, us, as)

			// Mock the context for valid / invalid authentication
			ctx := context.Background()
			if tt.expectError != codes.Unauthenticated {
				ctx = context.WithValue(ctx, auth.ContextKeyUserID, uint(1))
			}

			// Act
			resp, err := h.CreateArticle(ctx, tt.request)

			// Assert
			if err != nil {
				st, _ := status.FromError(err)
				if tt.expectError != st.Code() {
					t.Errorf("expected error code %v, got %v", tt.expectError, st.Code())
				}
			} else {
				if tt.expectError != codes.OK {
					t.Errorf("expected error code %v, got OK", tt.expectError)
				} else {
					// Verify the response
					if resp.Article.Title != tt.request.Article.Title {
						t.Errorf("expected title %v, got %v", tt.request.Article.Title, resp.Article.Title)
					}
					// Further assertions can be added here
				}
			}
		})
	}
}
