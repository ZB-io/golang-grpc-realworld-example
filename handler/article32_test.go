package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)





type MockArticleStore struct {
	articles map[uint]*model.Article
	err      error
}
type MockUserStore struct {
	users map[uint]*model.User
	err   error
}


/*
ROOST_METHOD_HASH=DeleteArticle_0347183038
ROOST_METHOD_SIG_HASH=DeleteArticle_b2585946c3

FUNCTION_DEF=func (h *Handler) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.Empty, error) 

 */
func (m *MockArticleStore) Delete(article *model.Article) error {
	if m.err != nil {
		return m.err
	}
	delete(m.articles, article.ID)
	return nil
}

func (m *MockArticleStore) GetByID(id uint) (*model.Article, error) {
	if m.err != nil {
		return nil, m.err
	}
	article, exists := m.articles[id]
	if !exists {
		return nil, errors.New("article not found")
	}
	return article, nil
}

func TestHandlerDeleteArticle(t *testing.T) {

	logger := zerolog.New(nil)

	tests := []struct {
		name          string
		setupContext  func() context.Context
		setupMocks    func() (*MockUserStore, *MockArticleStore)
		input         *pb.DeleteArticleRequest
		expectedCode  codes.Code
		expectedError string
	}{
		{
			name: "Successful Article Deletion",
			setupContext: func() context.Context {
				ctx := context.Background()

				return ctx
			},
			setupMocks: func() (*MockUserStore, *MockArticleStore) {
				us := &MockUserStore{
					users: map[uint]*model.User{
						1: {ID: 1, Username: "testuser"},
					},
				}
				as := &MockArticleStore{
					articles: map[uint]*model.Article{
						1: {
							ID:     1,
							Author: model.User{ID: 1},
						},
					},
				}
				return us, as
			},
			input:        &pb.DeleteArticleRequest{Slug: "1"},
			expectedCode: codes.OK,
		},
		{
			name: "Unauthenticated Request",
			setupContext: func() context.Context {
				return context.Background()
			},
			setupMocks: func() (*MockUserStore, *MockArticleStore) {
				return &MockUserStore{}, &MockArticleStore{}
			},
			input:         &pb.DeleteArticleRequest{Slug: "1"},
			expectedCode:  codes.Unauthenticated,
			expectedError: "unauthenticated",
		},
		{
			name: "Invalid Slug Format",
			setupContext: func() context.Context {
				ctx := context.Background()

				return ctx
			},
			setupMocks: func() (*MockUserStore, *MockArticleStore) {
				us := &MockUserStore{
					users: map[uint]*model.User{
						1: {ID: 1},
					},
				}
				return us, &MockArticleStore{}
			},
			input:         &pb.DeleteArticleRequest{Slug: "invalid-slug"},
			expectedCode:  codes.InvalidArgument,
			expectedError: "invalid article id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := tt.setupContext()
			us, as := tt.setupMocks()

			h := &Handler{
				logger: &logger,
				us:     us,
				as:     as,
			}

			_, err := h.DeleteArticle(ctx, tt.input)

			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}

				status, ok := status.FromError(err)
				if !ok {
					t.Errorf("expected gRPC status error, got %v", err)
					return
				}

				if status.Code() != tt.expectedCode {
					t.Errorf("expected status code %v, got %v", tt.expectedCode, status.Code())
				}

				if status.Message() != tt.expectedError {
					t.Errorf("expected error message %q, got %q", tt.expectedError, status.Message())
				}
			}
		})
	}
}

