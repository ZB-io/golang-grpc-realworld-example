package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	auth "github.com/raahii/golang-grpc-realworld-example/auth"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)

// Mock implementations and auxiliary structures for testing

type Handler struct {
	us *mockUserService
	as *mockArticleService
	logger mockLogger
}

type mockUserService struct {
	// Simulated methods
}

func (m *mockUserService) GetByID(id uint) (*model.User, error) {
	// Simulated logic
	return &model.User{ID: id}, nil
}

type mockArticleService struct {
	// Simulated methods
}

func (m *mockArticleService) GetByID(id uint) (*model.Article, error) {
	// Simulated logic
	return &model.Article{ID: id}, nil
}

func (m *mockArticleService) CreateComment(comment *model.Comment) error {
	// Simulated logic
	return nil
}

type mockLogger struct {}

func (m mockLogger) Info() *mockLogger { return &m }
func (m *mockLogger) Error() *mockLogger { return m }
func (m *mockLogger) Msgf(format string, args ...interface{}) {}
func (m *mockLogger) Msg(msg string) {}
func (m *mockLogger) Err(err error) *mockLogger { return m }

// TestCreateComment tests the CreateComment function for various scenarios
func TestCreateComment(t *testing.T) {
	h := &Handler{
		us: &mockUserService{},
		as: &mockArticleService{},
		logger: mockLogger{},
	}

	type args struct {
		ctx context.Context
		req *pb.CreateCommentRequest
	}

	tests := []struct {
		name       string
		args       args
		want       *pb.CommentResponse
		wantErr    error
		setupMocks func()
	}{
		{
			name: "Successful Comment Creation",
			args: args{
				ctx: context.WithValue(context.Background(), &auth.UserIDKey{}, uint(1)),
				req: &pb.CreateCommentRequest{
					Slug: "1",
					Comment: &pb.Comment{
						Body: "A comment",
					},
				},
			},
			want: &pb.CommentResponse{
				Comment: &pb.Comment{
					Body: "A comment",
					// Fields like author, created at etc., omitted for brevity
				},
			},
			setupMocks: func() {
				// No setup needed for success path
			},
		},
		{
			name: "Invalid User Authentication",
			args: args{
				ctx: context.Background(),
				req: &pb.CreateCommentRequest{
					Slug: "2",
					Comment: &pb.Comment{
						Body: "Another comment",
					},
				},
			},
			wantErr: status.Errorf(codes.Unauthenticated, "unauthenticated"),
			setupMocks: func() {
				// Mock error in authentication service
			},
		},
		{
			name: "Non-Existent User",
			args: args{
				ctx: context.WithValue(context.Background(), &auth.UserIDKey{}, uint(999)),
				req: &pb.CreateCommentRequest{
					Slug: "3",
					Comment: &pb.Comment{
						Body: "Yet another comment",
					},
				},
			},
			wantErr: status.Error(codes.NotFound, "user not found"),
			setupMocks: func() {
				h.us = &mockUserService{
					// Simulate user not found
				}
			},
		},
		{
			name: "Invalid Article Slug",
			args: args{
				ctx: context.WithValue(context.Background(), &auth.UserIDKey{}, uint(1)),
				req: &pb.CreateCommentRequest{
					Slug: "invalid-slug",
					Comment: &pb.Comment{
						Body: "Invalid slug comment",
					},
				},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid article id"),
			setupMocks: func() {
				// No setup needed for slug conversion error
			},
		},
		{
			name: "Non-Existent Article",
			args: args{
				ctx: context.WithValue(context.Background(), &auth.UserIDKey{}, uint(1)),
				req: &pb.CreateCommentRequest{
					Slug: "9999",
					Comment: &pb.Comment{
						Body: "Non-existent article comment",
					},
				},
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid article id"),
			setupMocks: func() {
				h.as = &mockArticleService{
					// Simulate article not found
				}
			},
		},
		{
			name: "Comment Validation Error",
			args: args{
				ctx: context.WithValue(context.Background(), &auth.UserIDKey{}, uint(1)),
				req: &pb.CreateCommentRequest{
					Slug: "4",
					Comment: &pb.Comment{
						Body: "",
					},
				},
			},
			wantErr: status.Error(codes.InvalidArgument, "validation error: invalid comment body"),
			setupMocks: func() {
				// No setup needed for validation error
			},
		},
		{
			name: "Comment Creation Failure",
			args: args{
				ctx: context.WithValue(context.Background(), &auth.UserIDKey{}, uint(1)),
				req: &pb.CreateCommentRequest{
					Slug: "5",
					Comment: &pb.Comment{
						Body: "This will fail to create",
					},
				},
			},
			wantErr: status.Error(codes.Aborted, "failed to create comment."),
			setupMocks: func() {
				h.as = &mockArticleService{
					// Simulate create comment error
					CreateComment: func(comment *model.Comment) error {
						return errors.New("database error")
					},
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			got, err := h.CreateComment(tt.args.ctx, tt.args.req)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				stat, _ := status.FromError(err)
				assert.Equal(t, tt.wantErr.Error(), stat.Message())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
