package handler

import (
	"context"
	"testing"

	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockArticleStore struct {
	GetByIDFunc      func(uint) (*model.Article, error)
	GetCommentsFunc  func(*model.Article) ([]model.Comment, error)
	CreateCommentFunc func(*model.Comment) error
	GetCommentByIDFunc func(uint) (*model.Comment, error)
	DeleteCommentFunc func(*model.Comment) error
}

type MockUserStore struct {
	GetByIDFunc     func(uint) (*model.User, error)
	IsFollowingFunc func(*model.User, *model.User) (bool, error)
}

type MockAuth struct {
	GetUserIDFunc func(context.Context) (uint, error)
}

func (m *MockArticleStore) GetByID(id uint) (*model.Article, error) {
	return m.GetByIDFunc(id)
}

func (m *MockArticleStore) GetComments(article *model.Article) ([]model.Comment, error) {
	return m.GetCommentsFunc(article)
}

func (m *MockArticleStore) CreateComment(comment *model.Comment) error {
	return m.CreateCommentFunc(comment)
}

func (m *MockArticleStore) GetCommentByID(id uint) (*model.Comment, error) {
	return m.GetCommentByIDFunc(id)
}

func (m *MockArticleStore) DeleteComment(comment *model.Comment) error {
	return m.DeleteCommentFunc(comment)
}

func (m *MockUserStore) GetByID(id uint) (*model.User, error) {
	return m.GetByIDFunc(id)
}

func (m *MockUserStore) IsFollowing(follower, followed *model.User) (bool, error) {
	return m.IsFollowingFunc(follower, followed)
}

func (m *MockAuth) GetUserID(ctx context.Context) (uint, error) {
	return m.GetUserIDFunc(ctx)
}

/*
ROOST_METHOD_HASH=DeleteComment_452af2f984
ROOST_METHOD_SIG_HASH=DeleteComment_27615e7d69
*/
func TestDeleteComment(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*MockUserStore, *MockArticleStore, *MockAuth)
		req            *pb.DeleteCommentRequest
		expectedError  error
		expectedResult *pb.Empty
	}{
		{
			name: "Successfully Delete a Comment",
			setupMocks: func(us *MockUserStore, as *MockArticleStore, ma *MockAuth) {
				us.GetByIDFunc = func(uint) (*model.User, error) {
					return &model.User{ID: 1}, nil
				}
				as.GetCommentByIDFunc = func(uint) (*model.Comment, error) {
					return &model.Comment{UserID: 1, ArticleID: 1}, nil
				}
				as.DeleteCommentFunc = func(*model.Comment) error {
					return nil
				}
				ma.GetUserIDFunc = func(context.Context) (uint, error) {
					return 1, nil
				}
			},
			req: &pb.DeleteCommentRequest{
				Slug: "1",
				Id:   "1",
			},
			expectedError:  nil,
			expectedResult: &pb.Empty{},
		},
		{
			name: "Unauthenticated User Attempt",
			setupMocks: func(us *MockUserStore, as *MockArticleStore, ma *MockAuth) {
				ma.GetUserIDFunc = func(context.Context) (uint, error) {
					return 0, status.Error(codes.Unauthenticated, "unauthenticated")
				}
			},
			req: &pb.DeleteCommentRequest{
				Slug: "1",
				Id:   "1",
			},
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zerolog.Nop()
			us := &MockUserStore{}
			as := &MockArticleStore{}
			ma := &MockAuth{}
			tt.setupMocks(us, as, ma)

			h := &Handler{
				logger: &logger,
				us:     us,
				as:     as,
				auth:   ma,
			}

			result, err := h.DeleteComment(context.Background(), tt.req)

			if err != nil {
				if tt.expectedError == nil {
					t.Errorf("expected no error, got %v", err)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if tt.expectedError != nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				}
			}

			if result != nil && tt.expectedResult == nil {
				t.Errorf("expected nil result, got %v", result)
			} else if result == nil && tt.expectedResult != nil {
				t.Errorf("expected result %v, got nil", tt.expectedResult)
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
		setupMocks     func(*MockArticleStore, *MockUserStore, *MockAuth)
		req            *pb.GetCommentsRequest
		ctx            context.Context
		expectedResp   *pb.CommentsResponse
		expectedErrMsg string
		expectedCode   codes.Code
	}{
		{
			name: "Successfully retrieve comments for a valid article",
			setupMocks: func(mas *MockArticleStore, mus *MockUserStore, ma *MockAuth) {
				mas.GetByIDFunc = func(uint) (*model.Article, error) {
					return &model.Article{ID: 1}, nil
				}
				mas.GetCommentsFunc = func(*model.Article) ([]model.Comment, error) {
					return []model.Comment{
						{ID: 1, Body: "Comment 1", Author: model.User{Username: "user1"}},
						{ID: 2, Body: "Comment 2", Author: model.User{Username: "user2"}},
					}, nil
				}
				mus.GetByIDFunc = func(uint) (*model.User, error) {
					return &model.User{ID: 1}, nil
				}
				mus.IsFollowingFunc = func(follower, followed *model.User) (bool, error) {
					return followed.Username == "user1", nil
				}
				ma.GetUserIDFunc = func(context.Context) (uint, error) {
					return 1, nil
				}
			},
			req: &pb.GetCommentsRequest{Slug: "1"},
			ctx: context.Background(),
			expectedResp: &pb.CommentsResponse{
				Comments: []*pb.Comment{
					{Id: "1", Body: "Comment 1", Author: &pb.Profile{Username: "user1", Following: true}},
					{Id: "2", Body: "Comment 2", Author: &pb.Profile{Username: "user2", Following: false}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArticleStore := &MockArticleStore{}
			mockUserStore := &MockUserStore{}
			mockAuth := &MockAuth{}
			tt.setupMocks(mockArticleStore, mockUserStore, mockAuth)

			h := &Handler{
				logger: zerolog.Nop(),
				as:     mockArticleStore,
				us:     mockUserStore,
				auth:   mockAuth,
			}

			resp, err := h.GetComments(tt.ctx, tt.req)

			if tt.expectedErrMsg != "" {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("expected gRPC status error, got %v", err)
				}
				if st.Code() != tt.expectedCode {
					t.Errorf("expected error code %v, got %v", tt.expectedCode, st.Code())
				}
				if st.Message() != tt.expectedErrMsg {
					t.Errorf("expected error message %q, got %q", tt.expectedErrMsg, st.Message())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if resp == nil {
					t.Fatalf("expected non-nil response, got nil")
				}
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
		setupMocks     func(*MockUserStore, *MockArticleStore, *MockAuth)
		req            *pb.CreateCommentRequest
		expectedResult *pb.CommentResponse
		expectedError  error
	}{
		{
			name: "Successfully Create a Comment",
			setupMocks: func(us *MockUserStore, as *MockArticleStore, ma *MockAuth) {
				us.GetByIDFunc = func(id uint) (*model.User, error) {
					return &model.User{Username: "testuser"}, nil
				}
				as.GetByIDFunc = func(id uint) (*model.Article, error) {
					return &model.Article{}, nil
				}
				as.CreateCommentFunc = func(comment *model.Comment) error {
					return nil
				}
				ma.GetUserIDFunc = func(ctx context.Context) (uint, error) {
					return 1, nil
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
			mockUS := &MockUserStore{}
			mockAS := &MockArticleStore{}
			mockAuth := &MockAuth{}
			tt.setupMocks(mockUS, mockAS, mockAuth)

			h := &Handler{
				logger: zerolog.Nop(),
				us:     mockUS,
				as:     mockAS,
				auth:   mockAuth,
			}

			result, err := h.CreateComment(context.Background(), tt.req)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Error("expected non-nil result, got nil")
				} else {
					if result.Comment.Body != tt.expectedResult.Comment.Body {
						t.Errorf("expected comment body %s, got %s", tt.expectedResult.Comment.Body, result.Comment.Body)
					}
					if result.Comment.Author.Username != tt.expectedResult.Comment.Author.Username {
						t.Errorf("expected author username %s, got %s", tt.expectedResult.Comment.Author.Username, result.Comment.Author.Username)
					}
				}
			}
		})
	}
}
