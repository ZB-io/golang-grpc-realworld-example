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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type mockUserStore struct {
	getByIDFunc func(uint) (*model.User, error)
}

func (m *mockUserStore) GetByID(id uint) (*model.User, error) {
	return m.getByIDFunc(id)
}

type mockArticleStore struct {
	getByIDFunc       func(uint) (*model.Article, error)
	getCommentByIDFunc func(uint) (*model.Comment, error)
	getCommentsFn     func(*model.Article) ([]model.Comment, error)
	deleteCommentFunc func(*model.Comment) error
	createCommentFunc func(*model.Comment) error
}

func (m *mockArticleStore) GetByID(id uint) (*model.Article, error) {
	return m.getByIDFunc(id)
}

func (m *mockArticleStore) GetCommentByID(id uint) (*model.Comment, error) {
	return m.getCommentByIDFunc(id)
}

func (m *mockArticleStore) GetComments(article *model.Article) ([]model.Comment, error) {
	return m.getCommentsFn(article)
}

func (m *mockArticleStore) DeleteComment(comment *model.Comment) error {
	return m.deleteCommentFunc(comment)
}

func (m *mockArticleStore) CreateComment(comment *model.Comment) error {
	return m.createCommentFunc(comment)
}

/*
ROOST_METHOD_HASH=DeleteComment_452af2f984
ROOST_METHOD_SIG_HASH=DeleteComment_27615e7d69
*/
func TestDeleteComment(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*mockUserStore, *mockArticleStore)
		req            *pb.DeleteCommentRequest
		expectedError  error
		expectedStatus codes.Code
	}{
		{
			name: "Successfully Delete a Comment",
			setupMocks: func(us *mockUserStore, as *mockArticleStore) {
				us.getByIDFunc = func(uint) (*model.User, error) {
					return &model.User{}, nil
				}
				as.getCommentByIDFunc = func(uint) (*model.Comment, error) {
					return &model.Comment{UserID: 1, ArticleID: 1}, nil
				}
				as.deleteCommentFunc = func(*model.Comment) error {
					return nil
				}
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
			setupMocks: func(us *mockUserStore, as *mockArticleStore) {
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
			setupMocks: func(us *mockUserStore, as *mockArticleStore) {
				us.getByIDFunc = func(uint) (*model.User, error) {
					return &model.User{}, nil
				}
				as.getCommentByIDFunc = func(uint) (*model.Comment, error) {
					return nil, errors.New("comment not found")
				}
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
			mockUS := &mockUserStore{}
			mockAS := &mockArticleStore{}
			tt.setupMocks(mockUS, mockAS)

			h := &Handler{
				logger: &zerolog.Logger{},
				us:     mockUS,
				as:     mockAS,
			}

			ctx := context.Background()
			if tt.name != "Attempt to Delete a Comment with Unauthenticated User" {
				ctx = context.WithValue(ctx, auth.UserIDKey, uint(1))
			}

			_, err := h.DeleteComment(ctx, tt.req)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else {
					statusErr, ok := status.FromError(err)
					if !ok {
						t.Errorf("expected gRPC status error, got %v", err)
					} else if statusErr.Code() != tt.expectedStatus {
						t.Errorf("expected status code %v, got %v", tt.expectedStatus, statusErr.Code())
					}
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
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
		setupMocks     func(*mockArticleStore, *mockUserStore)
		req            *pb.GetCommentsRequest
		ctx            context.Context
		expectedResult *pb.CommentsResponse
		expectedError  error
	}{
		{
			name: "Successfully retrieve comments for a valid article",
			setupMocks: func(as *mockArticleStore, us *mockUserStore) {
				as.getByIDFunc = func(uint) (*model.Article, error) {
					return &model.Article{Model: gorm.Model{ID: 1}}, nil
				}
				as.getCommentsFn = func(*model.Article) ([]model.Comment, error) {
					return []model.Comment{
						{Model: gorm.Model{ID: 1}, Body: "Comment 1", Author: model.User{Model: gorm.Model{ID: 1}, Username: "user1"}},
						{Model: gorm.Model{ID: 2}, Body: "Comment 2", Author: model.User{Model: gorm.Model{ID: 2}, Username: "user2"}},
					}, nil
				}
				us.getByIDFunc = func(uint) (*model.User, error) {
					return &model.User{Model: gorm.Model{ID: 1}}, nil
				}
			},
			req: &pb.GetCommentsRequest{Slug: "1"},
			ctx: context.WithValue(context.Background(), auth.UserIDKey, uint(1)),
			expectedResult: &pb.CommentsResponse{
				Comments: []*pb.Comment{
					{Id: "1", Body: "Comment 1", Author: &pb.Profile{Username: "user1"}},
					{Id: "2", Body: "Comment 2", Author: &pb.Profile{Username: "user2"}},
				},
			},
			expectedError: nil,
		},
		{
			name:           "Attempt to retrieve comments with an invalid slug",
			setupMocks:     func(as *mockArticleStore, us *mockUserStore) {},
			req:            &pb.GetCommentsRequest{Slug: "not-a-number"},
			ctx:            context.Background(),
			expectedResult: nil,
			expectedError:  status.Error(codes.InvalidArgument, "invalid article id"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArticleStore := &mockArticleStore{}
			mockUserStore := &mockUserStore{}
			tt.setupMocks(mockArticleStore, mockUserStore)

			h := &Handler{
				logger: &zerolog.Logger{},
				as:     mockArticleStore,
				us:     mockUserStore,
			}

			result, err := h.GetComments(tt.ctx, tt.req)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, but got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error %v, but got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.expectedResult != nil {
				if result == nil {
					t.Error("Expected non-nil result, but got nil")
				} else if len(result.Comments) != len(tt.expectedResult.Comments) {
					t.Errorf("Expected %d comments, but got %d", len(tt.expectedResult.Comments), len(result.Comments))
				} else {
					for i, expectedComment := range tt.expectedResult.Comments {
						if result.Comments[i].Id != expectedComment.Id {
							t.Errorf("Expected comment ID %s, but got %s", expectedComment.Id, result.Comments[i].Id)
						}
						if result.Comments[i].Body != expectedComment.Body {
							t.Errorf("Expected comment body %s, but got %s", expectedComment.Body, result.Comments[i].Body)
						}
						if result.Comments[i].Author.Username != expectedComment.Author.Username {
							t.Errorf("Expected author username %s, but got %s", expectedComment.Author.Username, result.Comments[i].Author.Username)
						}
					}
				}
			} else if result != nil {
				t.Error("Expected nil result, but got non-nil")
			}
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
		setupMocks     func(*mockUserStore, *mockArticleStore)
		req            *pb.CreateCommentRequest
		expectedResult *pb.CommentResponse
		expectedError  error
	}{
		{
			name: "Successfully Create a Comment",
			setupMocks: func(us *mockUserStore, as *mockArticleStore) {
				us.getByIDFunc = func(id uint) (*model.User, error) {
					return &model.User{Username: "testuser"}, nil
				}
				as.getByIDFunc = func(id uint) (*model.Article, error) {
					return &model.Article{}, nil
				}
				as.createCommentFunc = func(comment *model.Comment) error {
					return nil
				}
			},
			req: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.CreateCommentRequest_Comment{
					Body: "Test comment",
				},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUS := &mockUserStore{}
			mockAS := &mockArticleStore{}

			tt.setupMocks(mockUS, mockAS)

			h := &Handler{
				logger: &zerolog.Logger{},
				us:     mockUS,
				as:     mockAS,
			}

			origGetUserID := auth.GetUserID
			defer func() { auth.GetUserID = origGetUserID }()
			auth.GetUserID = func(ctx context.Context) (uint, error) {
				if tt.name == "Attempt to Create Comment with Unauthenticated User" {
					return 0, errors.New("unauthenticated")
				}
				return 1, nil
			}

			result, err := h.CreateComment(context.Background(), tt.req)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, but got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error %v, but got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got %v", err)
				}
				if result == nil {
					t.Error("Expected non-nil result, but got nil")
				} else {
					if result.Comment.Body != tt.expectedResult.Comment.Body {
						t.Errorf("Expected comment body %s, but got %s", tt.expectedResult.Comment.Body, result.Comment.Body)
					}
					if result.Comment.Author.Username != tt.expectedResult.Comment.Author.Username {
						t.Errorf("Expected author username %s, but got %s", tt.expectedResult.Comment.Author.Username, result.Comment.Author.Username)
					}
				}
			}
		})
	}
}
